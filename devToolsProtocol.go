package proton

import "encoding/json"

/*
BrowserClose Close browser gracefully.
*/
func (_this *Browser) BrowserClose() error {

	_, err := _this.send("Browser.close", h{})

	return err
}

/*
BrowserGetVersion Returns version information.
*/
func (_this *Browser) BrowserGetVersion() (Version, error) {
	result, err := _this.send("Browser.getVersion", h{})

	if err != nil {
		return Version{}, err
	}

	version := Version{}
	err = json.Unmarshal(result, &version)
	return version, err

}

/*
BrowserClearBrowserCache Clears browser cache.
*/
func (_this *Browser) BrowserClearBrowserCache() error {

	_, err := _this.send("Browser.clearBrowserCache", h{})

	return err

}

/*
BrowserClearBrowserCookies Clears browser cookies.
*/
func (_this *Browser) BrowserClearBrowserCookies() error {

	_, err := _this.send("Browser.clearBrowserCookies", h{})

	return err

}

/*
RuntimeEvaluate Evaluates expression on global object.
*/
func (_this *Browser) RuntimeEvaluate(Parameters RuntimeEvaluateParameters) (json.RawMessage, error) {

	return _this.send("Runtime.evaluate", structToMap(Parameters))
}

/* REVISADOS */
