package request

type (
	FileRequest struct {
		Path string `json:"path" yaml:"path" toml:"path"`
	}
	FileListRequest     FileRequest
	FileDownloadRequest FileRequest
	FileDestroyRequest  FileRequest
	FileZipRequest      FileRequest
)
