package response

type MessageBoardListResponse struct {
	MessageBoards []map[string]string `json:"messageBoards" xml:"messageBoards" yaml:"messageBoards" toml:"messageBoards" swaggertype:"array,object"`
}
