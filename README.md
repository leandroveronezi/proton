# proton

[![Go Report Card](https://goreportcard.com/badge/github.com/leandroveronezi/proton)](https://goreportcard.com/report/github.com/leandroveronezi/proton)
[![GoDoc](https://godoc.org/github.com/leandroveronezi/proton?status.svg)](https://godoc.org/github.com/leandroveronezi/proton)
![](https://img.shields.io/github/repo-size/leandroveronezi/proton.svg)
![MIT Licensed](https://img.shields.io/github/license/leandroveronezi/proton.svg)

<div>
  <p align="justify">
      A very small library for creating modern desktop apps in Go. Unlike Electron, 
      the browser is not bundled with the app by reusing the one that is already installed. 
      Proton uses the Chrome DevTools Protocol to interact with chromium-based browsers by providing a 
      UI layer allowing you to calling Go code from the UI and manipulating UI from Go in a seamless manner.
  </p>
</div>

## Requirements
  Requires Chrome/Chromium >= 70

## Features
* Pure Go library 
* Very simple API
* Small application size

## Example

```go
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

```

Also, see [examples](examples) for more details about binding functions and packaging binaries.

## Hello World

Here are the steps to run the hello world example.

```
cd examples/wikipedia
go run .
```

 


