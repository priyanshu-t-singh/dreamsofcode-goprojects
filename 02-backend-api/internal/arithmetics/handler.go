package arithmetics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct{}

type RequestBody struct {
	A int `json:"a"`
	B int `json:"b"`
}

type DivisionRequestBody struct {
	Dividend int `json:"dividend"`
	Divisor  int `json:"divisor"`
}

type ResponseBody struct {
	Result int `json:"result"`
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result int = body.A + body.B

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func (h *Handler) Subtract(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result int = body.A - body.B

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func (h *Handler) Multiply(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result int = body.A * body.B

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func (h *Handler) Divide(w http.ResponseWriter, r *http.Request) {
	var body DivisionRequestBody
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.Divisor == 0 {
		fmt.Fprint(w, "'division by zero' not allowed")
		return
	}

	var result int = body.Dividend / body.Divisor

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func (h *Handler) Sum(w http.ResponseWriter, r *http.Request) {
	var body []int
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result int = 0
	for _, num := range body {
		result += num
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}
