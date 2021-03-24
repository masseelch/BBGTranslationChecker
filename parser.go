package checker

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type (
	localizedText struct {
		Replacements Translations `xml:"Replace"`
		Rows         Translations `xml:"Row"`
	}
	gameData struct {
		XMLName       string        `xml:"GameData"`
		LocalizedText localizedText `xml:"LocalizedText"`
	}
)

func Parse(truth, dir, only string) (*File, []*File, error) {
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
		if f.Name() != truth && !strings.Contains(f.Name(), "_commented") && filepath.Ext(f.Name()) == ".xml" && (only == "" || f.Name() == only) {
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

	gd := new(gameData)
	err = xml.Unmarshal(b, gd)
	if err != nil {
		return &File{
			Filename: filepath.Base(filename),
			Error:    err,
		}, nil
	}

	trs := make(Translations, len(gd.LocalizedText.Replacements), len(gd.LocalizedText.Replacements)+len(gd.LocalizedText.Rows))
	_ = copy(trs, gd.LocalizedText.Replacements)
	trs = append(trs, gd.LocalizedText.Rows...)

	return &File{
		Filename:     filepath.Base(filename),
		Translations: trs,
		rows:         gd.LocalizedText.Rows,
		replacements: gd.LocalizedText.Replacements,
	}, nil
}
