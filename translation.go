package checker

type (
	File struct {
		Filename     string
		Translations Translations
	}
	Translation struct {
		Tag     string `xml:",attr"`
		Lang    string `xml:"Language,attr"`
		Message string `xml:"Text"`
	}
	Translations []Translation
)

func (ts Translations) LookupByTag(tag string) *Translation {
	for _, t := range ts {
		if t.Tag == tag {
			return &t
		}
	}

	return nil
}