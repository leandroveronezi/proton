package main

import (
	"fmt"
	"github.com/leandroveronezi/proton"
	"os"
	"os/signal"
)

func main() {

	conf := proton.Config{}

	conf.WindowState = proton.WindowStateFullscreen

	conf.Title = "Photon First Test"
	conf.Debug = false
	conf.Args = proton.DefaultBrowserArgs
	conf.UserDataDir = "./userdata"
	conf.UserDataDirKeep = true
	conf.Flavor = proton.Chrome

	browser := proton.Browser{}

	err := browser.Run(conf)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		browser.BrowserClose()
	}()

	browser.Bind("ola", func() string {
		return "mundo"
	})

	browser.Bind("goVersion", func() (proton.Version, error) {
		return browser.BrowserGetVersion()
	})

	browser.Bind("captureScreenshot", func() (string, error) {
		return browser.PageCaptureScreenshot(proton.PageCaptureScreenshotParameters{Format: proton.JPEG.Pointer()})
	})

	browser.Bind("printToPDF", func() (string, error) {
		return browser.PagePrintToPDF(proton.PrintToPDFParameters{})
	})

	browser.PageNavigate("https://www.wikipedia.org")

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-browser.Done():
	}

}
