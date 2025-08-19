package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/utils"
)

type RestaurantMenuCronHandler interface {
	CheckAndUpdateMenuCron() (bool, error)
}

// RestaurantScheduler handles periodic tasks related to restaurant menu
type RestaurantScheduler struct {
	restaurantHandler RestaurantMenuCronHandler
	stopChan          chan struct{}
	runningChan       chan struct{}
	stopOnce          sync.Once
}

// NewRestaurantScheduler creates a new restaurant scheduler
func NewRestaurantScheduler(handler RestaurantMenuCronHandler) *RestaurantScheduler {
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
				if s.isScheduledTimeAllowed(utils.Now()) {
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
	s.stopOnce.Do(func() {
		select {
		case <-s.runningChan:
			close(s.stopChan)
		default:
			utils.LogMessage(utils.LevelWarn, "Attempted to stop restaurant scheduler that wasn't running")
		}
	})
}

func (s *RestaurantScheduler) isScheduledTimeAllowed(t time.Time) bool {
	parisTime := utils.ToParisTime(t)
	weekday := parisTime.Weekday()
	fmt.Println("-------> Weekday:", weekday)
	fmt.Println("-------> Hour Paris time:", parisTime.Hour())
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := parisTime.Hour()
	return hour >= 8 && hour < 14
}

func (s *RestaurantScheduler) checkAndUpdateMenu() {
	now := utils.Now()

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
