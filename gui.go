package proton

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"os/exec"
	"regexp"
)

func (_this *Browser) Run(conf ...Config) error {

	if len(conf) > 0 {
		_this.config = conf[0]
	}

	if _this.config.BrowserBinary == "" {
		_this.config.BrowserBinary = _this.browserBinary()

		if _this.config.BrowserBinary == "" {
			return errors.New("Binary not found")
		}

	}

	if _this.config.Url == "" {
		_this.config.Url = _this.genEmptyHtml()
	}

	if _this.config.UserDataDir == "" {

		tempFolder, err := ioutil.TempDir("", "proton/userdata")
		if err != nil {
			return err
		}

		_this.config.UserDataDir = tempFolder

	}

	args := append(_this.config.Args, fmt.Sprintf("--app=%s", _this.config.Url))
	args = append(args, fmt.Sprintf("--user-data-dir=%s", _this.config.UserDataDir))

	if _this.config.Height <= 0 || _this.config.Width <= 0 {
		args = append(args, "--start-maximized")
	} else {
		args = append(args, fmt.Sprintf("--window-size=%d,%d", _this.config.Width, _this.config.Height))
	}

	if _this.config.Debug {
		args = append(args, "--auto-open-devtools-for-tabs")
	}

	args = append(args, "--remote-debugging-port=0")

	_this.config.Args = args

	return _this.makeBrowser()

}

func (_this *Browser) makeBrowser() error {

	// The first two IDs are used internally during the initialization
	_this.id = 2
	_this.pending = map[int]chan result{}
	_this.bindings = map[string]bindingFunc{}

	// Start chrome process
	_this.cmd = exec.Command(_this.config.BrowserBinary, _this.config.Args...)
	pipe, err := _this.cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := _this.cmd.Start(); err != nil {
		return err
	}

	// Wait for websocket address to be printed to stderr
	re := regexp.MustCompile(`^DevTools listening on (ws://.*?)\r?\n$`)
	m, err := readUntilMatch(pipe, re)
	if err != nil {
		_this.kill()
		return err
	}
	wsURL := m[1]

	// Open a websocket
	_this.ws, err = websocket.Dial(wsURL, "", "http://127.0.0.1")
	if err != nil {
		_this.kill()
		return err
	}

	// Find target and initialize session
	_this.target, err = _this.findTarget()
	if err != nil {
		_this.kill()
		return err
	}

	_this.session, err = _this.startSession(_this.target)
	if err != nil {
		_this.kill()
		return err
	}

	go _this.readLoop()

	for method, args := range map[string]h{
		"Page.enable":          nil,
		"Target.setAutoAttach": {"autoAttach": true, "waitForDebuggerOnStart": false},
		"Network.enable":       nil,
		"Runtime.enable":       nil,
		"Security.enable":      nil,
		"Performance.enable":   nil,
		"Log.enable":           nil,
	} {

		if _, err := _this.send(method, args); err != nil {
			_this.kill()
			_this.cmd.Wait()
			return err
		}

	}

	if !contains(_this.config.Args, "--headless") {
		win, err := _this.getWindowForTarget(_this.target)
		if err != nil {
			_this.kill()
			return err
		}
		_this.window = win.WindowID
	}

	done := make(chan struct{})
	if err != nil {
		return err
	}

	go func() {
		_this.cmd.Wait()
		close(done)
	}()

	_this.done = done

	return nil

}
