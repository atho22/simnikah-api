package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TestEvent represents a single event from `go test -json`
type TestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float64   `json:"Elapsed"`
}

// TestResult holds the aggregated result for a single test
type TestResult struct {
	Name     string
	Package  string
	Status   string // pass, fail, skip
	Duration time.Duration
	Output   []string
}

// CategoryResult groups tests by category
type CategoryResult struct {
	Name     string
	Tests    []TestResult
	Passed   int
	Failed   int
	Skipped  int
	Duration time.Duration
}

// ReportData holds all data for the HTML template
type ReportData struct {
	GeneratedAt    time.Time
	ProjectName    string
	TotalTests     int
	TotalPassed    int
	TotalFailed    int
	TotalSkipped   int
	TotalDuration  time.Duration
	PassRate       float64
	Categories     []CategoryResult
	SecurityTests  []TestResult
	RBACTests      []TestResult
	IDORTests      []TestResult
	AllTests       []TestResult
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║     SIPENA API - Test Report Generator          ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	// Run go test with JSON output
	fmt.Println("[1/3] Running tests...")
	cmd := exec.Command("go", "test", "-v", "-json", "-count=1", "./tests/...")
	cmd.Dir = getProjectRoot()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting tests: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON events
	var events []TestEvent
	scanner := bufio.NewScanner(stdout)
	// Increase buffer size for long output lines
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue // skip non-JSON lines
		}
		events = append(events, event)
	}

	if err := cmd.Wait(); err != nil {
		// Tests may fail, that's OK - we still have the events
		fmt.Println("  ⚠ Some tests failed (this is captured in the report)")
	}

	fmt.Printf("  ✓ Collected %d test events\n", len(events))

	// Aggregate results
	fmt.Println("[2/4] Processing results...")
	report := processEvents(events)

	// Generate HTML report
	fmt.Println("[3/4] Generating HTML report...")
	root := getProjectRoot()
	htmlPath := filepath.Join(root, "test_report.html")
	if err := generateHTMLReport(report, htmlPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HTML report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ HTML report: %s\n", htmlPath)

	// Generate PDF report
	fmt.Println("[4/4] Generating PDF report...")
	pdfPath := filepath.Join(root, "test_report.pdf")
	apiDataPath := filepath.Join(root, "test_api_data.json")
	if err := generatePDFReport(report, apiDataPath, pdfPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating PDF report: %v\n", err)
	} else {
		fmt.Printf("  ✓ PDF report:  %s\n", pdfPath)
	}

	fmt.Printf("\n✅ Reports generated successfully!\n")
	fmt.Printf("   Tests: %d total | %d passed | %d failed | %d skipped\n",
		report.TotalTests, report.TotalPassed, report.TotalFailed, report.TotalSkipped)
	fmt.Printf("   Pass Rate: %.1f%%\n", report.PassRate)
	fmt.Printf("   Duration: %s\n", report.TotalDuration.Round(time.Millisecond))
}

func getProjectRoot() string {
	// Walk up from current directory to find go.mod
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback
	dir, _ = os.Getwd()
	return filepath.Dir(filepath.Dir(dir))
}

func processEvents(events []TestEvent) ReportData {
	// Track test results by name
	testMap := make(map[string]*TestResult)

	for _, e := range events {
		if e.Test == "" {
			continue // package-level event
		}

		key := e.Package + "::" + e.Test
		switch e.Action {
		case "run":
			testMap[key] = &TestResult{
				Name:    e.Test,
				Package: e.Package,
				Status:  "running",
			}
		case "output":
			if tr, ok := testMap[key]; ok {
				tr.Output = append(tr.Output, e.Output)
			}
		case "pass":
			if tr, ok := testMap[key]; ok {
				tr.Status = "pass"
				tr.Duration = time.Duration(e.Elapsed * float64(time.Second))
			}
		case "fail":
			if tr, ok := testMap[key]; ok {
				tr.Status = "fail"
				tr.Duration = time.Duration(e.Elapsed * float64(time.Second))
			}
		case "skip":
			if tr, ok := testMap[key]; ok {
				tr.Status = "skip"
				tr.Duration = time.Duration(e.Elapsed * float64(time.Second))
			}
		}
	}

	// Collect all test results
	var allTests []TestResult
	for _, tr := range testMap {
		if tr.Status == "running" {
			continue // incomplete test
		}
		allTests = append(allTests, *tr)
	}

	// Sort by name for consistency
	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Name < allTests[j].Name
	})

	// Calculate totals
	report := ReportData{
		GeneratedAt: time.Now(),
		ProjectName: "SIPENA - Sistem Informasi Pernikahan",
		AllTests:    allTests,
	}

	for _, tr := range allTests {
		report.TotalTests++
		report.TotalDuration += tr.Duration
		switch tr.Status {
		case "pass":
			report.TotalPassed++
		case "fail":
			report.TotalFailed++
		case "skip":
			report.TotalSkipped++
		}
	}

	if report.TotalTests > 0 {
		report.PassRate = float64(report.TotalPassed) / float64(report.TotalTests) * 100
	}

	// Categorize tests
	report.Categories = categorizeTests(allTests)
	report.SecurityTests = filterTestsByPrefix(allTests, "TestJWT_", "TestInjection_", "TestHTTPMethods_", "Test404_")
	report.RBACTests = filterTestsByPrefix(allTests, "TestRBAC_")
	report.IDORTests = filterTestsByPrefix(allTests, "TestIDOR_")

	return report
}

func categorizeTests(tests []TestResult) []CategoryResult {
	categoryMap := map[string][]TestResult{
		"Authentication":     filterTestsByPrefix(tests, "TestRegister_", "TestLogin_", "TestGetProfile_"),
		"Health Check":       filterTestsByPrefix(tests, "TestHealth"),
		"RBAC":               filterTestsByPrefix(tests, "TestRBAC_"),
		"Notification CRUD":  filterTestsByPrefix(tests, "TestNotification_C", "TestNotification_G", "TestNotification_U", "TestNotification_M", "TestNotification_D", "TestNotification_S", "TestNotification_I", "TestNotification_P", "TestNotification_F"),
		"Notification IDOR":  filterTestsByPrefix(tests, "TestIDOR_Cannot", "TestIDOR_Mark"),
		"Registration (Catin)": filterTestsByPrefix(tests, "TestCatin_"),
		"Staff Operations":   filterTestsByPrefix(tests, "TestStaff_"),
		"Location & IDOR":    filterTestsByPrefix(tests, "TestIDOR_Update", "TestIDOR_Get", "TestLocation_"),
		"Penghulu":           filterTestsByPrefix(tests, "TestPenghulu_"),
		"Kepala KUA":         filterTestsByPrefix(tests, "TestKepalaKUA_"),
		"JWT Security":       filterTestsByPrefix(tests, "TestJWT_"),
		"Injection & XSS":    filterTestsByPrefix(tests, "TestInjection_"),
		"HTTP Methods":       filterTestsByPrefix(tests, "TestHTTPMethods_"),
		"Edge Cases":         filterTestsByPrefix(tests, "TestRegistration_Invalid", "TestNotification_Pagination", "TestNotification_Filter", "Test404_"),
	}

	var categories []CategoryResult
	// Maintain order
	order := []string{
		"Authentication", "Health Check", "RBAC",
		"Notification CRUD", "Notification IDOR",
		"Registration (Catin)", "Staff Operations",
		"Location & IDOR", "Penghulu", "Kepala KUA",
		"JWT Security", "Injection & XSS", "HTTP Methods", "Edge Cases",
	}

	for _, name := range order {
		tests := categoryMap[name]
		if len(tests) == 0 {
			continue
		}
		cat := CategoryResult{Name: name, Tests: tests}
		for _, t := range tests {
			switch t.Status {
			case "pass":
				cat.Passed++
			case "fail":
				cat.Failed++
			case "skip":
				cat.Skipped++
			}
			cat.Duration += t.Duration
		}
		categories = append(categories, cat)
	}

	return categories
}

func filterTestsByPrefix(tests []TestResult, prefixes ...string) []TestResult {
	var result []TestResult
	for _, t := range tests {
		for _, p := range prefixes {
			if strings.HasPrefix(t.Name, p) {
				result = append(result, t)
				break
			}
		}
	}
	return result
}

func generateHTMLReport(data ReportData, outputPath string) error {
	funcMap := template.FuncMap{
		"statusIcon": func(status string) string {
			switch status {
			case "pass":
				return "✅"
			case "fail":
				return "❌"
			case "skip":
				return "⏭️"
			default:
				return "⏳"
			}
		},
		"statusClass": func(status string) string {
			switch status {
			case "pass":
				return "pass"
			case "fail":
				return "fail"
			case "skip":
				return "skip"
			default:
				return "running"
			}
		},
		"formatDuration": func(d time.Duration) string {
			if d < time.Millisecond {
				return "<1ms"
			}
			if d < time.Second {
				return fmt.Sprintf("%dms", d.Milliseconds())
			}
			return fmt.Sprintf("%.2fs", d.Seconds())
		},
		"progressWidth": func(passed, total int) float64 {
			if total == 0 {
				return 0
			}
			return float64(passed) / float64(total) * 100
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"toUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"cleanOutput": func(lines []string) string {
			var cleaned []string
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l == "" {
					continue
				}
				// Skip test framework noise
				if strings.HasPrefix(l, "=== RUN") || strings.HasPrefix(l, "--- PASS") || strings.HasPrefix(l, "--- FAIL") {
					continue
				}
				cleaned = append(cleaned, l)
			}
			return strings.Join(cleaned, "\n")
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("file create error: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SIPENA API Test Report</title>
<style>
  :root {
    --pass: #22c55e;
    --fail: #ef4444;
    --skip: #f59e0b;
    --bg: #0f172a;
    --card: #1e293b;
    --card-hover: #334155;
    --text: #e2e8f0;
    --text-muted: #94a3b8;
    --border: #334155;
    --accent: #3b82f6;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
    background: var(--bg);
    color: var(--text);
    line-height: 1.6;
    padding: 2rem;
  }
  .container { max-width: 1200px; margin: 0 auto; }

  /* Header */
  .header {
    text-align: center;
    margin-bottom: 2rem;
    padding: 2rem;
    background: linear-gradient(135deg, #1e3a5f, #2d1b69);
    border-radius: 16px;
    border: 1px solid var(--border);
  }
  .header h1 { font-size: 2rem; margin-bottom: 0.5rem; }
  .header .subtitle { color: var(--text-muted); font-size: 0.95rem; }
  .header .timestamp { color: var(--text-muted); font-size: 0.85rem; margin-top: 0.5rem; }

  /* Summary Cards */
  .summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .stat-card {
    background: var(--card);
    border-radius: 12px;
    padding: 1.5rem;
    text-align: center;
    border: 1px solid var(--border);
    transition: transform 0.2s;
  }
  .stat-card:hover { transform: translateY(-2px); }
  .stat-card .number { font-size: 2.5rem; font-weight: 700; }
  .stat-card .label { color: var(--text-muted); font-size: 0.85rem; margin-top: 0.25rem; }
  .stat-card.pass .number { color: var(--pass); }
  .stat-card.fail .number { color: var(--fail); }
  .stat-card.skip .number { color: var(--skip); }
  .stat-card.total .number { color: var(--accent); }

  /* Progress Bar */
  .progress-bar {
    height: 12px;
    background: var(--border);
    border-radius: 6px;
    overflow: hidden;
    margin-bottom: 2rem;
  }
  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--pass), #16a34a);
    border-radius: 6px;
    transition: width 0.5s;
  }
  .progress-fill.has-failures {
    background: linear-gradient(90deg, var(--pass), var(--fail));
  }

  /* Category Section */
  .category {
    background: var(--card);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    border: 1px solid var(--border);
    overflow: hidden;
  }
  .category-header {
    padding: 1rem 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border);
    cursor: pointer;
    user-select: none;
  }
  .category-header:hover { background: var(--card-hover); }
  .category-header h3 { font-size: 1.1rem; }
  .category-badges { display: flex; gap: 0.5rem; align-items: center; }
  .badge {
    padding: 0.2rem 0.6rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 600;
  }
  .badge.pass { background: rgba(34,197,94,0.15); color: var(--pass); }
  .badge.fail { background: rgba(239,68,68,0.15); color: var(--fail); }
  .badge.skip { background: rgba(245,158,11,0.15); color: var(--skip); }
  .badge.duration { background: rgba(59,130,246,0.15); color: var(--accent); }

  /* Test Rows */
  .test-list { padding: 0; }
  .test-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1.5rem;
    border-bottom: 1px solid rgba(51,65,85,0.5);
    font-size: 0.9rem;
  }
  .test-row:last-child { border-bottom: none; }
  .test-row:hover { background: var(--card-hover); }
  .test-name { font-family: 'Cascadia Code', 'Fira Code', monospace; font-size: 0.85rem; }
  .test-meta { display: flex; gap: 0.75rem; align-items: center; white-space: nowrap; }
  .test-duration { color: var(--text-muted); font-size: 0.8rem; }

  /* Security Section */
  .security-section {
    background: var(--card);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    border: 1px solid var(--border);
    overflow: hidden;
  }
  .security-header {
    padding: 1rem 1.5rem;
    background: linear-gradient(135deg, rgba(239,68,68,0.1), rgba(245,158,11,0.1));
    border-bottom: 1px solid var(--border);
  }
  .security-header h3 { font-size: 1.1rem; }

  /* Debug Output */
  .debug-section {
    background: var(--card);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    border: 1px solid var(--border);
    overflow: hidden;
  }
  .debug-header {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border);
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
  }
  .debug-header:hover { background: var(--card-hover); }
  .debug-output {
    padding: 1rem 1.5rem;
    font-family: 'Cascadia Code', 'Fira Code', monospace;
    font-size: 0.8rem;
    white-space: pre-wrap;
    color: var(--text-muted);
    max-height: 300px;
    overflow-y: auto;
    display: none;
  }
  .debug-output.show { display: block; }
  .debug-output .error-line { color: var(--fail); }

  /* Footer */
  .footer {
    text-align: center;
    padding: 1.5rem;
    color: var(--text-muted);
    font-size: 0.85rem;
    border-top: 1px solid var(--border);
    margin-top: 2rem;
  }

  /* Collapsible */
  .collapsible-content { display: none; }
  .collapsible-content.show { display: block; }

  @media (max-width: 768px) {
    body { padding: 1rem; }
    .summary { grid-template-columns: repeat(2, 1fr); }
    .test-row { flex-direction: column; align-items: flex-start; gap: 0.5rem; }
  }
</style>
</head>
<body>
<div class="container">

  <div class="header">
    <h1>🧪 SIPENA API Test Report</h1>
    <div class="subtitle">Comprehensive API Testing & Security Audit Report</div>
    <div class="timestamp">Generated: {{.GeneratedAt.Format "Monday, 02 January 2006 - 15:04:05 MST"}}</div>
  </div>

  <!-- Summary -->
  <div class="summary">
    <div class="stat-card total">
      <div class="number">{{.TotalTests}}</div>
      <div class="label">Total Tests</div>
    </div>
    <div class="stat-card pass">
      <div class="number">{{.TotalPassed}}</div>
      <div class="label">Passed</div>
    </div>
    <div class="stat-card fail">
      <div class="number">{{.TotalFailed}}</div>
      <div class="label">Failed</div>
    </div>
    <div class="stat-card skip">
      <div class="number">{{.TotalSkipped}}</div>
      <div class="label">Skipped</div>
    </div>
    <div class="stat-card total">
      <div class="number">{{printf "%.1f%%" .PassRate}}</div>
      <div class="label">Pass Rate</div>
    </div>
    <div class="stat-card total">
      <div class="number">{{formatDuration .TotalDuration}}</div>
      <div class="label">Duration</div>
    </div>
  </div>

  <!-- Progress Bar -->
  <div class="progress-bar">
    <div class="progress-fill {{if gt .TotalFailed 0}}has-failures{{end}}" style="width: {{progressWidth .TotalPassed .TotalTests}}%"></div>
  </div>

  <!-- Security Tests Highlight -->
  {{if or (gt (len .SecurityTests) 0) (gt (len .RBACTests) 0) (gt (len .IDORTests) 0)}}
  <div class="security-section">
    <div class="security-header">
      <h3>🔒 Security Test Summary</h3>
    </div>
    <div class="test-list">
      {{range .SecurityTests}}
      <div class="test-row">
        <span class="test-name">{{statusIcon .Status}} {{.Name}}</span>
        <div class="test-meta">
          <span class="badge {{statusClass .Status}}">{{toUpper .Status}}</span>
          <span class="test-duration">{{formatDuration .Duration}}</span>
        </div>
      </div>
      {{end}}
      {{range .RBACTests}}
      <div class="test-row">
        <span class="test-name">{{statusIcon .Status}} {{.Name}}</span>
        <div class="test-meta">
          <span class="badge {{statusClass .Status}}">{{toUpper .Status}}</span>
          <span class="test-duration">{{formatDuration .Duration}}</span>
        </div>
      </div>
      {{end}}
      {{range .IDORTests}}
      <div class="test-row">
        <span class="test-name">{{statusIcon .Status}} {{.Name}}</span>
        <div class="test-meta">
          <span class="badge {{statusClass .Status}}">{{toUpper .Status}}</span>
          <span class="test-duration">{{formatDuration .Duration}}</span>
        </div>
      </div>
      {{end}}
    </div>
  </div>
  {{end}}

  <!-- Categories -->
  {{range .Categories}}
  <div class="category">
    <div class="category-header" onclick="toggleCategory(this)">
      <h3>{{.Name}} <span style="color:var(--text-muted);font-size:0.85rem">({{len .Tests}} tests)</span></h3>
      <div class="category-badges">
        {{if gt .Passed 0}}<span class="badge pass">{{.Passed}} PASS</span>{{end}}
        {{if gt .Failed 0}}<span class="badge fail">{{.Failed}} FAIL</span>{{end}}
        {{if gt .Skipped 0}}<span class="badge skip">{{.Skipped}} SKIP</span>{{end}}
        <span class="badge duration">{{formatDuration .Duration}}</span>
      </div>
    </div>
    <div class="test-list collapsible-content show">
      {{range .Tests}}
      <div class="test-row">
        <span class="test-name">{{statusIcon .Status}} {{.Name}}</span>
        <div class="test-meta">
          <span class="badge {{statusClass .Status}}">{{toUpper .Status}}</span>
          <span class="test-duration">{{formatDuration .Duration}}</span>
          {{if eq .Status "fail"}}<button onclick="toggleDebug(event, 'debug-{{.Name}}')" style="background:var(--fail);color:white;border:none;padding:2px 8px;border-radius:4px;cursor:pointer;font-size:0.75rem">DEBUG</button>{{end}}
        </div>
      </div>
      {{if eq .Status "fail"}}
      <div class="debug-output" id="debug-{{.Name}}">{{cleanOutput .Output}}</div>
      {{end}}
      {{end}}
    </div>
  </div>
  {{end}}

  <!-- Debug Output for Failed Tests -->
  {{if gt .TotalFailed 0}}
  <div class="debug-section">
    <div class="debug-header" onclick="toggleDebugSection(this)">
      <h3>🐛 Failed Tests - Debug Output</h3>
      <span style="color:var(--text-muted);font-size:0.85rem">Click to expand</span>
    </div>
    <div class="debug-output" id="debug-section">
      {{range .AllTests}}
      {{if eq .Status "fail"}}
═══════════════════════════════════════
FAILED: {{.Name}}
Package: {{.Package}}
Duration: {{formatDuration .Duration}}
═══════════════════════════════════════
{{cleanOutput .Output}}

      {{end}}
      {{end}}
    </div>
  </div>
  {{end}}

  <div class="footer">
    SIPENA API Test Report &bull; Generated automatically by Go Test Framework<br>
    Project: {{.ProjectName}}<br>
    Tests run with in-memory SQLite database &bull; No external dependencies
  </div>
</div>

<script>
function toggleCategory(header) {
  const content = header.nextElementSibling;
  content.classList.toggle('show');
}

function toggleDebug(event, id) {
  event.stopPropagation();
  const el = document.getElementById(id);
  if (el) el.classList.toggle('show');
}

function toggleDebugSection(header) {
  const content = header.nextElementSibling;
  content.classList.toggle('show');
}
</script>
</body>
</html>
`
