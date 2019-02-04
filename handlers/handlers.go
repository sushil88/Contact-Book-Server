package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"contactBook/conf"
	"contactBook/db"
	"contactBook/helpers"
	"contactBook/routes"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(conf.Get().Secret), nil)
}

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	formError := struct {
		Form string `json:"form"`
	}{}
	creds := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		formError.Form = "Invalid Credentials"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formError})
		return
	}

	user := struct {
		ID       int    `db:"id"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}{}

	err = db.Get().Get(&user, "select id, email, password from users where email = $1", creds.Email)
	if err != nil && err != sql.ErrNoRows {
		formError.Form = "Server error"
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"errors": formError})
		return
	}
	if user.Password != creds.Password {
		formError.Form = "Invalid Password"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formError})
		return
	}

	_, tokenString, err := tokenAuth.Encode(jwt.MapClaims{"user_id": user.ID, "email": user.Email})
	if err != nil {
		formError.Form = "Server error"
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"errors": formError})
		return
	}

	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"jwtToken": tokenString})
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	formErrors := struct {
		Form                 string `json:"form"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
		FirstName            string `json:"first_name"`
	}{}

	newUser := struct {
		Email                string `json:"email" db:"email"`
		Password             string `json:"password" db:"password"`
		PasswordConfirmation string `json:"password_confirmation" db:"-"`
		FirstName            string `json:"first_name" db:"first_name"`
		MiddleName           string `json:"middle_name" db:"middle_name"`
		LastName             string `json:"last_name" db:"last_name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		formErrors.Form = "Invalid Details"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	isSignupFormValid := true

	if newUser.FirstName == "" {
		isSignupFormValid = false
		formErrors.FirstName = "This is field is required"
	}
	if newUser.Email == "" {
		isSignupFormValid = false
		formErrors.Email = "This is field is required"
	}
	if newUser.Password == "" {
		isSignupFormValid = false
		formErrors.Password = "This is field is required"
	}
	if newUser.Password != newUser.PasswordConfirmation {
		isSignupFormValid = false
		formErrors.PasswordConfirmation = "Passwords must match"
	}

	if !isSignupFormValid {
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	email := ""
	err = db.Get().Get(&email, "select email from users where email = $1", newUser.Email)
	if err != nil && err != sql.ErrNoRows {
		formErrors.Form = "Server error"
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"errors": formErrors})
		helpers.Catch(err)
		return
	}
	if email != "" {
		formErrors.Form = "User already exists"
		helpers.RespondwithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": formErrors})
		return
	}

	queryInsert := `insert into users(
						first_name,
						middle_name,
						last_name,
						email,
						password,
						created_at,
						updated_at)
					values (
						$1,
						$2,
						$3,
						$4,
						$5,
						now() at time zone 'utc',
						now() at time zone 'utc')`

	_, err = db.Get().Exec(queryInsert, newUser.FirstName, newUser.MiddleName, newUser.LastName, newUser.Email, newUser.Password)
	if err != nil {
		helpers.RespondwithJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Server Error"})
		return
	}

	helpers.RespondwithJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Request processed in %s\n", time.Since(startTime))
		next(w, r)
	}
}

// AuthenticatedAPIHandlerFunc is our special func for authenticated JSON requests
type AuthenticatedAPIHandlerFunc func(http.ResponseWriter, *http.Request, string)

func (h *Handlers) AuthLogger(next AuthenticatedAPIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Request processed in %s\n", time.Since(startTime))
		_, claims, _ := jwtauth.FromContext(r.Context())
		s := fmt.Sprintf("%g", claims["user_id"].(float64))
		next(w, r, s)
	}
}

func (h *Handlers) SetupRoutes(r *chi.Mux) {
	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/contact/{email}", h.AuthLogger(routes.Contact))
		r.Get("/contacts", h.AuthLogger(routes.ContactIndex))
		r.Post("/contact", h.AuthLogger(routes.ContactCreate))
		r.Put("/contact", h.AuthLogger(routes.ContactUpdate))
		r.Delete("/contact/{id:[0-9]+}", h.AuthLogger(routes.ContactDelete))
	})

	// Public routes
	r.Get("/", h.Logger(h.Home))
	r.Post("/login", h.Logger(h.Login))
	r.Post("/register", h.Logger(h.RegisterUser))
	r.Get("/user/{email}", h.Logger(routes.User))

	// TODO Take care of version
	// TODO
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}
