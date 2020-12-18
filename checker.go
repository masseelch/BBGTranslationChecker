package checker

import "fmt"

type (
	Duplicates map[string]int
	// A report for a translation file.
	Report struct {
		// All tags that occur more than once.
		DuplicateTags Duplicates
		// All found language declarations.
		LanguageTags []string
		// Tags that do not appear in the given file but do in the truth.
		MissingTags []string
		// Tags do appear in the given file but are not translated yet (prefixed with "To Translate")
		MissingTranslations []string
	}
	// Key: filename
	Reports map[string]*Report
)

func Check(tf *File, cfs []*File) (Reports, error) {
	var err error
	rs := make(Reports)

	// Collect the report for every translation.
	rs[tf.Filename], err = report(nil, tf)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// t: File of truth
// c: File to check
func report(tf *File, cf *File) (*Report, error) {
	r := new(Report)

	// duplicate tags / lang
	r.duplicatesCheck(cf)

	// missing translations
	r.translationsCheck(tf, cf)

	return r, nil
}

func (r *Report) duplicatesCheck(f *File) {
	// duplicate tags / lang
	r.DuplicateTags = make(Duplicates)
	lang := make(Duplicates)

	// Count the occurrences of every tag.
	for _, t := range f.Translations {
		r.DuplicateTags[t.Tag]++
		lang[t.Lang]++
	}

	// Only keep those entries that occur more than once.
	r.DuplicateTags = r.DuplicateTags.cleaned()

	// Add every found lang to the report
	r.LanguageTags = lang.keys()
}

func (r *Report) translationsCheck(tf *File, cf *File) {
	// For every tag in the truth file check if there exists a translation in the file to check.
	// If found validate that the translation does no longer contain "To Translate".
	// If found check the numeric values in truth and checked file match.
	for _, tt := range tf.Translations {
		ct := cf.Translations.LookupByTag(tt.Tag)

		// Does a translation exist.
		if ct == nil {
			r.addMissingTag(tt.Tag)
		}

		// Translation tag found. Check if is yet untranslated.

	}
}

func (r *Report) addMissingTag(tag string) {
	if r.MissingTags == nil {
		r.MissingTags = make([]string, 0)
	}

	r.MissingTags = append(r.MissingTags, tag)
}

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

	fmt.Printf("map: %v, len: %d\n", d, len(d))

	i := 0
	for k := range d {
		fmt.Printf("key is %s\n", k)
		ks[i] = k
		i++
	}

	return ks
}
