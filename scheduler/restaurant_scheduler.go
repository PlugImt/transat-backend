package scheduler

import (
	"log"
	"time"

	"github.com/plugimt/transat-backend/handlers/restaurant"
	"github.com/plugimt/transat-backend/utils"
)

// RestaurantScheduler handles periodic tasks related to restaurant menu
type RestaurantScheduler struct {
	restaurantHandler *restaurant.RestaurantHandler
	stopChan          chan struct{}
	runningChan       chan struct{}
}

// NewRestaurantScheduler creates a new restaurant scheduler
func NewRestaurantScheduler(handler *restaurant.RestaurantHandler) *RestaurantScheduler {
	return &RestaurantScheduler{
		restaurantHandler: handler,
		stopChan:          make(chan struct{}),
		runningChan:       make(chan struct{}),
	}
}

func (s *RestaurantScheduler) Start() {
	utils.LogMessage(utils.LevelInfo, "Starting restaurant menu scheduler")

	go s.checkAndUpdateMenu()

	go func() {
		close(s.runningChan)

		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if s.isScheduledTimeAllowed(time.Now()) {
					go s.checkAndUpdateMenu()
				} else {
					utils.LogMessage(utils.LevelDebug, "Scheduler: Outside allowed time window (weekdays 9h-14h), skipping check")
				}
			case <-s.stopChan:
				utils.LogMessage(utils.LevelInfo, "Stopping restaurant menu scheduler")
				return
			}
		}
	}()
}

// Stop halts the scheduler
func (s *RestaurantScheduler) Stop() {
	// Only stop if it's running
	select {
	case <-s.runningChan:
		// It's running, so we can stop it
		close(s.stopChan)
	default:
		// It's not running
		utils.LogMessage(utils.LevelWarn, "Attempted to stop restaurant scheduler that wasn't running")
	}
}

func (s *RestaurantScheduler) isScheduledTimeAllowed(t time.Time) bool {
	// Check if it's a weekend (Saturday = 6, Sunday = 0)
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := t.Hour()
	return hour >= 9 && hour < 14
}

func (s *RestaurantScheduler) checkAndUpdateMenu() {
	now := time.Now()

	if !s.isScheduledTimeAllowed(now) {
		utils.LogMessage(utils.LevelInfo, "Scheduler: Outside allowed time window (weekdays 9h-14h), skipping menu check")
		return
	}

	utils.LogMessage(utils.LevelInfo, "Running scheduled menu check")

	updated, err := s.restaurantHandler.CheckAndUpdateMenuCron()
	if err != nil {
		log.Printf("Error in scheduled menu check: %v", err)
	} else if updated {
		utils.LogMessage(utils.LevelInfo, "Menu was updated during scheduled check")
	} else {
		utils.LogMessage(utils.LevelInfo, "No menu update needed during scheduled check")
	}
}
