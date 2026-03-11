package initialize

var New app

type app struct{}

func (*app) Config() *ConfigInitialize           { return &ConfigInitialize{} }
func (*app) Zap() *ZapInitialize                 { return &ZapInitialize{} }
func (*app) Timezone() *TimezoneInitialize       { return &TimezoneInitialize{} }
func (*app) FileManager() *FileManagerInitialize { return &FileManagerInitialize{} }
func (*app) DB() *SDB                            { return &SDB{} }
