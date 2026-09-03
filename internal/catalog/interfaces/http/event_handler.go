package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EventHandler struct {
	createUc  CreateEventUseCase
	getUc     GetEventsUseCase
	getByIdUc GetEventUseCase
	updateUc  UpdateEventUseCase
	deleteUc  DeleteEventUseCase
	getAllUc  GetAllEventsIncludingDeletedUseCase
}

func NewEventHandler(createUc CreateEventUseCase,
	getUc GetEventsUseCase,
	getByIdUc GetEventUseCase,
	updateUc UpdateEventUseCase,
	deleteUc DeleteEventUseCase,
	getAllUc GetAllEventsIncludingDeletedUseCase) *EventHandler {
	return &EventHandler{
		createUc:  createUc,
		getUc:     getUc,
		getByIdUc: getByIdUc,
		updateUc:  updateUc,
		deleteUc:  deleteUc,
		getAllUc:  getAllUc,
	}
}

// CreateEvent godoc
//
//		@Summary		Yangi event (banner) yaratish
//		@Description	Eyebrow, title, subtitle, cta, category_id va rasm orqali yangi event yaratadi
//		@Tags			events
//		@Accept			multipart/form-data
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			eyebrow_uz	formData	string	false	"Kichik ustki matn (o'zbekcha)"
//		@Param			eyebrow_eng	formData	string	false	"Kichik ustki matn (inglizcha)"
//		@Param			eyebrow_ru	formData	string	false	"Kichik ustki matn (ruscha)"
//		@Param			title_uz	formData	string	true	"Sarlavha (o'zbekcha)"
//		@Param			title_eng	formData	string	true	"Sarlavha (inglizcha)"
//		@Param			title_ru	formData	string	true	"Sarlavha (ruscha)"
//		@Param			subtitle_uz	formData	string	false	"Kichik matn (o'zbekcha)"
//		@Param			subtitle_eng	formData	string	false	"Kichik matn (inglizcha)"
//		@Param			subtitle_ru	formData	string	false	"Kichik matn (ruscha)"
//		@Param			cta_uz		formData	string	false	"Tugma matni (o'zbekcha)"
//		@Param			cta_eng		formData	string	false	"Tugma matni (inglizcha)"
//		@Param			cta_ru		formData	string	false	"Tugma matni (ruscha)"
//		@Param			category_id	formData	string	true	"Kategoriya ID"
//		@Param			is_root		formData	boolean	false	"true bo'lsa, ro'yxatda birinchi chiqadi"
//		@Param			image		formData	file	true	"Event rasmi"
//		@Success		201			{object}	response.Envelope{data=application.EventOutput}	"Event yaratildi"
//		@Failure		400			{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//		@Failure		401			{object}	response.Envelope	"Autentifikatsiyadan o'tilmagan"
//		@Failure		403			{object}	response.Envelope	"Huquq yetarli emas"
//		@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/events [post]
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "fayl hajmi juda katta yoki noto'g'ri format")
		return
	}

	isRoot, _ := strconv.ParseBool(r.FormValue("is_root"))

	file, header, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "rasm yuklanmadi")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "rasmni o'qishda xatolik")
		return
	}

	input := application.CreateEventInput{
		EyebrowUz:   r.FormValue("eyebrow_uz"),
		EyebrowEng:  r.FormValue("eyebrow_eng"),
		EyebrowRu:   r.FormValue("eyebrow_ru"),
		TitleUz:     r.FormValue("title_uz"),
		TitleEng:    r.FormValue("title_eng"),
		TitleRu:     r.FormValue("title_ru"),
		SubtitleUz:  r.FormValue("subtitle_uz"),
		SubtitleEng: r.FormValue("subtitle_eng"),
		SubtitleRu:  r.FormValue("subtitle_ru"),
		CTAUz:       r.FormValue("cta_uz"),
		CTAEng:      r.FormValue("cta_eng"),
		CTARu:       r.FormValue("cta_ru"),
		CategoryID:  r.FormValue("category_id"),
		IsRoot:      isRoot,
		Image: media.UploadInput{
			FileName:    "event-images/" + uuid.NewString(),
			ContentType: header.Header.Get("Content-Type"),
			Data:        data,
		},
	}

	output, err := h.createUc.Execute(r.Context(), input)
	if err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, output)
}

// GetEvents godoc
//
//	@Summary		Eventlarni olish
//	@Description	Faol (o'chirilmagan) eventlar ro'yxati. is_root=true bo'lganlar birinchi chiqadi. lang bo'yicha localized javob qaytadi.
//	@Tags			events
//	@Produce		json
//	@Param			lang	query		string	false	"Til: uz (default), eng yoki ru"
//	@Success		200	{object}	response.Envelope{data=[]application.EventOutput}	"Eventlar"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/events [get]
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	lang := parseLang(r)
	events, err := h.getUc.Execute(r.Context(), lang)
	if err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusOK, events)
}

// GetEvent godoc
//
//	@Summary		Eventni olish
//	@Description	Eventni ID bo'yicha olish
//	@Tags			events
//	@Produce		json
//	@Param			id		path		string	true	"Event ID"
//	@Param			lang	query		string	false	"Til: uz (default), eng yoki ru"
//	@Success		200	{object}	response.Envelope{data=application.EventOutput}	"Event"
//	@Failure		404	{object}	response.Envelope	"Event topilmadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/events/{id} [get]
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lang := parseLang(r)
	event, err := h.getByIdUc.Execute(r.Context(), id, lang)
	if err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusOK, event)
}

// UpdateEvent godoc
//
//		@Summary		Eventni yangilash
//		@Description	Eventni ID bo'yicha yangilash (JSON body, rasm bu yerda o'zgartirilmaydi)
//		@Tags			events
//		@Accept			json
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id		path		string					true	"Event ID"
//		@Param			event	body		application.UpdateEventInput	true	"Yangilash uchun event ma'lumotlari"
//		@Success		200		{object}	response.Envelope	"Event yangilash muvaffaqiyatli"
//		@Failure		400		{object}	response.Envelope	"Validatsiya xatosi"
//		@Failure		404		{object}	response.Envelope	"Event topilmadi"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/events/{id} [put]
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input application.UpdateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.ID = id

	if err := h.updateUc.Execute(r.Context(), input); err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil)
}

// DeleteEvent godoc
//
//		@Summary		Eventni o'chirish
//		@Description	Eventni ID bo'yicha o'chirish (soft delete)
//		@Tags			events
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id	path		string	true	"Event ID"
//		@Success		200	{object}	response.Envelope	"Event o'chirish muvaffaqiyatli"
//		@Failure		404	{object}	response.Envelope	"Event topilmadi"
//		@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/events/{id} [delete]
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.deleteUc.Execute(r.Context(), id); err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Event o'chirildi")
}

// GetAllEvents godoc
//
//		@Summary		Eventlarni olish - admin (o'chirilganlar bilan birga)
//		@Description	Barcha eventlarni (o'chirilganlarni ham) qaytaradi
//		@Tags			events
//	 @Security		BearerAuth
//		@Produce		json
//		@Success		200	{object}	response.Envelope{data=[]application.EventOutputForAdmin}	"Eventlar"
//		@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/events/admin [get]
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.getAllUc.Execute(r.Context())
	if err != nil {
		writeEventError(w, err)
		return
	}

	response.Success(w, http.StatusOK, events)
}
