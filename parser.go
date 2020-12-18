package checker

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
)

type (
	localizedText struct {
		Translations Translations `xml:"Replace"`
	}
	gameDate struct {
		LocalizedText localizedText `xml:"LocalizedText"`
	}
)

func Parse(truth string, dir string) (*File, []*File, error) {
	// Truth
	t, err := parseFile(filepath.Join(dir, truth))
	if err != nil {
		return nil, nil, err
	}

	// Translations
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	ts := make([]*File, 0)
	for _, f := range fs {
		if f.Name() != truth && filepath.Ext(f.Name()) == ".xml" {
			t, err := parseFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, nil, err
			}

			ts = append(ts, t)
		}
	}

	return t, ts, nil
}

func parseFile(filename string) (*File, error) {
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
