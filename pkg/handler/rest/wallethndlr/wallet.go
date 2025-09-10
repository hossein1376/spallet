package wallethndlr

import (
	"net/http"

	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/handler/rest/serde"
	"github.com/hossein1376/spallet/pkg/service"
)

const UserID = "user_id"

type WalletHandler struct {
	services *service.Services
}

func New(services *service.Services) *WalletHandler {
	return &WalletHandler{services: services}
}

func (h WalletHandler) TopUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := bindTopUpRequest(r)
	if err != nil {
		resp := serde.Response{Message: err.Error()}
		serde.WriteJson(ctx, w, http.StatusBadRequest, resp)
		return
	}

	err = h.services.Wallets.TopUpService(
		ctx, req.userID, req.amount, req.releaseDate, req.description,
	)
	if err != nil {
		status, resp := serde.ExtractFromErr(ctx, err)
		serde.WriteJson(ctx, w, status, resp)
		return
	}

	serde.WriteJson(ctx, w, http.StatusNoContent, nil)
}

func (h WalletHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := bindWithdrawRequest(r)
	if err != nil {
		resp := serde.Response{Message: err.Error()}
		serde.WriteJson(ctx, w, http.StatusBadRequest, resp)
		return
	}

	refId, err := h.services.Wallets.WithdrawalService(ctx, req.userID, req.Amount)
	if err != nil {
		status, resp := serde.ExtractFromErr(ctx, err)
		serde.WriteJson(ctx, w, status, resp)
		return
	}

	resp := serde.Response{Data: withdrawResponse{RefID: refId}}
	serde.WriteJson(ctx, w, http.StatusOK, resp)
}

func (h WalletHandler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64(r.PathValue(UserID))
	if err != nil {
		resp := serde.Response{Message: err.Error()}
		serde.WriteJson(ctx, w, http.StatusBadRequest, resp)
		return
	}
	balance, err := h.services.Wallets.BalanceService(ctx, model.UserID(userID))
	if err != nil {
		status, resp := serde.ExtractFromErr(ctx, err)
		serde.WriteJson(ctx, w, status, resp)
		return
	}

	serde.WriteJson(ctx, w, http.StatusOK, serde.Response{Data: balance})
}

func (h WalletHandler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := bindHistoryRequest(r)
	if err != nil {
		resp := serde.Response{Message: err.Error()}
		serde.WriteJson(ctx, w, http.StatusBadRequest, resp)
		return
	}
	history, err := h.services.Wallets.HistoryService(
		ctx, req.userID, req.count, req.threshold,
	)
	if err != nil {
		status, resp := serde.ExtractFromErr(ctx, err)
		serde.WriteJson(ctx, w, status, resp)
		return
	}

	serde.WriteJson(ctx, w, http.StatusOK, serde.Response{Data: history})
}
