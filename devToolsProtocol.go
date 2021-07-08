package proton

import "encoding/json"

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
