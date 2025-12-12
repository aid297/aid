package regexp

type (
	Attributer interface{ Register(regexp *Regexp) }

	AttrTargetString  struct{ target string }
	AttrTargetsString struct{ targets []string }
	AttrTargetError   struct{ target error }
	AttrTargetsError  struct{ targets []error }
)

func TargetString(target string) AttrTargetString   { return AttrTargetString{target} }
func (my AttrTargetString) Register(regexp *Regexp) { regexp.target = my.target }

func TargetsString(targets ...string) AttrTargetsString { return AttrTargetsString{targets} }
func (my AttrTargetsString) Register(regexp *Regexp)    { regexp.targets = my.targets }

func TargetError(target error) AttrTargetError     { return AttrTargetError{target} }
func (my AttrTargetError) Register(regexp *Regexp) { regexp.target = my.target.Error() }

func TargetsError(targets ...error) AttrTargetsError { return AttrTargetsError{targets} }
func (my AttrTargetsError) Register(regexp *Regexp) {
	if len(my.targets) == 0 {
		return
	}

	ret := make([]string, 0, len(my.targets))
	for idx := range my.targets {
		ret = append(ret, my.targets[idx].Error())
	}
	regexp.targets = ret
}
