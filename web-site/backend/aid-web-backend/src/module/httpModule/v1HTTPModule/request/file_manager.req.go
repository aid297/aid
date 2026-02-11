package request

type (
	FileRequest struct {
		Path string `json:"path" yaml:"path" toml:"path"`
	}
	FileListRequest        FileStoreFolderRequest
	FileStoreFolderRequest struct {
		FileRequest
		Name string `json:"name" yaml:"name" toml:"name"`
	}
	FileDestroyRequest FileStoreFolderRequest
	FileZipRequest     FileStoreFolderRequest
)
