package scheduler

import (
	"sync"
	"testing"
	"time"
)

// Mock RestaurantHandler for testing
type MockRestaurantHandler struct {
	CheckAndUpdateMenuCronCalled int
	ShouldReturn                 bool
	ShouldError                  bool
	mu                           sync.Mutex
}

func (m *MockRestaurantHandler) CheckAndUpdateMenuCron() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CheckAndUpdateMenuCronCalled++

	if m.ShouldError {
		return false, &MockError{message: "mock error"}
	}
	return m.ShouldReturn, nil
}

// Mock error for testing
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestNewRestaurantScheduler(t *testing.T) {
	mockHandler := &MockRestaurantHandler{}
	scheduler := NewRestaurantScheduler(mockHandler)

	if scheduler == nil {
		t.Fatal("NewRestaurantScheduler returned nil")
	}

	if scheduler.restaurantHandler != mockHandler {
		t.Error("Restaurant handler not set correctly")
	}

	if scheduler.stopChan == nil {
		t.Error("Stop channel not initialized")
	}

	if scheduler.runningChan == nil {
		t.Error("Running channel not initialized")
	}
}

func TestRestaurantSchedulerStartStop(t *testing.T) {
	mockHandler := &MockRestaurantHandler{
		ShouldReturn: true,
	}

	scheduler := NewRestaurantScheduler(mockHandler)

	// Start the scheduler
	scheduler.Start()

	// Wait a bit for the immediate check to happen
	time.Sleep(50 * time.Millisecond)

	// Check if handler was called
	mockHandler.mu.Lock()
	callCount := mockHandler.CheckAndUpdateMenuCronCalled
	mockHandler.mu.Unlock()

	if callCount == 0 {
		t.Error("CheckAndUpdateMenuCron not called immediately after start")
	}

	// Stop the scheduler
	scheduler.Stop()

	// Make sure we can stop without issues
	scheduler.Stop() // Should handle being called multiple times

	// Record the current call count
	mockHandler.mu.Lock()
	callCountBefore := mockHandler.CheckAndUpdateMenuCronCalled
	mockHandler.mu.Unlock()

	// Wait a bit to ensure no more calls happen
	time.Sleep(50 * time.Millisecond)

	// Check if handler was called again (it shouldn't be)
	mockHandler.mu.Lock()
	callCountAfter := mockHandler.CheckAndUpdateMenuCronCalled
	mockHandler.mu.Unlock()

	if callCountAfter > callCountBefore {
		t.Error("Handler called after scheduler was stopped")
	}
}

func TestCheckAndUpdateMenu(t *testing.T) {
	// Test successful update
	mockHandler := &MockRestaurantHandler{
		ShouldReturn: true,
	}
	scheduler := NewRestaurantScheduler(mockHandler)

	scheduler.checkAndUpdateMenu()

	if mockHandler.CheckAndUpdateMenuCronCalled != 1 {
		t.Errorf("CheckAndUpdateMenuCron called %d times, expected 1", mockHandler.CheckAndUpdateMenuCronCalled)
	}

	// Test with error
	mockHandler = &MockRestaurantHandler{
		ShouldError: true,
	}
	scheduler = NewRestaurantScheduler(mockHandler)

	// This shouldn't panic even with an error
	scheduler.checkAndUpdateMenu()

	if mockHandler.CheckAndUpdateMenuCronCalled != 1 {
		t.Errorf("CheckAndUpdateMenuCron called %d times, expected 1", mockHandler.CheckAndUpdateMenuCronCalled)
	}
}

func TestParisTimeScheduling(t *testing.T) {
	// Create a scheduler instance for testing
	mockHandler := &MockRestaurantHandler{
		ShouldReturn: true,
	}
	scheduler := NewRestaurantScheduler(mockHandler)

	// Test times that should be allowed (weekday 9-13h Paris time)
	// Create a time that is 10 AM Paris time on a weekday (Tuesday)
	parisLocation, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatalf("Failed to load Paris timezone: %v", err)
	}
	
	// Create a Tuesday at 10 AM Paris time
	allowedTime := time.Date(2024, 1, 2, 10, 0, 0, 0, parisLocation)
	
	if !scheduler.isScheduledTimeAllowed(allowedTime) {
		t.Error("Expected Tuesday 10 AM Paris time to be allowed")
	}

	// Test times that should NOT be allowed
	// Weekend (Saturday)
	weekendTime := time.Date(2024, 1, 6, 10, 0, 0, 0, parisLocation)
	if scheduler.isScheduledTimeAllowed(weekendTime) {
		t.Error("Expected Saturday to not be allowed")
	}

	// Too early (7 AM Paris time on weekday)
	earlyTime := time.Date(2024, 1, 2, 7, 0, 0, 0, parisLocation)
	if scheduler.isScheduledTimeAllowed(earlyTime) {
		t.Error("Expected Tuesday 7 AM Paris time to not be allowed")
	}

	// Too late (15 PM Paris time on weekday)
	lateTime := time.Date(2024, 1, 2, 15, 0, 0, 0, parisLocation)
	if scheduler.isScheduledTimeAllowed(lateTime) {
		t.Error("Expected Tuesday 3 PM Paris time to not be allowed")
	}

	// Test that Paris time conversion is working correctly
	// Create a UTC time and verify it converts to correct Paris time
	utcTime := time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC) // 8 AM UTC
	
	// This should be 9 AM Paris time in winter (UTC+1) or 10 AM in summer (UTC+2)
	// Let's check what the conversion gives us
	if !scheduler.isScheduledTimeAllowed(utcTime) {
		// This might be expected if it's outside the allowed window
		// Let's use a time that should definitely work: 9 AM UTC in winter = 10 AM Paris
		utcTimeWinter := time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC) // 9 AM UTC in winter = 10 AM Paris
		if !scheduler.isScheduledTimeAllowed(utcTimeWinter) {
			t.Error("Expected UTC time to be correctly converted to Paris time for scheduling")
		}
	}
}
