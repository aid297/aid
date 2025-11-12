package str

type (
	RegexpRegexpAttributer interface{ Register(regexp *Regexp) }

	AttrRegexpTargetString  struct{ target string }
	AttrRegexpTargetsString struct{ targets []string }
	AttrRegexpTargetError   struct{ target error }
	AttrRegexpTargetsError  struct{ targets []error }
)

func RegexpTargetString(target string) AttrRegexpTargetString {
	return AttrRegexpTargetString{target: target}
}
func (my AttrRegexpTargetString) Register(regexp *Regexp) { regexp.target = my.target }

func RegexpTargetsString(targets ...string) AttrRegexpTargetsString {
	return AttrRegexpTargetsString{targets: targets}
}
func (my AttrRegexpTargetsString) Register(regexp *Regexp) { regexp.targets = my.targets }

func RegexpTargetError(target error) AttrRegexpTargetError {
	return AttrRegexpTargetError{target: target}
}
func (my AttrRegexpTargetError) Register(regexp *Regexp) { regexp.target = my.target.Error() }

func RegexpTargetsError(targets ...error) AttrRegexpTargetsError {
	return AttrRegexpTargetsError{targets: targets}
}
func (my AttrRegexpTargetsError) Register(regexp *Regexp) {
	if len(my.targets) == 0 {
		return
	}

	ret := make([]string, 0, len(my.targets))
	for idx := range my.targets {
		ret = append(ret, my.targets[idx].Error())
	}
	regexp.targets = ret
}
