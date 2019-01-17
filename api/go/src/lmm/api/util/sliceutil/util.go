package sliceutil

// ContainsString returns true if strings contains target
func ContainsString(target interface{}, strings []string) bool {
	return AnyString(strings, func(s string) bool {
		return s == target
	})
}

// AnyString maps string items to judge and return true if any of result judge is true
func AnyString(strings []string, judge func(s string) bool) bool {
	for _, s := range strings {
		if judge(s) {
			return true
		}
	}
	return false
}
