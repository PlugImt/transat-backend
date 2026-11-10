package photobooth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Request signing for the booth's private routes.
//
// Every booth has its own key: an id and a secret of at least 32 random bytes, listed on the server in
// PHOTOBOOTH_DEVICE_KEYS ("booth-1:<secret>,booth-2:<secret>") and typed once into the booth's admin panel as
// "booth-1.<secret>". The secret never travels: each request carries
//
//	X-Photobooth-Key:       booth-1
//	X-Photobooth-Timestamp: unix seconds
//	X-Photobooth-Nonce:     16 to 64 random [A-Za-z0-9_-]
//	X-Photobooth-Signature: base64url(HMAC-SHA256(secret, canonical))
//
// with canonical = "PB1\n" + METHOD + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + hex(sha256(body)).
// A request is accepted within MaxClockSkew of the server clock and only once (the nonce is remembered), so a
// captured request can neither be altered nor replayed, and a leaked log line holds nothing reusable.
const (
	HeaderKey       = "X-Photobooth-Key"
	HeaderTimestamp = "X-Photobooth-Timestamp"
	HeaderNonce     = "X-Photobooth-Nonce"
	HeaderSignature = "X-Photobooth-Signature"

	MaxClockSkew   = 5 * time.Minute
	minSecretBytes = 32
)

var (
	keyIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,32}$`)
	nonceRe = regexp.MustCompile(`^[A-Za-z0-9_-]{16,64}$`)

	ErrUnauthorized = errors.New("unauthorized")
	ErrClockSkew    = errors.New("request timestamp outside the accepted window")
)

// Verifier checks signed booth requests.
type Verifier struct {
	keys   map[string][]byte
	now    func() time.Time
	mu     sync.Mutex
	nonces map[string]time.Time
	sweep  time.Time
}

// ParseDeviceKeys reads "id:secret,id:secret". Secrets are base64url (padding optional) or raw text, and must
// decode to at least 32 bytes.
func ParseDeviceKeys(raw string) (map[string][]byte, error) {
	keys := map[string][]byte{}
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		id, secret, ok := strings.Cut(entry, ":")
		if !ok || !keyIDRe.MatchString(id) {
			return nil, fmt.Errorf("invalid photobooth key entry %q: expected id:secret", truncate(entry, 12))
		}
		decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(secret, "="))
		if err != nil {
			decoded = []byte(secret)
		}
		if len(decoded) < minSecretBytes {
			return nil, fmt.Errorf("photobooth key %q is too short: use at least %d random bytes", id, minSecretBytes)
		}
		if _, dup := keys[id]; dup {
			return nil, fmt.Errorf("photobooth key %q is listed twice", id)
		}
		keys[id] = decoded
	}
	return keys, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func NewVerifier(keys map[string][]byte) *Verifier {
	return &Verifier{keys: keys, now: time.Now, nonces: map[string]time.Time{}}
}

// Canonical is the string a booth signs.
func Canonical(method, path, timestamp, nonce string, body []byte) string {
	sum := sha256.Sum256(body)
	return "PB1\n" + strings.ToUpper(method) + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + hex.EncodeToString(sum[:])
}

// Sign is what a booth computes; exposed for tests and tooling.
func Sign(secret []byte, canonical string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(canonical))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// Verify returns the key id of a correctly signed, fresh, never-seen request.
func (v *Verifier) Verify(method, path, keyID, timestamp, nonce, signature string, body []byte) (string, error) {
	secret, ok := v.keys[keyID]
	if !ok || !keyIDRe.MatchString(keyID) || !nonceRe.MatchString(nonce) || signature == "" {
		return "", ErrUnauthorized
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", ErrUnauthorized
	}
	expected := Sign(secret, Canonical(method, path, timestamp, nonce, body))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return "", ErrUnauthorized
	}
	now := v.now()
	at := time.Unix(ts, 0)
	if at.Before(now.Add(-MaxClockSkew)) || at.After(now.Add(MaxClockSkew)) {
		return "", ErrClockSkew
	}
	if !v.remember(keyID+"/"+nonce, now) {
		return "", ErrUnauthorized
	}
	return keyID, nil
}

// remember records a nonce for twice the clock window; false when it was already used.
func (v *Verifier) remember(nonce string, now time.Time) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if now.Sub(v.sweep) > time.Minute {
		for n, exp := range v.nonces {
			if now.After(exp) {
				delete(v.nonces, n)
			}
		}
		v.sweep = now
	}
	if exp, seen := v.nonces[nonce]; seen && now.Before(exp) {
		return false
	}
	v.nonces[nonce] = now.Add(2 * MaxClockSkew)
	return true
}
