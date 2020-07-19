package helpers

func NSet(flags ...bool) int {
	var n int
	for _, flag := range flags {
		if flag {
			n++
		}
	}
	return n
}
