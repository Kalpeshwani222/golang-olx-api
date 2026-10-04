package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/kalpeshWani222/olx-api/internal/httpx"
	"github.com/kalpeshWani222/olx-api/internal/middleware"
	"github.com/kalpeshWani222/olx-api/internal/token"
)

type AuthHandler struct {
	db        *sql.DB
	logger    *slog.Logger
	jwtSecret []byte
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		logger:    logger,
		jwtSecret: []byte(jwtSecret),
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIdFromContext(ctx)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("register: decode failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	//validate the body
	if err := req.Validate(); err != nil {
		var valErr ValidationError
		errors.As(err, &valErr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, valErr.Field)
		return
	}

	//check if email exist
	var exists bool
	if err := ah.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, req.Email).Scan(&exists); err != nil {
		ah.logger.Error("register: email check failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if exists {
		// Email is already taken return it
		httpx.ValidationError(w, http.StatusConflict, "email already registered", httpx.CodeConflict, "email")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ah.logger.Error("register: bcrypt failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// Insert the new user into DB
	var user User
	row := ah.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, created_at
	`, req.Email, string(hashedPassword))

	if err := row.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
		ah.logger.Error("register: insert failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ah.logger.Info("user registered", "requestId", requestID, "user_id", user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIdFromContext(ctx)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("login: decode failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	//Validate
	if err := req.Validate(); err != nil {
		var valErr ValidationError
		errors.As(err, &valErr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, valErr.Field)
		return
	}

	//Find the user by email
	var user User
	row := ah.db.QueryRowContext(ctx, `
		SELECT id, email, password FROM users WHERE email = $1
	`, req.Email)

	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.Error(w, http.StatusUnauthorized, "invalid credentials", httpx.CodeUnauthenticated)
			return
		}
		ah.logger.Error("login: query failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	//Compare the provided password with the stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid credentials", httpx.CodeUnauthenticated)
		return
	}

	//Generate a JWT token
	tokenStr, err := token.Generate(user.ID, ah.jwtSecret)
	if err != nil {
		ah.logger.Error("login: token generation failed", "requestId", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ah.logger.Info("user logged in", "requestId", requestID, "user_id", user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}