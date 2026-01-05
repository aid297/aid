package filesystemV3

var APP struct {
	Dir          Dir
	DirOperation DirOperation
	DirAttr      struct {
		Path  AttrDirPath
		IsRel AttrDirIsRel
	}
	File     File
	FileAttr struct {
		Path  AttrFilePath
		IsRel AttrFileIsRel
	}
}
