package response

import "github.com/aid297/aid/filesystem/filesystemV4"

type (
	FileUploadResponse struct {
		FileName    string `json:"filename" swaggertype:"string"`
		Size        int64  `json:"size" swaggertype:"integer"`
		SavedAs     string `json:"savedAs" swaggertype:"string"`
		SavedPath   string `json:"savedPath" swaggertype:"string"`
		ContentType string `json:"contentType" swaggertype:"string"`
	}

	FileListResponse struct {
		Filesystemers []filesystemV4.Filesystemer `json:"filesystemers" swaggertype:"array,object"`
		CurrentPath   string                      `json:"currentPath" swaggertype:"string"`
	}
)
