package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/lib/utils"
)

// Required middleware requires a request a valid jwt token in the authorization header.
func Required(endpoint func(http.ResponseWriter, *http.Request), secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if len(authHeader) <= 0 {
			utils.RespondError(w, errs.ErrInvalidAuth.Code, errs.ErrInvalidAuth.Msg, "")
			return
		}

		code, err := parseToken(r, authHeader, secret)
		if err != nil {
			utils.RespondError(w, code, err.Error(), "")
			return
		}

		endpoint(w, r)
	})
}

// Optional middleware will only evaluate a jwt if the authorization header is set.
func Optional(endpoint func(http.ResponseWriter, *http.Request), secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) <= 0 {
			endpoint(w, r)
			return
		}

		code, err := parseToken(r, authHeader, secret)
		if err != nil {
			utils.RespondError(w, code, err.Error(), "")
			return
		}

		endpoint(w, r)
	})
}

func parseToken(r *http.Request, auth, secret string) (int, error) {
	splitToken := strings.Split(auth, "Bearer ")
	if len(splitToken) != 2 {
		return errs.ErrInvalidAuth.Code, errs.ErrInvalidAuth
	}

	reqToken := splitToken[1]
	if reqToken == "" {
		return errs.ErrInvalidAuth.Code, errs.ErrInvalidAuth
	}

	token, err := jwt.ParseWithClaims(reqToken, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("error authorizing user: %s\n", err.Error())
		return errs.ErrInvalidAuth.Code, errs.ErrInvalidAuth
	}

	if !token.Valid {
		return errs.ErrForbidden.Code, errs.ErrForbidden
	}

	claims, ok := token.Claims.(*utils.CustomClaims)
	if !ok {
		return errs.ErrInternal.Code, errs.ErrInternal
	}
	if !token.Valid {
		return errs.ErrForbidden.Code, errs.ErrForbidden
	}

	context.Set(r, "claims", claims)

	return 0, nil
}
