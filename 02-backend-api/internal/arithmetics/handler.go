package arithmetics

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/models"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/middleware"
)

type Handler struct {
	Repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{Repository: repository}
}

type RequestBody struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type DivisionRequestBody struct {
	Dividend float64 `json:"dividend"`
	Divisor  float64 `json:"divisor"`
}

type ResponseBody struct {
	Result float64 `json:"result"`
}

func (h *Handler) saveCalculations(r *http.Request, calc *models.Calculation) error {
	username, ok := middleware.GetAuthenticatedUserID(r)
	if !ok {
		return fmt.Errorf("failed to get authenticated user ID")
	}

	calc.Username = username
	if err := h.Repository.SaveCalculation(calc); err != nil {
		return err
	}
	return nil
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result float64 = body.A + body.B

	calc := &models.Calculation{
		Input:     []float64{body.A, body.B},
		Operation: "add",
		Result:    result,
	}
	if err := h.saveCalculations(r, calc); err != nil {
		sendErrorResponse(w, err)
		return
	}

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

	var result float64 = body.A - body.B

	calc := &models.Calculation{
		Input:     []float64{body.A, body.B},
		Operation: "subtract",
		Result:    result,
	}
	if err := h.saveCalculations(r, calc); err != nil {
		sendErrorResponse(w, err)
		return
	}
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

	var result float64 = body.A * body.B

	calc := &models.Calculation{
		Input:     []float64{body.A, body.B},
		Operation: "multiply",
		Result:    result,
	}
	if err := h.saveCalculations(r, calc); err != nil {
		sendErrorResponse(w, err)
		return
	}

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

	var result float64 = body.Dividend / body.Divisor

	calc := &models.Calculation{
		Input:     []float64{body.Dividend, body.Divisor},
		Operation: "divide",
		Result:    result,
	}
	if err := h.saveCalculations(r, calc); err != nil {
		sendErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func (h *Handler) Sum(w http.ResponseWriter, r *http.Request) {
	var body []float64
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var result float64 = 0
	for _, num := range body {
		result += float64(num)
	}

	calc := &models.Calculation{
		Input:     body,
		Operation: "sum",
		Result:    result,
	}
	if err := h.saveCalculations(r, calc); err != nil {
		sendErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}
