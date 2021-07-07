package main

import (
	"github.com/leandroveronezi/proton"
	"log"
	"os"
	"os/signal"
)

func main() {

	conf := proton.Config{}

	conf.WindowState = proton.WindowStateMaximized
	conf.Title = "Photon"
	conf.Args = proton.DefaultBrowserArgs
	conf.UserDataDir = "./userdata"
	conf.UserDataDirKeep = true
	conf.Flavor = proton.Edge

	browser := proton.Browser{}

	err := browser.Run(conf)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer browser.BrowserClose()

	browser.PageNavigate(proton.PageNavigateParameters{Url: "https://www.wikipedia.org"})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-browser.Done():
	}

}
