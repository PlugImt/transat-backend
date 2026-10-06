package models

import (
	"fmt"
	"strconv"
)

// NavigationVersion is bumped on breaking changes of the NavigationTarget wire format.
const NavigationVersion = 1

// NavigationType identifies the kind of app destination a notification opens.
// To add one: declare it here, register it in navigationRules, and map it in the
// app's destination registry (src/services/notifications/destinations).
type NavigationType string

const (
	NavigationEvent       NavigationType = "event"
	NavigationClub        NavigationType = "club"
	NavigationRestaurant  NavigationType = "restaurant"
	NavigationReservation NavigationType = "reservation"
	NavigationService     NavigationType = "service"
)

// navigationRule describes what a destination type needs and how pre-navigation-context apps open it.
type navigationRule struct {
	requiresID   bool
	legacyScreen string
}

var navigationRules = map[NavigationType]navigationRule{
	NavigationEvent:       {requiresID: true, legacyScreen: "Events"},
	NavigationClub:        {requiresID: true, legacyScreen: "Clubs"},
	NavigationRestaurant:  {requiresID: false, legacyScreen: "Restaurant"},
	NavigationReservation: {requiresID: true, legacyScreen: "Reservation"},
	NavigationService:     {requiresID: true},
}

// NavigationTarget is the structured destination carried by a push notification.
type NavigationTarget struct {
	Version int            `json:"v"`
	Type    NavigationType `json:"type"`
	// ID identifies the entity (event id, club id, reservation item id, service key...).
	ID string `json:"id,omitempty"`
	// Params holds optional extras the destination can use (e.g. "date", "title").
	Params map[string]string `json:"params,omitempty"`
}

func newNavigationTarget(t NavigationType, id string, params map[string]string) *NavigationTarget {
	return &NavigationTarget{Version: NavigationVersion, Type: t, ID: id, Params: params}
}

func EventNavigation(eventID int) *NavigationTarget {
	return newNavigationTarget(NavigationEvent, strconv.Itoa(eventID), nil)
}

func ClubNavigation(clubID int) *NavigationTarget {
	return newNavigationTarget(NavigationClub, strconv.Itoa(clubID), nil)
}

func RestaurantNavigation() *NavigationTarget {
	return newNavigationTarget(NavigationRestaurant, "", nil)
}

// ReservationNavigation opens the reservation page of an item; params may carry "date" (YYYY-MM-DD) and "title".
func ReservationNavigation(itemID int, params map[string]string) *NavigationTarget {
	return newNavigationTarget(NavigationReservation, strconv.Itoa(itemID), params)
}

// ServiceNavigation opens the landing page of a service by key (e.g. "laundry").
func ServiceNavigation(service string) *NavigationTarget {
	return newNavigationTarget(NavigationService, service, nil)
}

// Validate checks a target, typically one received through the API, before it is sent to devices.
func (n *NavigationTarget) Validate() error {
	rule, ok := navigationRules[n.Type]
	if !ok {
		return fmt.Errorf("unknown navigation type %q", n.Type)
	}
	if rule.requiresID && n.ID == "" {
		return fmt.Errorf("navigation type %q requires an id", n.Type)
	}
	return nil
}

// normalized returns a copy stamped with the current wire version.
func (n *NavigationTarget) normalized() *NavigationTarget {
	return newNavigationTarget(n.Type, n.ID, n.Params)
}

// ExpoData builds the `data` object delivered with the push, or nil without navigation.
// "screen" is kept only so apps predating the navigation context still open the right tab.
func (p NotificationPayload) ExpoData() map[string]interface{} {
	if p.Navigation == nil {
		return nil
	}

	data := map[string]interface{}{"navigation": p.Navigation.normalized()}
	if rule, ok := navigationRules[p.Navigation.Type]; ok && rule.legacyScreen != "" {
		data["screen"] = rule.legacyScreen
	}
	return data
}
