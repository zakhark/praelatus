// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"

	"github.com/gorilla/mux"
)

func userRouter(router *mux.Router) {
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/tokens", login).Methods("POST")

	router.HandleFunc("/users/me", loggedInUser)
	router.HandleFunc("/users/{username}", singleUser)
	router.HandleFunc("/users/{username}/avatar", avatar)
	router.HandleFunc("/users/{username}/leadof", leadOf)
}

func loggedInUser(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusNotFound, "no session found")
		return
	}

	utils.SendJSON(w, *u)
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginReq)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := Repo.Users().Get(&models.User{}, loginReq.Username)
	if err != nil {
		utils.Error(w, err)
		return
	}

	if user.CheckPw([]byte(loginReq.Password)) {
		err := middleware.SetUserSession(user, w)
		if err != nil {
			utils.Error(w, err)
			return
		}

		utils.SendJSON(w, user)
		return
	}

	utils.APIErr(w, http.StatusForbidden, "invalid password")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	loggedInUser := middleware.GetUserSession(r)
	if loggedInUser == nil {
		loggedInUser = &models.User{}
	}

	var u *models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.ValidateModel(u); err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if !loggedInUser.IsAdmin {
		u.IsAdmin = false
	}

	u, err = models.NewUser(u.Username, u.Password, u.FullName, u.Email, u.IsAdmin)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = middleware.SetUserSession(*u, w)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, u)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)

	}

	var users models.Users
	var err error

	users, err = Repo.Users().Search(u, q)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, users)
}

func singleUser(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var user models.User
	var err error

	username := mux.Vars(r)["username"]

	switch r.Method {
	case "GET":
		user, err = Repo.Users().Get(u, username)
	case "DELETE":
		err = Repo.Users().Delete(u, username)
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&user)
		if err != nil {
			break
		}

		err = Repo.Users().Update(u, username, user)
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	user.Password = ""
	utils.SendJSON(w, user)
}

func avatar(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var user models.User
	var err error

	username := mux.Vars(r)["username"]

	user, err = Repo.Users().Get(u, username)
	if err != nil {
		utils.Error(w, err)
		return
	}

	w.Write([]byte(user.ProfilePic))
}

func leadOf(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	username := mux.Vars(r)["username"]

	projects, err := Repo.Projects().HasLead(u, models.User{Username: username})
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, projects)
}
