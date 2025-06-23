// main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// Colors
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

type Violation struct {
	Rule        string `json:"rule"`
	Message     string `json:"message"`
	Line        int    `json:"line"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type FileResult struct {
	Filename   string      `json:"filename"`
	Violations []Violation `json:"violations"`
	Score      float64     `json:"score"`
	LineCount  int         `json:"line_count"`
}

type Report struct {
	Files            []FileResult `json:"files"`
	TotalScore       float64      `json:"total_score"`
	TotalFiles       int          `json:"total_files"`
	TotalLines       int          `json:"total_lines"`
	TotalViolations  int          `json:"total_violations"`
	CleanFiles       int          `json:"clean_files"`
}

func main() {
	var (
		pathFlag    = flag.String("path", "", "Path to file or directory to analyze")
		verboseFlag = flag.Bool("verbose", false, "Verbose output")
		jsonFlag    = flag.Bool("json", false, "JSON output format")
		silentFlag  = flag.Bool("silent", false, "Silent mode (exit code only)")
		levelFlag   = flag.Int("level", 1, "Verification level (1=basic, 2=advanced)")
	)
	flag.Parse()

	// Get path from flag or argument
	path := *pathFlag
	if path == "" && len(flag.Args()) > 0 {
		path = flag.Args()[0]
	}

	if path == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file_or_directory>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	analyzer := NewAnalyzer(*levelFlag)
	report, err := analyzer.AnalyzePath(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *silentFlag {
		if report.TotalViolations > 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *jsonFlag {
		output, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(output))
	} else {
		printReport(report, *verboseFlag)
	}

	if report.TotalViolations > 0 {
		os.Exit(1)
	}
}

type Analyzer struct {
	level int
	rules map[string]Rule
}

type Rule struct {
	Code        string
	Name        string
	Description string
	Severity    string
	Level       int
	Check       func(*FileAnalysis, string, int) []Violation
}

type FileAnalysis struct {
	Filename  string
	Lines     []string
	Functions []FunctionInfo
}

type FunctionInfo struct {
	Name      string
	StartLine int
	EndLine   int
	ParamCount int
}

func NewAnalyzer(level int) *Analyzer {
	a := &Analyzer{
		level: level,
		rules: make(map[string]Rule),
	}
	a.initRules()
	return a
}

func (a *Analyzer) initRules() {
	// Level 1 rules (basic)
	a.rules["C-L1"] = Rule{
		Code: "C-L1", Name: "Line Length", Description: "Line too long (80 chars max)",
		Severity: "major", Level: 1, Check: checkLineLength,
	}
	a.rules["C-L2"] = Rule{
		Code: "C-L2", Name: "Empty Lines", Description: "Forbidden empty lines",
		Severity: "minor", Level: 1, Check: checkEmptyLines,
	}
	a.rules["C-L3"] = Rule{
		Code: "C-L3", Name: "Indentation", Description: "TAB indentation only",
		Severity: "major", Level: 1, Check: checkIndentation,
	}
	a.rules["C-L4"] = Rule{
		Code: "C-L4", Name: "Variable Declaration", Description: "One variable per line",
		Severity: "major", Level: 1, Check: checkVariableDeclaration,
	}
	a.rules["C-V1"] = Rule{
		Code: "C-V1", Name: "Variable Position", Description: "Variables at function start",
		Severity: "major", Level: 1, Check: checkVariablePosition,
	}
	a.rules["C-O1"] = Rule{
		Code: "C-O1", Name: "Filename", Description: "Filename in snake_case",
		Severity: "major", Level: 1, Check: checkFilename,
	}
	a.rules["C-O2"] = Rule{
		Code: "C-O2", Name: "Function Count", Description: "Max 3 functions per file",
		Severity: "major", Level: 1, Check: checkFunctionCount,
	}
	a.rules["C-F1"] = Rule{
		Code: "C-F1", Name: "Function Name", Description: "Function name in snake_case",
		Severity: "major", Level: 1, Check: checkFunctionNames,
	}
	a.rules["C-F2"] = Rule{
		Code: "C-F2", Name: "Macro Name", Description: "Macro in SCREAMING_SNAKE_CASE",
		Severity: "major", Level: 1, Check: checkMacroNames,
	}
	a.rules["C-F3"] = Rule{
		Code: "C-F3", Name: "Function Length", Description: "Function max 25 lines",
		Severity: "major", Level: 1, Check: checkFunctionLength,
	}

	// Level 2 rules (advanced)
	if a.level >= 2 {
		a.rules["C-C1"] = Rule{
			Code: "C-C1", Name: "Comment Format", Description: "/* */ comments only",
			Severity: "minor", Level: 2, Check: checkCommentFormat,
		}
		a.rules["C-C2"] = Rule{
			Code: "C-C2", Name: "Function Comment", Description: "Function comment required",
			Severity: "minor", Level: 2, Check: checkFunctionComment,
		}
		a.rules["C-G1"] = Rule{
			Code: "C-G1", Name: "Global Variables", Description: "No non-const globals",
			Severity: "major", Level: 2, Check: checkGlobalVariables,
		}
		a.rules["C-F4"] = Rule{
			Code: "C-F4", Name: "Function Parameters", Description: "Max 4 parameters",
			Severity: "major", Level: 2, Check: checkFunctionParameters,
		}
		a.rules["C-L5"] = Rule{
			Code: "C-L5", Name: "For Loop Declaration", Description: "No declaration in for loops",
			Severity: "major", Level: 2, Check: checkForLoopDeclaration,
		}
	}
}

func (a *Analyzer) AnalyzePath(path string) (*Report, error) {
	var files []string
	
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(p, ".c") || strings.HasSuffix(p, ".h") {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(path, ".c") || strings.HasSuffix(path, ".h") {
		files = append(files, path)
	}

	report := &Report{
		Files: make([]FileResult, 0, len(files)),
	}

	for _, file := range files {
		result, err := a.analyzeFile(file)
		if err != nil {
			continue
		}
		report.Files = append(report.Files, *result)
		report.TotalFiles++
		report.TotalLines += result.LineCount
		report.TotalViolations += len(result.Violations)
		if len(result.Violations) == 0 {
			report.CleanFiles++
		}
	}

	// Calculate total score
	if report.TotalFiles > 0 {
		totalScore := 0.0
		for _, file := range report.Files {
			totalScore += file.Score
		}
		report.TotalScore = totalScore / float64(report.TotalFiles)
	}

	return report, nil
}

func (a *Analyzer) analyzeFile(filename string) (*FileResult, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	analysis := &FileAnalysis{
		Filename: filename,
		Lines:    lines,
		Functions: extractFunctions(lines),
	}

	var violations []Violation
	for _, rule := range a.rules {
		if rule.Level <= a.level {
			ruleViolations := rule.Check(analysis, filename, 0)
			violations = append(violations, ruleViolations...)
		}
	}

	// Calculate score (100 - penalty per violation)
	score := 100.0
	for _, v := range violations {
		penalty := 5.0 // major violations
		if v.Severity == "minor" {
			penalty = 2.0
		}
		score -= penalty
	}
	if score < 0 {
		score = 0
	}

	return &FileResult{
		Filename:   filepath.Base(filename),
		Violations: violations,
		Score:      score,
		LineCount:  len(lines),
	}, nil
}

// Rule checking functions
func checkLineLength(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		if len(line) > 80 {
			violations = append(violations, Violation{
				Rule:        "C-L1",
				Message:     "Line too long",
				Line:        i + 1,
				Severity:    "major",
				Description: fmt.Sprintf("Line contains %d characters (max 80)", len(line)),
			})
		}
	}
	return violations
}

func checkEmptyLines(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	lines := analysis.Lines
	
	// Check first line
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		violations = append(violations, Violation{
			Rule:        "C-L2",
			Message:     "Empty line at beginning of file",
			Line:        1,
			Severity:    "minor",
			Description: "File should not start with empty line",
		})
	}
	
	// Check last line
	if len(lines) > 1 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		violations = append(violations, Violation{
			Rule:        "C-L2",
			Message:     "Empty line at end of file",
			Line:        len(lines),
			Severity:    "minor",
			Description: "File should not end with empty line",
		})
	}
	
	// Check consecutive empty lines
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" && strings.TrimSpace(lines[i-1]) == "" {
			violations = append(violations, Violation{
				Rule:        "C-L2",
				Message:     "Consecutive empty lines",
				Line:        i + 1,
				Severity:    "minor",
				Description: "Multiple consecutive empty lines are forbidden",
			})
		}
	}
	
	return violations
}

func checkIndentation(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		if len(line) > 0 && line[0] == ' ' {
			violations = append(violations, Violation{
				Rule:        "C-L3",
				Message:     "Space indentation",
				Line:        i + 1,
				Severity:    "major",
				Description: "Use TAB for indentation, not spaces",
			})
		}
	}
	return violations
}

func checkVariableDeclaration(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		trimmed := strings.TrimSpace(line)
		// Simple check for multiple variable declarations
		if strings.Contains(trimmed, "int ") || strings.Contains(trimmed, "char ") || 
		   strings.Contains(trimmed, "float ") || strings.Contains(trimmed, "double ") {
			if strings.Count(trimmed, ",") > 0 && !strings.Contains(trimmed, "for") {
				violations = append(violations, Violation{
					Rule:        "C-L4",
					Message:     "Multiple variable declaration",
					Line:        i + 1,
					Severity:    "major",
					Description: "Declare only one variable per line",
				})
			}
		}
	}
	return violations
}

func checkVariablePosition(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	// This is a simplified check - would need proper C parsing for accuracy
	return []Violation{}
}

func checkFilename(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	
	if !isSnakeCase(name) {
		violations = append(violations, Violation{
			Rule:        "C-O1",
			Message:     "Invalid filename format",
			Line:        0,
			Severity:    "major",
			Description: "Filename must be in snake_case",
		})
	}
	return violations
}

func checkFunctionCount(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	funcCount := 0
	
	for _, line := range analysis.Lines {
		trimmed := strings.TrimSpace(line)
		// Simple function detection
		if strings.Contains(trimmed, "(") && strings.Contains(trimmed, ")") && 
		   strings.Contains(trimmed, "{") && !strings.HasPrefix(trimmed, "//") &&
		   !strings.HasPrefix(trimmed, "/*") && !strings.Contains(trimmed, "if") &&
		   !strings.Contains(trimmed, "while") && !strings.Contains(trimmed, "for") {
			if !strings.Contains(trimmed, "main") {
				funcCount++
			}
		}
	}
	
	if funcCount > 3 {
		violations = append(violations, Violation{
			Rule:        "C-O2",
			Message:     "Too many functions",
			Line:        0,
			Severity:    "major",
			Description: fmt.Sprintf("File contains %d functions (max 3 excluding main)", funcCount),
		})
	}
	return violations
}

func checkFunctionNames(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for _, fn := range analysis.Functions {
		if !isSnakeCase(fn.Name) && fn.Name != "main" {
			violations = append(violations, Violation{
				Rule:        "C-F1",
				Message:     "Invalid function name",
				Line:        fn.StartLine,
				Severity:    "major",
				Description: fmt.Sprintf("Function '%s' must be in snake_case", fn.Name),
			})
		}
	}
	return violations
}

func checkMacroNames(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#define ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				macroName := parts[1]
				if !isScreamingSnakeCase(macroName) {
					violations = append(violations, Violation{
						Rule:        "C-F2",
						Message:     "Invalid macro name",
						Line:        i + 1,
						Severity:    "major",
						Description: fmt.Sprintf("Macro '%s' must be in SCREAMING_SNAKE_CASE", macroName),
					})
				}
			}
		}
	}
	return violations
}

func checkFunctionLength(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for _, fn := range analysis.Functions {
		length := fn.EndLine - fn.StartLine + 1
		if length > 25 {
			violations = append(violations, Violation{
				Rule:        "C-F3",
				Message:     "Function too long",
				Line:        fn.StartLine,
				Severity:    "major",
				Description: fmt.Sprintf("Function '%s' has %d lines (max 25)", fn.Name, length),
			})
		}
	}
	return violations
}

// Level 2 checks
func checkCommentFormat(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		if strings.Contains(line, "//") {
			violations = append(violations, Violation{
				Rule:        "C-C1",
				Message:     "Invalid comment format",
				Line:        i + 1,
				Severity:    "minor",
				Description: "Use /* */ comments only, not // comments",
			})
		}
	}
	return violations
}

func checkFunctionComment(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	// Simplified check - would need better parsing
	return []Violation{}
}

func checkGlobalVariables(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	// Simplified check - would need proper C parsing
	return []Violation{}
}

func checkFunctionParameters(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for _, fn := range analysis.Functions {
		if fn.ParamCount > 4 {
			violations = append(violations, Violation{
				Rule:        "C-F4",
				Message:     "Too many parameters",
				Line:        fn.StartLine,
				Severity:    "major",
				Description: fmt.Sprintf("Function '%s' has %d parameters (max 4)", fn.Name, fn.ParamCount),
			})
		}
	}
	return violations
}

func checkForLoopDeclaration(analysis *FileAnalysis, filename string, lineNum int) []Violation {
	var violations []Violation
	for i, line := range analysis.Lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "for") && strings.Contains(trimmed, "int ") {
			violations = append(violations, Violation{
				Rule:        "C-L5",
				Message:     "Variable declaration in for loop",
				Line:        i + 1,
				Severity:    "major",
				Description: "Do not declare variables in for loop initialization",
			})
		}
	}
	return violations
}

// Helper functions
func isSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			return false
		}
		if r == '_' && (i == 0 || i == len(s)-1) {
			return false
		}
	}
	return true
}

func isScreamingSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r >= 'a' && r <= 'z' {
			return false
		}
		if r == '_' && (i == 0 || i == len(s)-1) {
			return false
		}
	}
	return true
}

func extractFunctions(lines []string) []FunctionInfo {
	var functions []FunctionInfo
	var currentFunc *FunctionInfo
	braceCount := 0
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Simple function detection
		if strings.Contains(trimmed, "(") && strings.Contains(trimmed, ")") && 
		   (strings.Contains(trimmed, "{") || (i+1 < len(lines) && strings.Contains(strings.TrimSpace(lines[i+1]), "{"))) {
			if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") &&
			   !strings.Contains(trimmed, "if") && !strings.Contains(trimmed, "while") &&
			   !strings.Contains(trimmed, "for") && !strings.Contains(trimmed, "switch") {
				
				// Extract function name
				parenPos := strings.Index(trimmed, "(")
				if parenPos > 0 {
					funcPart := trimmed[:parenPos]
					parts := strings.Fields(funcPart)
					if len(parts) > 0 {
						funcName := parts[len(parts)-1]
						if strings.Contains(funcName, "*") {
							funcName = strings.TrimLeft(funcName, "*")
						}
						
						// Count parameters
						paramPart := trimmed[parenPos+1:]
						closeParenPos := strings.Index(paramPart, ")")
						if closeParenPos > 0 {
							params := paramPart[:closeParenPos]
							paramCount := 0
							if strings.TrimSpace(params) != "" && strings.TrimSpace(params) != "void" {
								paramCount = strings.Count(params, ",") + 1
							}
							
							currentFunc = &FunctionInfo{
								Name:       funcName,
								StartLine:  i + 1,
								ParamCount: paramCount,
							}
						}
					}
				}
			}
		}
		
		// Count braces to find function end
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")
		
		if currentFunc != nil && braceCount == 0 && strings.Contains(line, "}") {
			currentFunc.EndLine = i + 1
			functions = append(functions, *currentFunc)
			currentFunc = nil
		}
	}
	
	return functions
}

func printReport(report *Report, verbose bool) {
	// Print header
	fmt.Println(ColorBold + "╔══════════════════════════════════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorBold + "║                         EPICSTYLE - RAPPORT D'ANALYSE                        ║" + ColorReset)
	fmt.Println(ColorBold + "╚══════════════════════════════════════════════════════════════════════════════╝" + ColorReset)
	fmt.Println()

	// Print summary
	fmt.Printf("📊 %sRÉSUMÉ GLOBAL%s\n", ColorBold, ColorReset)
	fmt.Printf("   • Fichiers analysés: %d\n", report.TotalFiles)
	fmt.Printf("   • Lignes de code: %d\n", report.TotalLines)
	fmt.Printf("   • Violations totales: %d\n", report.TotalViolations)
	fmt.Printf("   • Fichiers propres: %d/%d\n", report.CleanFiles, report.TotalFiles)
	
	cleanPercent := 0.0
	if report.TotalFiles > 0 {
		cleanPercent = float64(report.CleanFiles) / float64(report.TotalFiles) * 100
	}
	fmt.Printf("   • Propreté: %.1f%% %s\n", cleanPercent, getProgressBar(cleanPercent))
	fmt.Println()

	// Sort files by score (descending)
	sort.Slice(report.Files, func(i, j int) bool {
		return report.Files[i].Score > report.Files[j].Score
	})

	// Print file results
	for _, file := range report.Files {
		if len(file.Violations) == 0 {
			fmt.Printf("%s✅ %s%s (%.1f%% - %d lignes)\n", 
				ColorGreen, file.Filename, ColorReset, file.Score, file.LineCount)
		} else {
			fmt.Printf("%s❌ %s%s (%.1f%% - %d lignes - %d violations)\n", 
				ColorRed, file.Filename, ColorReset, file.Score, file.LineCount, len(file.Violations))
		}
		
		if verbose && len(file.Violations) > 0 {
			for _, v := range file.Violations {
				severity := ColorYellow + "MINOR" + ColorReset
				if v.Severity == "major" {
					severity = ColorRed + "MAJOR" + ColorReset
				}
				fmt.Printf("    [%s] Line %d: %s - %s\n", severity, v.Line, v.Rule, v.Message)
				if v.Description != "" {
					fmt.Printf("         %s\n", v.Description)
				}
			}
		}
	}
	
	fmt.Println()

	// Print final score
	scoreColor := ColorRed
	scoreMessage := "❌ ÉCHEC! Beaucoup de travail nécessaire."
	if report.TotalScore >= 90 {
		scoreColor = ColorGreen
		scoreMessage = "🎉 EXCELLENT! Code très propre."
	} else if report.TotalScore >= 75 {
		scoreColor = ColorYellow
		scoreMessage = "🎉 TRÈS BIEN! Quelques petits détails à corriger."
	} else if report.TotalScore >= 50 {
		scoreColor = ColorYellow
		scoreMessage = "⚠️  CORRECT! Plusieurs améliorations nécessaires."
	}

	fmt.Println(ColorBold + "╔══════════════════════════════════════════════════════════════════════════════╗" + ColorReset)
	fmt.Printf("║%s                             SCORE GLOBAL: %.1f%%                              %s ║\n", 
		scoreColor, report.TotalScore, ColorReset)
	fmt.Printf("║           %s%.1f%%           ║\n", getProgressBar(report.TotalScore), report.TotalScore)
	fmt.Printf("║                   %s                  ║\n", scoreMessage)
	fmt.Println(ColorBold + "╚══════════════════════════════════════════════════════════════════════════════╝" + ColorReset)
}

func getProgressBar(percentage float64) string {
	barLength := 50
	filled := int(percentage / 100 * float64(barLength))
	empty := barLength - filled
	
	bar := ColorGreen + strings.Repeat("█", filled) + ColorReset + strings.Repeat("░", empty)
	return "[" + bar + "]"
}