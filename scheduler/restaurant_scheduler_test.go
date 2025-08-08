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
