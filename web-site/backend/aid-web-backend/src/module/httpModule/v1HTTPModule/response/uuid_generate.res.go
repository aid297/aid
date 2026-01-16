package response

type UUIDResponse struct {
	UUID string `json:"uuid" yaml:"uuid" toml:"uuid"`
}

var UUID UUIDResponse

func (UUIDResponse) New(uuid string) UUIDResponse { return UUIDResponse{UUID: uuid} }

type UUIDGenerateResponse struct {
	UUIDs []UUIDResponse `json:"uuids" yaml:"uuids" toml:"uuids"`
}
