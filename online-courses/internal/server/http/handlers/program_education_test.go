package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestCreateProgram(t *testing.T) {

	url := "http://localhost:8080/api/v1/programeducation"
	jsonBody :=
		`{
    "name_prof_education": "Курсы по ujdyeasdf базовый уровень",
    "time_education": 31,
    "price": 16000,
    "id_educationtype": "0774eba8-e911-4080-a82b-88539f73519b",
    "id_divisionseducation": "05df98bb-f9bc-48ff-8e4a-e69cf9e511aa"
	}`

	req, _ := http.NewRequest("POST", url, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MiwidXNlcm5hbWUiOiIxIiwicm9sZSI6IndvcmtlciIsInN1YiI6IjEiLCJleHAiOjE3NzY4NTk1MDksImlhdCI6MTc3NjgxNjMwOSwianRpIjoiNGViNzllZTItNThkZi00NTQ4LTk0Y2EtMjFjMDU4MTczNWY5In0.itVPbZLh4XuLeSkRN0lK6SIk9-y7lyQOynzZWyZF-xs")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("got %d", resp.StatusCode)
	}

}
