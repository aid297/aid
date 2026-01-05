package filesystemV4

var APP struct {
	Dir          Dir
	DirOperation DirOperation
	DirOperAttr  struct {
		Flag AttrDirOperationFlag
		Mode AttrDirOperationMode
	}
	File          File
	FileOperation FileOperation
	FileOperAttr  struct {
		Flag AttrFileOperationFlag
		Mode AttrFileOperationMode
	}
}
