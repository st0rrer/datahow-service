package log

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	Service *Service
}


func (h *Handler) ProcessMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	message := &Message{}
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		http.Error(w, fmt.Errorf("could not deserialize message. %w", err).Error(), 400)
		return
	}

	err = h.Service.ProcessMessage(message)
	if err != nil {
		http.Error(w, fmt.Errorf("could not process message. %w", err).Error(), 400)
		return
	}
}
