package filesystemV33

var APP struct {
	Dir     Dir
	DirAttr struct {
		Path  AttrDirPath
		IsRel AttrDirIsRel
	}
	File     File
	FileAttr struct {
		Path  AttrFilePath
		IsRel AttrFileIsRel
	}
}
