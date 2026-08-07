package http

import (
	"encoding/json"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

type IdentityHandler struct {
	registerUc RegisterUseCase
	loginUc    LoginUseCase
	refreshUc  RefreshUseCase
}

func NewIdentityHandler(registerUc RegisterUseCase, loginUc LoginUseCase, refreshUc RefreshUseCase) *IdentityHandler {
	return &IdentityHandler{
		registerUc: registerUc,
		loginUc:    loginUc,
		refreshUc:  refreshUc,
	}
}

func (h *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.registerUc.Execute(r.Context(), application.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		writeIdentityError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, RegisterResponse{
		UserID: output.UserID,
		Email:  output.Email,
	})
}

func (h *IdentityHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.loginUc.Execute(r.Context(), application.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		writeIdentityError(w, err)
		return
	}

	response.Success(w, http.StatusOK, TokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

func (h *IdentityHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.refreshUc.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	response.Success(w, http.StatusOK, TokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
