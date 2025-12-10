package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActivateNodeBackup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/cdpbackup/activate/node123/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"status":              "active",
			"detail":              "Backup activated successfully",
			"last_recovery_point": "2025-01-01T10:00:00Z",
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Backup activated successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	config := &BackupConfig{
		PlanID:          1,
		Type:            "DAILY",
		CompressionType: "GZip",
	}

	result, _, err := client.Backup.ActivateNodeBackup(context.Background(), "node123", config)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.Status != "active" {
		t.Errorf("Expected Status active, got %s", result.Status)
	}

	if result.Detail != "Backup activated successfully" {
		t.Errorf("Expected Detail 'Backup activated successfully', got %s", result.Detail)
	}
}

func TestActivateNodeBackup_EmptyNodeID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "", &BackupConfig{})
	assertError(t, err, "")
}

func TestActivateNodeBackup_NilConfig(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", nil)
	assertError(t, err, "")
}

func TestActivateNodeBackup_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestActivateNodeBackup_BadRequest(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Invalid backup configuration")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "INVALID"})
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/cdpbackup/deactivate/node123/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "node123")
	assertNoError(t, err)
}

func TestDeactivateNodeBackup_EmptyNodeID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "")
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_NotFound(t *testing.T) {
	server := newErrorServer(t, http.StatusNotFound, "Node not found")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "notfound")
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetNodeBackupStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/cdpbackup/node123/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"status":              "active",
			"detail":              "Backup is configured",
			"last_recovery_point": "2025-01-01T10:00:00Z",
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "success", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.Status != "active" {
		t.Errorf("Expected Status active, got %s", result.Status)
	}

	if result.Detail != "Backup is configured" {
		t.Errorf("Expected Detail 'Backup is configured', got %s", result.Detail)
	}
}

func TestGetNodeBackupStatus_EmptyNodeID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetNodeBackupStatus(context.Background(), "")
	assertError(t, err, "")
}

func TestGetNodeBackupStatus_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetNodeBackupStatus(context.Background(), "notfound")
	assertNoError(t, err)

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestGetNodeBackupStatus_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetBackupAgentStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/cdpbackup/node123/cdp-agent/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"status":  "running",
			"version": "1.0.0",
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "success", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetBackupAgentStatus(context.Background(), "node123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.Status != "running" {
		t.Errorf("Expected Status running, got %s", result.Status)
	}

	if result.Version != "1.0.0" {
		t.Errorf("Expected Version 1.0.0, got %s", result.Version)
	}
}

func TestGetBackupAgentStatus_EmptyNodeID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetBackupAgentStatus(context.Background(), "")
	assertError(t, err, "")
}

func TestGetBackupAgentStatus_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetBackupAgentStatus(context.Background(), "notfound")
	assertNoError(t, err)

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestGetBackupAgentStatus_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetBackupAgentStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestListNodeBackupStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/cdpbackup/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := []interface{}{
			map[string]interface{}{
				"status":              "active",
				"detail":              "Backup is configured",
				"last_recovery_point": "2025-01-01T10:00:00Z",
			},
			map[string]interface{}{
				"status": "inactive",
				"detail": "Backup not configured",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.ListNodeBackupStatus(context.Background())
	assertNoError(t, err)

	if len(result) != 2 {
		t.Errorf("Expected 2 backup statuses, got %d", len(result))
	}

	if result[0].Status != "active" {
		t.Errorf("Expected first status 'active', got %s", result[0].Status)
	}

	if result[1].Status != "inactive" {
		t.Errorf("Expected second status 'inactive', got %s", result[1].Status)
	}
}

func TestListNodeBackupStatus_EmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.ListNodeBackupStatus(context.Background())
	assertNoError(t, err)

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestListNodeBackupStatus_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ListNodeBackupStatus(context.Background())
	assertError(t, err, "")
}

func TestParseHoursOfDay_Success(t *testing.T) {
	hours := []string{"9", "12", "18"}
	result := ParseHoursOfDay(hours)
	expected := "9,12,18"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestParseHoursOfDay_EmptySlice(t *testing.T) {
	hours := []string{}
	result := ParseHoursOfDay(hours)
	expected := ""
	if result != expected {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestParseHoursOfDay_SingleHour(t *testing.T) {
	hours := []string{"9"}
	result := ParseHoursOfDay(hours)
	expected := "9"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestParseHoursOfDayFromString_Success(t *testing.T) {
	hoursStr := "9,12,18"
	result := ParseHoursOfDayFromString(hoursStr)
	expected := []string{"9", "12", "18"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d hours, got %d", len(expected), len(result))
	}
	for i, h := range result {
		if h != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], h)
		}
	}
}

func TestParseHoursOfDayFromString_EmptyString(t *testing.T) {
	result := ParseHoursOfDayFromString("")
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
}

func TestParseHoursOfDayFromString_WithSpaces(t *testing.T) {
	hoursStr := "9, 12, 18"
	result := ParseHoursOfDayFromString(hoursStr)
	expected := []string{"9", "12", "18"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d hours, got %d", len(expected), len(result))
	}
	for i, h := range result {
		if h != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], h)
		}
	}
}

func TestParseHoursOfDayFromString_EmptyElements(t *testing.T) {
	hoursStr := "9,,12,18"
	result := ParseHoursOfDayFromString(hoursStr)
	expected := []string{"9", "12", "18"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d hours, got %d", len(expected), len(result))
	}
}

func TestBackupConfig_WithCompression(t *testing.T) {
	config := &BackupConfig{
		PlanID:           1,
		Type:             "DAILY",
		CompressionType:  "GZip",
		CompressionLevel: "High",
	}

	if config.PlanID != 1 {
		t.Errorf("Expected PlanID 1, got %d", config.PlanID)
	}
	if config.CompressionType != "GZip" {
		t.Errorf("Expected CompressionType GZip, got %s", config.CompressionType)
	}
	if config.CompressionLevel != "High" {
		t.Errorf("Expected CompressionLevel High, got %s", config.CompressionLevel)
	}
}

func TestBackupConfig_WithEncryption(t *testing.T) {
	config := &BackupConfig{
		PlanID:               1,
		Type:                 "DAILY",
		IsEncryptionRequired: true,
		EncryptionPassphrase: "secret123",
	}

	if !config.IsEncryptionRequired {
		t.Error("Expected encryption to be required")
	}
	if config.EncryptionPassphrase != "secret123" {
		t.Errorf("Expected passphrase 'secret123', got %s", config.EncryptionPassphrase)
	}
}

func TestBackupConfig_WithDatabase(t *testing.T) {
	config := &BackupConfig{
		PlanID:     1,
		Type:       "DAILY",
		DBEnabled:  true,
		DBUsername: "admin",
		DBPassword: "dbpass",
	}

	if !config.DBEnabled {
		t.Error("Expected DB to be enabled")
	}
	if config.DBUsername != "admin" {
		t.Errorf("Expected DBUsername 'admin', got %s", config.DBUsername)
	}
	if config.DBPassword != "dbpass" {
		t.Errorf("Expected DBPassword 'dbpass', got %s", config.DBPassword)
	}
}

func TestBackupConfig_WithExcludePaths(t *testing.T) {
	config := &BackupConfig{
		PlanID:            1,
		Type:              "DAILY",
		ExcludeFileFolder: "/tmp,/var/log",
	}

	if config.ExcludeFileFolder != "/tmp,/var/log" {
		t.Errorf("Expected ExcludeFileFolder '/tmp,/var/log', got %s", config.ExcludeFileFolder)
	}
}

func TestBackupConfig_WithSchedule(t *testing.T) {
	config := &BackupConfig{
		PlanID:         1,
		Type:           "HOURLY",
		HoursOfDay:     "9,12,18",
		StartingMinute: 30,
	}

	if config.HoursOfDay != "9,12,18" {
		t.Errorf("Expected HoursOfDay '9,12,18', got %s", config.HoursOfDay)
	}
	if config.StartingMinute != 30 {
		t.Errorf("Expected StartingMinute 30, got %d", config.StartingMinute)
	}
}

func TestBackupConfig_WithBackupNow(t *testing.T) {
	config := &BackupConfig{
		PlanID:    1,
		Type:      "DAILY",
		BackupNow: true,
	}

	if !config.BackupNow {
		t.Error("Expected BackupNow to be true")
	}
}

func TestBackupStatus_AllFields(t *testing.T) {
	status := &BackupStatus{
		Status:            "active",
		Detail:            "Backup is configured",
		LastRecoveryPoint: "2025-01-01T10:00:00Z",
	}

	if status.Status != "active" {
		t.Errorf("Expected Status 'active', got %s", status.Status)
	}
	if status.Detail != "Backup is configured" {
		t.Errorf("Expected Detail 'Backup is configured', got %s", status.Detail)
	}
	if status.LastRecoveryPoint != "2025-01-01T10:00:00Z" {
		t.Errorf("Expected LastRecoveryPoint '2025-01-01T10:00:00Z', got %s", status.LastRecoveryPoint)
	}
}

func TestBackupAgentStatus_AllFields(t *testing.T) {
	status := &BackupAgentStatus{
		Status:  "running",
		Version: "1.0.0",
	}

	if status.Status != "running" {
		t.Errorf("Expected Status 'running', got %s", status.Status)
	}
	if status.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got %s", status.Version)
	}
}

func TestGetNodeBackupStatus_Non404Error(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetBackupAgentStatus_Non404Error(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetBackupAgentStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetNodeBackupStatus_BadGateway(t *testing.T) {
	server := newErrorServer(t, http.StatusBadGateway, "Bad gateway")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetBackupAgentStatus_BadGateway(t *testing.T) {
	server := newErrorServer(t, http.StatusBadGateway, "Bad gateway")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetBackupAgentStatus(context.Background(), "node123")
	assertError(t, err, "")
}

// Additional HTTP Error Status Code Tests (per CLAUDE.md guidelines)

func TestActivateNodeBackup_Unauthorized(t *testing.T) {
	server := newErrorServer(t, http.StatusUnauthorized, "Unauthorized")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestActivateNodeBackup_Forbidden(t *testing.T) {
	server := newErrorServer(t, http.StatusForbidden, "Forbidden")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestActivateNodeBackup_Conflict(t *testing.T) {
	server := newErrorServer(t, http.StatusConflict, "Backup already active")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestActivateNodeBackup_ServiceUnavailable(t *testing.T) {
	server := newErrorServer(t, http.StatusServiceUnavailable, "Service unavailable")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestActivateNodeBackup_GatewayTimeout(t *testing.T) {
	server := newErrorServer(t, http.StatusGatewayTimeout, "Gateway timeout")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_Unauthorized(t *testing.T) {
	server := newErrorServer(t, http.StatusUnauthorized, "Unauthorized")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "node123")
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_Forbidden(t *testing.T) {
	server := newErrorServer(t, http.StatusForbidden, "Forbidden")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "node123")
	assertError(t, err, "")
}

func TestDeactivateNodeBackup_Conflict(t *testing.T) {
	server := newErrorServer(t, http.StatusConflict, "Backup operation in progress")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Backup.DeactivateNodeBackup(context.Background(), "node123")
	assertError(t, err, "")
}

// Response Parsing Tests (per CLAUDE.md guidelines)

func TestActivateNodeBackup_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ActivateNodeBackup(context.Background(), "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
	assertError(t, err, "")
}

func TestGetNodeBackupStatus_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestGetNodeBackupStatus_MissingRequiredFields(t *testing.T) {
	server := newMissingFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"detail": "Backup is configured",
		},
	})
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertNoError(t, err)
	if result == nil {
		t.Fatal("Expected result object even with missing fields")
	}
	// Status should be empty string (zero value)
	if result.Status != "" {
		t.Errorf("Expected empty status, got %s", result.Status)
	}
}

func TestGetNodeBackupStatus_NullFieldValues(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"status":              "active",
			"detail":              nil,
			"last_recovery_point": nil,
		},
	})
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Backup.GetNodeBackupStatus(context.Background(), "node123")
	assertNoError(t, err)
	if result == nil {
		t.Fatal("Expected result object")
	}
	if result.Status != "active" {
		t.Errorf("Expected Status 'active', got %s", result.Status)
	}
	// Null fields should be empty strings (zero values)
	if result.Detail != "" {
		t.Errorf("Expected empty Detail for null, got %s", result.Detail)
	}
}

func TestGetBackupAgentStatus_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.GetBackupAgentStatus(context.Background(), "node123")
	assertError(t, err, "")
}

func TestListNodeBackupStatus_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Backup.ListNodeBackupStatus(context.Background())
	assertError(t, err, "")
}

func TestListNodeBackupStatus_EmptyDataField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Missing "data" field - should handle gracefully
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "success", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, _ = client.Backup.ListNodeBackupStatus(context.Background())
	// Should handle missing data field gracefully - error is acceptable
}

// Network/Timeout Tests (per CLAUDE.md guidelines)

func TestActivateNodeBackup_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.Backup.ActivateNodeBackup(ctx, "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
		return err
	}
	testNetworkError(t, testFunc)
}

func TestActivateNodeBackup_ContextTimeout(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.Backup.ActivateNodeBackup(ctx, "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
		return err
	}
	testContextTimeout(t, testFunc)
}

func TestActivateNodeBackup_ContextCancellation(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.Backup.ActivateNodeBackup(ctx, "node123", &BackupConfig{PlanID: 1, Type: "DAILY"})
		return err
	}
	testContextCancellation(t, testFunc)
}

func TestGetNodeBackupStatus_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.Backup.GetNodeBackupStatus(ctx, "node123")
		return err
	}
	testNetworkError(t, testFunc)
}

func TestGetNodeBackupStatus_ContextTimeout(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.Backup.GetNodeBackupStatus(ctx, "node123")
		return err
	}
	testContextTimeout(t, testFunc)
}

func TestDeactivateNodeBackup_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, err := client.Backup.DeactivateNodeBackup(ctx, "node123")
		return err
	}
	testNetworkError(t, testFunc)
}

func TestDeactivateNodeBackup_ContextCancellation(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, err := client.Backup.DeactivateNodeBackup(ctx, "node123")
		return err
	}
	testContextCancellation(t, testFunc)
}
