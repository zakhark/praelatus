// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

// Event represents an event happening on a given ticket, Data contains
// additional data about the event for example if it is a transition event then
// the transition will be in Data, if it is a comment added event then Data
// will be the comment itself so on and so forth
type Event struct {
	Ticket Ticket      `json:"ticket"`
	Data   interface{} `json:"data"`
}
