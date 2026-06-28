package testhelper

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"sync"
	"time"
)

// APIRecord captures a single API request/response pair
type APIRecord struct {
	TestName       string            `json:"test_name"`
	Category       string            `json:"category"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	RequestBody    string            `json:"request_body"`
	StatusCode     int               `json:"status_code"`
	ResponseBody   string            `json:"response_body"`
	ResponseHeader map[string]string `json:"response_header"`
	Duration       time.Duration     `json:"duration"`
	Timestamp      time.Time         `json:"timestamp"`
}

// ReportData holds the complete test report
type ReportData struct {
	GeneratedAt time.Time   `json:"generated_at"`
	APIRecords  []APIRecord `json:"api_records"`
}

// ReportRecorder captures API interactions during tests
type ReportRecorder struct {
	mu       sync.Mutex
	records  []APIRecord
	current  string
	category string
	enabled  bool
}

// Global recorder instance
var Recorder = &ReportRecorder{enabled: true}

// StartTest marks the beginning of a named test for response capture
func StartTest(name, category string) {
	Recorder.mu.Lock()
	defer Recorder.mu.Unlock()
	Recorder.current = name
	Recorder.category = category
}

// RecordAPI captures an API request/response pair
func RecordAPI(method, path string, reqBody interface{}, w *httptest.ResponseRecorder, dur time.Duration) {
	if !Recorder.enabled {
		return
	}

	Recorder.mu.Lock()
	defer Recorder.mu.Unlock()

	record := APIRecord{
		TestName:       Recorder.current,
		Category:       Recorder.category,
		Method:         method,
		Path:           path,
		StatusCode:     w.Code,
		Timestamp:      time.Now(),
		Duration:       dur,
		ResponseHeader: make(map[string]string),
	}

	// Capture request body
	if reqBody != nil {
		if b, err := json.MarshalIndent(reqBody, "", "  "); err == nil {
			record.RequestBody = string(b)
		}
	}

	// Capture response body (truncate if too long)
	body := w.Body.String()
	if len(body) > 4000 {
		body = body[:4000] + "\n... (truncated)"
	}
	record.ResponseBody = body

	// Capture response headers
	for k, v := range w.Header() {
		if len(v) > 0 {
			record.ResponseHeader[k] = v[0]
		}
	}

	Recorder.records = append(Recorder.records, record)
}

// SaveReport writes the captured data to a JSON file
func SaveReport(path string) error {
	Recorder.mu.Lock()
	defer Recorder.mu.Unlock()

	data := ReportData{
		GeneratedAt: time.Now(),
		APIRecords:  Recorder.records,
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create report file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// GetRecords returns a copy of all captured records
func GetRecords() []APIRecord {
	Recorder.mu.Lock()
	defer Recorder.mu.Unlock()
	result := make([]APIRecord, len(Recorder.records))
	copy(result, Recorder.records)
	return result
}

// Reset clears all recorded data
func ResetRecorder() {
	Recorder.mu.Lock()
	defer Recorder.mu.Unlock()
	Recorder.records = nil
	Recorder.current = ""
	Recorder.category = ""
}

// StatusText returns a human-readable status description
func StatusText(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "Success"
	case code >= 300 && code < 400:
		return "Redirect"
	case code == 400:
		return "Bad Request"
	case code == 401:
		return "Unauthorized"
	case code == 403:
		return "Forbidden"
	case code == 404:
		return "Not Found"
	case code == 409:
		return "Conflict"
	case code >= 500:
		return "Server Error"
	default:
		return fmt.Sprintf("Status %d", code)
	}
}
