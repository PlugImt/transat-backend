package scheduler

import (
	"github.com/plugimt/transat-backend/handlers/restaurant"
)

// Scheduler manages all scheduled tasks in the application
type Scheduler struct {
	restaurantScheduler *RestaurantScheduler
	// Add other schedulers here as needed
}

// NewScheduler creates a new main scheduler
func NewScheduler(restaurantHandler *restaurant.RestaurantHandler) *Scheduler {
	return &Scheduler{
		restaurantScheduler: NewRestaurantScheduler(restaurantHandler),
		// Initialize other schedulers here
	}
}

// StartAll starts all schedulers
func (s *Scheduler) StartAll() {
	// Start restaurant scheduler
	s.restaurantScheduler.Start()
	
	// Start other schedulers as needed
}

// StopAll stops all schedulers
func (s *Scheduler) StopAll() {
	// Stop restaurant scheduler
	s.restaurantScheduler.Stop()
	
	// Stop other schedulers as needed
} 