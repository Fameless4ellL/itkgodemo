package v1

import (
	"itkdemo/internal/domain"
	"itkdemo/internal/usecase"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	uc *usecase.WalletUseCase
}

func NewHandler(uc *usecase.WalletUseCase) *Handler {
	return &Handler{uc: uc}
}

// CreateWallet creates a new wallet.
func (h *Handler) CreateWallet(c echo.Context) error {
	wallet, err := h.uc.CreateWallet()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, wallet)
	return nil
}

func (h *Handler) DeleteWallet(c echo.Context) error {
	id, err := uuid.Parse(c.QueryParam("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrInvalidID)
	}
	if err := h.uc.DeleteWallet(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound)
	}
	c.NoContent(http.StatusNoContent)
	return nil
}

func (h *Handler) GetBalance(c echo.Context) error {
	id, err := uuid.Parse(c.QueryParam("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrInvalidID)
	}
	wallet, err := h.uc.GetWallet(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound)
	}
	c.JSON(http.StatusOK, wallet)
	return nil
}

func (h *Handler) Operation(c echo.Context) error {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrInvalidID)
	}
	u := new(domain.Operation)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	u.WalletID = uid
	if err := h.uc.UpdateWallet(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, u)
	return nil
}
