package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// APIRecord from test capture
type APIRecord struct {
	TestName       string            `json:"test_name"`
	Category       string            `json:"category"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	RequestBody    string            `json:"request_body"`
	StatusCode     int               `json:"status_code"`
	ResponseBody   string            `json:"response_body"`
	ResponseHeader map[string]string `json:"response_header"`
	Duration       int64             `json:"duration"`
	Timestamp      string            `json:"timestamp"`
}

// APIData holds the captured API data
type APIData struct {
	GeneratedAt string      `json:"generated_at"`
	APIRecords  []APIRecord `json:"api_records"`
}

func loadAPIData(path string) (*APIData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var apiData APIData
	if err := json.Unmarshal(data, &apiData); err != nil {
		return nil, err
	}
	return &apiData, nil
}

func generatePDFReport(report ReportData, apiDataPath, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 25)
	pdf.SetMargins(20, 20, 20)

	// Load API data
	var apiData *APIData
	if data, err := loadAPIData(apiDataPath); err == nil {
		apiData = data
	}

	// === COVER PAGE ===
	drawCoverPage(pdf, report)

	// === TABLE OF CONTENTS ===
	drawTableOfContents(pdf, report, apiData)

	// === TEST RESULTS BY CATEGORY ===
	for _, cat := range report.Categories {
		drawCategoryPage(pdf, cat)
	}

	// === API ENDPOINT DOCUMENTATION ===
	if apiData != nil && len(apiData.APIRecords) > 0 {
		drawAPIDocumentation(pdf, apiData)
	}

	// === SECURITY AUDIT ===
	drawSecurityAuditPage(pdf, report, apiData)

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return fmt.Errorf("PDF write error: %w", err)
	}
	return nil
}

// ============== COVER PAGE ==============

func drawCoverPage(pdf *gofpdf.Fpdf, report ReportData) {
	pdf.AddPage()

	// Title area
	pdf.SetFillColor(30, 58, 95)
	pdf.Rect(0, 0, 210, 120, "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 28)
	pdf.SetY(30)
	pdf.CellFormat(0, 12, "SIPENA API", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 16)
	pdf.CellFormat(0, 10, "Test & Security Report", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetY(70)
	pdf.CellFormat(0, 7, "Sistem Informasi Pernikahan", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "Comprehensive API Testing & Security Audit", "", 1, "C", false, 0, "")

	// Stats cards
	pdf.SetY(135)
	pdf.SetTextColor(60, 60, 60)
	cardW := 38.0
	gap := 4.0
	startX := (210 - (5*cardW + 4*gap)) / 2

	cards := []struct {
		label string
		value string
		r, g, b int
	}{
		{"Total Tests", fmt.Sprintf("%d", report.TotalTests), 59, 130, 246},
		{"Passed", fmt.Sprintf("%d", report.TotalPassed), 34, 197, 94},
		{"Failed", fmt.Sprintf("%d", report.TotalFailed), 239, 68, 68},
		{"Skipped", fmt.Sprintf("%d", report.TotalSkipped), 245, 158, 11},
		{"Pass Rate", fmt.Sprintf("%.1f%%", report.PassRate), 59, 130, 246},
	}

	for i, card := range cards {
		x := startX + float64(i)*(cardW+gap)
		pdf.SetFillColor(card.r, card.g, card.b)
		pdf.SetXY(x, 135)
		pdf.Rect(x, 135, cardW, 28, "F")

		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 18)
		pdf.SetXY(x, 137)
		pdf.CellFormat(cardW, 10, card.value, "", 1, "C", false, 0, "")
		pdf.SetFont("Helvetica", "", 7)
		pdf.SetXY(x, 149)
		pdf.CellFormat(cardW, 6, card.label, "", 1, "C", false, 0, "")
	}

	// Info section
	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetY(180)
	pdf.CellFormat(0, 7, fmt.Sprintf("Generated: %s", report.GeneratedAt.Format("Monday, 02 January 2006 - 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Duration: %s", report.TotalDuration.Round(time.Millisecond)), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "Test Environment: In-Memory SQLite (Isolated)", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "Framework: Go Testing + Gin httptest", "", 1, "C", false, 0, "")
}

// ============== TABLE OF CONTENTS ==============

func drawTableOfContents(pdf *gofpdf.Fpdf, report ReportData, apiData *APIData) {
	pdf.AddPage()
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(0, 12, "Table of Contents", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetTextColor(60, 60, 60)
	items := []string{
		"1. Cover Page",
		"2. Table of Contents",
		"3. Test Results Summary",
	}

	for i, cat := range report.Categories {
		items = append(items, fmt.Sprintf("   3.%d %s (%d tests)", i+1, cat.Name, len(cat.Tests)))
	}

	items = append(items, "4. API Endpoint Documentation")

	if apiData != nil {
		categories := getAPICategories(apiData)
		for i, cat := range categories {
			items = append(items, fmt.Sprintf("   4.%d %s (%d endpoints)", i+1, cat, countAPIsInCategory(apiData, cat)))
		}
	}

	items = append(items, "5. Security Audit Results")

	pdf.SetFont("Helvetica", "", 11)
	for _, item := range items {
		pdf.CellFormat(0, 8, item, "", 1, "L", false, 0, "")
	}
}

// ============== CATEGORY PAGES ==============

func drawCategoryPage(pdf *gofpdf.Fpdf, cat CategoryResult) {
	pdf.AddPage()

	// Category header
	pdf.SetFillColor(30, 58, 95)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 12, fmt.Sprintf("  %s", cat.Name), "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Summary bar
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(60, 60, 60)
	summary := fmt.Sprintf("Total: %d  |  Passed: %d  |  Failed: %d  |  Duration: %s",
		len(cat.Tests), cat.Passed, cat.Failed, formatDuration(cat.Duration))
	pdf.CellFormat(0, 6, summary, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Table header
	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(10, 7, " #", "1", 0, "C", true, 0, "")
	pdf.CellFormat(100, 7, " Test Name", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 7, "Status", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 7, "Duration", "1", 1, "C", true, 0, "")

	// Test rows
	for i, test := range cat.Tests {
		if pdf.GetY() > 260 {
			pdf.AddPage()
			pdf.SetFillColor(240, 240, 240)
			pdf.SetFont("Helvetica", "B", 8)
			pdf.CellFormat(10, 7, " #", "1", 0, "C", true, 0, "")
			pdf.CellFormat(100, 7, " Test Name", "1", 0, "L", true, 0, "")
			pdf.CellFormat(25, 7, "Status", "1", 0, "C", true, 0, "")
			pdf.CellFormat(35, 7, "Duration", "1", 1, "C", true, 0, "")
		}

		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Status color
		switch test.Status {
		case "pass":
			pdf.SetTextColor(34, 197, 94)
		case "fail":
			pdf.SetTextColor(239, 68, 68)
		case "skip":
			pdf.SetTextColor(245, 158, 11)
		}

		pdf.SetFont("Helvetica", "", 8)
		pdf.CellFormat(10, 6, fmt.Sprintf(" %d", i+1), "1", 0, "C", true, 0, "")

		pdf.SetTextColor(60, 60, 60)
		name := test.Name
		if len(name) > 55 {
			name = name[:52] + "..."
		}
		pdf.CellFormat(100, 6, " "+name, "1", 0, "L", true, 0, "")

		switch test.Status {
		case "pass":
			pdf.SetTextColor(34, 197, 94)
		case "fail":
			pdf.SetTextColor(239, 68, 68)
		case "skip":
			pdf.SetTextColor(245, 158, 11)
		}
		pdf.SetFont("Helvetica", "B", 8)
		pdf.CellFormat(25, 6, strings.ToUpper(test.Status), "1", 0, "C", true, 0, "")

		pdf.SetTextColor(100, 100, 100)
		pdf.SetFont("Helvetica", "", 8)
		pdf.CellFormat(35, 6, formatDuration(test.Duration), "1", 1, "C", true, 0, "")
	}

	// Failed test details
	if cat.Failed > 0 {
		pdf.Ln(5)
		pdf.SetTextColor(239, 68, 68)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(0, 8, "Failed Tests - Debug Output", "", 1, "L", false, 0, "")

		for _, test := range cat.Tests {
			if test.Status != "fail" {
				continue
			}
			pdf.SetTextColor(60, 60, 60)
			pdf.SetFont("Helvetica", "B", 9)
			pdf.CellFormat(0, 6, test.Name, "", 1, "L", false, 0, "")

			pdf.SetFont("Courier", "", 7)
			pdf.SetTextColor(100, 100, 100)
			output := cleanTestOutput(test.Output)
			if output != "" {
				pdf.MultiCell(170, 4, output, "", "", false)
			}
			pdf.Ln(2)
		}
	}
}

// ============== API DOCUMENTATION PAGES ==============

func drawAPIDocumentation(pdf *gofpdf.Fpdf, apiData *APIData) {
	// Section title page
	pdf.AddPage()
	pdf.SetFillColor(30, 58, 95)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetY(40)
	pdf.CellFormat(0, 12, "API Endpoint Documentation", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(200, 200, 200)
	pdf.CellFormat(0, 8, "Complete request/response capture from test execution", "", 1, "C", false, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 7, fmt.Sprintf("%d API interactions recorded", len(apiData.APIRecords)), "", 1, "C", false, 0, "")

	// Group by category
	categories := getAPICategories(apiData)

	for _, cat := range categories {
		records := getAPIsByCategory(apiData, cat)
		for _, rec := range records {
			drawAPIEndpoint(pdf, rec)
		}
	}
}

func drawAPIEndpoint(pdf *gofpdf.Fpdf, rec APIRecord) {
	pdf.AddPage()

	// Endpoint header
	pdf.SetFillColor(methodColor(rec.Method))
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 10)

	methodW := 22.0
	pdf.CellFormat(methodW, 8, " "+rec.Method, "1", 0, "L", true, 0, "")

	pdf.SetFillColor(245, 245, 245)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetFont("Courier", "B", 9)
	pdf.CellFormat(148, 8, " "+rec.Path, "1", 1, "L", true, 0, "")

	pdf.Ln(2)

	// Test info
	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("Test: %s  |  Category: %s  |  Status: %d %s",
		rec.TestName, rec.Category, rec.StatusCode, statusText(rec.StatusCode)), "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Request Body
	if rec.RequestBody != "" {
		pdf.SetTextColor(30, 58, 95)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(0, 7, "Request Body", "", 1, "L", false, 0, "")

		pdf.SetFillColor(248, 248, 248)
		pdf.SetFont("Courier", "", 7)
		pdf.SetTextColor(60, 60, 60)
		body := truncate(rec.RequestBody, 800)
		pdf.MultiCell(170, 4, body, "1", "", false)
		pdf.Ln(3)
	}

	// Response Status
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(0, 7, fmt.Sprintf("Response  [%d %s]", rec.StatusCode, statusText(rec.StatusCode)), "", 1, "L", false, 0, "")

	// Response Headers (selected)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 100, 100)
	if ct, ok := rec.ResponseHeader["Content-Type"]; ok {
		pdf.CellFormat(0, 5, fmt.Sprintf("Content-Type: %s", ct), "", 1, "L", false, 0, "")
	}
	pdf.Ln(1)

	// Response Body
	pdf.SetFillColor(248, 248, 248)
	pdf.SetFont("Courier", "", 7)
	pdf.SetTextColor(60, 60, 60)
	respBody := rec.ResponseBody
	if respBody == "" {
		respBody = "(empty response)"
	}

	// Try to pretty-print JSON
	var prettyJSON map[string]interface{}
	if err := json.Unmarshal([]byte(respBody), &prettyJSON); err == nil {
		if b, err := json.MarshalIndent(prettyJSON, "", "  "); err == nil {
			respBody = string(b)
		}
	}

	respBody = truncate(respBody, 2000)
	pdf.MultiCell(170, 4, respBody, "1", "", false)
}

// ============== SECURITY AUDIT PAGE ==============

func drawSecurityAuditPage(pdf *gofpdf.Fpdf, report ReportData, apiData *APIData) {
	pdf.AddPage()

	pdf.SetFillColor(139, 0, 0)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 12, "  Security Audit Results", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// JWT Security
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "JWT Security Tests", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	secTests := []struct {
		name   string
		desc   string
		status string
	}{
		{"Expired Token", "Token past expiration date is rejected", findStatus(report, "TestJWT_ExpiredToken")},
		{"Invalid Signature", "Token signed with wrong key is rejected", findStatus(report, "TestJWT_InvalidSignature")},
		{"Algorithm None Attack", "JWT alg:none attack is blocked", findStatus(report, "TestJWT_NoneAlgorithm")},
		{"Malformed Token", "Invalid JWT format is rejected", findStatus(report, "TestJWT_MalformedToken")},
		{"Role Escalation", "Forged role with wrong key is rejected", findStatus(report, "TestJWT_RoleEscalation_ForgedRole")},
		{"Invalid Role", "Non-existent role value is rejected", findStatus(report, "TestJWT_InvalidRole")},
	}

	drawSecurityTable(pdf, secTests)
	pdf.Ln(5)

	// RBAC
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "Role-Based Access Control (RBAC)", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	rbacTests := []struct {
		name   string
		desc   string
		status string
	}{
		{"Kepala KUA Endpoints", "user_biasa blocked from kepala_kua routes", findStatus(report, "TestRBAC_KepalaKUAOnlyEndpoints")},
		{"Staff Endpoints", "user_biasa blocked from staff routes", findStatus(report, "TestRBAC_StaffEndpoints")},
		{"Multi-Role Endpoints", "user_biasa blocked from staff+kepala routes", findStatus(report, "TestRBAC_StaffAndKepalaKUAEndpoints")},
		{"Penghulu Endpoints", "staff blocked from penghulu routes", findStatus(report, "TestRBAC_PenghuluEndpoint_RejectsOtherRoles")},
		{"No Token Rejection", "All authenticated endpoints reject missing token", findStatus(report, "TestRBAC_AuthenticatedEndpointsRejectNoToken")},
	}

	drawSecurityTable(pdf, rbacTests)
	pdf.Ln(5)

	// IDOR
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "IDOR Prevention", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	idorTests := []struct {
		name   string
		desc   string
		status string
	}{
		{"Notification Read", "User B cannot read User A notification", findStatus(report, "TestIDOR_CannotReadOtherUsersNotification")},
		{"Notification Update", "User B cannot update User A notification", findStatus(report, "TestIDOR_CannotUpdateOtherUsersNotification")},
		{"Notification Delete", "User B cannot delete User A notification", findStatus(report, "TestIDOR_CannotDeleteOtherUsersNotification")},
		{"Notification List", "User B sees only own notifications", findStatus(report, "TestIDOR_CannotSeeOtherUsersNotifications")},
		{"Mark All Read Scope", "Mark-all only affects own notifications", findStatus(report, "TestIDOR_MarkAllRead_OnlyOwnNotifications")},
		{"Location Update", "Non-owner blocked from updating location", findStatus(report, "TestIDOR_UpdateLocation_OtherUserBlocked")},
		{"Location View", "Non-owner blocked from viewing location", findStatus(report, "TestIDOR_GetLocationDetail_OtherUserBlocked")},
	}

	drawSecurityTable(pdf, idorTests)
	pdf.Ln(5)

	// Injection
	if pdf.GetY() > 220 {
		pdf.AddPage()
	}
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "Injection & Input Validation", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	injTests := []struct {
		name   string
		desc   string
		status string
	}{
		{"SQL Injection", "Malicious SQL in registration ID blocked", findStatus(report, "TestInjection_SQLInjection_RegistrationID")},
		{"XSS Prevention", "Content-Type header prevents XSS sniffing", findStatus(report, "TestInjection_XSSInNotification")},
		{"HTTP Method Abuse", "Wrong HTTP methods rejected", findStatus(report, "TestHTTPMethods_NoUnexpectedMethods")},
	}

	drawSecurityTable(pdf, injTests)

	// API Response examples from captured data
	if apiData != nil {
		drawSecurityAPIExamples(pdf, apiData)
	}
}

func drawSecurityTable(pdf *gofpdf.Fpdf, tests []struct {
	name   string
	desc   string
	status string
}) {
	// Table header
	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(50, 6, " Test", "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 6, " Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 6, "Result", "1", 1, "C", true, 0, "")

	for i, test := range tests {
		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(60, 60, 60)
		pdf.CellFormat(50, 5, " "+test.name, "1", 0, "L", true, 0, "")
		pdf.CellFormat(95, 5, " "+test.desc, "1", 0, "L", true, 0, "")

		if test.status == "pass" {
			pdf.SetFillColor(34, 197, 94)
			pdf.SetTextColor(255, 255, 255)
		} else if test.status == "fail" {
			pdf.SetFillColor(239, 68, 68)
			pdf.SetTextColor(255, 255, 255)
		} else {
			pdf.SetFillColor(200, 200, 200)
			pdf.SetTextColor(100, 100, 100)
		}
		pdf.SetFont("Helvetica", "B", 8)
		pdf.CellFormat(25, 5, strings.ToUpper(test.status), "1", 1, "C", true, 0, "")
	}
}

func drawSecurityAPIExamples(pdf *gofpdf.Fpdf, apiData *APIData) {
	pdf.AddPage()
	pdf.SetTextColor(30, 58, 95)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 10, "Security API Response Examples", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Find security-related API records
	securityRecords := []APIRecord{}
	for _, rec := range apiData.APIRecords {
		if rec.Category == "Security" || rec.Category == "Authentication" {
			securityRecords = append(securityRecords, rec)
		}
	}

	// Show up to 6 examples
	limit := 6
	if len(securityRecords) < limit {
		limit = len(securityRecords)
	}

	for i := 0; i < limit; i++ {
		rec := securityRecords[i]
		if pdf.GetY() > 230 {
			pdf.AddPage()
		}

		// Header
		pdf.SetFillColor(methodColor(rec.Method))
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 8)
		pdf.CellFormat(18, 6, " "+rec.Method, "1", 0, "L", true, 0, "")

		pdf.SetFillColor(245, 245, 245)
		pdf.SetTextColor(60, 60, 60)
		pdf.SetFont("Courier", "", 8)
		pdf.CellFormat(122, 6, " "+rec.Path, "1", 0, "L", true, 0, "")

		pdf.SetFillColor(statusColor(rec.StatusCode))
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 8)
		pdf.CellFormat(30, 6, fmt.Sprintf(" %d", rec.StatusCode), "1", 1, "C", true, 0, "")

		// Request
		if rec.RequestBody != "" {
			pdf.SetFont("Helvetica", "B", 7)
			pdf.SetTextColor(100, 100, 100)
			pdf.CellFormat(0, 4, "Request:", "", 1, "L", false, 0, "")
			pdf.SetFont("Courier", "", 6)
			pdf.SetFillColor(248, 248, 248)
			pdf.MultiCell(170, 3, truncate(rec.RequestBody, 300), "1", "", false)
		}

		// Response
		pdf.SetFont("Helvetica", "B", 7)
		pdf.SetTextColor(100, 100, 100)
		pdf.CellFormat(0, 4, "Response:", "", 1, "L", false, 0, "")
		pdf.SetFont("Courier", "", 6)
		pdf.SetFillColor(248, 248, 248)

		respBody := rec.ResponseBody
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal([]byte(respBody), &prettyJSON); err == nil {
			if b, err := json.MarshalIndent(prettyJSON, "", "  "); err == nil {
				respBody = string(b)
			}
		}
		pdf.MultiCell(170, 3, truncate(respBody, 500), "1", "", false)
		pdf.Ln(4)
	}
}

// ============== HELPERS ==============

func methodColor(method string) (int, int, int) {
	switch method {
	case "GET":
		return 34, 197, 94
	case "POST":
		return 59, 130, 246
	case "PUT":
		return 245, 158, 11
	case "DELETE":
		return 239, 68, 68
	default:
		return 100, 100, 100
	}
}

func statusColor(code int) (int, int, int) {
	switch {
	case code >= 200 && code < 300:
		return 34, 197, 94
	case code >= 300 && code < 400:
		return 59, 130, 246
	case code >= 400 && code < 500:
		return 245, 158, 11
	default:
		return 239, 68, 68
	}
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 409:
		return "Conflict"
	case 500:
		return "Internal Server Error"
	default:
		return fmt.Sprintf("Status %d", code)
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return "<1ms"
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... (truncated)"
}

func cleanTestOutput(lines []string) string {
	var cleaned []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "=== RUN") || strings.HasPrefix(l, "--- PASS") || strings.HasPrefix(l, "--- FAIL") {
			continue
		}
		cleaned = append(cleaned, l)
	}
	return strings.Join(cleaned, "\n")
}

func findStatus(report ReportData, prefix string) string {
	for _, t := range report.AllTests {
		if strings.HasPrefix(t.Name, prefix) {
			return t.Status
		}
	}
	return "n/a"
}

func getAPICategories(apiData *APIData) []string {
	seen := make(map[string]bool)
	var cats []string
	for _, rec := range apiData.APIRecords {
		cat := rec.Category
		if cat == "" {
			cat = "General"
		}
		if !seen[cat] {
			seen[cat] = true
			cats = append(cats, cat)
		}
	}
	return cats
}

func getAPIsByCategory(apiData *APIData, category string) []APIRecord {
	var records []APIRecord
	for _, rec := range apiData.APIRecords {
		cat := rec.Category
		if cat == "" {
			cat = "General"
		}
		if cat == category {
			records = append(records, rec)
		}
	}
	return records
}

func countAPIsInCategory(apiData *APIData, category string) int {
	return len(getAPIsByCategory(apiData, category))
}
