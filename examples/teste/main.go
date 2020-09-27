package main

import (
	"fmt"
	"github.com/leandroveronezi/proton"
	"os"
	"os/signal"
)

func main() {

	conf := proton.Config{}

	conf.Title = "Photon First Test"
	conf.Debug = false
	conf.Args = proton.DefaultBrowserArgs
	conf.UserDataDir = "./userdata"
	conf.UserDataDirKeep = true
	conf.Flavor = proton.Edge

	browser := proton.Browser{}

	err := browser.Run(conf)

	if err != nil {
		fmt.Println(err)
		return
	}

	browser.Bind("ola", func() string {
		return "mundo"
	})

	browser.Bind("goVersion", func() (proton.Version, error) {
		return browser.GetVersion()
	})

	browser.Bind("captureScreenshot", func() (string, error) {
		return browser.CaptureScreenshot(proton.ScreenshotParameters{Format: proton.JPEG.Pointer()})
	})

	browser.Navigate("https://www.wikipedia.org")

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-browser.Done():
	}

}
