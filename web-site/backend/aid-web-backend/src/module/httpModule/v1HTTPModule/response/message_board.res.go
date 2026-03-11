package response

import "github.com/aid297/aid/simpleDB/simpleDBDriver"

type MessageBoardListResponse struct {
	MessageBoards []simpleDBDriver.Row `json:"messageBoards" xml:"messageBoards" yaml:"messageBoards" toml:"messageBoards" swaggertype:"array,object"`
}
