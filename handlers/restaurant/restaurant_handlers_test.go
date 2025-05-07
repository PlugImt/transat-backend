package restaurant

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/plugimt/transat-backend/models"
)

// crée un serveur HTTP de mock
func setupTestServer(t *testing.T, filePath string) *httptest.Server {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read test data file '%s': %v", filePath, err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}))
	return server
}

// Helpers de test

// ça permet d'override le menuSourceURL pour le test
func (h *RestaurantHandler) SetMenuSourceURLForTesting(url string) {
	h.menuSourceURL = url
}

// accès direct à fetchMenuFromAPI
func (h *RestaurantHandler) FetchMenuFromAPIForTesting() (*models.MenuData, error) {
	return h.fetchMenuFromAPI()
}

func TestFetchMenuFromAPI(t *testing.T) {
	// Setup mock server
	testDataPath := filepath.Join("../../tests/handlers", "data.html") // Relative path within tests/handlers
	server := setupTestServer(t, testDataPath)
	defer server.Close()

	// crée une instance du handler (les dépendances peuvent être nil pour ce test spécifique)
	// On passe nil pour DB, TransService, NotifService car elles ne sont pas utilisées par fetchMenuFromAPI directement
	// faudra peut-être rajouter des mocks mais pour l'instant ça suffit
	handler := NewRestaurantHandler(nil, nil, nil)
	handler.SetMenuSourceURLForTesting(server.URL) // On doit ajouter cette méthode helper

	// Call the function to test
	menuData, err := handler.FetchMenuFromAPIForTesting() // Need to expose fetchMenuFromAPI or use a test helper

	// --- Assertions ---
	if err != nil {
		t.Fatalf("fetchMenuFromAPI failed: %v", err)
	}

	if menuData == nil {
		t.Fatal("fetchMenuFromAPI returned nil menuData, expected valid data")
	}

	// c'est dans le data.html
	expectedGrilladesMidi := []string{
		"Suprème de poulet",
		"Saucisse de Toulouse sauce moutarde et thym",
		"Effiloché de poulet basquaise",
		"Couscous merguez",
	}
	expectedAccompMidi := []string{
		"Semoule bio",
		"Pommes de terre boulangère",
		"Brocolis",
		"Poélée de carottes",
	}
	expectedCibo := []string{
		"Oeufs brouillés légumes provençale", // Note the trailing space
	}
	// Use reflect.DeepEqual for slice comparison (ignores order if we sort, but let's stick to order for simplicity now)
	// Note: processRawMenuItems appends, so order should match the file if no duplicates were removed (which there aren't in this case)
	// Let's use compareStringSlicesIgnoreOrder helper for robustness (needs exposure or duplication)

	if !compareStringSlicesIgnoreOrder(menuData.GrilladesMidi, expectedGrilladesMidi) {
		t.Errorf(`GrilladesMidi mismatch:
Got:      %v
Expected: %v`, menuData.GrilladesMidi, expectedGrilladesMidi)
	}
	if !compareStringSlicesIgnoreOrder(menuData.AccompMidi, expectedAccompMidi) {
		t.Errorf(`AccompMidi mismatch:
Got:      %v
Expected: %v`, menuData.AccompMidi, expectedAccompMidi)
	}
	if !compareStringSlicesIgnoreOrder(menuData.Cibo, expectedCibo) {
		t.Errorf(`Cibo mismatch:
Got:      %v
Expected: %v`, menuData.Cibo, expectedCibo)
	}

	// Check other categories are empty
	if len(menuData.Migrateurs) != 0 {
		t.Errorf("Expected empty Migrateurs, got: %v", menuData.Migrateurs)
	}
	if len(menuData.GrilladesSoir) != 0 {
		t.Errorf("Expected empty GrilladesSoir, got: %v", menuData.GrilladesSoir)
	}
	if len(menuData.AccompSoir) != 0 {
		t.Errorf("Expected empty AccompSoir, got: %v", menuData.AccompSoir)
	}

}

// --- Mocks (if needed for other tests) ---
// Example: Mock TranslationService
type MockTranslationService struct{}

func (m *MockTranslationService) TranslateMenu(menu *models.MenuData, targetLang string) (*models.MenuData, error) {
	// Simple mock: return original menu or an error
	if targetLang == "error" {
		return nil, fmt.Errorf("mock translation error")
	}
	// Return a copy to avoid modifying the original in tests
	translatedMenu := *menu
	// Modify slightly to show it's "translated"
	if len(translatedMenu.GrilladesMidi) > 0 {
		translatedMenu.GrilladesMidi[0] = translatedMenu.GrilladesMidi[0] + " (" + targetLang + ")"
	}
	return &translatedMenu, nil
}

// Example: Mock NotificationService
type MockNotificationService struct {
	SendDailyMenuCalled bool
	LastMenuSent        *models.MenuData
	ShouldError         bool
}

func (m *MockNotificationService) SendDailyMenuNotification(menu *models.MenuData) error {
	m.SendDailyMenuCalled = true
	m.LastMenuSent = menu
	if m.ShouldError {
		return fmt.Errorf("mock notification error")
	}
	return nil
}

// TODO: ajouter des tests pour GetRestaurantMenu (en gros tout, donc faut mocker db/context/translation/notifs)
// TODO: ajouter des tests pour CheckAndUpdateMenuCron
// TODO: tester les traductions
