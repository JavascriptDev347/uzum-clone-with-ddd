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

// Register godoc
//
//	@Summary		Yangi foydalanuvchi ro'yxatdan o'tkazish
//	@Description	Email va parol orqali yangi akkaunt yaratadi
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Ro'yxatdan o'tish ma'lumotlari"
//	@Success		201		{object}	response.Envelope{data=RegisterResponse}	"Foydalanuvchi yaratildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov yoki email formati"
//	@Failure		409		{object}	response.Envelope	"Email allaqachon band"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/auth/register [post]
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

// Login godoc
//
//	@Summary		Tizimga kirish
//	@Description	Email va parolni tekshirib, access va refresh tokenlarni qaytaradi
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Kirish ma'lumotlari"
//	@Success		200		{object}	response.Envelope{data=TokenResponse}	"Tokenlar muvaffaqiyatli berildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov tanasi"
//	@Failure		401		{object}	response.Envelope	"Email yoki parol noto'g'ri"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/auth/login [post]
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

// Refresh godoc
//
//	@Summary		Tokenni yangilash
//	@Description	Amaldagi refresh token orqali yangi access va refresh token juftligini oladi
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	response.Envelope{data=TokenResponse}	"Yangi tokenlar"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov tanasi"
//	@Failure		401		{object}	response.Envelope	"Refresh token yaroqsiz yoki muddati o'tgan"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/auth/refresh [post]
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
