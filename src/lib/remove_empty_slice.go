package lib

func RemoveEmptySlice(in []string) []string {

	var s []string

	for _, item := range in {
		if item != "" {
			s = append(s, item)
		}
	}

	return s
}
