package response

import "github.com/aid297/aid/filesystem/filesystemV4"

type (
	FileUploadResponse struct {
		FileName    string `json:"filename" swaggertype:"string"`
		Size        int64  `json:"size" swaggertype:"integer"`
		ContentType string `json:"contentType" swaggertype:"string"`
	}

	FileListResponse struct {
		Items       []filesystemV4.Filesystemer `json:"items" swaggertype:"array,object"`
		CurrentPath string                      `json:"currentPath" swaggertype:"string"`
	}

	FileZipResponse struct {
		Name string `json:"name" swaggertype:"string"`
	}
)
