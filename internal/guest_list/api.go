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
