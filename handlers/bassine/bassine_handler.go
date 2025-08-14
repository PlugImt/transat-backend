package bassine

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/bassine/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type BassineHandler struct {
	BassineRepository *repository.BassineRepository
}

func NewBassineHandler(db *sql.DB) *BassineHandler {
	return &BassineHandler{
		BassineRepository: repository.NewBassineRepository(db),
	}
}
