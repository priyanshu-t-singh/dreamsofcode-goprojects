package arithmetics

import (
	"net/http"
)

func GetRoutes() *http.ServeMux {
	handler := &Handler{}

	router := http.NewServeMux()
	router.HandleFunc("POST /add", handler.Add)
	router.HandleFunc("POST /subtract", handler.Subtract)
	router.HandleFunc("POST /multiply", handler.Multiply)
	router.HandleFunc("POST /divide", handler.Divide)
	router.HandleFunc("POST /sum", handler.Sum)

	return router
}
