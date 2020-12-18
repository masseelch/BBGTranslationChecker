package checker

type Duplicates map[string]int

func (d Duplicates) cleaned() Duplicates {
	if len(d) == 0 {
		return nil
	}

	r := make(Duplicates)
	for k, c := range d {
		if c > 1 {
			r[k] = c
		}
	}

	// If there are no entries left return nil map.
	if len(r) == 0 {
		return nil
	}

	return r
}

func (d Duplicates) keys() []string {
	ks := make([]string, len(d))

	i := 0
	for k := range d {
		ks[i] = k
		i++
	}

	return ks
}
