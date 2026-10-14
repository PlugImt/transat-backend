package push

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the timings and limits of the delivery pipeline. Every field has an environment variable, see
// ConfigFromEnv.
type Config struct {
	// PollInterval is how often the worker looks for work when nothing wakes it up.
	PollInterval time.Duration
	// BatchSize is how many messages go into one request to the provider (capped by the provider's own limit).
	BatchSize int
	// ClaimSize is how many deliveries a worker takes in one round.
	ClaimSize int
	// MaxAttempts is how many times one delivery is tried in all, the first try included.
	MaxAttempts int
	// RetryBase and RetryMax bound the exponential backoff between two attempts.
	RetryBase time.Duration
	RetryMax  time.Duration
	// Lease is how long a worker owns the deliveries it claimed; a crashed worker's claims are taken back after it.
	Lease time.Duration
	// DefaultTTL is how long a message may wait to be pushed when its sender gives no deadline. The provider is
	// told the same limit, so a message queued at APNs or FCM does not appear hours later.
	DefaultTTL time.Duration
	// ReceiptDelay is how long to wait after the provider accepted a message before asking for its receipt.
	ReceiptDelay time.Duration
	// ReceiptRetry is the delay before asking again for a receipt that is not available yet.
	ReceiptRetry time.Duration
	// ReceiptTTL is how long the provider keeps receipts; an accepted message still without one is unconfirmed.
	ReceiptTTL time.Duration
	// MaxPerSecond caps the messages handed to the provider each second; 0 means no cap.
	MaxPerSecond int
	// Retention is how long finished messages and their delivery log are kept.
	Retention time.Duration
}

// DefaultConfig returns the values used when nothing is configured.
func DefaultConfig() Config {
	return Config{
		PollInterval: 5 * time.Second,
		BatchSize:    100,
		ClaimSize:    500,
		MaxAttempts:  5,
		RetryBase:    30 * time.Second,
		RetryMax:     15 * time.Minute,
		Lease:        2 * time.Minute,
		DefaultTTL:   2 * time.Hour,
		ReceiptDelay: 15 * time.Minute,
		ReceiptRetry: 5 * time.Minute,
		ReceiptTTL:   24 * time.Hour,
		MaxPerSecond: 300,
		Retention:    30 * 24 * time.Hour,
	}
}

// ConfigFromEnv reads PUSH_POLL_INTERVAL, PUSH_BATCH_SIZE, PUSH_MAX_ATTEMPTS, PUSH_RETRY_BASE, PUSH_RETRY_MAX,
// PUSH_LEASE, PUSH_DEFAULT_TTL, PUSH_RECEIPT_DELAY, PUSH_MAX_PER_SECOND and PUSH_RETENTION. Durations are Go
// durations such as "30s" or "2h". Invalid values are ignored.
func ConfigFromEnv() Config {
	cfg := DefaultConfig()
	cfg.PollInterval = envDuration("PUSH_POLL_INTERVAL", cfg.PollInterval)
	cfg.BatchSize = envInt("PUSH_BATCH_SIZE", cfg.BatchSize, 1)
	cfg.MaxAttempts = envInt("PUSH_MAX_ATTEMPTS", cfg.MaxAttempts, 1)
	cfg.RetryBase = envDuration("PUSH_RETRY_BASE", cfg.RetryBase)
	cfg.RetryMax = envDuration("PUSH_RETRY_MAX", cfg.RetryMax)
	cfg.Lease = envDuration("PUSH_LEASE", cfg.Lease)
	cfg.DefaultTTL = envDuration("PUSH_DEFAULT_TTL", cfg.DefaultTTL)
	cfg.ReceiptDelay = envDuration("PUSH_RECEIPT_DELAY", cfg.ReceiptDelay)
	cfg.MaxPerSecond = envInt("PUSH_MAX_PER_SECOND", cfg.MaxPerSecond, 0)
	cfg.Retention = envDuration("PUSH_RETENTION", cfg.Retention)
	if cfg.RetryMax < cfg.RetryBase {
		log.Printf("Warning: PUSH_RETRY_MAX is below PUSH_RETRY_BASE, using the defaults for both")
		d := DefaultConfig()
		cfg.RetryBase, cfg.RetryMax = d.RetryBase, d.RetryMax
	}
	return cfg
}

func envDuration(name string, def time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		log.Printf("Warning: ignoring %s=%q", name, raw)
		return def
	}
	return d
}

func envInt(name string, def, min int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < min {
		log.Printf("Warning: ignoring %s=%q", name, raw)
		return def
	}
	return n
}
