package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

const baseURL = "http://localhost:8080"

func getToken(t *testing.T, role string) string {
	url := fmt.Sprintf("%s/dummyLogin", baseURL)
	payload := map[string]string{"role": role}
	jsonBody, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("failed to get jwt for %s: %v", role, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)

	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(bodyBytes))
}

func createPVZ(t *testing.T, token string) string {
	url := fmt.Sprintf("%s/pvz", baseURL)
	body := map[string]string{"city": "Москва"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create pvz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var res struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.ID
}

func createReception(t *testing.T, token, pvzID string) string {
	url := fmt.Sprintf("%s/receptions", baseURL)
	body := map[string]string{"pvzId": pvzID}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create reception: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var res struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.ID
}

func addProduct(t *testing.T, token, pvzID string, i int) {
	url := fmt.Sprintf("%s/products", baseURL)
	typeProduct := "электроника"
	if i >= 20 {
		typeProduct = "обувь"
	} else if i >= 35 {
		typeProduct = "одежда"
	}
	body := map[string]string{
		"type":  typeProduct,
		"pvzId": pvzID,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to add product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func closeReception(t *testing.T, token, pvzID string) {
	url := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to close reception: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var res struct {
		Status string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	if res.Status != "close" {
		t.Errorf("expected status 'close', got '%s'", res.Status)
	}
}

func TestFullReceptionFlow(t *testing.T) {
	moderatorJWT := getToken(t, "moderator")
	employeeJWT := getToken(t, "employee")

	pvzID := createPVZ(t, moderatorJWT)
	_ = createReception(t, employeeJWT, pvzID)

	for i := 1; i <= 50; i++ {
		addProduct(t, employeeJWT, pvzID, 1)

	}

	closeReception(t, employeeJWT, pvzID)
}
