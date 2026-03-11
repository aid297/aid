package v1HTTPRoute

var New app

type app struct{}

func (*app) Rezip() *RezipRoute               { return &RezipRoute{} }
func (*app) UUID() *UUIDRoute                 { return &UUIDRoute{} }
func (*app) Upload() *FileManagerRoute        { return &FileManagerRoute{} }
func (*app) MessageBoard() *MessageBoardRoute { return &MessageBoardRoute{} }
