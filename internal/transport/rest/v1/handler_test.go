package v1

import (
	"bytes"
	"encoding/json"
	"itkdemo/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletUseCase struct {
	mock.Mock
}

func (m *MockWalletUseCase) CreateWallet() (*domain.Wallet, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) DeleteWallet(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWalletUseCase) GetWallet(id uuid.UUID) (*domain.Wallet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) UpdateWallet(op *domain.Operation) error {
	args := m.Called(op)
	return args.Error(0)
}

func TestHandler(t *testing.T) {
	e := echo.New()
	mockID := uuid.New()

	t.Run("CreateWalletSuccess", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}

		m.On("CreateWallet").Return(&domain.Wallet{ID: mockID, Balance: 0}, nil)

		req := httptest.NewRequest(http.MethodPost, "/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateWallet(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("CreateWalletError", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("CreateWallet").Return(nil, assert.AnError)

		c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
		err := h.CreateWallet(c)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, he.Code)
	})

	t.Run("GetBalanceSuccess", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("GetWallet", mockID).Return(&domain.Wallet{ID: mockID, Balance: 100}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(mockID.String())

		assert.NoError(t, h.GetBalance(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GetBalanceInvalidID", func(t *testing.T) {
		h := &Handler{uc: nil}
		req := httptest.NewRequest(http.MethodGet, "/uuid", nil)
		c := e.NewContext(req, httptest.NewRecorder())

		err := h.GetBalance(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("GetBalanceNotFound", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("GetWallet", mockID).Return(nil, domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		c := e.NewContext(req, httptest.NewRecorder())
		c.SetParamNames("id")
		c.SetParamValues(mockID.String())

		err := h.GetBalance(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, he.Code)
	})

	t.Run("OperationSuccess", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}

		op := domain.Operation{WalletID: mockID, Amount: 1000, Type: domain.Withdraw}
		m.On("UpdateWallet", mock.MatchedBy(func(o *domain.Operation) bool {
			return o.WalletID == mockID && o.Amount == 1000 && o.Type == domain.Withdraw
		})).Return(nil)

		body, _ := json.Marshal(op)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		assert.NoError(t, h.Operation(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("OperationBindError", func(t *testing.T) {
		h := &Handler{uc: nil}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := h.Operation(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("OperationBadRequest", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}

		op := domain.Operation{WalletID: mockID, Amount: 1000, Type: domain.Withdraw}
		m.On("UpdateWallet", mock.MatchedBy(func(o *domain.Operation) bool {
			return o.WalletID == mockID && o.Amount == 1000 && o.Type == domain.Withdraw
		})).Return(domain.ErrInsufficientBalance)

		body, _ := json.Marshal(op)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := h.Operation(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("DeleteBadRequest", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}

		op := domain.Operation{Amount: 1000, Type: domain.Withdraw}
		m.On("UpdateWallet", mock.MatchedBy(func(o *domain.Operation) bool {
			return o.WalletID == mockID && o.Amount == 1000
		})).Return(domain.ErrNotFound)

		body, _ := json.Marshal(op)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(mockID.String())

		err := h.DeleteWallet(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("DeleteWalletSuccess", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("DeleteWallet", mockID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/?id="+mockID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, h.DeleteWallet(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("DeleteWalletInvalidID", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("DeleteWallet", mockID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/?id=uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.DeleteWallet(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("DeleteWalletNotFound", func(t *testing.T) {
		m := new(MockWalletUseCase)
		h := &Handler{uc: m}
		m.On("DeleteWallet", mockID).Return(domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/?id="+mockID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.DeleteWallet(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, he.Code)
	})
}
