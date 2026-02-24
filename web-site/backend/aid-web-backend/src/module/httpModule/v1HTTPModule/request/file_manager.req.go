package request

type (
	FileRequest struct {
		Path string `json:"path" yaml:"path" toml:"path" swaggertype:"string"`
	}
	FileListRequest        FileStoreFolderRequest
	FileStoreFolderRequest struct {
		FileRequest
		Name string `json:"name" yaml:"name" toml:"name" swaggertype:"string"`
	}
	FileDestroyRequest FileStoreFolderRequest
	FileZipRequest     FileStoreFolderRequest
)
