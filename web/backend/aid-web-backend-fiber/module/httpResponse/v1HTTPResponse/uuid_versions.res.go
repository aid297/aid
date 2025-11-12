package v1HTTPResponse

import "hr-fiber/module/httpRequest/v1HTTPRequest"

type UUIDVersionsResponse struct {
	Versions map[string]v1HTTPRequest.UUIDVersion `json:"versions"`
}
