package models

import (
	"fmt"
	"regexp"
	"strings"
)

type AudienceType string

const (
	AudienceAll    AudienceType = "all"
	AudienceClub   AudienceType = "club"
	AudienceCampus AudienceType = "campus"
	AudienceUsers  AudienceType = "users"
	// AudienceCohort targets a formation (FIL, FISE...), a graduation year, or both.
	AudienceCohort AudienceType = "cohort"
	// AudienceAppVersion targets the devices running an app older than BelowVersion (to ask them to update).
	// Only those devices are notified, not the users' other, up-to-date devices.
	AudienceAppVersion AudienceType = "app_version"
)

var appVersionRe = regexp.MustCompile(`^[0-9]{1,4}(\.[0-9]{1,4}){0,3}$`)

var audiencePlatforms = map[string]bool{"": true, "ios": true, "android": true}

var audienceCampuses = map[string]bool{"NANTES": true, "BREST": true, "RENNES": true}

// Matches the newf.formation_name constraint.
var audienceFormations = map[string]bool{"FISE": true, "FIL": true, "FIP": true, "FIT": true, "FID": true}

const (
	minGraduationYear = 2000
	maxGraduationYear = 2100
)

// NotificationAudience selects the users an admin notification is sent to.
type NotificationAudience struct {
	Type           AudienceType `json:"type"`
	ClubID         int          `json:"clubId,omitempty"`
	Campus         string       `json:"campus,omitempty"`
	Emails         []string     `json:"emails,omitempty"`
	Formation      string       `json:"formation,omitempty"`
	GraduationYear int          `json:"graduationYear,omitempty"`
	// BelowVersion and Platform go with AudienceAppVersion; an empty platform means every platform.
	BelowVersion string `json:"belowVersion,omitempty"`
	Platform     string `json:"platform,omitempty"`
}

func (a NotificationAudience) Validate() error {
	switch a.Type {
	case AudienceAll:
	case AudienceClub:
		if a.ClubID <= 0 {
			return fmt.Errorf("audience club requires a clubId")
		}
	case AudienceCampus:
		if !audienceCampuses[strings.ToUpper(a.Campus)] {
			return fmt.Errorf("unknown campus %q", a.Campus)
		}
	case AudienceUsers:
		if len(a.Emails) == 0 {
			return fmt.Errorf("audience users requires at least one email")
		}
	case AudienceCohort:
		if a.Formation == "" && a.GraduationYear == 0 {
			return fmt.Errorf("audience cohort requires a formation or a graduationYear")
		}
		if a.Formation != "" && !audienceFormations[strings.ToUpper(a.Formation)] {
			return fmt.Errorf("unknown formation %q", a.Formation)
		}
		if a.GraduationYear != 0 && (a.GraduationYear < minGraduationYear || a.GraduationYear > maxGraduationYear) {
			return fmt.Errorf("graduationYear out of range")
		}
	case AudienceAppVersion:
		if !appVersionRe.MatchString(a.BelowVersion) {
			return fmt.Errorf("audience app_version requires a belowVersion such as 2.4.0")
		}
		if !audiencePlatforms[strings.ToLower(a.Platform)] {
			return fmt.Errorf("unknown platform %q", a.Platform)
		}
	default:
		return fmt.Errorf("unknown audience type %q", a.Type)
	}
	return nil
}
