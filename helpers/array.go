package helpers


func Contains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}