package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const packageName = "main"

func FileLine(depth int) (string, int) {

	for i := depth; ; i++ {

		_, file, line, ok := runtime.Caller(i)

		if !ok {
			break
		}

		if strings.Contains(file, packageName) {
			//continue
		}

		aux, _ := filepath.Abs(file)
		return aux, line
	}

	return "", 0

}

func Log(v ...interface{}) {

	//log.SetOutput(os.Stdout)

	filename, fileline := FileLine(2)

	titulo := filename + ":" + strconv.Itoa(fileline)

	fmt.Print("[" + time.Now().Format("02/01/2006 15:04:05") + "] ")
	fmt.Println(titulo)
	fmt.Println(v)
	fmt.Println("")

}
