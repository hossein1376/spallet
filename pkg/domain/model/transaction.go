package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          TxID       `json:"id"`
	UserID      UserID     `json:"user_id"`
	Amount      int64      `json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Type        TxType     `json:"type"`
	Status      *TxStatus  `json:"status,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Description *string    `json:"description,omitempty"`
	RefID       *uuid.UUID `json:"ref_id,omitempty"`
}

type TxID int64

func (id TxID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type TxType int

const (
	invalidTxType TxType = iota
	TxTypeDeposit
	TxTypeWithdrawal
)

func (t TxType) String() string {
	switch t {
	case TxTypeDeposit:
		return "deposit"
	case TxTypeWithdrawal:
		return "withdrawal"
	default:
		panic(fmt.Errorf("unknown transaction type: %d", t))
	}
}

func (t TxType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TxType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("invalid transaction type: %T", src)
	}
	txType, err := ParseTxType(s)
	if err != nil {
		return fmt.Errorf("parse transaction type: %w", err)
	}
	*t = txType
	return nil
}

func ParseTxType(s string) (TxType, error) {
	switch s {
	case "deposit":
		return TxTypeDeposit, nil
	case "withdrawal":
		return TxTypeWithdrawal, nil
	default:
		return invalidTxType, fmt.Errorf("unknown value: %s", s)
	}
}

type TxStatus int

const (
	invalidTxStatus TxStatus = iota
	TxStatusPending
	TxStatusCompleted
	TxStatusFailed
)

func (s TxStatus) String() string {
	switch s {
	case TxStatusPending:
		return "pending"
	case TxStatusCompleted:
		return "completed"
	case TxStatusFailed:
		return "failed"
	default:
		panic(fmt.Errorf("unknown transaction status: %d", s))
	}
}

func (s TxStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *TxStatus) Scan(src any) error {
	t, ok := src.(string)
	if !ok {
		return fmt.Errorf("invalid transaction status: %T", src)
	}
	status, err := ParseTxStatus(t)
	if err != nil {
		return fmt.Errorf("parse transaction status: %w", err)
	}
	*s = status
	return nil
}

func ParseTxStatus(s string) (TxStatus, error) {
	switch s {
	case "pending":
		return TxStatusPending, nil
	case "completed":
		return TxStatusCompleted, nil
	case "failed":
		return TxStatusFailed, nil
	default:
		return invalidTxStatus, fmt.Errorf("unknown value: %s", s)
	}
}

type InsertTxOption struct {
	Status      *TxStatus
	ReleaseDate *time.Time
	Description *string
	RefID       *uuid.UUID
}

func Ptr[T any](v T) *T {
	return &v
}
