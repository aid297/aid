package httpClientV2

type (
	ContentType string
	Accept      string
)

var (
	ContentTypeJSON               ContentType = "json"
	ContentTypeXML                ContentType = "xml"
	ContentTypeXWwwFormURLencoded ContentType = "form"
	ContentTypeFormData           ContentType = "form-data"
	ContentTypePlain              ContentType = "plain"
	ContentTypeHTML               ContentType = "html"
	ContentTypeCSS                ContentType = "css"
	ContentTypeJavascript         ContentType = "javascript"
	ContentTypeSteam              ContentType = "steam"

	ContentTypes = map[ContentType]string{
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

	AcceptJSON       Accept = "json"
	AcceptXML        Accept = "xml"
	AcceptPlain      Accept = "plain"
	AcceptHTML       Accept = "html"
	AcceptCSS        Accept = "css"
	AcceptJavascript Accept = "javascript"
	AcceptSteam      Accept = "steam"
	AcceptAny        Accept = "any"

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
