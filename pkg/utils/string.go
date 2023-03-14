package utils

func ContainString(sli []string, elem string) bool {
	for _, e := range sli {
		if elem == e {
			return true
		}
	}
	return false
}
