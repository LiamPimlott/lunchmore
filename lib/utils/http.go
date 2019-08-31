package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LiamPimlott/lunchmore/lib/errs"
)

// ErrorResponse represents a json error response body.
type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// CustomClaims represents custom jwt claims.
type CustomClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

// NotImplementedHandler writes a not implemented response.
var NotImplementedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

// Decode decodes a json request body.
func Decode(w http.ResponseWriter, r *http.Request, body interface{}) {
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println("error decoding json.")
		WriteError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg)
		return
	}
}

// Respond encodes and writes a json response body.
func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error encoding json: %s", err.Error())
		WriteError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg)
	}
}

// RespondError encodes and writes a error json response body.
func RespondError(w http.ResponseWriter, code int, status, msg string) {
	w.WriteHeader(code)

	Respond(w, ErrorResponse{
		Code:    code,
		Status:  status,
		Message: msg,
	})
}

// WriteError writes an error response to a ResponseWriter.
func WriteError(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	w.Write([]byte(status))
}

// GenerateToken generates a signed jwt token string for a user.
func GenerateToken(id uint, secret string) (string, error) {
	claims := CustomClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ts, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("error generating token: %s", err.Error())
		return "", err
	}

	return ts, nil
}
