package v1HTTPAPI

var New app

type app struct{}

func (*app) Rezip() *RezipAPI               { return &RezipAPI{} }
func (*app) UUID() *UUIDAPI                 { return &UUIDAPI{} }
func (*app) FileManager() *FileManagerAPI   { return &FileManagerAPI{} }
func (*app) MessageBoard() *MessageBoardAPI { return &MessageBoardAPI{} }
