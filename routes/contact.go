package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"contactBook/db"
	"contactBook/helpers"

	"github.com/go-chi/chi"
)

type contactPayload struct {
	ID           int    `json:"id" db:"id"`
	EmailAddress string `json:"email_address" db:"email_address"`
	FirstName    string `json:"first_name" db:"first_name"`
	MiddleName   string `json:"middle_name" db:"middle_name"`
	LastName     string `json:"last_name" db:"last_name"`
	PhoneNumber  string `json:"phone_number" db:"phone_number"`
}

//Contact Returns details for given contact id
func Contact(w http.ResponseWriter, r *http.Request, user_id string) {
	email := chi.URLParam(r, "email")
	if email == "" {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "Invalid user"})
		return
	}

	contact := contactPayload{}

	err := db.Get().Get(&contact, `select id, email_address, first_name, middle_name, last_name, phone_number  from contacts where email_address = $1 and user_id = $2`, email, user_id)
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "User not found"})
		return
	}
	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"contact": contact})
}

// ContactIndex - Returns list of contacts
func ContactIndex(w http.ResponseWriter, r *http.Request, user_id string) {
	contacts := []contactPayload{}

	p := r.FormValue("page")
	q := r.FormValue("query")
	var err error
	if n, e := strconv.Atoi(p); e == nil && n > 0 {
		err = db.Get().Select(&contacts, `select id, email_address, first_name, middle_name, last_name, phone_number  from contacts where user_id = $1 and (last_name ilike $2 or first_name ilike $2 or middle_name ilike $2 or email_address ilike $2) limit $3 offset $4`, user_id, "%"+q+"%", 10, (n-1)*10)
	} else { // if page number is not provided, first 10
		err = db.Get().Select(&contacts, `select id, email_address, first_name, middle_name, last_name, phone_number  from contacts where user_id = $1 and (last_name ilike $2 or first_name ilike $2 or middle_name ilike $2 or email_address ilike $2) limit $3`, user_id, "%"+q+"%", 10)
	}
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Server error"})
		return
	}
	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"contacts": contacts})
}

// ContactCreate - Creates new contact
func ContactCreate(w http.ResponseWriter, r *http.Request, user_id string) {
	formErrors := struct {
		Form      string `json:"form"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
	}{}

	contact := struct {
		EmailAddress string `json:"email_address" db:"email_address"`
		FirstName    string `json:"first_name" db:"first_name"`
		MiddleName   string `json:"middle_name" db:"middle_name"`
		LastName     string `json:"last_name" db:"last_name"`
		PhoneNumber  string `json:"phone_number" db:"phone_number"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		formErrors.Form = "Invalid Details"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	isContactFormValid := true

	if contact.FirstName == "" {
		isContactFormValid = false
		formErrors.FirstName = "This is field is required"
	}
	if contact.EmailAddress == "" {
		isContactFormValid = false
		formErrors.Email = "This is field is required"
	}

	if !isContactFormValid {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	id := 0
	err = db.Get().Get(&id, "select id from contacts where email_address = $1 and user_id = $2", contact.EmailAddress, user_id)
	if err != nil && err != sql.ErrNoRows {
		formErrors.Form = "Server error"
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"errors": formErrors})
		helpers.Catch(err)
		return
	}
	if id != 0 {
		formErrors.Form = "Contact already exists"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	queryInsert := `insert into contacts(
						first_name,
						middle_name,
						last_name,
						email_address,
						phone_number,
						created_at,
						updated_at,
                  		user_id)
					values (
						$1,
						$2,
						$3,
						$4,
						$5,
						now() at time zone 'utc',
						now() at time zone 'utc',
					    $6)`

	_, err = db.Get().Exec(queryInsert, contact.FirstName, contact.MiddleName, contact.LastName, contact.EmailAddress, contact.PhoneNumber, user_id)
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Server Error"})
		helpers.Catch(err)
		return
	}

	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

// ContactUpdate - Updates contact
func ContactUpdate(w http.ResponseWriter, r *http.Request, user_id string) {
	formErrors := struct {
		Form      string `json:"form"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
	}{}

	contact := contactPayload{}

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		formErrors.Form = "Invalid Details"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	isContactFormValid := true

	if contact.FirstName == "" {
		isContactFormValid = false
		formErrors.FirstName = "This is field is required"
	}
	if contact.EmailAddress == "" {
		isContactFormValid = false
		formErrors.Email = "This is field is required"
	}

	if !isContactFormValid {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	queryUpdate := `update contacts set
						first_name = $1,
						middle_name = $2,
						last_name = $3,
						email_address = $4,
						phone_number = $5,
						updated_at = now() at time zone 'utc' where id = $6`

	_, err = db.Get().Exec(queryUpdate, contact.FirstName, contact.MiddleName, contact.LastName, contact.EmailAddress, contact.PhoneNumber, contact.ID)
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Server Error"})
		return
	}

	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

// ContactDelete - Deletes contact
func ContactDelete(w http.ResponseWriter, r *http.Request, user_id string) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "Invalid contact"})
		return
	}
	_, err := db.Get().Exec(`delete from contacts where id = $1`, id)
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "User not found"})
		return
	}
	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}
