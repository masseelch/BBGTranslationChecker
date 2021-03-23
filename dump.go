package checker

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (rs Reports) DumpSummary(w io.Writer) error {
	for _, r := range rs {
		if _, err := w.Write(r.dumpSummary()); err != nil {
			return err
		}
	}

	return nil
}

func (rs Reports) DumpWithComments(tf *File) error {
	for _, r := range rs {
		// Skip the truth
		if r.File.Filename != tf.Filename {
			if err := r.dumpWithComments(tf); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r Report) dumpSummary() []byte {
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
		for t, c := range r.LanguageTags {
			buf.WriteString(fmt.Sprintf("\t\t- %s: %d times\n", t, c))
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
			buf.WriteString(fmt.Sprintf("\t\t- %s:\n\t\t\tShould Have:\t%v\n\t\t\tDoes Have:\t%v\n", d.Tag, d.Truth, d.Translation))
		}
		buf.WriteString("\n")
	}

	return buf.Bytes()
}

func (r Report) dumpWithComments(tf *File) error {
	ext := filepath.Ext(r.File.Filename)

	f, err := os.Create(strings.TrimSuffix(r.File.Filename, ext) + "_commented" + ext)
	if err != nil {
		return err
	}
	defer f.Close()

	// If there is nothing do do for this file just say it.
	if r.perfect() {
		_, _ = f.WriteString(fmt.Sprintf("Nothing To Report. You Did a GREAT Job!\n"))
		return nil
	}

	// If the file had an XML-Error while parsing say so.
	if r.XMLError != nil {
		_, _ = f.WriteString(fmt.Sprintf("Invalid XML:\n\t%s\n", r.XMLError))
		return nil
	}

	// DuplicateTags
	if r.DuplicateTags != nil {
		for t, c := range r.DuplicateTags {
			for _, tr := range r.File.Translations.AllByTag(t) {
				tr.AddReportToComment(fmt.Sprintf("- This tag exists %d times", c))
			}
		}
	}

	// LanguageTags
	// In this case we loop over all translations. If they happen to have a language-tag which is
	// different from most of the other language tags it most likely is wrong.
	if r.LanguageTags != nil && len(r.LanguageTags) > 1 {
		h := r.LanguageTags.highest()
		for _, tr := range r.File.Translations {
			if tr.Lang() != h {
				tr.AddReportToComment(fmt.Sprintf("- The 'Language'-tag (%s) differs from most of the other entries (%s)", tr.Lang(), h))
			}
		}
	}

	// ObsoleteTags
	if r.ObsoleteTags != nil {
		for _, t := range r.ObsoleteTags {
			for _, tr := range r.File.Translations.AllByTag(t) {
				tr.AddReportToComment("- This tag does no longer exist in the original file")
			}
		}
	}

	// MissingTags
	// If there a translations missing, we add them to the translation. The original Text gets
	// prefixed with 'TO TRANSLATE' and a comment is added, that this entry still needs translation.
	if r.MissingTags != nil {
		for _, t := range r.MissingTags {
			// Lookup the original translation and add a copy of the original to
			// the file we are editing.
			tr := tf.rows.LookupByTag(t)
			if tr != nil {
				r.File.rows = append(r.File.rows, tr)
			} else {
				tr = tf.replacements.LookupByTag(t)
				r.File.replacements = append(r.File.replacements, tr)
			}
			r.File.Translations = append(r.File.Translations, tr)

			// Prefix the message with a 'TO TRANSLATE' notice. Add a comment.
			tr.Message = "[COLOR_RED]TO_TRANSLATE: " + tr.Message
			tr.AddReportToComment("- This tag has no translation yet")
		}
	}

	// MissingTranslations
	if r.MissingTranslations != nil {
		for _, t := range r.MissingTranslations {
			tr := r.File.Translations.LookupByTag(t)
			tr.AddReportToComment("- This tag has no translation yet")
		}
	}

	// NumericDifferences
	if r.NumericDifferences != nil {
		for _, d := range r.NumericDifferences {
			tr := r.File.Translations.LookupByTag(d.Tag)
			tr.AddReportToComment("- Difference In Values Detected (Translation Might Still Be Correct)")
			tr.AddReportToComment(fmt.Sprintf("\tShould Have %v", d.Truth))
			tr.AddReportToComment(fmt.Sprintf("\tDoes Have   %v", d.Translation))
		}
	}

	// Marshal to xml and write to file.
	x, err := xml.MarshalIndent(gameData{
		LocalizedText: localizedText{
			Replacements: r.File.replacements,
			Rows:         r.File.rows,
		},
	}, "", "\t")
	if err != nil {
		panic(err)
	}

	// One thing we do is to replace all "language"-tags with "Language" to keep the file consistent.
	x = bytes.ReplaceAll(x, []byte("language"), []byte("Language"))

	_, err = f.WriteString(xml.Header)
	if err != nil {
		return err
	}

	_, err = f.Write(x)
	return err
}

func (r Report) perfect() bool {
	return r.DuplicateTags == nil &&
		r.LanguageTags == nil &&
		r.MissingTags == nil &&
		r.MissingTranslations == nil &&
		r.NumericDifferences == nil &&
		r.XMLError == nil
}
