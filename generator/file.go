package main

import (
	"io"
	"os"
	"strings"
	"time"
)

type goFile struct {
	Name        string
	Description string
	Major       string
	Minor       string
	value       string
}

func NewFile() goFile {
	return goFile{}
}

func (_this *goFile) Add(str string) {
	_this.value += str
}

func (_this goFile) Generate() string {

	aux := `package proton

import "encoding/json"

// {DATA} Generated from the latest (tip-of-tree) of protocol {MINOR}.{MAJOR}

`
	aux = strings.ReplaceAll(aux, "{DATA}", time.Now().Format("02/01/2006 15:04:05"))

	aux = strings.ReplaceAll(aux, "{MINOR}", _this.Minor)
	aux = strings.ReplaceAll(aux, "{MAJOR}", _this.Major)

	if len(strings.Trim(_this.Description, " ")) > 0 {
		aux += "/*" + _this.Description + "*/" + lineBreak + lineBreak
	}

	aux += _this.value

	return aux

}

func (_this goFile) Save() error {

	fo, err := os.Create("./files/dtp_" + _this.Name + ".go")
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(_this.Generate()))
	if err != nil {
		return err
	}

	return nil
}
