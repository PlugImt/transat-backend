package scheduler

import (
	"sync"
	"time"

	"github.com/plugimt/transat-backend/utils"
)

type UserScheduleCronHandler interface {
	SyncAll() error
}

type UserScheduleScheduler struct {
	handler     UserScheduleCronHandler
	stopChan    chan struct{}
	runningChan chan struct{}
	stopOnce    sync.Once
	syncMutex   sync.Mutex
}

func NewUserScheduleScheduler(handler UserScheduleCronHandler) *UserScheduleScheduler {
	return &UserScheduleScheduler{
		handler:     handler,
		stopChan:    make(chan struct{}),
		runningChan: make(chan struct{}),
	}
}

func (s *UserScheduleScheduler) Start() {
	utils.LogMessage(utils.LevelInfo, "Starting user schedule ICS scheduler")

	go func() {
		close(s.runningChan)

		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go s.syncAllSchedules()
			case <-s.stopChan:
				utils.LogMessage(utils.LevelInfo, "Stopping user schedule ICS scheduler")
			}
		}
	}()
}

func (s *UserScheduleScheduler) Stop() {
	s.stopOnce.Do(func() {
		select {
		case <-s.runningChan:
			close(s.stopChan)
		default:
			utils.LogMessage(utils.LevelWarn, "Attempted to stop user schedule scheduler that wasn't running")
		}
	})
}

func (s *UserScheduleScheduler) syncAllSchedules() {
	utils.LogMessage(utils.LevelInfo, "Running scheduled ICS sync for all users")

	if err := s.handler.SyncAll(); err != nil {
		utils.LogMessage(utils.LevelError, "Error in scheduled ICS sync")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		utils.LogMessage(utils.LevelInfo, "Scheduled ICS sync completed successfully")
	}
}
