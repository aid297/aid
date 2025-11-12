package response

type UUIDVersionsResponse struct {
	Versions map[string]string `json:"versions" yaml:"versions" toml:"versions"`
}
