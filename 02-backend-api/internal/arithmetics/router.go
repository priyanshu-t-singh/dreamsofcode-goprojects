package arithmetics

import (
	"database/sql"
	"net/http"
)

func GetRoutes(db *sql.DB) *http.ServeMux {
	handler := &Handler{
		Repository: NewRepository(db),
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /add", handler.Add)
	router.HandleFunc("POST /subtract", handler.Subtract)
	router.HandleFunc("POST /multiply", handler.Multiply)
	router.HandleFunc("POST /divide", handler.Divide)
	router.HandleFunc("POST /sum", handler.Sum)

	return router
}
