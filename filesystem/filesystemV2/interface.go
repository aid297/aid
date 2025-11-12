package filesystemV2

type FilesystemV2Type string

const (
	FilesystemV2File FilesystemV2Type = "FILE"
	FilesystemV2Dir  FilesystemV2Type = "DIR"
)

type IFilesystemV2 interface{ WhoAmI() FilesystemV2Type }
