package httpClientV4

type (
	ContentType string
	Accept      string
)

const (
	ContentTypeJSON               ContentType = "application/json"
	ContentTypeXML                ContentType = "application/xml"
	ContentTypeXWwwFormURLencoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData           ContentType = "multipart/form-data"
	ContentTypePlain              ContentType = "text/plain"
	ContentTypeHTML               ContentType = "text/html"
	ContentTypeCSS                ContentType = "text/css"
	ContentTypeJavascript         ContentType = "text/javascript"
	ContentTypeSteam              ContentType = "application/octet-stream"

	AcceptJSON       Accept = "application/json"
	AcceptXML        Accept = "application/xml"
	AcceptPlain      Accept = "text/plain"
	AcceptHTML       Accept = "text/html"
	AcceptCSS        Accept = "text/css"
	AcceptJavascript Accept = "text/javascript"
	AcceptSteam      Accept = "application/octet-stream"
	AcceptAny        Accept = "*/*"
)
