package controller

import (
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/adapter/inbound/dto"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/core/port"
)

type ProductController struct {
	productService port.ProductService
	logger         *slog.Logger
}

func NewProductController(productService port.ProductService, logger *slog.Logger) *ProductController {
	return &ProductController{
		productService: productService,
		logger:         logger,
	}
}

func (p *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.ProductDto

	err := json.NewDecoder(r.Body).Decode(&productDto)

	if err != nil {
		p.logger.Error("Unable to decode the request body.  %v", err)
	}

	errValidate := dto.ValidateProduct(productDto)

	fmt.Println(errValidate)

	if len(errValidate.Errors) > 0 {
		p.logger.Error("validate error.  %v", errValidate)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errValidate)
		return
	}

	productDomain := dto.ConvertDtoToDomain(productDto)

	product, err := p.productService.Create(productDomain)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (p *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.ProductDto

	err := json.NewDecoder(r.Body).Decode(&productDto)

	if err != nil {
		p.logger.Error("Unable to decode the request body.  %v", err)
	}

	errValidate := dto.ValidateProduct(productDto)

	fmt.Println(errValidate)

	if len(errValidate.Errors) > 0 {
		p.logger.Error("validate error.  %v", errValidate)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errValidate)
		return
	}

	productDomain := dto.ConvertDtoToDomain(productDto)

	product, err := p.productService.Update(productDomain)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (p *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIdParam, ok := vars["productId"]

	if !ok {
		fmt.Println("id is missing in parameters")
	}

	productId, err := strconv.Atoi(productIdParam)

	if err != nil {
		p.logger.Error("Error to convert productId to int.  %v", err)
	}

	err = p.productService.Delete(productId)
	if err != nil {
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		return
	}
}
