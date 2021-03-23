package checker

// key: translation tag
type Duplicates map[string]int

// Only keep those entries that occur more than once.
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

// Return the key with the highest count
func (d Duplicates) highest() string {
	var h string
	for t, c := range d {
		if h == "" || c > d[h] {
			h = t
		}
	}

	return h
}
