package filesystemV4

var APP struct {
	Dir          Dir
	DirOperation DirOperation
	DirOperAttr  struct {
		Flag AttrDirFlag
		Mode AttrDirMode
	}
	File          File
	FileOperation FileOperation
	FileOperAttr  struct {
		Flag AttrFileFlag
		Mode AttrFileMode
	}
}
