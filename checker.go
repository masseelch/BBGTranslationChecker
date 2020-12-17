package checker

type (
	// Key: tag, Value: count
	Duplicates map[string]int
	Report struct {
		// Key: filename
		Duplicates map[string]Duplicates

	}
)


func Check(truth *File, translations []*File) *Report {
	r := &Report{
		Duplicates: make(map[string]Duplicates),
	}

	// Check for duplicate tags in all translations.
	td := duplicates(truth.Translations)
	r.Duplicates[truth.Filename] = td

	// Duplicates
	// Missing translations
	// Wrong of values
	return r
}

func duplicates(ts []Translation) Duplicates {
	d := make(Duplicates)

	// Count the occurrences of every tag.
	for _, t := range ts {
		d[t.Tag]++
	}

	// Only keep those entries that occur more than once.
	r := make(Duplicates)
	for k, c := range d {
		if c > 1 {
			r[k] = c
		}
	}

	return r
}
