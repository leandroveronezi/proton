package proton

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type h = map[string]interface{}

// Result is a struct for the resulting value of the JS expression or an error.
type result struct {
	Value json.RawMessage
	Err   error
}

type bindingFunc func(args []json.RawMessage) (interface{}, error)

// Msg is a struct for incoming messages (results and async events)
type msg struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  json.RawMessage `json:"error"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Browser struct {
	config Config
	done   chan struct{}
	sync.Mutex
	cmd      *exec.Cmd
	ws       *websocket.Conn
	id       int32
	target   string
	session  string
	window   int
	pending  map[int]chan result
	bindings map[string]bindingFunc
}

func (_this *Browser) findTarget() (string, error) {
	err := websocket.JSON.Send(_this.ws, h{
		"id": 0, "method": "Target.setDiscoverTargets", "params": h{"discover": true},
	})
	if err != nil {
		return "", err
	}
	for {
		m := msg{}
		if err = websocket.JSON.Receive(_this.ws, &m); err != nil {
			return "", err
		} else if m.Method == "Target.targetCreated" {
			target := struct {
				TargetInfo struct {
					Type string `json:"type"`
					ID   string `json:"targetId"`
				} `json:"targetInfo"`
			}{}
			if err := json.Unmarshal(m.Params, &target); err != nil {
				return "", err
			} else if target.TargetInfo.Type == "page" {
				return target.TargetInfo.ID, nil
			}
		}
	}
}

func (_this *Browser) startSession(target string) (string, error) {
	err := websocket.JSON.Send(_this.ws, h{
		"id": 1, "method": "Target.attachToTarget", "params": h{"targetId": target},
	})
	if err != nil {
		return "", err
	}
	for {
		m := msg{}
		if err = websocket.JSON.Receive(_this.ws, &m); err != nil {
			return "", err
		} else if m.ID == 1 {
			if m.Error != nil {
				return "", errors.New("Target error: " + string(m.Error))
			}
			session := struct {
				ID string `json:"sessionId"`
			}{}
			if err := json.Unmarshal(m.Result, &session); err != nil {
				return "", err
			}
			return session.ID, nil
		}
	}
}

func (_this *Browser) getWindowForTarget(target string) (windowTargetMessage, error) {
	var m windowTargetMessage
	msg, err := _this.send("Browser.getWindowForTarget", h{"targetId": target})
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(msg, &m)
	return m, err
}

func (_this *Browser) readLoop() {
	for {
		m := msg{}
		if err := websocket.JSON.Receive(_this.ws, &m); err != nil {
			return
		}

		if m.Method == "Target.receivedMessageFromTarget" {
			params := struct {
				SessionID string `json:"sessionId"`
				Message   string `json:"message"`
			}{}
			json.Unmarshal(m.Params, &params)
			if params.SessionID != _this.session {
				continue
			}
			res := targetMessage{}
			json.Unmarshal([]byte(params.Message), &res)

			if res.ID == 0 && res.Method == "Runtime.consoleAPICalled" || res.Method == "Runtime.exceptionThrown" {

				if _this.config.Debug {
					log.Println(params.Message)
				}

			} else if res.ID == 0 && res.Method == "Runtime.bindingCalled" {
				payload := struct {
					Name string            `json:"name"`
					Seq  int               `json:"seq"`
					Args []json.RawMessage `json:"args"`
				}{}
				json.Unmarshal([]byte(res.Params.Payload), &payload)

				_this.Lock()
				binding, ok := _this.bindings[res.Params.Name]
				_this.Unlock()
				if ok {
					jsString := func(v interface{}) string { b, _ := json.Marshal(v); return string(b) }
					go func() {
						result, error := "", `""`
						if r, err := binding(payload.Args); err != nil {
							error = jsString(err.Error())
						} else if b, err := json.Marshal(r); err != nil {
							error = jsString(err.Error())
						} else {
							result = string(b)
						}
						expr := fmt.Sprintf(`
							if (%[4]s) {
								window['%[1]s']['errors'].get(%[2]d)(%[4]s);
							} else {
								window['%[1]s']['callbacks'].get(%[2]d)(%[3]s);
							}
							window['%[1]s']['callbacks'].delete(%[2]d);
							window['%[1]s']['errors'].delete(%[2]d);
							`, payload.Name, payload.Seq, result, error)
						_this.send("Runtime.evaluate", h{"expression": expr, "contextId": res.Params.ID})
					}()
				}
				continue
			}

			_this.Lock()
			resc, ok := _this.pending[res.ID]
			delete(_this.pending, res.ID)
			_this.Unlock()

			if !ok {
				continue
			}

			if res.Error.Message != "" {
				resc <- result{Err: errors.New(res.Error.Message)}
			} else if res.Result.Exception.Exception.Value != nil {
				resc <- result{Err: errors.New(string(res.Result.Exception.Exception.Value))}
			} else if res.Result.Result.Type == "object" && res.Result.Result.Subtype == "error" {
				resc <- result{Err: errors.New(res.Result.Result.Description)}
			} else if res.Result.Result.Type != "" {
				resc <- result{Value: res.Result.Result.Value}
			} else {
				res := targetMessageTemplate{}
				json.Unmarshal([]byte(params.Message), &res)
				resc <- result{Value: res.Result}
			}
		} else if m.Method == "Target.targetDestroyed" {
			params := struct {
				TargetID string `json:"targetId"`
			}{}
			json.Unmarshal(m.Params, &params)
			if params.TargetID == _this.target {
				_this.kill(true)
				return
			}
		}
	}
}

func (_this *Browser) send(method string, params h) (json.RawMessage, error) {
	id := atomic.AddInt32(&_this.id, 1)
	b, err := json.Marshal(h{"id": int(id), "method": method, "params": params})
	if err != nil {
		return nil, err
	}
	resc := make(chan result)
	_this.Lock()
	_this.pending[int(id)] = resc
	_this.Unlock()

	if _this.config.Debug {
		log.Println(string(b))
	}

	if err := websocket.JSON.Send(_this.ws, h{
		"id":     int(id),
		"method": "Target.sendMessageToTarget",
		"params": h{"message": string(b), "sessionId": _this.session},
	}); err != nil {
		return nil, err
	}
	res := <-resc
	return res.Value, res.Err
}

func (_this *Browser) bind(name string, f bindingFunc) error {
	_this.Lock()
	// check if binding already exists
	_, exists := _this.bindings[name]

	_this.bindings[name] = f
	_this.Unlock()

	if exists {
		// Just replace callback and return, as the binding was already added to js
		// and adding it again would break it.
		return nil
	}

	if _, err := _this.send("Runtime.addBinding", h{"name": name}); err != nil {
		return err
	}
	script := fmt.Sprintf(`(() => {
	const bindingName = '%s';
	const binding = window[bindingName];
	window[bindingName] = async (...args) => {
		const me = window[bindingName];
		let errors = me['errors'];
		let callbacks = me['callbacks'];
		if (!callbacks) {
			callbacks = new Map();
			me['callbacks'] = callbacks;
		}
		if (!errors) {
			errors = new Map();
			me['errors'] = errors;
		}
		const seq = (me['lastSeq'] || 0) + 1;
		me['lastSeq'] = seq;
		const promise = new Promise((resolve, reject) => {
			callbacks.set(seq, resolve);
			errors.set(seq, reject);
		});
		binding(JSON.stringify({name: bindingName, seq, args}));
		return promise;
	}})();
	`, name)
	_, err := _this.send("Page.addScriptToEvaluateOnNewDocument", h{"source": script})
	if err != nil {
		return err
	}

	awaitPromise := true
	returnByValue := true

	_, err = _this.RuntimeEvaluate(RuntimeEvaluateParameters{Expression: script, AwaitPromise: &awaitPromise, ReturnByValue: &returnByValue})

	return err
}

func (_this *Browser) setBounds(b Bounds) error {
	if b.WindowState == "" {
		b.WindowState = WindowStateNormal
	}
	param := h{"windowId": _this.window, "bounds": b}
	if b.WindowState != WindowStateNormal {
		param["bounds"] = h{"windowState": b.WindowState}
	}
	_, err := _this.send("Browser.setWindowBounds", param)
	return err
}

func (_this *Browser) bounds() (Bounds, error) {
	result, err := _this.send("Browser.getWindowBounds", h{"windowId": _this.window})
	if err != nil {
		return Bounds{}, err
	}
	bounds := struct {
		Bounds Bounds `json:"bounds"`
	}{}
	err = json.Unmarshal(result, &bounds)
	return bounds.Bounds, err
}

func (_this *Browser) kill(exited bool) error {

	if _this.ws != nil {

		if err := _this.ws.Close(); err != nil {
			return err
		}

	}

	if exited {
		return nil
	}

	// TODO: cancel all pending requests
	if state := _this.cmd.ProcessState; state == nil || !state.Exited() {

		sig := os.Interrupt

		if runtime.GOOS == "windows" {
			sig = os.Kill
		}

		if err := _this.cmd.Process.Signal(sig); err != nil {
			return err
		}

		return _this.cmd.Process.Kill()
	}

	return nil
}

func readUntilMatch(r io.ReadCloser, re *regexp.Regexp) ([]string, error) {
	br := bufio.NewReader(r)
	for {
		if line, err := br.ReadString('\n'); err != nil {
			r.Close()
			return nil, err
		} else if m := re.FindStringSubmatch(line); m != nil {
			go io.Copy(ioutil.Discard, br)
			return m, nil
		}
	}
}

func contains(arr []string, x string) bool {
	for _, n := range arr {
		if x == n {
			return true
		}
	}
	return false
}

func structToMap(s interface{}) h {

	result, err := json.Marshal(s)

	if err != nil {
		return h{}
	}

	aux := h{}

	err = json.Unmarshal(result, &aux)

	if err != nil {
		return h{}
	}

	for key, val := range aux {

		if val == nil {
			delete(aux, key)
		}

	}

	return aux
}

func (_this *Browser) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	// f must return either value and error or just error
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}

	return _this.bind(name, func(raw []json.RawMessage) (interface{}, error) {
		if len(raw) != v.Type().NumIn() {
			return nil, errors.New("function arguments mismatch")
		}
		args := []reflect.Value{}
		for i := range raw {
			arg := reflect.New(v.Type().In(i))
			if err := json.Unmarshal(raw[i], arg.Interface()); err != nil {
				return nil, err
			}
			args = append(args, arg.Elem())
		}
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		res := v.Call(args)
		switch len(res) {
		case 0:
			// No results from the function, just return nil
			return nil, nil
		case 1:
			// One result may be a value, or an error
			if res[0].Type().Implements(errorType) {
				if res[0].Interface() != nil {
					return nil, res[0].Interface().(error)
				}
				return nil, nil
			}
			return res[0].Interface(), nil
		case 2:
			// Two results: first one is value, second is error
			if !res[1].Type().Implements(errorType) {
				return nil, errors.New("second return value must be an error")
			}
			if res[1].Interface() == nil {
				return res[0].Interface(), nil
			}
			return res[0].Interface(), res[1].Interface().(error)
		default:
			return nil, errors.New("unexpected number of return values")
		}
	})
}

func (_this *Browser) Eval(js string) Value {

	awaitPromise := true
	returnByValue := true

	v, err := _this.RuntimeEvaluate(RuntimeEvaluateParameters{Expression: js, AwaitPromise: &awaitPromise, ReturnByValue: &returnByValue})
	return value{err: err, raw: v}
}

func (_this *Browser) SetBounds(b Bounds) error {
	return _this.setBounds(b)
}

func (_this *Browser) Bounds() (Bounds, error) {
	return _this.bounds()
}

func (_this *Browser) browserBinary() string {

	if _this.config.BrowserBinary != "" {
		return _this.config.BrowserBinary
	}

	if _this.config.Flavor == Chrome {
		return _this.browserBinaryChrome()
	}

	return _this.browserBinaryEdge()

}

func (_this *Browser) browserBinaryEdge() string {

	var paths []string

	switch runtime.GOOS {
	case "darwin":

		paths = []string{
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}

	case "windows":

		paths = []string{
			os.Getenv("LocalAppData") + "/Microsoft/Edge/Application/msedge.exe",
			os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe",
			os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		}

	default:

		paths = []string{
			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-dev",
		}

	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}

	return ""

}

func (_this *Browser) browserBinaryChrome() string {

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	case "windows":
		paths = []string{
			"C:/Program Files/Google/Chrome/Application/chrome.exe",
			os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe",
		}
	default:
		paths = []string{
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""

}

func (_this *Browser) Done() <-chan struct{} {
	return _this.done
}

func (_this *Browser) Close() error {

	// ignore err, as the chrome process might be already dead, when user close the window.
	_this.kill(false)
	<-_this.done

	if !_this.config.UserDataDirKeep {
		if err := os.RemoveAll(_this.config.UserDataDir); err != nil {
			return err
		}
	}

	return nil
}

func (_this *Browser) genEmptyHtml() string {

	template := `<!DOCTYPE html><html><head><title>{{title}}</title></head><body></body></html>`

	/* * /
	template = `
		<!DOCTYPE html>
		<html>
			<head>
		    <meta charset="utf-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1">
		    <title>{{title}}</title>

		    <style>
		        body, html {
		            margin: 0;
		            padding: 0;
		        }
		        html,body{height:100%;overflow:hidden}body{background:linear-gradient(to left, #141E30, #243B55);transform:scale(1.2, 1.2)}body>div{border-radius:50%;border:1px solid #fff;transform-style:preserve-3d;transform:rotateX(80deg) rotateY(20deg);position:absolute;left:50%;top:50%;margin-left:-100px;margin-top:-100px}body>div:first-of-type:after{content:"";position:absolute;height:40px;width:40px;background:#fff;border-radius:50%;transform:rotateX(-80deg) rotateY(0);box-shadow:0 0 25px #fff;animation:nucleus_ 2s infinite linear;left:50%;top:50%;margin-top:-20px;margin-left:-20px}body>div:nth-of-type(2){transform:rotateX(-80deg) rotateY(20deg)}body>div:nth-of-type(2)>div,body>div:nth-of-type(2)>div:after{animation-delay:-.5s}body>div:nth-of-type(3){transform:rotateX(-70deg) rotateY(60deg)}body>div:nth-of-type(3)>div,body>div:nth-of-type(3)>div:after{animation-delay:-1s}body>div:nth-of-type(4){transform:rotateX(70deg) rotateY(60deg)}body>div:nth-of-type(4)>div,body>div:nth-of-type(4)>div:after{animation-delay:-1.5s}body>div>div{width:200px;height:200px;position:relative;transform-style:preserve-3d;animation:trail_ 2s infinite linear}body>div>div:after{content:"";position:absolute;top:-5px;box-shadow:0 0 12px #fff;left:50%;margin-left:-5px;width:5px;height:5px;border-radius:50%;background-color:#fff;animation:particle_ 2s infinite linear}@keyframes trail_{from{transform:rotateZ(0deg)}to{transform:rotateZ(360deg)}}@keyframes particle_{from{transform:rotateX(90deg) rotateY(0deg)}to{transform:rotateX(90deg) rotateY(-360deg)}}@keyframes nucleus_{0%, 100%{box-shadow:0 0 0 transparent}50%{box-shadow:0 0 25px #fff}}
		    </style>
		    </head>
		<body>
			<div>
			  <div></div>
			</div>
			<div>
			  <div></div>
			</div>
			<div>
			  <div></div>
			</div>
			<div>
			  <div></div>
			</div>
		</body></html>
		`
	/**/

	template = strings.ReplaceAll(template, "{{title}}", _this.config.Title)

	template = "data:text/html;base64," + base64.StdEncoding.EncodeToString([]byte(template))

	return template

}
