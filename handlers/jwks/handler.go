package jwks

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zalando/planb-tokeninfo/keyloader"
)

type jwksHandler struct {
	loader keyloader.KeyLoader
}

// NewHandler creates an http.Handler that provides JWKS responses from the KeyLoader kl
func NewHandler(kl keyloader.KeyLoader) http.Handler {
	return &jwksHandler{loader: kl}
}

// ServeHTTP serializes the current snapshot of Keys from the KeyLoader as a JSON Web Key Set
func (h *jwksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper := &jwksWrapper{keys: h.loader.Keys()}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wrapper); err != nil {
		log.Println("Failed to finish JWKS response: ", err)
	}
}
