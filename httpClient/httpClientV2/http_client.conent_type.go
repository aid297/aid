package httpClientV2

type (
	ContentType = string
	Accept      string
)

var (
	ContentTypeJSON               ContentType = "application/json"
	ContentTypeXML                ContentType = "application/xml"
	ContentTypeXWwwFormURLencoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData           ContentType = "form-data"
	ContentTypePlain              ContentType = "text/plain"
	ContentTypeHTML               ContentType = "text/html"
	ContentTypeCSS                ContentType = "text/css"
	ContentTypeJavascript         ContentType = "text/javascript"
	ContentTypeSteam              ContentType = "application/octet-stream"
	ContentTypes                              = map[ContentType]string{
		ContentTypeJSON:               "application/json",
		ContentTypeXML:                "application/xml",
		ContentTypeXWwwFormURLencoded: "application/x-www-form-urlencoded",
		ContentTypeFormData:           "form-data",
		ContentTypePlain:              "text/plain",
		ContentTypeHTML:               "text/html",
		ContentTypeCSS:                "text/css",
		ContentTypeJavascript:         "text/javascript",
		ContentTypeSteam:              "application/octet-stream",
	}

	AcceptJSON       Accept = "application/json"
	AcceptXML        Accept = "application/xml"
	AcceptPlain      Accept = "text/plain"
	AcceptHTML       Accept = "text/html"
	AcceptCSS        Accept = "text/css"
	AcceptJavascript Accept = "text/javascript"
	AcceptSteam      Accept = "application/octet-stream"
	AcceptAny        Accept = "*/*"

	Accepts = map[Accept]string{
		AcceptJSON:       "application/json",
		AcceptXML:        "application/xml",
		AcceptPlain:      "text/plain",
		AcceptHTML:       "text/html",
		AcceptCSS:        "text/css",
		AcceptJavascript: "text/javascript",
		AcceptSteam:      "application/octet-stream",
		AcceptAny:        "*/*",
	}
)
