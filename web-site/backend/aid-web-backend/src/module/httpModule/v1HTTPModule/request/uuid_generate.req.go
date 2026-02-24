package request

type (
	UUIDGenerateRequest struct {
		Number         uint8       `json:"number" yaml:"number" toml:"number" swaggertype:"integer" v-rule:"(required)(min>0)(max<256)" v-name:"生成数量"`
		NoSubsTractKey bool        `json:"noSubsTractKey" yaml:"noSubsTractKey" swaggertype:"boolean" toml:"noSubsTractKey" v-rule:"bool" v-name:"不减去key"`
		IsUpper        bool        `json:"isUpper" yaml:"isUpper" toml:"isUpper" swaggertype:"boolean" v-rule:"bool" v-name:"是否大写"`
		Version        UUIDVersion `json:"version" yaml:"version" toml:"version" swaggertype:"string" v-rule:"(required)(in:v1,v4,v6,v7)" v-name:"uuid版本"` // Ensure to convert UUIDVersion to string when asserting
	}

	UUIDVersion = string
)

const (
	UUIDVersionV1 UUIDVersion = "v1"
	UUIDVersionV4 UUIDVersion = "v4"
	UUIDVersionV6 UUIDVersion = "v6"
	UUIDVersionV7 UUIDVersion = "v7"
)
