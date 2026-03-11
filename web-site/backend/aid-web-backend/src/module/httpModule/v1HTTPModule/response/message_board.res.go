package response

import "github.com/aid297/aid/simpleDB/kernal"

type MessageBoardListResponse struct {
	MessageBoards []kernal.Row `json:"messageBoards" xml:"messageBoards" yaml:"messageBoards" toml:"messageBoards" swaggertype:"array,object"`
}
