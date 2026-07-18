package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/mouzong/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		// expect the ID in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI unable to capture number")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, w, r)
		return
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to un marshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

// getProducts returns the products from the data store
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal the json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to un arshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)

}
