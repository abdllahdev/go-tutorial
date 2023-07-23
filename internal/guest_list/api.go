package guest_list

import (
	"encoding/json"
	"net/http"

	"github.com/getground/tech-tasks/backend/internal/entity"
	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router, service GuestListService) {
	h := handler{service}
	r.HandleFunc("/tables", h.createTable).Methods(http.MethodPost)
	r.HandleFunc("/guest_list", h.getAllGuests).Methods(http.MethodGet)
	r.HandleFunc("/guest_list/{name}", h.addGuest).Methods(http.MethodPost)
	r.HandleFunc("/guests/{name}", h.checkInGuest).Methods(http.MethodPut)
	r.HandleFunc("/guests", h.getAllCheckedInGuests).Methods(http.MethodGet)
	r.HandleFunc("/seats_empty", h.countEmptySeat).Methods(http.MethodGet)
	r.HandleFunc("/guests/{name}", h.checkoutGuest).Methods(http.MethodDelete)
}

type handler struct {
	service GuestListService
}

func (h handler) createTable(w http.ResponseWriter, r *http.Request) {
	var table entity.Table
	err := json.NewDecoder(r.Body).Decode(&table)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTable, err := h.service.CreateTable(&table)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTable)
}

func (h handler) addGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var requestBody entity.AddGuestRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var guest entity.Guest
	guest.Name = vars["name"]
	guest.TableID = requestBody.Table
	guest.AccompanyingGuests = requestBody.AccompanyingGuests

	newGuest, err := h.service.AddGuest(&guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGuest)
}

func (h handler) getAllGuests(w http.ResponseWriter, r *http.Request) {
	guests, err := h.service.GetAllGuests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := entity.GetAllGuestsResponseBody{
		Guests: guests,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func (h handler) checkInGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var requestBody entity.CheckInGuestRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var guest entity.Guest
	guest.Name = vars["name"]
	guest.AccompanyingGuests = requestBody.AccompanyingGuests

	_, err = h.service.CheckInGuest(&guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := entity.CheckInGuestResponseBody{
		Name: guest.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func (h handler) getAllCheckedInGuests(w http.ResponseWriter, r *http.Request) {
	checkedInGuests, err := h.service.GetAllCheckedInGuests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := entity.GetAllCheckedInGuestsResponseBody{
		Guests: checkedInGuests,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func (h handler) countEmptySeat(w http.ResponseWriter, r *http.Request) {
	emptySeats, err := h.service.CountEmptySeats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody := entity.CountEmptySeatsResponseBody{
		SeatsEmpty: emptySeats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func (h handler) checkoutGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var guest entity.Guest
	guest.Name = vars["name"]
	err := h.service.CheckoutGuest(&guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w)
}
