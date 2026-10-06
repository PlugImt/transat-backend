package models

import (
	"fmt"
	"strings"
)

type AudienceType string

const (
	AudienceAll    AudienceType = "all"
	AudienceClub   AudienceType = "club"
	AudienceCampus AudienceType = "campus"
	AudienceUsers  AudienceType = "users"
)

var audienceCampuses = map[string]bool{"NANTES": true, "BREST": true, "RENNES": true}

// NotificationAudience selects the users an admin notification is sent to.
type NotificationAudience struct {
	Type   AudienceType `json:"type"`
	ClubID int          `json:"clubId,omitempty"`
	Campus string       `json:"campus,omitempty"`
	Emails []string     `json:"emails,omitempty"`
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
	default:
		return fmt.Errorf("unknown audience type %q", a.Type)
	}
	return nil
}
