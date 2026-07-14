package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mouzong/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProdcut(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	d, err := json.Marshal(lp)

	if err != nil {
		http.Error(w, "Unable to marshal the json", http.StatusInternalServerError)
	}

	w.Write(d)
}
