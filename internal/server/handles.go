package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/astrorick/seekret/internal/api"
	"github.com/astrorick/seekret/internal/database"
	"github.com/astrorick/seekret/pkg/jwt"
)

// helper func for error responses
func writeErrorResponse(w http.ResponseWriter, httpStatus int, errorResponse *api.ServerErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(errorResponse)
}

// implement username checks here
func isValidUsername(username string) bool {
	// TODO: add strict username characters checks
	if len(username) < 4 || len(username) > 32 {
		return false
	}

	return true
}

// implement password checks here
func isValidPassword(password string) bool {
	// TODO: add strict password characters checks
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	return true
}

// CreateUserRequestHandler handles client request for the creation of a new user in the database.
// An error response is returned to the client if the request body cannot be parsed, if the username/password is not valid or if the user already exists.
// If the request is successful, a new JWT will be returned in the response body to be used by the client.
func (srv *Server) CreateUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodPost {
			writeErrorResponse(w, http.StatusMethodNotAllowed, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "method not allowed",
			})
			return
		}

		// read request body
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to read request body",
			})
			return
		}

		// parse request body
		var createUserRequest api.CreateUserRequest
		if err := json.Unmarshal(reqBody, &createUserRequest); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to parse request body",
			})
			return
		}

		// check parsed data
		if !isValidUsername(createUserRequest.Username) {
			writeErrorResponse(w, http.StatusBadRequest, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid username",
			})
			return
		}
		if !isValidPassword(createUserRequest.Password) {
			writeErrorResponse(w, http.StatusBadRequest, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid password",
			})
			return
		}

		// check if specified username exists
		userExists, err := srv.Database.UserExists(createUserRequest.Username)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database read error",
			})
			return
		}
		if userExists {
			writeErrorResponse(w, http.StatusConflict, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "username not available",
			})
			return
		}

		// generate new srp salt and verifier
		newUserSalt, err := srv.SRP.NewSalt()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal srp salt generation error",
			})
			return
		}
		newUserVerifier, err := srv.SRP.GetVerifier(newUserSalt, createUserRequest.Username, createUserRequest.Password)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal srp verifier generation error",
			})
			return
		}

		// add new user
		if err := srv.Database.CreateUser(&database.NewUser{
			Username: createUserRequest.Username,
			Salt:     newUserSalt,
			Verifier: newUserVerifier,
		}); err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// generate a jwt for the newly created user
		signedJWTString, err := srv.JWT.NewSignedJWTString(&jwt.MapClaims{
			"username": createUserRequest.Username,
			"exp":      time.Now().Add(24 * time.Second).Unix(),
		})
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, &api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal jwt generation error",
			})
			return
		}

		// send feedback
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(api.CreateUserResponse{
			JWT: signedJWTString,
		})
	}
}

// DeleteUserRequestHandler handles client request for removing an existing user from the database.
// An error response is returned to the client if the request body cannot be parsed, TODO.
func (srv *Server) DeleteUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodDelete {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "method not allowed",
			})
			return
		}

		// read request body
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to read request body",
			})
			return
		}

		// parse request body
		var deleteUserRequest api.DeleteUserRequest
		if err := json.Unmarshal(reqBody, &deleteUserRequest); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to parse request content",
			})
			return
		}

		// check if bearer token present
		authString := r.Header.Get("Authorization")
		if authString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "no authorization token provided",
			})
			return
		}
		authParts := strings.Split(authString, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "bad authentication token format",
			})
			return
		}

		// parse jwt and validate
		token, err := srv.JWT.ParseJWTString(authParts[1])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "bad authentication token format",
			})
			return
		}
		if !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid authorization token",
			})
			return
		}

		// check jwt claims to see if they match the user being deleted
		mapClaims, err := srv.JWT.GetClaims(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid token claims",
			})
			return
		}
		jwtUsername, ok := (*mapClaims)["username"].(string)
		if !ok || jwtUsername != deleteUserRequest.Username {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid or missing username claim",
			})
			return
		}
		jwtExp, ok := (*mapClaims)["exp"].(float64)
		if !ok || time.Unix(int64(jwtExp), 0).Before(time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "missing exp claim or jwt expired",
			})
			return
		}

		// check if specified username exists
		userExists, err := srv.Database.UserExists(deleteUserRequest.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database read error",
			})
			return
		}
		if !userExists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "specified user does not exist",
			})
			return
		}

		// delete user
		if err := srv.Database.DeleteUser(deleteUserRequest.Username); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// send feedback
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
