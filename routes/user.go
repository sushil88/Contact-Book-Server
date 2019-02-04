package routes

import (
	"database/sql"
	"github.com/go-chi/chi"
	"contactBook/db"
	"contactBook/helpers"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
		return
	}
	id := 0
	err := db.Get().Get(&id, "select id from users where email = $1", email)
	if err != nil && err != sql.ErrNoRows {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Server Error"})
		return
	}
	if err == sql.ErrNoRows {
		helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"exists": false})
		return
	}
	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"exists": true})
}
