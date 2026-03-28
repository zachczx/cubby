package server

import (
	"encoding/json"
	"net/http"

	"strings"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/market"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func (s *Service) LogMarketPriceHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	familyID, err := user.GetUserFamilyID(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input market.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	input.ItemName = strings.TrimSpace(input.ItemName)
	if input.ItemName == "" {
		response.RespondWithError(w, http.StatusBadRequest, "item name is required")
		return
	}
	if input.Price < 0 {
		response.RespondWithError(w, http.StatusBadRequest, "price cannot be negative")
		return
	}
	if input.Quantity != nil && *input.Quantity < 0 {
		response.RespondWithError(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}

	loggedBy := userID
	p := market.MarketPrice{
		FamilyID: familyID,
		LoggedBy: &loggedBy,
		ItemName: input.ItemName,
		Category: input.Category,
		Country:  input.Country,
		Store:    input.Store,
		Unit:     input.Unit,
		Quantity: input.Quantity,
		Price:    input.Price,
		IsPromo:  input.IsPromo,
		Remarks:  input.Remarks,
	}

	result, err := market.LogPrice(s.DB, p)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, result)
}

func (s *Service) GetMarketPricesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	prices, err := market.GetPrices(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, prices)
}

func (s *Service) UpdateMarketPriceHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	priceID, err := uuid.Parse(r.PathValue("priceID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input market.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	input.ItemName = strings.TrimSpace(input.ItemName)
	if input.ItemName == "" {
		response.RespondWithError(w, http.StatusBadRequest, "item name is required")
		return
	}
	if input.Price < 0 {
		response.RespondWithError(w, http.StatusBadRequest, "price cannot be negative")
		return
	}
	if input.Quantity != nil && *input.Quantity < 0 {
		response.RespondWithError(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}

	p := market.MarketPrice{
		ID:       priceID,
		ItemName: input.ItemName,
		Category: input.Category,
		Country:  input.Country,
		Store:    input.Store,
		Unit:     input.Unit,
		Quantity: input.Quantity,
		Price:    input.Price,
		IsPromo:  input.IsPromo,
		Remarks:  input.Remarks,
	}

	if err := market.UpdatePrice(s.DB, p, userID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetMarketInsightsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	insights, err := market.GetInsights(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, insights)
}

func (s *Service) DeleteMarketPriceHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	priceID, err := uuid.Parse(r.PathValue("priceID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := market.DeletePrice(s.DB, userID, priceID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
