package proton

import "encoding/json"

//BrowserClose Close browser gracefully.
func (_this *Browser) BrowserClose() error {

	_, err := _this.send("Browser.close", h{})

	return err
}

//BrowserGetVersion Returns version information.
func (_this *Browser) BrowserGetVersion() (BrowserGetVersionReturn, error) {
	result, err := _this.send("Browser.getVersion", h{})

	data := BrowserGetVersionReturn{}

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(result, &data)

	return data, err
}
