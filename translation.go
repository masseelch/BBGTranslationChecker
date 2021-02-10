package checker

type (
	File struct {
		Filename     string
		Error        error
		Translations Translations
	}
	Translation struct {
		Tag       string `xml:",attr"`
		LangUpper string `xml:"Language,attr"`
		LangLower string `xml:"language,attr"`
		Message   string `xml:"Text"`
	}
	Translations []Translation
)

func (t Translation) Lang() string {
	if t.LangUpper != "" {
		return t.LangUpper
	}

	return t.LangLower
}

func (ts Translations) LookupByTag(tag string) *Translation {
	for _, t := range ts {
		if t.Tag == tag {
			return &t
		}
	}

	return nil
}
