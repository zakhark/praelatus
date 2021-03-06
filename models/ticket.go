// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Ticket represents a ticket
type Ticket struct {
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
	Key         string    `bson:"_id" json:"key"`
	Summary     string    `json:"summary" required:"true"`
	Description string    `json:"description" required:"true"`
	Status      string    `json:"status"`
	Reporter    string    `json:"reporter" required:"true"`
	Assignee    string    `json:"assignee"`
	Type        string    `json:"type" required:"true"`
	Labels      []string  `json:"labels"`

	Fields   []Field   `json:"fields"`
	Comments []Comment `json:"comments,omitempty"`

	Workflow bson.ObjectId `json:"workflow"`
	Project  string        `json:"project" required:"true"`
}

func (t Ticket) String() string {
	return jsonString(t)
}

// Transition searches through the available transitions for the ticket
// returning a boolean indicating success or failure and the transition
func (t Ticket) Transition(db *mgo.Database, name string) (Transition, bool) {
	var workflow Workflow

	err := db.C("workflows").FindId(t.Workflow).One(&workflow)
	if err != nil {
		log.Println(err.Error())
		return Transition{}, false
	}

	for _, transition := range workflow.Transitions {
		if transition.Name == name && t.Status == transition.FromStatus {
			return transition, true
		}
	}

	return Transition{}, false
}

// Comment is a comment on an issue / ticket.
type Comment struct {
	UpdatedDate time.Time `json:"updatedDate"`
	CreatedDate time.Time `json:"createdDate"`
	Body        string    `json:"body" required:"true"`
	Author      string    `json:"author" required:"true"`
}

func (c *Comment) String() string {
	return jsonString(c)
}
