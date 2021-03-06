// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

func projectRouter(router *mux.Router) {
	router.HandleFunc("/projects", GetAllProjects).Methods("GET")
	router.HandleFunc("/projects", CreateProject).Methods("POST")

	router.HandleFunc("/projects/{key}", SingleProject)
}

// SingleProject will get a project by it's project key
func SingleProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)

	var p models.Project
	var err error

	key := mux.Vars(r)["key"]

	switch r.Method {
	case "GET":
		p, err = Repo.Projects().Get(u, key)
	case "DELETE":
		err = Repo.Projects().Delete(u, key)
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&p)
		if err != nil {
			break
		}

		err = Repo.Projects().Update(u, key, p)
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, p)
}

// GetAllProjects will get all the projects on this instance that the user has
// permissions to
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
	}

	projects, err := Repo.Projects().Search(u, q)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, projects)
}

// CreateProject will create a project in the database. This route requires
// admin.
func CreateProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.Error(w, repo.ErrUnauthorized)
		return
	}

	var p models.Project

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.ValidateModel(p); err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err = Repo.Projects().Create(u, p)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, p)
}
