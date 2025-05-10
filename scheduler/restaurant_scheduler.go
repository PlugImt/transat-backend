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

// Start begins the periodic menu check schedule
func (s *RestaurantScheduler) Start() {
	utils.LogMessage(utils.LevelInfo, "Starting restaurant menu scheduler")
	
	// Run immediately at startup
	go s.checkAndUpdateMenu()
	
	go func() {
		// Signal that the scheduler is running
		close(s.runningChan)
		
		// Create a ticker that ticks every 10 minutes
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				go s.checkAndUpdateMenu() // Run in goroutine to avoid blocking the ticker
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

// checkAndUpdateMenu calls the restaurant handler to check for menu updates
func (s *RestaurantScheduler) checkAndUpdateMenu() {
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