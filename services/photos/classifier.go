// Package photos classifies the folders of the shared Photos Drive into albums.
// It is pure (no network, no database): the Drive sync feeds it a folder tree
// and persists the returned Results, so rules can be replayed without calling Drive.
package photos

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ClassifierVersion is stored with every result so stale classifications can be recomputed.
const ClassifierVersion = 1

type FolderKind string

const (
	KindStructuralYear     FolderKind = "structural_year"
	KindStructuralCategory FolderKind = "structural_category"
	KindGenericContainer   FolderKind = "generic_container"
	KindEvent              FolderKind = "event"
	KindEventSection       FolderKind = "event_section"
	KindUnassigned         FolderKind = "unassigned"
	KindIgnored            FolderKind = "ignored"
)

type Status string

const (
	StatusCertain   Status = "certain"
	StatusProbable  Status = "probable"
	StatusUncertain Status = "uncertain"
)

type Rule string

const (
	RuleYear                 Rule = "YEAR_FOLDER"
	RuleCategoryName         Rule = "CATEGORY_NAME"
	RuleGenericName          Rule = "GENERIC_NAME"
	RuleEventDatedStrict     Rule = "EVENT_DATED_STRICT"
	RuleEventDatedLoose      Rule = "EVENT_DATED_LOOSE"
	RuleEventDatedNested     Rule = "EVENT_DATED_NESTED"
	RuleEventUndatedInStruct Rule = "EVENT_UNDATED_IN_STRUCTURE"
	RuleEventUndatedInGener  Rule = "EVENT_UNDATED_IN_GENERIC"
	RuleEventUndatedRoot     Rule = "EVENT_UNDATED_AT_ROOT"
	RulePromoteContainer     Rule = "PROMOTE_TO_CATEGORY"
	RuleSectionInherit       Rule = "SECTION_INHERITS_EVENT"
	RuleUnassignedTechnical  Rule = "UNASSIGNED_TECHNICAL_NAME"
	RuleUnassignedNoEvent    Rule = "UNASSIGNED_UNDER_UNASSIGNED"
	RuleShortcutAlias        Rule = "SHORTCUT_ALIAS_IGNORED"
	RuleShortcutUnresolved   Rule = "SHORTCUT_TARGET_UNRESOLVED"
)

type DateSource string

const (
	DateNone       DateSource = ""
	DateFolderName DateSource = "folder_name"
	DateNormalized DateSource = "folder_name_normalized"
	DatePartial    DateSource = "folder_name_partial"
)

// Node is one Drive folder. When reached through a shortcut, ID is the target's ID
// and Name is the shortcut's name (the label chosen by the organisers).
type Node struct {
	ID               string
	Name             string
	Children         []*Node
	ShortcutID       string
	TargetUnresolved bool // shortcut whose target cannot be read
	AliasOfKnown     bool // shortcut to a folder already in the tree (cycle or duplicate)
}

type Result struct {
	FolderID          string
	Kind              FolderKind
	Rule              Rule
	Confidence        float64
	Status            Status
	AcademicYear      string
	Category          string
	AlbumRootID       string // event folder owning this folder; empty when no album
	ParentAlbumRootID string // enclosing event for a dated event nested in another event
	Title             string
	EventDate         *time.Time // UTC midnight, date-only semantics
	EventDateSource   DateSource
	SectionPath       string // path below the event root, e.g. "Roman / Scène"
	Flags             []string
}

// Publishable tells whether the folder is the root of an album that may be listed.
func (r Result) Publishable() bool {
	return r.Kind == KindEvent && r.Status != StatusUncertain
}

var (
	strictEventRe = regexp.MustCompile(`^(\d{4})\.(\d{2})\.(\d{2}) - (\S.*)$`)
	looseEventRe  = regexp.MustCompile(`^(\d{4})\s*[.\-/]\s*(\d{1,2})\s*[.\-/]\s*(\d{1,2})(?:\s*[-–—:]\s*|\s+|$)(.*)$`)
	yearRe        = regexp.MustCompile(`^(\d{4})\s*[-–/]\s*(\d{4})$`)
	partialDateRe = regexp.MustCompile(`(?:^|\s)(\d{1,2})[-/.](\d{1,2})$`)
	cameraDirRe   = regexp.MustCompile(`^\d{3}[a-z0-9_]{3,}$`)
	leadingOrder  = regexp.MustCompile(`^\d+[\s.\-_]+`)
	spaces        = regexp.MustCompile(`\s+`)

	categoryNames = set("sport", "sports")
	genericNames  = set("autre", "autres", "divers", "vrac", "other", "misc")
	technicalName = set(
		"photobooth", "photo", "photos", "autres photos", "selection", "la selection",
		"retouchee", "retouchees", "retouche", "retouches", "gopro", "a trier", "raw",
		"jpg", "jpeg", "original", "originaux", "originals", "export", "exports",
		"brut", "bruts", "video", "videos", "dcim",
	)
)

func set(values ...string) map[string]bool {
	m := make(map[string]bool, len(values))
	for _, v := range values {
		m[v] = true
	}
	return m
}

var accentStripper = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// normalize lowercases, strips accents and collapses spaces for name comparisons.
func normalize(name string) string {
	s, _, err := transform.String(accentStripper, name)
	if err != nil {
		s = name
	}
	return spaces.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), " ")
}

func isTechnical(name string) bool {
	n := normalize(name)
	if technicalName[n] || cameraDirRe.MatchString(n) {
		return true
	}
	return technicalName[leadingOrder.ReplaceAllString(n, "")]
}

func validDate(y, m, d int) (time.Time, bool) {
	if y < 2000 || y > 2100 {
		return time.Time{}, false
	}
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return t, t.Year() == y && int(t.Month()) == m && t.Day() == d
}

type dated struct {
	title  string
	date   time.Time
	strict bool
}

func parseDated(name string) (dated, bool) {
	name = strings.TrimSpace(name)
	if m := strictEventRe.FindStringSubmatch(name); m != nil {
		y, _ := strconv.Atoi(m[1])
		mo, _ := strconv.Atoi(m[2])
		d, _ := strconv.Atoi(m[3])
		if t, ok := validDate(y, mo, d); ok {
			return dated{title: m[4], date: t, strict: true}, true
		}
		return dated{}, false
	}
	if m := looseEventRe.FindStringSubmatch(name); m != nil {
		y, _ := strconv.Atoi(m[1])
		mo, _ := strconv.Atoi(m[2])
		d, _ := strconv.Atoi(m[3])
		if t, ok := validDate(y, mo, d); ok {
			return dated{title: m[4], date: t}, true
		}
	}
	return dated{}, false
}

func parseYear(name string) (string, int, bool) {
	m := yearRe.FindStringSubmatch(strings.TrimSpace(name))
	if m == nil {
		return "", 0, false
	}
	a, _ := strconv.Atoi(m[1])
	b, _ := strconv.Atoi(m[2])
	if b != a+1 {
		return "", 0, false
	}
	return m[1] + "-" + m[2], a, true
}

// academicYearOf maps a date to its university year (August to July).
func academicYearOf(t time.Time) (string, int) {
	start := t.Year()
	if t.Month() < time.August {
		start--
	}
	return strconv.Itoa(start) + "-" + strconv.Itoa(start+1), start
}

func hasDatedDescendant(n *Node) bool {
	for _, c := range n.Children {
		if _, ok := parseDated(c.Name); ok || hasDatedDescendant(c) {
			return true
		}
	}
	return false
}

func statusOf(confidence float64) Status {
	switch {
	case confidence >= 0.9:
		return StatusCertain
	case confidence >= 0.6:
		return StatusProbable
	default:
		return StatusUncertain
	}
}

type eventCtx struct {
	rootID     string
	confidence float64
	status     Status
	path       []string
}

type ctx struct {
	year       string
	yearStart  int
	category   string
	parentKind FolderKind
	event      *eventCtx
}

// Classify returns one Result per node reachable from roots (ignored aliases are not descended).
func Classify(roots []*Node) map[string]Result {
	out := make(map[string]Result)
	for _, r := range roots {
		walk(r, ctx{}, out)
	}
	return out
}

func walk(n *Node, c ctx, out map[string]Result) {
	res := Result{FolderID: n.ID, AcademicYear: c.year, Category: c.category}
	child := c

	switch {
	case n.AliasOfKnown:
		res.Kind, res.Rule, res.Confidence = KindIgnored, RuleShortcutAlias, 1
		res.Status = StatusCertain
		out[n.ID] = res
		return

	case c.event != nil:
		if d, ok := parseDated(n.Name); ok {
			rule := RuleEventDatedNested
			fillEvent(&res, n, c, d, 0.7, rule)
			res.ParentAlbumRootID = c.event.rootID
			res.Flags = append(res.Flags, "nested_dated_event")
			child = ctx{year: c.year, yearStart: c.yearStart, category: c.category, parentKind: KindEvent, event: &eventCtx{rootID: n.ID, confidence: res.Confidence, status: res.Status}}
		} else {
			path := append(append([]string{}, c.event.path...), strings.TrimSpace(n.Name))
			res.Kind, res.Rule = KindEventSection, RuleSectionInherit
			res.Confidence, res.Status = c.event.confidence, c.event.status
			res.AlbumRootID = c.event.rootID
			res.SectionPath = strings.Join(path, " / ")
			child.event = &eventCtx{rootID: c.event.rootID, confidence: c.event.confidence, status: c.event.status, path: path}
			child.parentKind = KindEventSection
		}

	default:
		if d, ok := parseDated(n.Name); ok {
			rule, conf := RuleEventDatedStrict, 1.0
			if !d.strict {
				rule, conf = RuleEventDatedLoose, 0.95
			}
			fillEvent(&res, n, c, d, conf, rule)
			child.event = &eventCtx{rootID: n.ID, confidence: res.Confidence, status: res.Status}
			child.parentKind = KindEvent
			break
		}
		if label, start, ok := parseYear(n.Name); ok {
			res.Kind, res.Rule, res.Confidence, res.Status = KindStructuralYear, RuleYear, 1, StatusCertain
			res.AcademicYear = label
			child.year, child.yearStart, child.parentKind = label, start, KindStructuralYear
			break
		}
		key := normalize(n.Name)
		if categoryNames[key] {
			res.Kind, res.Rule, res.Confidence, res.Status = KindStructuralCategory, RuleCategoryName, 1, StatusCertain
			res.Category = strings.TrimSpace(n.Name)
			child.category, child.parentKind = res.Category, KindStructuralCategory
			break
		}
		if genericNames[key] {
			res.Kind, res.Rule, res.Confidence, res.Status = KindGenericContainer, RuleGenericName, 1, StatusCertain
			child.parentKind = KindGenericContainer
			break
		}
		if isTechnical(n.Name) {
			res.Kind, res.Rule, res.Confidence, res.Status = KindUnassigned, RuleUnassignedTechnical, 0.3, StatusUncertain
			child.parentKind = KindUnassigned
			break
		}
		if c.parentKind == KindUnassigned {
			res.Kind, res.Rule, res.Confidence, res.Status = KindUnassigned, RuleUnassignedNoEvent, 0.3, StatusUncertain
			child.parentKind = KindUnassigned
			break
		}
		if hasDatedDescendant(n) {
			res.Kind, res.Rule, res.Confidence, res.Status = KindStructuralCategory, RulePromoteContainer, 0.7, StatusProbable
			res.Category = strings.TrimSpace(n.Name)
			child.category, child.parentKind = res.Category, KindStructuralCategory
			break
		}
		rule, conf := RuleEventUndatedRoot, 0.4
		switch c.parentKind {
		case KindStructuralYear, KindStructuralCategory:
			rule, conf = RuleEventUndatedInStruct, 0.7
		case KindGenericContainer:
			rule, conf = RuleEventUndatedInGener, 0.4
		}
		fillEvent(&res, n, c, dated{title: strings.TrimSpace(n.Name)}, conf, rule)
		child.event = &eventCtx{rootID: n.ID, confidence: res.Confidence, status: res.Status}
		child.parentKind = KindEvent
	}

	if n.TargetUnresolved {
		res.Rule, res.Confidence, res.Status = RuleShortcutUnresolved, 0.3, StatusUncertain
		res.Flags = append(res.Flags, "target_unresolved")
		if child.event != nil && child.event.rootID == n.ID {
			child.event.confidence, child.event.status = res.Confidence, res.Status
		}
	}
	if n.ShortcutID != "" {
		res.Flags = append(res.Flags, "via_shortcut")
	}
	out[n.ID] = res
	for _, k := range n.Children {
		walk(k, child, out)
	}
}

// fillEvent sets the event-specific fields (title, date, year) of res.
func fillEvent(res *Result, n *Node, c ctx, d dated, conf float64, rule Rule) {
	res.Kind, res.Rule, res.Confidence = KindEvent, rule, conf
	res.AlbumRootID = n.ID
	res.Title = strings.TrimSpace(spaces.ReplaceAllString(d.title, " "))

	switch {
	case !d.date.IsZero():
		t := d.date
		res.EventDate = &t
		res.EventDateSource = DateFolderName
		if !d.strict {
			res.EventDateSource = DateNormalized
			res.Flags = append(res.Flags, "name_normalized")
		}
		if res.Title == "" {
			res.Title = strings.TrimSpace(n.Name)
			res.Flags = append(res.Flags, "title_from_date_only")
		}
		year, start := academicYearOf(t)
		if c.year == "" {
			res.AcademicYear = year
			res.Flags = append(res.Flags, "year_from_date")
		} else if start != c.yearStart {
			res.Flags = append(res.Flags, "year_mismatch")
		}
	case c.year != "":
		if m := partialDateRe.FindStringSubmatch(res.Title); m != nil {
			day, _ := strconv.Atoi(m[1])
			month, _ := strconv.Atoi(m[2])
			y := c.yearStart + 1
			if month >= 8 {
				y = c.yearStart
			}
			if t, ok := validDate(y, month, day); ok {
				res.EventDate = &t
				res.EventDateSource = DatePartial
				res.Flags = append(res.Flags, "date_partial")
				if res.Confidence > 0.65 {
					res.Confidence = 0.65
				}
			}
		}
	}
	res.Status = statusOf(res.Confidence)
}
