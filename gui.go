package proton

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type gui struct {
	done    chan struct{}
	config  Config
	browser *browser
}

func (_this *gui) genEmptyHtml() string {

	template := `data:text/html,<!DOCTYPE html><html><head><title>{{title}}</title></head><body></body></html>`

	template = strings.ReplaceAll(template, "{{title}}", _this.config.Title)

	return template

}

func (_this *gui) Navigate(url string) error {
	return _this.browser.navigate(url)
}

func (_this *gui) Run() error {

	chrome, err := _this.newBrowser()
	done := make(chan struct{})
	if err != nil {
		return err
	}

	go func() {
		chrome.cmd.Wait()
		close(done)
	}()

	_this.browser = chrome
	_this.done = done

	return nil

}

func (_this *gui) Done() <-chan struct{} {
	return _this.done
}

func (_this *gui) Close() error {
	// ignore err, as the chrome process might be already dead, when user close the window.
	_this.browser.kill()
	<-_this.done

	if !_this.config.UserDataDirKeep {
		if err := os.RemoveAll(_this.config.UserDataDir); err != nil {
			return err
		}
	}

	return nil
}

func (_this *gui) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	// f must return either value and error or just error
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}

	return _this.browser.bind(name, func(raw []json.RawMessage) (interface{}, error) {
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

func (_this *gui) Eval(js string) Value {
	v, err := _this.browser.eval(js)
	return value{err: err, raw: v}
}

func (_this *gui) SetBounds(b Bounds) error {
	return _this.browser.setBounds(b)
}

func (_this *gui) Bounds() (Bounds, error) {
	return _this.browser.bounds()
}

func (_this *gui) browserBinary() string {

	if _this.config.BrowserBinary != "" {
		return _this.config.BrowserBinary
	}

	if _this.config.Flavor == Chrome {
		return _this.browserBinaryChrome()
	}

	return _this.browserBinaryEdge()

}

func (_this *gui) browserBinaryEdge() string {

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		return ""
	case "windows":
		paths = []string{
			os.Getenv("LocalAppData") + "/Microsoft/Edge/Application/msedge.exe",
			os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe",
			os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		}
	default:
		return ""
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""

}

func (_this *gui) browserBinaryChrome() string {

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

func (_this *gui) newBrowser() (*browser, error) {

	// The first two IDs are used internally during the initialization
	c := &browser{
		id:       2,
		pending:  map[int]chan result{},
		bindings: map[string]bindingFunc{},
	}

	// Start chrome process
	c.cmd = exec.Command(_this.config.BrowserBinary, _this.config.Args...)
	pipe, err := c.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := c.cmd.Start(); err != nil {
		return nil, err
	}

	// Wait for websocket address to be printed to stderr
	re := regexp.MustCompile(`^DevTools listening on (ws://.*?)\r?\n$`)
	m, err := readUntilMatch(pipe, re)
	if err != nil {
		c.kill()
		return nil, err
	}
	wsURL := m[1]

	// Open a websocket
	c.ws, err = websocket.Dial(wsURL, "", "http://127.0.0.1")
	if err != nil {
		c.kill()
		return nil, err
	}

	// Find target and initialize session
	c.target, err = c.findTarget()
	if err != nil {
		c.kill()
		return nil, err
	}

	c.session, err = c.startSession(c.target)
	if err != nil {
		c.kill()
		return nil, err
	}
	go c.readLoop()
	for method, args := range map[string]h{
		"Page.enable":          nil,
		"Target.setAutoAttach": {"autoAttach": true, "waitForDebuggerOnStart": false},
		"Network.enable":       nil,
		"Runtime.enable":       nil,
		"Security.enable":      nil,
		"Performance.enable":   nil,
		"Log.enable":           nil,
	} {
		if _, err := c.send(method, args); err != nil {
			c.kill()
			c.cmd.Wait()
			return nil, err
		}
	}

	if !contains(_this.config.Args, "--headless") {
		win, err := c.getWindowForTarget(c.target)
		if err != nil {
			c.kill()
			return nil, err
		}
		c.window = win.WindowID
	}

	return c, nil
}

func New(conf ...Config) (*gui, error) {

	newGui := gui{}

	if len(conf) > 0 {
		newGui.config = conf[0]
	}

	if newGui.config.BrowserBinary == "" {
		newGui.config.BrowserBinary = newGui.browserBinary()
	}

	if newGui.config.Url == "" {
		newGui.config.Url = newGui.genEmptyHtml()
	}

	if newGui.config.UserDataDir == "" {

		tempFolder, err := ioutil.TempDir("", "proton/userdata")
		if err != nil {
			return nil, err
		}

		newGui.config.UserDataDir = tempFolder

	}

	args := append(newGui.config.Args, fmt.Sprintf("--app=%s", newGui.config.Url))
	args = append(args, fmt.Sprintf("--user-data-dir=%s", newGui.config.UserDataDir))

	if newGui.config.Height == 0 || newGui.config.Width == 0 {
		args = append(args, "--start-maximized")
	} else {
		args = append(args, fmt.Sprintf("--window-size=%d,%d", newGui.config.Width, newGui.config.Height))
	}

	if newGui.config.Debug {
		args = append(args, "--auto-open-devtools-for-tabs")
	}

	args = append(args, "--remote-debugging-port=0")

	newGui.config.Args = args

	return &newGui, nil

}
