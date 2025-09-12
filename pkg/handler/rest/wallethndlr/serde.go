package wallethndlr

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/handler/rest/serde"
)

type topUpRequest struct {
	Amount      int64   `json:"amount"`
	ReleaseDate string  `json:"release_date"`
	Description *string `json:"description"`
}

type topUpDTO struct {
	amount      int64
	releaseDate *time.Time
	description *string
	userID      model.UserID
}

func bindTopUpRequest(r *http.Request) (topUpDTO, error) {
	var (
		req topUpRequest
		dto topUpDTO
	)

	err := serde.ReadJson(r, &req)
	if err != nil {
		return dto, fmt.Errorf("read request: %w", err)
	}
	userID, err := parseInt64(r.PathValue(UserID))
	if err != nil {
		return dto, fmt.Errorf("invalid user id: %w", err)
	}
	dto.userID = model.UserID(userID)
	dto.description = req.Description

	if req.Amount <= 0 {
		return dto, fmt.Errorf("invalid amount: %d", req.Amount)
	}
	dto.amount = req.Amount

	releaseDateStr := req.ReleaseDate
	if releaseDateStr == "" {
		return dto, nil
	}
	releaseDate, err := time.Parse(time.RFC3339, releaseDateStr)
	switch {
	case err != nil:
		return dto, fmt.Errorf("parse release date: %w", err)
	case releaseDate.Before(time.Now()):
		return dto, fmt.Errorf("release date is in the past: %s", releaseDate)
	default:
		dto.releaseDate = &releaseDate
	}

	return dto, nil
}

type withdrawRequest struct {
	Amount      int64   `json:"amount"`
	Description *string `json:"description"`

	userID model.UserID
}

type withdrawResponse struct {
	RefID *uuid.UUID `json:"ref_id"`
}

func bindWithdrawRequest(r *http.Request) (withdrawRequest, error) {
	var req withdrawRequest
	err := serde.ReadJson(r, &req)
	if err != nil {
		return req, fmt.Errorf("read request: %w", err)
	}
	userID, err := parseInt64(r.PathValue(UserID))
	if err != nil {
		return req, fmt.Errorf("invalid user id: %w", err)
	}
	req.userID = model.UserID(userID)
	if req.Amount <= 0 {
		return req, fmt.Errorf("invalid amount: %d", req.Amount)
	}
	return req, nil
}

type historyRequest struct {
	count     int64
	threshold int64
	userID    model.UserID
}

func bindHistoryRequest(r *http.Request) (*historyRequest, error) {
	userID, err := parseInt64(r.PathValue(UserID))
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	query := r.URL.Query()

	count, err := serde.ValueOrDefault(query.Get("count"), parseInt64)
	if err != nil {
		return nil, fmt.Errorf("invalid count query: %w", err)
	}
	if count < 0 || count > 100 {
		return nil, fmt.Errorf("invalid count value: %d", count)
	}
	if count == 0 {
		count = 20
	}

	threshold, err := serde.ValueOrDefault(query.Get("threshold"), parseInt64)
	if err != nil {
		return nil, fmt.Errorf("invalid threshold query: %w", err)
	}
	if threshold < 0 {
		return nil, fmt.Errorf("invalid threshold value: %d", threshold)
	}

	return &historyRequest{
		count:     count,
		threshold: threshold,
		userID:    model.UserID(userID),
	}, nil
}

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
