package restaurant

import (
	"fmt"

	"github.com/plugimt/transat-backend/handlers/restaurant/repository"
	"github.com/plugimt/transat-backend/handlers/restaurant/service"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type Scheduler struct {
	MenuRepository          *repository.MenuRepository
	MenuService             *service.MenuService
	NotifService            *services.NotificationService
	lastNotificationDate    string
	menuSimilarityThreshold float64
}

func NewScheduler(menuRepo *repository.MenuRepository, menuService *service.MenuService, notifService *services.NotificationService) *Scheduler {
	return &Scheduler{
		MenuRepository:          menuRepo,
		MenuService:             menuService,
		NotifService:            notifService,
		menuSimilarityThreshold: 0.7,
	}
}

func (s *Scheduler) CheckAndUpdateMenuCron() (bool, error) {
	utils.LogHeader("🤖 Cron: Check & Update Restaurant Menu")

	// 1. Fetch latest base menu from API
	baseMenuData, err := s.MenuService.FetchMenuFromAPI()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Cron: Failed to fetch base menu from API")
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error: %v", err))
		utils.LogFooter()
		return false, err
	}

	// 2. Convert to FetchedItems format
	fetchedItems := s.MenuService.ConvertMenuDataToFetchedItems(baseMenuData)
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Converted %d menu items for synchronization", len(fetchedItems)))

	// 3. Synchronize with database
	err = s.MenuRepository.SyncTodaysMenu(fetchedItems)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to sync menu: %v", err))
		utils.LogFooter()
		return false, err
	}

	// Note: Notifications are automatically handled within SyncTodaysMenu if conditions are met

	utils.LogMessage(utils.LevelInfo, "Menu synchronization completed successfully")
	utils.LogFooter()
	return true, nil
}
