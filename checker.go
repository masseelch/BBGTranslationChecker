package checker

import (
	"regexp"
	"sort"
	"strings"
)

const (
	translationMarker = "to translate"
)

var (
	extractNumericValuesRegex = regexp.MustCompile("[0-9]+")
)

type (
	// A report for a translation file.
	Report struct {
		File *File
		// Parsing XML error
		XMLError error
		// All tags that occur more than once.
		DuplicateTags Duplicates
		// All tags that occur in the translation file but not the truth.
		ObsoleteTags []string
		// All found language declarations. [tag -> count]
		LanguageTags Duplicates
		// Tags that do not appear in the given file but do in the truth.
		MissingTags []string
		// Tags do appear in the given file but are not translated yet (prefixed with "To Translate")
		MissingTranslations []string
		// Where do the numeric values in the translation differ from the truth?
		NumericDifferences []NumericDifference
	}
	// Key: filename
	Reports map[string]*Report
)

func Check(tf *File, cfs []*File) (Reports, error) {
	var err error
	rs := make(Reports)

	// Collect (limited) reports for the truth
	rs[tf.Filename], err = report(nil, tf)
	if err != nil {
		return nil, err
	}

	// Collect the report for every translation.
	for _, cf := range cfs {
		rs[cf.Filename], err = report(tf, cf)
		if err != nil {
			return nil, err
		}
	}

	return rs, nil
}

// t: File of truth
// c: File to check
func report(tf *File, cf *File) (*Report, error) {
	r := new(Report)
	r.File = cf

	// If there has been error while parsing report it.
	if r.File.Error != nil {
		r.XMLError = cf.Error
		return r, nil
	}

	// duplicate tags / lang
	r.duplicatesCheck(cf)

	// If not truth file is given skip the parts where is is mandatory.
	if tf == nil {
		return r, nil
	}

	// obsolete translations
	r.obsoletesCheck(tf, cf)

	// missing translations
	r.translationsCheck(tf, cf)

	return r, nil
}

func (r *Report) duplicatesCheck(f *File) {
	// duplicate tags / lang
	r.DuplicateTags = make(Duplicates)
	r.LanguageTags = make(Duplicates)

	// Count the occurrences of every tag.
	for _, t := range f.Translations {
		r.DuplicateTags[t.Tag]++
		r.LanguageTags[t.Lang()]++
	}

	// Only keep those entries that occur more than once.
	r.DuplicateTags = r.DuplicateTags.cleaned()
}

func (r *Report) obsoletesCheck(tf *File, cf *File) {
	r.ObsoleteTags = make([]string, 0)

	// For every tag in the check file check if there exists an entry in the truth file.
	// If there is no tag found mark the tag as obsolete.
	for _, ct := range cf.Translations {
		if tf.Translations.LookupByTag(ct.Tag) == nil {
			r.ObsoleteTags = append(r.ObsoleteTags, ct.Tag)
		}
	}

	// Clean up
	if len(r.ObsoleteTags) == 0 {
		r.ObsoleteTags = nil
	}
}

func (r *Report) translationsCheck(tf *File, cf *File) {
	r.MissingTags = make([]string, 0)
	r.MissingTranslations = make([]string, 0)
	r.NumericDifferences = make([]NumericDifference, 0)

	// For every tag in the truth file check if there exists a translation in the file to check.
	// If found validate that the translation does no longer contain "To Translate".
	// If found check the numeric values in truth and checked file match.
	for _, tt := range tf.Translations {
		ct := cf.Translations.LookupByTag(tt.Tag)

		// If the tag is not found mark it and skip to the next entry.
		if ct == nil {
			r.MissingTags = append(r.MissingTags, tt.Tag)
			continue
		}

		// If the translation exists but has the translationMarker still present we report that and
		// won't check for numeric differences.
		if strings.Contains(strings.ToLower(ct.Message), translationMarker) {
			r.MissingTranslations = append(r.MissingTranslations, tt.Tag)
		}

		// Check if the numeric values in check match the truth.
		tn := extractNumericValues(tt.Message)
		cn := extractNumericValues(ct.Message)

		// If the amount of values extracted differ we can report a difference without further checking.
		if len(tn) != len(cn) {
			r.NumericDifferences = append(r.NumericDifferences, NumericDifference{
				Tag:         tt.Tag,
				Truth:       tn,
				Translation: cn,
			})
		} else {
			for i, v := range tn {
				if v != cn[i] {
					r.NumericDifferences = append(r.NumericDifferences, NumericDifference{
						Tag:         tt.Tag,
						Truth:       tn,
						Translation: cn,
					})
					// If there is at least one difference we do not need to check the remaining values.
					break
				}
			}
		}
	}

	// Clean up
	if len(r.MissingTags) == 0 {
		r.MissingTags = nil
	}
	if len(r.MissingTranslations) == 0 {
		r.MissingTranslations = nil
	}
	if len(r.NumericDifferences) == 0 {
		r.NumericDifferences = nil
	}
}

// Extracts all non overlapping numeric values. The resulting slice is sorted lexicographically.
func extractNumericValues(msg string) []string {
	nvs := extractNumericValuesRegex.FindAllString(msg, -1)
	sort.Slice(nvs, func(i, j int) bool {
		return nvs[i] < nvs[j]
	})

	return nvs
}
