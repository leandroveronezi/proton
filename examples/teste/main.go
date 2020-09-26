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
	conf.UserDataDir = "./proton/userdata"
	conf.UserDataDirKeep = true
	conf.Flavor = proton.Edge

	gui, err := proton.New(conf)

	if err != nil {
		fmt.Println(err)
		return
	}

	gui.Run()

	gui.Navigate("http://www.uol.com.br")

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-gui.Done():
	}

}
