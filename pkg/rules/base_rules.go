package rules

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// LineLengthRule vérifie la longueur des lignes (max 80 caractères)
type LineLengthRule struct{}

func (r *LineLengthRule) Name() string        { return "C-L1" }
func (r *LineLengthRule) Description() string { return "Une ligne ne doit pas dépasser 80 caractères" }
func (r *LineLengthRule) Level() int          { return 1 }

func (r *LineLengthRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	for i, line := range ctx.Lines {
		if len(line) > 80 {
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Ligne trop longue",
				Line:     i + 1,
				Severity: "major",
				Description: "La ligne contient plus de 80 caractères",
			})
		}
	}
	
	return violations
}

// EmptyLinesRule vérifie les lignes vides en début/fin et consécutives
type EmptyLinesRule struct{}

func (r *EmptyLinesRule) Name() string        { return "C-L2" }
func (r *EmptyLinesRule) Description() string { return "Pas de lignes vides en début/fin de fichier ni consécutives" }
func (r *EmptyLinesRule) Level() int          { return 1 }

func (r *EmptyLinesRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	lines := ctx.Lines
	
	if len(lines) == 0 {
		return violations
	}
	
	// Ligne vide en début
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		violations = append(violations, Violation{
			Rule:     r.Name(),
			Message:  "Ligne vide en début de fichier",
			Line:     1,
			Severity: "major",
		})
	}
	
	// Ligne vide en fin
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		violations = append(violations, Violation{
			Rule:     r.Name(),
			Message:  "Ligne vide en fin de fichier",
			Line:     len(lines),
			Severity: "major",
		})
	}
	
	// Lignes vides consécutives
	for i := 0; i < len(lines)-1; i++ {
		if strings.TrimSpace(lines[i]) == "" && strings.TrimSpace(lines[i+1]) == "" {
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Lignes vides consécutives",
				Line:     i + 2,
				Severity: "major",
			})
		}
	}
	
	return violations
}

// IndentationRule vérifie l'indentation en TAB uniquement
type IndentationRule struct{}

func (r *IndentationRule) Name() string        { return "C-L3" }
func (r *IndentationRule) Description() string { return "Indentation en TAB uniquement" }
func (r *IndentationRule) Level() int          { return 1 }

func (r *IndentationRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	for i, line := range ctx.Lines {
		if strings.Contains(line, "    ") { // 4 espaces
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Utilisation d'espaces au lieu de tabulations",
				Line:     i + 1,
				Severity: "minor",
			})
		}
	}
	
	return violations
}

// VariableDeclarationRule vérifie une variable par ligne
type VariableDeclarationRule struct{}

func (r *VariableDeclarationRule) Name() string        { return "C-L4" }
func (r *VariableDeclarationRule) Description() string { return "Une seule déclaration de variable par ligne" }
func (r *VariableDeclarationRule) Level() int          { return 1 }

func (r *VariableDeclarationRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	// Regex pour détecter les déclarations multiples
	multiDeclRegex := regexp.MustCompile(`^\s*(int|char|float|double|long|short|unsigned)\s+\w+\s*,\s*\w+`)
	
	for i, line := range ctx.Lines {
		if multiDeclRegex.MatchString(line) {
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Plusieurs variables déclarées sur une ligne",
				Line:     i + 1,
				Severity: "major",
			})
		}
	}
	
	return violations
}

// FilenameRule vérifie le nom de fichier en snake_case
type FilenameRule struct{}

func (r *FilenameRule) Name() string        { return "C-O1" }
func (r *FilenameRule) Description() string { return "Nom de fichier en snake_case" }
func (r *FilenameRule) Level() int          { return 1 }

func (r *FilenameRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	filename := filepath.Base(ctx.Filename)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	
	if !isSnakeCase(nameWithoutExt) {
		violations = append(violations, Violation{
			Rule:     r.Name(),
			Message:  "Nom de fichier non conforme au snake_case",
			Line:     1,
			Severity: "major",
			Description: "Le nom de fichier doit être en snake_case (ex: mon_fichier.c)",
		})
	}
	
	return violations
}

// FunctionNamingRule vérifie les noms de fonction en snake_case
type FunctionNamingRule struct{}

func (r *FunctionNamingRule) Name() string        { return "C-F1" }
func (r *FunctionNamingRule) Description() string { return "Nom de fonction en snake_case" }
func (r *FunctionNamingRule) Level() int          { return 1 }

func (r *FunctionNamingRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	// Regex pour les déclarations de fonction
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\(`)
	
	for i, line := range ctx.Lines {
		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			funcName := matches[1]
			if funcName != "main" && !isSnakeCase(funcName) {
				violations = append(violations, Violation{
					Rule:     r.Name(),
					Message:  "Nom de fonction non conforme au snake_case",
					Line:     i + 1,
					Severity: "major",
					Description: "Le nom de fonction '" + funcName + "' doit être en snake_case",
				})
			}
		}
	}
	
	return violations
}

// MacroNamingRule vérifie les noms de macro en SCREAMING_SNAKE_CASE
type MacroNamingRule struct{}

func (r *MacroNamingRule) Name() string        { return "C-F2" }
func (r *MacroNamingRule) Description() string { return "Nom de macro en SCREAMING_SNAKE_CASE" }
func (r *MacroNamingRule) Level() int          { return 1 }

func (r *MacroNamingRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	// Regex pour #define
	defineRegex := regexp.MustCompile(`^\s*#define\s+(\w+)`)
	
	for i, line := range ctx.Lines {
		matches := defineRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			macroName := matches[1]
			if !isScreamingSnakeCase(macroName) {
				violations = append(violations, Violation{
					Rule:     r.Name(),
					Message:  "Nom de macro non conforme au SCREAMING_SNAKE_CASE",
					Line:     i + 1,
					Severity: "major",
					Description: "Le nom de macro '" + macroName + "' doit être en SCREAMING_SNAKE_CASE",
				})
			}
		}
	}
	
	return violations
}

// FunctionLengthRule vérifie la longueur des fonctions (max 25 lignes)
type FunctionLengthRule struct{}

func (r *FunctionLengthRule) Name() string        { return "C-F3" }
func (r *FunctionLengthRule) Description() string { return "Fonction de maximum 25 lignes" }
func (r *FunctionLengthRule) Level() int          { return 1 }

func (r *FunctionLengthRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\([^)]*\)\s*$`)
	inFunction := false
	funcStart := 0
	funcName := ""
	braceCount := 0
	
	for i, line := range ctx.Lines {
		// Début de fonction
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName = matches[1]
			funcStart = i + 1
			inFunction = true
			braceCount = 0
		}
		
		if inFunction {
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			
			// Fin de fonction
			if braceCount == 0 && strings.Contains(line, "}") {
				funcLength := i + 1 - funcStart + 1
				if funcLength > 25 {
					violations = append(violations, Violation{
						Rule:     r.Name(),
						Message:  "Fonction trop longue",
						Line:     funcStart,
						Severity: "major",
						Description: "La fonction '" + funcName + "' fait " + 
							strings.Repeat("", funcLength) + " lignes (max: 25)",
					})
				}
				inFunction = false
			}
		}
	}
	
	return violations
}

// Fonctions utilitaires
func isSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	
	for _, r := range s {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	
	// Ne doit pas commencer ou finir par _
	return !strings.HasPrefix(s, "_") && !strings.HasSuffix(s, "_")
}

func isScreamingSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	
	return !strings.HasPrefix(s, "_") && !strings.HasSuffix(s, "_")
}