package checker

import (
	"fmt"
	"strings"
)

const (
	commentPrefix = " BBG-TRANSLATION-CHECKER-NOTES\n\t\t\t"
)

type (
	File struct {
		Filename     string
		Error        error
		Translations Translations
		rows         Translations
		replacements Translations
	}
	Translation struct {
		Comment   string `xml:",comment"`
		Tag       string `xml:",attr"`
		LangUpper string `xml:"Language,attr,omitempty"`
		LangLower string `xml:"language,attr,omitempty"`
		Message   string `xml:"Text"`
	}
	Translations []*Translation
)

func (t Translation) Lang() string {
	if t.LangUpper != "" {
		return t.LangUpper
	}

	return t.LangLower
}

func (t *Translation) AddReportToComment(r string) {
	// Make sure there is the starting part.
	if !strings.HasPrefix(t.Comment, commentPrefix) {
		t.Comment = commentPrefix
	}

	// xml.MarshalIndent does not indent multiline-comments. We do know the indent of the comment
	// since our xml-structure is well-defined. Still, this is just a workaround.
	t.Comment += fmt.Sprintf("\t\t%s\n\t\t\t", r)
}

func (t *Translation) Copy() *Translation {
	return &Translation{
		Comment:   t.Comment,
		Tag:       t.Tag,
		LangUpper: t.LangUpper,
		LangLower: t.LangLower,
		Message:   t.Message,
	}
}

func (ts Translations) LookupByTag(tag string) *Translation {
	for _, t := range ts {
		if t.Tag == tag {
			return t
		}
	}

	return nil
}

func (ts Translations) AllByTag(tag string) []*Translation {
	rs := make([]*Translation, 0)

	for _, t := range ts {
		if t.Tag == tag {
			rs = append(rs, t)
		}
	}

	return rs
}