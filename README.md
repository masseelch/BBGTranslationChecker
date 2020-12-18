# BBG Translations Checker
This tool checks the translation of [CPLs](https://cpl.gg/) [BBG-Mod](https://steamcommunity.com/sharedfiles/filedetails/?id=2312050357).

Example usage:
```bash
bbg-translation-checker --source ./_example --output ./_example/report 

// Output:
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

Fell free to report errors or ask questions in the [issue tracker](https://github.com/masseelch/BBGTranslationChecker/issues).