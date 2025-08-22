package scheduler

import (
	"github.com/robfig/cron/v3"
)

// Scheduler manages all scheduled tasks in the application
type Scheduler struct {
	restaurantScheduler *RestaurantScheduler
	// Add other schedulers here as needed
	cron         *cron.Cron
	cronEntryIDs map[string]cron.EntryID
}

// NewScheduler creates a new main scheduler
func NewScheduler(restaurantHandler RestaurantMenuCronHandler) *Scheduler {
	return &Scheduler{
		restaurantScheduler: NewRestaurantScheduler(restaurantHandler),
		// Initialize other schedulers here
		cron:         cron.New(),
		cronEntryIDs: make(map[string]cron.EntryID),
	}
}

// StartAll starts all schedulers
func (s *Scheduler) StartAll() {
	// Start restaurant scheduler
	s.restaurantScheduler.Start()

	// Start other schedulers as needed

	// Start cron scheduler
	s.cron.Start()
}

// StopAll stops all schedulers
func (s *Scheduler) StopAll() {
	// Stop restaurant scheduler
	s.restaurantScheduler.Stop()

	// Stop other schedulers as needed

	// Stop cron scheduler
	s.cron.Stop()
}

// AddCronJob adds a new cron job with a name for identification
func (s *Scheduler) AddCronJob(name string, spec string, job func()) error {
	id, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return err
	}

	s.cronEntryIDs[name] = id
	return nil
}

// RemoveCronJob removes a cron job by name
func (s *Scheduler) RemoveCronJob(name string) {
	if id, exists := s.cronEntryIDs[name]; exists {
		s.cron.Remove(id)
		delete(s.cronEntryIDs, name)
	}
}

// GetCronEntries returns all registered cron entries
func (s *Scheduler) GetCronEntries() []cron.Entry {
	return s.cron.Entries()
}
