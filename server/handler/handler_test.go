package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/api"
	"server/model"
	"server/persistence"
	"testing"

	"gorm.io/datatypes"
)

var database *persistence.Database

func TestSteps(t *testing.T) {
	installer := api.NewApp(database)
	req, err := http.NewRequest("GET", "/step", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	hndlr := http.HandlerFunc(installer.GetAllSteps)

	hndlr.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("Get all steps expected code 200, got %d", rec.Code)
	}
	var response []map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	if len(response) != 4 {
		t.Errorf("Expected 4 steps, got %d", len(response))
	}
	for _, step := range response {
		switch step["name"] {
		case "step1", "step2", "step3", "step4":
			continue
		default:
			t.Errorf("Unexpected step in get all steps: %s", step["name"])
		}
	}
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	database = persistence.NewInMemoryDB()

	step1 := model.Step{}
	step1.Name = "step1"
	step1.Template = datatypes.JSONMap{}
	step1.Priority = 0

	step2 := model.Step{}
	step2.Name = "step2"
	step2.Template = datatypes.JSONMap{}
	step2.Priority = 1

	step3 := model.Step{}
	step3.Name = "step3"
	step3.Template = datatypes.JSONMap{}
	step3.Priority = 1

	step4 := model.Step{}
	step4.Name = "step4"
	step4.Template = datatypes.JSONMap{}
	step4.Priority = 2

	database.Instance.Save(&step1)
	database.Instance.Save(&step2)
	database.Instance.Save(&step3)
	database.Instance.Save(&step4)

	m.Run()
}
