// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"encoding/json"
	"testing"

	"github.com/praelatus/praelatus/models"
)

func projectFromJSON(jsn []byte) (interface{}, error) {
	var tk models.Project
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func projectsFromJSON(jsn []byte) (interface{}, error) {
	var tk []models.Project
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func toProjects(v interface{}) []models.Project {
	return v.([]models.Project)
}

func toProject(v interface{}) models.Project {
	return v.(models.Project)
}

var projectRouteTests = []routeTest{
	{
		Name:      "Get Project",
		Endpoint:  "/api/v1/projects/TEST2",
		Converter: projectFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			project := toProject(v)

			if project.Key != "TEST2" {
				t.Errorf("Expected TEST2 Got: %s", project.Key)
			}
		},
	},
	{
		Name:      "Get All Projects",
		Endpoint:  "/api/v1/projects",
		Converter: projectsFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			projects := toProjects(v)

			if len(projects) <= 1 {
				t.Errorf("Expected More than 1 Project Got: %d", len(projects))
				return
			}

			if projects[0].Key != "TEST" {
				t.Errorf("Expected TEST Got: %s", projects[0].Key)
			}
		},
	},

	{
		Name:      "Create Project",
		Admin:     true,
		Method:    "POST",
		Endpoint:  "/api/v1/projects",
		Converter: projectFromJSON,
		Body: models.Project{
			Key:  "FAKEPROJ",
			Name: "Fake Project",
			Lead: "testuser",
		},
		Validator: func(v interface{}, t *testing.T) {
			project := toProject(v)

			if project.Key != "FAKEPROJ" {
				t.Errorf("Expected FAKEPROJ Got: %s", project.Key)
			}
		},
	},

	{
		Name:     "Remove Project",
		Endpoint: "/api/v1/projects/TEST2",
		Method:   "DELETE",
		Admin:    true,
	},
}

func TestProjectRoutes(t *testing.T) {
	testRoutes(projectRouteTests, t)
}
