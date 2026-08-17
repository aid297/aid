package validations

import (
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/aid297/aid/v2/anySlices"
	"github.com/aid297/aid/v2/points"
)

func getRuleRequired(rules anySlices.AnySlicer[string]) bool { return rules.In("required", "!") }

func getRuleExFnNames(rule string) (exFnNames []string) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "ex:"); ok {
		exFnNames = strings.Split(value, string(validatorExIns.defaultSliceSplitChar))
		return
	}

	return
}

func getRuleUintSize(rule string) (size *uint, eq bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "size=="); ok {
		size = points.New(cast.ToUint(value))
		eq = true
		return
	}

	if value, ok = strings.CutPrefix(rule, "size!="); ok {
		size = points.New(cast.ToUint(value))
		eq = false
		return
	}

	return
}

func getRuleIntSize(rule string) (size *int, eq bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "size=="); ok {
		size = points.New(cast.ToInt(value))
		eq = true
		return
	}

	if value, ok = strings.CutPrefix(rule, "size!="); ok {
		size = points.New(cast.ToInt(value))
		eq = false
		return
	}

	return
}

func getRuleFloatSize(rule string) (size *float64, eq bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "size=="); ok {
		size = points.New(cast.ToFloat64(value))
		eq = true
		return
	}

	if value, ok = strings.CutPrefix(rule, "size!="); ok {
		size = points.New(cast.ToFloat64(value))
		eq = false
		return
	}

	return
}

func getRuleUintMin(rule string) (size *uint, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		size = points.New(cast.ToUint(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		size = points.New(cast.ToUint(value))
		return
	}

	return
}

func getRuleUintMax(rule string) (size *uint, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		size = points.New(cast.ToUint(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		size = points.New(cast.ToUint(value))
		include = false
		return
	}

	return
}

func getRuleStrTimeMin(rule string) (*string, bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "str-time>="); ok {
		return &value, true
	}

	if value, ok = strings.CutPrefix(rule, "str-time>"); ok {
		return &value, false
	}

	return nil, false
}

func getRuleStrTimeMax(rule string) (*string, bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "str-time<="); ok {
		return &value, true
	}

	if value, ok = strings.CutPrefix(rule, "str-time<"); ok {
		return &value, false
	}

	return nil, false
}

func getRuleIntMin(rule string) (size *int, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		size = points.New(cast.ToInt(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		size = points.New(cast.ToInt(value))
		return
	}

	return
}

func getRuleIntMax(rule string) (size *int, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		size = points.New(cast.ToInt(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		size = points.New(cast.ToInt(value))
		include = false
		return
	}

	return
}

func getRuleFloatMin(rule string) (size *float64, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		size = points.New(cast.ToFloat64(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		size = points.New(cast.ToFloat64(value))
		include = false
		return
	}

	return
}

func getRuleFloatMax(rule string) (size *float64, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		size = points.New(cast.ToFloat64(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		size = points.New(cast.ToFloat64(value))
		include = false
		return
	}

	return
}

func getRuleIn(rule string) (in []string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "in=="); ok {
		in = strings.Split(value, string(validatorExIns.defaultSliceSplitChar))
		return
	}

	return
}

func getRuleRegexEq(rule string) (pattern string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "regex=="); ok {
		pattern = value
		return
	}

	return
}

func getRuleRegexNotEq(rule string) (pattern string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "regex!="); ok {
		pattern = value
		return
	}

	return
}

func getRuleNotIn(rule string) (notIn []string) {
	var (
		value string
		ok    bool
	)
	if value, ok = strings.CutPrefix(rule, "in!="); ok {
		notIn = strings.Split(value, string(validatorExIns.defaultSliceSplitChar))
		return
	}

	return
}

func getRuleTimeMin(rule string) (t *time.Time, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "min>="); ok {
		t = points.New(cast.ToTime(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "min>"); ok {
		t = points.New(cast.ToTime(value))
		include = false
		return
	}

	return
}

func getRuleTimeMax(rule string) (t *time.Time, include bool) {
	var (
		value string
		ok    bool
	)

	if value, ok = strings.CutPrefix(rule, "max<="); ok {
		t = points.New(cast.ToTime(value))
		include = true
		return
	}
	if value, ok = strings.CutPrefix(rule, "max<"); ok {
		t = points.New(cast.ToTime(value))
		include = false
		return
	}

	return
}

func getRuleTimeIn(rule string) (in []time.Time) {
	var (
		value string
		ok    bool
		times []string
	)
	if value, ok = strings.CutPrefix(rule, "in=="); ok {
		times = strings.Split(value, string(validatorExIns.defaultSliceSplitChar))
	}

	if len(times) > 0 {
		in = make([]time.Time, 0, len(times))
		for idx := range times {
			in = append(in, cast.ToTime(times[idx]))
		}
	}

	return
}

func getRuleTimeNotIn(rule string) (notIn []time.Time) {
	var (
		value string
		ok    bool
		times []string
	)
	if value, ok = strings.CutPrefix(rule, "in!="); ok {
		times = strings.Split(value, string(validatorExIns.defaultSliceSplitChar))
	}

	if len(times) > 0 {
		notIn = make([]time.Time, 0, len(times))
		for idx := range times {
			notIn = append(notIn, cast.ToTime(times[idx]))
		}
	}

	return
}
