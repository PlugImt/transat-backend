package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/plugimt/transat-backend/services/reservations"
)

// ReservationScheduler sends the reservation reminders: one hour before a slot, and one, five and seven days
// into a loan of an item without a slot. Which reminder is due, and whether it was already sent, is decided by
// the reservations service from the database, so a restart or a second instance never duplicates one.
type ReservationScheduler struct {
	service  *reservations.Service
	interval time.Duration
	stopChan chan struct{}
}

// NewReservationScheduler creates the reservation scheduler.
func NewReservationScheduler(service *reservations.Service) *ReservationScheduler {
	return &ReservationScheduler{service: service, interval: time.Minute, stopChan: make(chan struct{})}
}

// Start starts the scheduler loop.
func (rs *ReservationScheduler) Start() {
	go rs.run()
	log.Println("✅ Reservation scheduler started")
}

// Stop stops the scheduler loop.
func (rs *ReservationScheduler) Stop() {
	close(rs.stopChan)
	log.Println("🛑 Reservation scheduler stopped")
}

func (rs *ReservationScheduler) run() {
	ticker := time.NewTicker(rs.interval)
	defer ticker.Stop()
	for {
		select {
		case <-rs.stopChan:
			return
		case <-ticker.C:
			rs.tick()
		}
	}
}

func (rs *ReservationScheduler) tick() {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	sent, err := rs.service.Tick(ctx)
	if err != nil {
		log.Printf("Reservation reminders failed: %v", err)
		return
	}
	if sent > 0 {
		log.Printf("Queued %d reservation reminder(s)", sent)
	}
}
