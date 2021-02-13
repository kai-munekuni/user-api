package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/kai-munekuni/user-api/internal/domain/model"
	"github.com/kai-munekuni/user-api/internal/domain/repository"
)

var (
	_ http.Handler = (*signupHandler)(nil)
	_ http.Handler = (*healthCheckHandler)(nil)
	_ http.Handler = (*getUserHandler)(nil)
	_ http.Handler = (*patchUserHandler)(nil)
	_ http.Handler = (*deleteUserHandler)(nil)
)

var userIDRe = regexp.MustCompile("^[a-zA-Z0-9]{6,20}$")
var passwordRe = regexp.MustCompile("^[^\r\n\t\f\v\a\v ]{8,20}$")
var nicknameRe = regexp.MustCompile("^[^\a\f\r\v]{0,30}$")
var commentRe = regexp.MustCompile("^[^\a\f\r\v]{0,100}$")

type healthCheckHandler struct {
}

func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	success(w,
		map[string]interface{}{
			"message": "healthy",
		})
}

type signupHandler struct {
	userRepo repository.User
}

func (h *signupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		internalServerError(w, "marshal error")
		return
	}
	if params["user_id"] == "" {
		badRequest(w, "Account creation failed", "required user_id is missing")
		return
	}

	if params["password"] == "" {
		badRequest(w, "Account creation failed", "required password is missing")
		return
	}
	userID := params["user_id"]
	password := params["password"]

	if !userIDRe.Match([]byte(userID)) {
		badRequest(w, "Account creation failed", "invalid user_id pattern")
		return
	}

	if !passwordRe.Match([]byte(password)) {
		badRequest(w, "Account creation failed", "invalid password pattern")
		return
	}

	err := h.userRepo.CreateUser(r.Context(), userID, password)
	if err != nil {
		log.Printf("%+v", err)
		if errors.Is(err, model.ErrUserAlreadyExists) {
			badRequest(w, "Account creation failed", "already same user_id is used")
		} else {
			internalServerError(w, "something went wrong")
		}
		return
	}

	success(w,
		map[string]interface{}{
			"message": "Account successfully created",
			"user":    signupResponse{UserID: userID, NickName: userID},
		})
}

type signupResponse struct {
	UserID   string `json:"user_id"`
	NickName string `json:"nickname"`
}

type getUserHandler struct {
	userRepo repository.User
}

func (h *getUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := h.userRepo.Find(r.Context(), vars["id"])

	if err != nil {
		log.Printf("%+v", err)
		notFound(w)
		return
	}
	userBody := make(map[string]string)
	userBody["user_id"] = user.ID
	userBody["nickname"] = user.Nickname
	if user.Comment != "" {
		userBody["comment"] = user.Comment
	}

	success(w,
		map[string]interface{}{
			"message": "User details by user_id",
			"user":    userBody,
		})
}

type patchUserHandler struct {
	userRepo repository.User
}

func (h *patchUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := getUserID(ctx)
	if userID != vars["id"] {
		forbidden(w)
		return
	}
	var params map[string]*string
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		badRequest(w, "User updation failed", "invalid")
		return
	}

	if params["user_id"] != nil || params["password"] != nil {
		badRequest(w, "User updation failed", "not updatable user_id and password")
		return
	}

	nickname := params["nickname"]
	comment := params["comment"]
	if nickname == nil && comment == nil {
		badRequest(w, "User updation failed", "required nickname or comment")
		return
	}

	if nickname != nil && !nicknameRe.Match([]byte(string(*nickname))) {
		badRequest(w, "User updation failed", "invalid nickname pattern")
		return
	}

	if comment != nil && !commentRe.Match([]byte(string(*comment))) {
		badRequest(w, "User updation failed", "invalid comment pattern")
		return
	}

	user, err := h.userRepo.UpdateField(ctx, userID, nickname, comment)

	if err != nil {
		log.Printf("%+v", err)
		notFound(w)
		return
	}

	success(w,
		map[string]interface{}{
			"message": "User successfully updated",
			"recipe": []patchUserResponse{
				{Nickname: user.Nickname, Comment: user.Comment},
			},
		})
}

type patchUserResponse struct {
	Nickname string `json:"nickname"`
	Comment  string `json:"comment"`
}

type deleteUserHandler struct {
	userRepo repository.User
}

func (h *deleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserID(ctx)
	err := h.userRepo.Delete(ctx, userID)
	if err != nil {
		log.Printf("%+v", err)
		internalServerError(w, "db error")
		return
	}

	success(w,
		map[string]interface{}{
			"message": "Account and user successfully removed",
		})
}
