package httpClients

type (
	HTTPContentType = string
	HTTPAccept      string
)

var (
	ContentTypeJSON               HTTPContentType = "application/json"
	ContentTypeXML                HTTPContentType = "application/xml"
	ContentTypeXWWWFormURLEncoded HTTPContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData           HTTPContentType = "form-data"
	ContentTypePlain              HTTPContentType = "text/plain"
	ContentTypeHTML               HTTPContentType = "text/html"
	ContentTypeCSS                HTTPContentType = "text/css"
	ContentTypeJavascript         HTTPContentType = "text/javascript"
	ContentTypeSteam              HTTPContentType = "application/octet-stream"
	HTTPContentTypes                              = map[HTTPContentType]string{
		ContentTypeJSON:               "application/json",
		ContentTypeXML:                "application/xml",
		ContentTypeXWWWFormURLEncoded: "application/x-www-form-urlencoded",
		ContentTypeFormData:           "form-data",
		ContentTypePlain:              "text/plain",
		ContentTypeHTML:               "text/html",
		ContentTypeCSS:                "text/css",
		ContentTypeJavascript:         "text/javascript",
		ContentTypeSteam:              "application/octet-stream",
	}

	AcceptJSON       HTTPAccept = "application/json"
	AcceptXML        HTTPAccept = "application/xml"
	AcceptPlain      HTTPAccept = "text/plain"
	AcceptHTML       HTTPAccept = "text/html"
	AcceptCSS        HTTPAccept = "text/css"
	AcceptJavascript HTTPAccept = "text/javascript"
	AcceptSteam      HTTPAccept = "application/octet-stream"
	AcceptAny        HTTPAccept = "*/*"

	HTTPAccepts = map[HTTPAccept]string{
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
