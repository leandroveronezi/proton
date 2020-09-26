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

	gui, err := proton.New(conf)

	if err != nil {
		fmt.Println(err)
		return
	}

	gui.Run()

	gui.Bind("ola", func() string {
		return "mundo"
	})

	gui.Navigate("https://www.wikipedia.org")

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-gui.Done():
	}

}
