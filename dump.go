package checker

import (
	"bytes"
	"fmt"
	"io"
)

func (rs Reports) FDump(w io.Writer) error {
	for _, r := range rs {
		if _, err := w.Write(r.dump()); err != nil {
			return err
		}
	}

	return nil
}

func (r Report) dump() []byte {
	buf := new(bytes.Buffer)

	// Header
	buf.WriteString(fmt.Sprintf("%s:\n", r.File.Filename))

	// If there is nothing do do for this file just say it.
	if r.perfect() {
		buf.WriteString(fmt.Sprintf("\tNothing To Report. You Did a GREAT Job!\n\n"))
		return buf.Bytes()
	}

	// Report XML Error
	if r.XMLError != nil {
		buf.WriteString(fmt.Sprintf("\tInvalid XML:\n\t\t%s\n\n", r.XMLError))
		return buf.Bytes()
	}

	// Duplicates
	if r.DuplicateTags != nil {
		buf.WriteString("\tDuplicate Entries:\n")
		for t, c := range r.DuplicateTags {
			buf.WriteString(fmt.Sprintf("\t\t- %s: %d times\n", t, c))
		}
		buf.WriteString("\n")
	}

	// LanguageTags
	if r.LanguageTags != nil && len(r.LanguageTags) > 1 {
		buf.WriteString("\tMultiple Language Tags Found:\n")
		for _, t := range r.LanguageTags {
			buf.WriteString(fmt.Sprintf("\t\t- %s\n", t))
		}
		buf.WriteString("\n")
	}

	// ObsoleteTags
	if r.ObsoleteTags != nil {
		buf.WriteString("\tObsolete Translations:\n")
		for _, t := range r.ObsoleteTags {
			buf.WriteString(fmt.Sprintf("\t\t- %s\n", t))
		}
		buf.WriteString("\n")
	}

	// MissingTags
	if r.MissingTags != nil {
		buf.WriteString("\tMissing Translations (No Entry):\n")
		for _, t := range r.MissingTags {
			buf.WriteString(fmt.Sprintf("\t\t- %s\n", t))
		}
		buf.WriteString("\n")
	}

	// MissingTranslations
	if r.MissingTranslations != nil {
		buf.WriteString("\tMissing Translations (Untranslated Entry):\n")
		for _, t := range r.MissingTranslations {
			buf.WriteString(fmt.Sprintf("\t\t- %s\n", t))
		}
		buf.WriteString("\n")
	}

	// NumericDifferences
	if r.NumericDifferences != nil {
		buf.WriteString("\tDifferences In Values (Translation Might Still Be Correct):\n")
		for _, d := range r.NumericDifferences {
			buf.WriteString(fmt.Sprintf("\t\t- %s:\n\t\t\tShould Have:%v\n\t\t\tDoes Have:%v\n", d.Tag, d.Truth, d.Translation))
		}
		buf.WriteString("\n")
	}

	return buf.Bytes()
}

func (r Report) perfect() bool {
	return r.DuplicateTags == nil &&
		r.LanguageTags == nil &&
		r.MissingTags == nil &&
		r.MissingTranslations == nil &&
		r.NumericDifferences == nil &&
		r.XMLError == nil
}
