package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"be-project-monitoring/internal/domain/model"
)

type (
	createUserReq struct {
		Email          string
		Username       string
		FirstName      string
		LastName       string
		Group          string
		GithubUsername string
		Password       string
	}
	authUserReq struct {
		Username string
		Password string
	}

	createUserResp struct {
		*model.User
		Token string
	}
)

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	userReq := &createUserReq{}
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		s.response.ErrorWithCode(w, err, http.StatusBadRequest)
	}

	newUser := &model.User{
		Role:           model.Student,
		Email:          userReq.Email,
		Username:       userReq.Username,
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
		Group:          userReq.Group,
		GithubUsername: userReq.GithubUsername,
		HashedPassword: hashPass(userReq.Password),
	}

	user, token, err := s.svc.CreateUser(r.Context(), newUser)
	if err != nil {
		s.response.ErrorWithCode(w, err, http.StatusInternalServerError)
	}

	s.response.Created(w, createUserResp{
		User:  user,
		Token: token,
	})

}

func hashPass(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
