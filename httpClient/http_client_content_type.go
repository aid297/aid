package httpClient

type (
	ContentType string
	Accept      string
)

var (
	ContentTypeJson       ContentType = "json"
	ContentTypeXml        ContentType = "xml"
	ContentTypeForm       ContentType = "form"
	ContentTypeFormData   ContentType = "form-data"
	ContentTypePlain      ContentType = "plain"
	ContentTypeHtml       ContentType = "html"
	ContentTypeCss        ContentType = "css"
	ContentTypeJavascript ContentType = "javascript"
	ContentTypeSteam      ContentType = "steam"

	ContentTypes = map[ContentType]string{
		ContentTypeJson:       "application/json",
		ContentTypeXml:        "application/xml",
		ContentTypeForm:       "application/x-www-form-urlencoded",
		ContentTypeFormData:   "form-data",
		ContentTypePlain:      "text/plain",
		ContentTypeHtml:       "text/html",
		ContentTypeCss:        "text/css",
		ContentTypeJavascript: "text/javascript",
		ContentTypeSteam:      "application/octet-stream",
	}

	AcceptJson       Accept = "json"
	AcceptXml        Accept = "xml"
	AcceptPlain      Accept = "plain"
	AcceptHtml       Accept = "html"
	AcceptCss        Accept = "css"
	AcceptJavascript Accept = "javascript"
	AcceptSteam      Accept = "steam"
	AcceptAny        Accept = "any"

	Accepts = map[Accept]string{
		AcceptJson:       "application/json",
		AcceptXml:        "application/xml",
		AcceptPlain:      "text/plain",
		AcceptHtml:       "text/html",
		AcceptCss:        "text/css",
		AcceptJavascript: "text/javascript",
		AcceptSteam:      "application/octet-stream",
		AcceptAny:        "*/*",
	}
)
