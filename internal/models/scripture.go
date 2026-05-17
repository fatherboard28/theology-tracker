package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ScriptureEntityType string

const (
	ScriptureEntityWorkItem ScriptureEntityType = "WorkItem"
	ScriptureEntityTopic    ScriptureEntityType = "Topic"
	ScriptureEntitySession  ScriptureEntityType = "Session"
)

type ScriptureTag struct {
	ID         int64               `db:"id"`
	Reference  string              `db:"reference"`
	EntityType ScriptureEntityType `db:"entity_type"`
	EntityID   int64               `db:"entity_id"`
}

// ── Reference validation ──────────────────────────────────────────────────────

// validBooks is the canonical set of 3-character book abbreviations.
var validBooks = map[string]bool{
	// Old Testament
	"Gen": true, "Exo": true, "Lev": true, "Num": true, "Deu": true,
	"Jos": true, "Jdg": true, "Rut": true, "1Sa": true, "2Sa": true,
	"1Ki": true, "2Ki": true, "1Ch": true, "2Ch": true, "Ezr": true,
	"Neh": true, "Est": true, "Job": true, "Psa": true, "Pro": true,
	"Ecc": true, "Sol": true, "Isa": true, "Jer": true, "Lam": true,
	"Eze": true, "Dan": true, "Hos": true, "Joe": true, "Amo": true,
	"Oba": true, "Jon": true, "Mic": true, "Nah": true, "Hab": true,
	"Zep": true, "Hag": true, "Zec": true, "Mal": true,
	// New Testament
	"Mat": true, "Mrk": true, "Luk": true, "Jhn": true, "Act": true,
	"Rom": true, "1Co": true, "2Co": true, "Gal": true, "Eph": true,
	"Php": true, "Col": true, "1Th": true, "2Th": true, "1Ti": true,
	"2Ti": true, "Tit": true, "Phm": true, "Heb": true, "Jas": true,
	"1Pe": true, "2Pe": true, "1Jn": true, "2Jn": true, "3Jn": true,
	"Jud": true, "Rev": true,
}

// refPattern matches: "Rom 8", "Rom 8:28", "Rev 22:1–5", "Psa 119"
// The en-dash (–) in ranges is also handled.
var refPattern = regexp.MustCompile(
	`^([A-Z][a-z0-9]{2})\s+(\d+)(?::(\d+)(?:[–\-](\d+))?)?$`,
)

// ValidateScriptureReference checks format and book abbreviation.
// Returns a normalised reference string on success, or an error.
func ValidateScriptureReference(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	m := refPattern.FindStringSubmatch(ref)
	if m == nil {
		return "", fmt.Errorf("invalid format: expected e.g. "Rom 8:28" or "Gen 1:1–3"")
	}

	book := m[1]
	if !validBooks[book] {
		return "", fmt.Errorf("unknown book abbreviation: %q", book)
	}

	// Validate verse numbers are positive integers.
	if m[3] != "" {
		v, _ := strconv.Atoi(m[3])
		if v < 1 {
			return "", fmt.Errorf("verse number must be a positive integer")
		}
	}
	if m[4] != "" {
		v, _ := strconv.Atoi(m[4])
		if v < 1 {
			return "", fmt.Errorf("end verse must be a positive integer")
		}
	}

	// Normalise: replace ASCII hyphen in ranges with en-dash.
	normalised := strings.ReplaceAll(ref, "-", "–")
	return normalised, nil
}

// ChapterPrefix returns the "Book Chapter" prefix used for range queries.
// e.g. "Rom 8:28" → "Rom 8" so a LIKE 'Rom 8%' query finds all of Romans 8.
func ChapterPrefix(ref string) string {
	parts := strings.SplitN(ref, ":", 2)
	return strings.TrimSpace(parts[0])
}
