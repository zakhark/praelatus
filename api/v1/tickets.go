// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"
)

func ticketRouter(router *mux.Router) {
	router.HandleFunc("/tickets", getAllTickets).Methods("GET")
	router.HandleFunc("/tickets", createTicket).Methods("POST")
	router.HandleFunc("/tickets/{key}", singleTicket)

	router.HandleFunc("/tickets/{key}/addComment", addComment).Methods("POST")
}

func createTicket(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var t models.Ticket

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.ValidateModel(t); err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err = Repo.Tickets().Create(u, t)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, t)
}

func singleTicket(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var t models.Ticket
	var err error

	id := mux.Vars(r)["key"]

	switch r.Method {
	case "GET":
		t, err = Repo.Tickets().Get(u, id)
	case "DELETE":
		err = Repo.Tickets().Delete(u, id)
	case "PUT":
		var ticket models.Ticket

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&ticket)
		if err != nil {
			break
		}

		err = Repo.Tickets().Update(u, id, ticket)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, t)
}

// getAllTickets will return all tickets which the user has permissions to.
func getAllTickets(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	q := r.FormValue("q")
	tickets, err := Repo.Tickets().Search(u, q)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, tickets)
}

func addComment(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusForbidden, "you must be logged in to comment")
		return
	}

	var c models.Comment

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.ValidateModel(c)
	if err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	key := mux.Vars(r)["key"]

	ticket, err := Repo.Tickets().AddComment(u, key, c)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, ticket)
}
