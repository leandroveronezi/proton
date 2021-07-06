package proton

import "encoding/json"

// WindowState defines the state of the Chrome window, possible values are
// "normal", "maximized", "minimized" and "fullscreen".
type WindowState string

const (
	// WindowStateNormal defines a normal state of the browser window
	WindowStateNormal WindowState = "normal"
	// WindowStateMaximized defines a maximized state of the browser window
	WindowStateMaximized WindowState = "maximized"
	// WindowStateMinimized defines a minimized state of the browser window
	WindowStateMinimized WindowState = "minimized"
	// WindowStateFullscreen defines a fullscreen state of the browser window
	WindowStateFullscreen WindowState = "fullscreen"
)

// Bounds defines settable window properties.
type Bounds struct {
	Left        int         `json:"left"`
	Top         int         `json:"top"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	WindowState WindowState `json:"windowState"`
}

type windowTargetMessage struct {
	WindowID int    `json:"windowId"`
	Bounds   Bounds `json:"bounds"`
}

type targetMessageTemplate struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params struct {
		Name    string `json:"name"`
		Payload string `json:"payload"`
		ID      int    `json:"executionContextId"`
		Args    []struct {
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"args"`
	} `json:"params"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	Result json.RawMessage `json:"result"`
}

type targetMessage struct {
	targetMessageTemplate
	Result struct {
		Result struct {
			Type        string          `json:"type"`
			Subtype     string          `json:"subtype"`
			Description string          `json:"description"`
			Value       json.RawMessage `json:"value"`
			ObjectID    string          `json:"objectId"`
		} `json:"result"`
		Exception struct {
			Exception struct {
				Value json.RawMessage `json:"value"`
			} `json:"exception"`
		} `json:"exceptionDetails"`
	} `json:"result"`
}

type Version struct {
	ProtocolVersion string `json:"protocolVersion"`
	Product         string `json:"product"`
	Revision        string `json:"revision"`
	UserAgent       string `json:"userAgent"`
	JsVersion       string `json:"jsVersion"`
}

type ScreenshotFormat string

const (
	JPEG ScreenshotFormat = "jpeg"
	PNG  ScreenshotFormat = "png"
)

func (_this ScreenshotFormat) Pointer() *ScreenshotFormat {
	return &_this
}

type Viewport struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Scale  int `json:"scale"`
}

type PageCaptureScreenshotParameters struct {
	Format  *ScreenshotFormat `json:"format"`
	Quality *int              `json:"quality"`
	Clip    *Viewport         `json:"clip"`
}

type PageReloadParameters struct {
	IgnoreCache            *bool   `json:"ignoreCache"`
	ScriptToEvaluateOnLoad *string `json:"scriptToEvaluateOnLoad"`
}

type PrintToPDFParameters struct {
	Landscape               *bool   `json:"landscape"`
	DisplayHeaderFooter     *bool   `json:"displayHeaderFooter"`
	PrintBackground         *bool   `json:"printBackground"`
	Scale                   *int    `json:"scale"`
	PaperWidth              *int    `json:"paperWidth"`
	PaperHeight             *int    `json:"paperHeight"`
	MarginTop               *int    `json:"marginTop"`
	MarginBottom            *int    `json:"marginBottom"`
	MarginLeft              *int    `json:"marginLeft"`
	MarginRight             *int    `json:"marginRight"`
	PageRanges              *string `json:"pageRanges"`
	IgnoreInvalidPageRanges *bool   `json:"ignoreInvalidPageRanges"`

	/*
		date: formatted print date
		title: document title
		url: document location
		pageNumber: current page number
		totalPages: total pages in the document
	*/

	HeaderTemplate    *string `json:"headerTemplate"`
	FooterTemplate    *string `json:"footerTemplate"`
	PreferCSSPageSize *bool   `json:"preferCSSPageSize"`
}
