# BBG Translations Checker
This tool checks the translation of [CPLs](https://cpl.gg/) [BBG-Mod](https://steamcommunity.com/sharedfiles/filedetails/?id=2312050357).

### How to use
- Download the correct executable for your platform. 
- Put the executable in a folder containing the xml-files you want to check.
- Run the executable.

You will then find a new file `reports.txt` containing information about all issued found and some files ending in `_commented.xml`.
These files are copies of the translations checked annotated with the same information found in `reports.txt`. 

Now just `STRG+F` for `BBG-TRANSLATION-CHECKER-NOTES` and fix the issues.

#### Configuration
The checker has some configuration options. The following are the defaults.
```yaml
# config.yaml

source: .          # path/to/folder/containing/translation/files
truth: english.xml # filename (relative to source) of the file to use as reference
only: ~            # if set only checks filename given (relative to source)
```

#### Example Output
```
// reports.txt:
german.xml:
	Duplicate Entries:
		- LOC_UNIT_INQUISITOR_EXPANSION2_DESCRIPTION: 2 times
		[...]

	Obsolete Translations:
		- LOC_BUILDING_STABLE_DESCRIPTION
		- LOC_BUILDING_MILITARY_ACADEMY_DESCRIPTION
		[...]

	Missing Translations (No Entry):
		- LOC_UNIT_COLOMBIAN_LLANERO_DESCRIPTION
		- LOC_ABILITY_LLANERO_ADJACENCY_STRENGTH_DESCRIPTION
		[...]

	Missing Translations (Untranslated Entry):
		- LOC_LEADER_TRAIT_FEZ_DESCRIPTION
		[...]

	Differences In Values (Translation Might Still Be Correct):
		- LOC_TRAIT_LEADER_AMBIORIX_DESCRIPTION:
			Should Have:[1]
			Does Have:[2]
		[...]
```
```xml
<!-- german_commented.xml: -->
<Replace Tag="LOC_TRAIT_LEADER_RELIGION_CITY_STATES_DESCRIPTION" Language="de_DE">
    <!-- BBG-TRANSLATION-CHECKER-NOTES
            - Difference In Values Detected (Translation Might Still Be Correct)
                Should Have [50]
                Does Have   [2 50]
    -->
    <Text>[...]</Text>
</Replace>
<Replace Tag="LOC_IMPROVEMENT_FISHING_BOATS_DESCRIPTION" Language="de_DE">
    <!-- BBG-TRANSLATION-CHECKER-NOTES
            - This tag exists 2 times
    -->
    <Text>[...]</Text>
</Replace>
<Replace Tag="LOC_IMPROVEMENT_FISHING_BOATS_DESCRIPTION" Language="de_DE">
    <!-- BBG-TRANSLATION-CHECKER-NOTES
            - This tag exists 2 times
    -->
    <Text>[...]</Text>
</Replace>
<Replace Tag="LOC_BELIEF_JUST_WAR_DESCRIPTIONS" Language="de">
<!-- BBG-TRANSLATION-CHECKER-NOTES
        - The 'Language'-tag (de) differs from most of the other entries (de_DE)
        - This tag does no longer exist in the original file
-->
<Text>[...]</Text>
</Replace>
```

Fell free to report errors or ask questions in the [issue tracker](https://github.com/masseelch/BBGTranslationChecker/issues).