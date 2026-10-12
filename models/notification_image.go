package models

import (
	"fmt"
	"net/url"
	"strings"
)

// MaxNotificationImageURLLength keeps the link small: it travels inside the 4 KiB push payload.
const MaxNotificationImageURLLength = 512

// ValidateImageURL checks the link of a notification image. The phone downloads it, so it must be a public https
// address; the server never fetches it and cannot tell whether it is reachable or how large the file is.
//
// What happens to the image depends on the platform, and a notification is always delivered as text when the
// image cannot be shown:
//   - Android: sent as the provider's richContent.image, which FCM and the OS fetch and display on their own.
//     Size and format limits are theirs.
//   - iOS: shown only by the app's notification service extension (transat-app/targets/notification-service),
//     which needs the push to be mutable-content (set whenever an image is sent). The extension downloads
//     JPEG, PNG or GIF files up to 5 MB over https and drops anything else; an app build without the
//     extension shows the text only.
//   - Both: the link counts toward the 4 KiB payload limit.
func ValidateImageURL(raw string) error {
	if len(raw) > MaxNotificationImageURLLength {
		return fmt.Errorf("image url is too long")
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || strings.ContainsAny(raw, " \t\r\n") {
		return fmt.Errorf("image url must be an absolute https link")
	}
	return nil
}

// OwnedImageURL returns the link when it is a valid notification image hosted on the same host as base (our public
// storage domain), and "" otherwise. A link chosen by a user (an event poster) is only trusted when it points to
// our own storage: every recipient's phone fetches it, which would otherwise show their address to whoever
// controls that host. A producer that gets "" sends a text notification.
func OwnedImageURL(raw, base string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || ValidateImageURL(raw) != nil {
		return ""
	}
	owner, err := url.Parse(strings.TrimSpace(base))
	if err != nil || owner.Hostname() == "" {
		return ""
	}
	parsed, _ := url.Parse(raw)
	if !strings.EqualFold(parsed.Host, owner.Host) {
		return ""
	}
	return raw
}
