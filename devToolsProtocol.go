package proton

import "encoding/json"

/*
Capture page screenshot
*/
func (_this *Browser) PageCaptureScreenshot(Parameters PageCaptureScreenshotParameters) (string, error) {

	result, err := _this.send("Page.captureScreenshot", structToMap(Parameters))

	if err != nil {
		return "", err
	}

	data := struct {
		Data string `json:"data"`
	}{}

	err = json.Unmarshal(result, &data)
	return data.Data, err

}

/*
Reloads given page optionally ignoring the cache.
*/
func (_this *Browser) PageReload() error {

	_, err := _this.send("Page.reload", h{})

	return err

}

/*
Resets navigation history for the current page.
*/
func (_this *Browser) PageResetNavigationHistory() error {

	_, err := _this.send("Page.resetNavigationHistory", h{})

	return err

}

/*
Force the page stop all navigations and pending resource fetches.
*/
func (_this *Browser) PageStopLoading() error {

	_, err := _this.send("Page.stopLoading", h{})

	return err

}

/*
Brings page to front (activates tab).
*/
func (_this *Browser) PageBringToFront() error {

	_, err := _this.send("Page.bringToFront", h{})

	return err

}

/*
Print page as PDF.
*/
func (_this *Browser) PagePrintToPDF(Parameters PrintToPDFParameters) (string, error) {

	result, err := _this.send("Page.printToPDF", structToMap(Parameters))

	if err != nil {
		return "", err
	}

	data := struct {
		Data string `json:"data"`
	}{}

	err = json.Unmarshal(result, &data)
	return data.Data, err

}

/*
Navigates current page to the given URL.
*/
func (_this *Browser) PageNavigate(url string) error {
	_, err := _this.send("Page.navigate", h{"url": url})
	return err
}

/*
Close browser gracefully.
*/
func (_this *Browser) BrowserClose() error {

	_, err := _this.send("Browser.close", h{})

	return err
}

/*
Returns version information.
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
Clears browser cache.
*/
func (_this *Browser) BrowserClearBrowserCache() error {

	_, err := _this.send("Browser.clearBrowserCache", h{})

	return err

}

/*
Clears browser cookies.
*/
func (_this *Browser) BrowserClearBrowserCookies() error {

	_, err := _this.send("Browser.clearBrowserCookies", h{})

	return err

}
