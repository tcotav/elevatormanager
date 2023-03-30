package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

// ref - https://gin-gonic.com/docs/testing/

func TestPushDestination(t *testing.T) {
	router := setupRouter()

	// first call elevator to our floor
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callElevator/2/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	retMap := make(map[string]int)
	err := json.Unmarshal(w.Body.Bytes(), &retMap)
	if err != nil {
		t.Errorf("Expected no error on json unmarshal of callelevator, got %s", err.Error())
	}
	elevatorID := retMap["elevatorID"]

	// then push a destination button
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/pushDestination/%d/4", elevatorID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// test out of bounds elevator
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/pushDestination/10/2", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}

	// test out of bounds floor
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/pushDestination/%d/102", elevatorID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}
}
func TestCallElevator(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callElevator/7/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// test out of bounds floor
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/callElevator/102/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}

	// test out of bounds direction
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/callElevator/2/3", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}
}

func TestResetElevator(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/resetElevator/0", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// test out of bounds elevator
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/resetElevator/10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}
}

func TestGetAllElevatorState(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/getAllElevatorState", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

func TestMaintenanceCallOverride(t *testing.T) {
	router := setupRouter()

	// move elevator to another floor
	// first call elevator to our floor
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callElevator/3/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	retMap := make(map[string]int)
	err := json.Unmarshal(w.Body.Bytes(), &retMap)
	if err != nil {
		t.Errorf("Expected no error on json unmarshal of callelevator, got %s", err.Error())
	}
	elevatorID := retMap["elevatorID"]

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/maintenanceCallOverride/%d/6/1", elevatorID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// test out of bounds elevator
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/maintenanceCallOverride/100/2/1", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}
}