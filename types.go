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

type ScreenshotFormatType string

const (
	JPEG ScreenshotFormatType = "jpeg"
	PNG  ScreenshotFormatType = "png"
	WEBP ScreenshotFormatType = "webp"
)

func (_this ScreenshotFormatType) Pointer() *ScreenshotFormatType {
	return &_this
}

type ViewportType struct {
	X      int `json:"x"`      //X offset in device independent pixels (dip).
	Y      int `json:"y"`      //Y offset in device independent pixels (dip).
	Width  int `json:"width"`  //Rectangle width in device independent pixels (dip).
	Height int `json:"height"` //Rectangle height in device independent pixels (dip).
	Scale  int `json:"scale"`  //Page scale factor.
}

//Page.captureScreenshot Parameters
type PageCaptureScreenshotParameters struct {
	Format  *ScreenshotFormatType `json:"format"`  //Image compression format (defaults to png). Allowed Values: jpeg, png, webp
	Quality *int                  `json:"quality"` //Compression quality from range [0..100] (jpeg only).
	Clip    *ViewportType         `json:"clip"`    //Capture the screenshot of a given region only.
}

//Page.captureScreenshot Return
type PageCaptureScreenshotReturn struct {
	Data string `json:"data"` //Base64-encoded image data. (Encoded as a base64 string when passed over JSON)
}

type PageReloadParameters struct {
	IgnoreCache            *bool   `json:"ignoreCache"`
	ScriptToEvaluateOnLoad *string `json:"scriptToEvaluateOnLoad"`
}

//Page.printToPDF Parameters
type PrintToPDFParameters struct {
	Landscape               *bool   `json:"landscape"`               //Paper orientation. Defaults to false.
	DisplayHeaderFooter     *bool   `json:"displayHeaderFooter"`     //Display header and footer. Defaults to false.
	PrintBackground         *bool   `json:"printBackground"`         //Print background graphics. Defaults to false.
	Scale                   *int    `json:"scale"`                   //Scale of the webpage rendering. Defaults to 1.
	PaperWidth              *int    `json:"paperWidth"`              //Paper width in inches. Defaults to 8.5 inches.
	PaperHeight             *int    `json:"paperHeight"`             //Paper height in inches. Defaults to 11 inches.
	MarginTop               *int    `json:"marginTop"`               //Top margin in inches. Defaults to 1cm (~0.4 inches).
	MarginBottom            *int    `json:"marginBottom"`            //Bottom margin in inches. Defaults to 1cm (~0.4 inches).
	MarginLeft              *int    `json:"marginLeft"`              //Left margin in inches. Defaults to 1cm (~0.4 inches).
	MarginRight             *int    `json:"marginRight"`             //Right margin in inches. Defaults to 1cm (~0.4 inches).
	PageRanges              *string `json:"pageRanges"`              //Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages.
	IgnoreInvalidPageRanges *bool   `json:"ignoreInvalidPageRanges"` //Whether to silently ignore invalid but successfully parsed page ranges, such as '3-2'. Defaults to false.
	HeaderTemplate          *string `json:"headerTemplate"`          //HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: date: formatted print date title: document title, url: document location, pageNumber: current page number, totalPages: total pages in the document
	FooterTemplate          *string `json:"footerTemplate"`          //HTML template for the print footer. Should use the same format as the
	PreferCSSPageSize       *bool   `json:"preferCSSPageSize"`       //Whether or not to prefer page size as defined by css. Defaults to false, in which case the content will be scaled to fit the paper size.
}

//Page.printToPDF Return
type PrintToPDFReturn struct {
	Data string `json:"data"` //Base64-encoded image data. (Encoded as a base64 string when passed over JSON)
}

type RuntimeEvaluateParameters struct {
	Expression            string  `json:"expression"`
	ObjectGroup           *string `json:"objectGroup"`
	IncludeCommandLineAPI *bool   `json:"includeCommandLineAPI"`
	Silent                *bool   `json:"silent"`
	ContextId             *int    `json:"contextId"`
	ReturnByValue         *bool   `json:"returnByValue"`
	GeneratePreview       *bool   `json:"generatePreview"`
	UserGesture           *bool   `json:"userGesture"`
	AwaitPromise          *bool   `json:"awaitPromise"`
	ThrowOnSideEffect     *bool   `json:"throwOnSideEffect"`
	Timeout               *int    `json:"timeout"`
	DisableBreaks         *bool   `json:"disableBreaks"`
}

type TransitionType string

const (
	Link             TransitionType = "link"
	Typed            TransitionType = "typed"
	AddressBar       TransitionType = "address_bar"
	AutoBookmark     TransitionType = "auto_bookmark"
	AutoSubframe     TransitionType = "auto_subframe"
	ManualSubframe   TransitionType = "manual_subframe"
	Generated        TransitionType = "generated"
	AutoToplevel     TransitionType = "auto_toplevel"
	FormSubmit       TransitionType = "form_submit"
	Reload           TransitionType = "reload"
	Keyword          TransitionType = "keyword"
	KeywordGenerated TransitionType = "keyword_generated"
	Other            TransitionType = "other"
)

func (_this TransitionType) Pointer() *TransitionType {
	return &_this
}

//Page.navigate Parameters
type PageNavigateParameters struct {
	Url            string          `json:"url"`            //URL to navigate the page to.
	Referrer       *string         `json:"referrer"`       //Referrer URL.
	TransitionType *TransitionType `json:"transitionType"` //Intended transition type.
	FrameId        *string         `json:"frameId"`        //Frame id to navigate, if not specified navigates the top frame.
}

//Page.navigate Return
type PageNavigateReturn struct {
	FrameId   string  `json:"frameId"`   //Frame id that has navigated (or failed to navigate)
	LoaderId  *string `json:"loaderId"`  //Loader identifier.
	ErrorText *string `json:"errorText"` //User friendly error message, present if and only if navigation has failed.
}

//Page.setDocumentContent Parameters
type PageSetDocumentContentParameters struct {
	FrameId *string `json:"frameId"` //Frame id to set HTML for.
	Html    string  `json:"html"`    //HTML content to set.
}

type ScriptIdentifierType string //Unique script identifier.

type PageremoveScriptToEvaluateOnNewDocumentParameters struct {
	Identifier ScriptIdentifierType `json:"identifier"` //Unique script identifier.
}

//Page.handleJavaScriptDialog Parameters
type PageHandleJavaScriptDialogParameters struct {
	Accept     bool    `json:"accept"`     //Whether to accept or dismiss the dialog.
	PromptText *string `json:"promptText"` //The text to enter into the dialog prompt before accepting. Used only if this is a prompt dialog.
}
