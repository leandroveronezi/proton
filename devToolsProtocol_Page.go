package proton

import "encoding/json"

//TODO: Page.addScriptToEvaluateOnNewDocument

//PageBringToFront Brings page to front (activates tab).
func (_this *Browser) PageBringToFront() error {

	_, err := _this.send("Page.bringToFront", h{})

	return err

}

//PageCaptureScreenshot Capture page screenshot
func (_this *Browser) PageCaptureScreenshot(Parameters PageCaptureScreenshotParameters) (PageCaptureScreenshotReturn, error) {

	result, err := _this.send("Page.captureScreenshot", structToMap(Parameters))

	data := PageCaptureScreenshotReturn{}

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(result, &data)

	return data, err

}

//TODO: Page.createIsolatedWorld

//PageDisable Disables page domain notifications.
func (_this *Browser) PageDisable() error {

	_, err := _this.send("Page.disable", h{})

	return err

}

//PageEnable Enables page domain notifications.
func (_this *Browser) PageEnable() error {

	_, err := _this.send("Page.enable", h{})

	return err

}

//TODO: Page.getAppManifest
//TODO: Page.getFrameTree
//TODO: Page.getLayoutMetrics
//TODO: Page.getNavigationHistory

//PageHandleJavaScriptDialog Accepts or dismisses a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload).
func (_this *Browser) PageHandleJavaScriptDialog(Parameters PageHandleJavaScriptDialogParameters) error {

	_, err := _this.send("Page.handleJavaScriptDialog", structToMap(Parameters))

	return err
}

//PageNavigate Navigates current page to the given URL.
func (_this *Browser) PageNavigate(Parameters PageNavigateParameters) (PageNavigateReturn, error) {

	result, err := _this.send("Page.navigate", structToMap(Parameters))

	data := PageNavigateReturn{}

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(result, &data)

	return data, err
}

//TODO: Page.navigateToHistoryEntry

//PagePrintToPDF Print page as PDF.
func (_this *Browser) PagePrintToPDF(Parameters PrintToPDFParameters) (PrintToPDFReturn, error) {

	result, err := _this.send("Page.printToPDF", structToMap(Parameters))

	data := PrintToPDFReturn{}

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(result, &data)

	return data, err

}

//PageReload Reloads given page optionally ignoring the cache.
func (_this *Browser) PageReload() error {

	_, err := _this.send("Page.reload", h{})

	return err

}

//Page.removeScriptToEvaluateOnNewDocument Removes given script from the list.
func (_this *Browser) PageremoveScriptToEvaluateOnNewDocument(Parameters PageremoveScriptToEvaluateOnNewDocumentParameters) error {

	_, err := _this.send("Page.removeScriptToEvaluateOnNewDocument", structToMap(Parameters))

	return err

}

//PageResetNavigationHistory Resets navigation history for the current page.
func (_this *Browser) PageResetNavigationHistory() error {

	_, err := _this.send("Page.resetNavigationHistory", h{})

	return err

}

//PageSetDocumentContent Sets given markup as the document's HTML.
func (_this *Browser) PageSetDocumentContent(Parameters PageSetDocumentContentParameters) error {

	_, err := _this.send("Page.navigate", structToMap(Parameters))

	return err

}

//PageStopLoading Force the page stop all navigations and pending resource fetches.
func (_this *Browser) PageStopLoading() error {

	_, err := _this.send("Page.stopLoading", h{})

	return err

}
