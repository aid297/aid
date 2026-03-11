package v1HTTPService

var New app

type app struct{}

func (*app) UUID() *UUIDService                 { return &UUIDService{} }
func (*app) MessageBoard() *MessageBoardService { return &MessageBoardService{} }
