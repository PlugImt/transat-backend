package scheduler

import "github.com/robfig/cron/v3"

// Scheduler manages all scheduled tasks in the application
type Scheduler struct {
	cron                  *cron.Cron
	cronEntryIDs          map[string]cron.EntryID
	restaurantScheduler   *RestaurantScheduler
	userScheduleScheduler *UserScheduleScheduler
}

// NewScheduler creates a new main scheduler
func NewScheduler(restaurantHandler RestaurantMenuCronHandler, userScheduleHandler UserScheduleCronHandler) *Scheduler {
	return &Scheduler{
		cron:                  cron.New(),
		cronEntryIDs:          make(map[string]cron.EntryID),
		restaurantScheduler:   NewRestaurantScheduler(restaurantHandler),
		userScheduleScheduler: NewUserScheduleScheduler(userScheduleHandler),
	}
}

// StartAll starts all schedulers
func (s *Scheduler) StartAll() {
	s.restaurantScheduler.Start()
	s.userScheduleScheduler.Start()
	s.cron.Start()
}

// StopAll stops all schedulers
func (s *Scheduler) StopAll() {
	s.restaurantScheduler.Stop()
	s.userScheduleScheduler.Stop()
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
