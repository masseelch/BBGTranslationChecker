package checker

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
)

type (
	File struct {
		Filename     string
		Translations []Translation
	}
	Translation struct {
		Tag     string `xml:",attr"`
		Lang    string `xml:"Language,attr"`
		Message string `xml:"Text"`
	}
	localizedText struct {
		Translations []Translation `xml:"Replace"`
	}
	gameDate struct {
		LocalizedText localizedText `xml:"LocalizedText"`
	}
)

func ParseFile(filename string) (*File, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	gd := new(gameDate)
	return &File{
		Filename:     filepath.Base(filename),
		Translations: gd.LocalizedText.Translations,
	}, xml.Unmarshal(b, gd)
}
