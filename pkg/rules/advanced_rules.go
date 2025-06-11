package rules

import (
	"fmt"
	"regexp"
	"strings"
)

// CommentFormatRule vérifie le format des commentaires
type CommentFormatRule struct{}

func (r *CommentFormatRule) Name() string        { return "C-C1" }
func (r *CommentFormatRule) Description() string { return "Format de commentaire correct (/* */ pour blocs)" }
func (r *CommentFormatRule) Level() int          { return 2 }

func (r *CommentFormatRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	for i, line := range ctx.Lines {
		// Vérifier les commentaires //
		if strings.Contains(line, "//") {
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Utilisation de // interdit",
				Line:     i + 1,
				Severity: "major",
				Description: "Utiliser /* */ pour les commentaires",
			})
		}
	}
	
	return violations
}

// FunctionCommentRule vérifie les commentaires de fonction
type FunctionCommentRule struct{}

func (r *FunctionCommentRule) Name() string        { return "C-C2" }
func (r *FunctionCommentRule) Description() string { return "Commentaire de fonction obligatoire" }
func (r *FunctionCommentRule) Level() int          { return 2 }

func (r *FunctionCommentRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\([^)]*\)\s*$`)
	
	for i, line := range ctx.Lines {
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName := matches[1]
			
			// Ignorer main
			if funcName == "main" {
				continue
			}
			
			// Vérifier s'il y a un commentaire avant la fonction
			hasComment := false
			if i > 0 {
				prevLine := strings.TrimSpace(ctx.Lines[i-1])
				if strings.HasPrefix(prevLine, "/**") || strings.HasPrefix(prevLine, "/*") {
					hasComment = true
				}
			}
			
			if !hasComment {
				violations = append(violations, Violation{
					Rule:     r.Name(),
					Message:  "Commentaire de fonction manquant",
					Line:     i + 1,
					Severity: "major",
					Description: "La fonction '" + funcName + "' doit avoir un commentaire",
				})
			}
		}
	}
	
	return violations
}

// GlobalVariableRule vérifie les déclarations globales non const
type GlobalVariableRule struct{}

func (r *GlobalVariableRule) Name() string        { return "C-G1" }
func (r *GlobalVariableRule) Description() string { return "Pas de déclaration globale non const" }
func (r *GlobalVariableRule) Level() int          { return 2 }

func (r *GlobalVariableRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	globalVarRegex := regexp.MustCompile(`^\s*(int|char|float|double|long|short|unsigned)\s+\w+\s*[=;]`)
	inFunction := false
	braceLevel := 0
	
	for i, line := range ctx.Lines {
		// Suivre les niveaux de braces pour savoir si on est dans une fonction
		braceLevel += strings.Count(line, "{") - strings.Count(line, "}")
		
		// Si on trouve une fonction, on est dans du code local
		if strings.Contains(line, "(") && 
		   strings.Contains(line, ")") && 
		   (strings.Contains(line, "{") || (i < len(ctx.Lines)-1 && strings.Contains(ctx.Lines[i+1], "{"))) {
			inFunction = true
		}
		
		// Si braceLevel revient à 0, on sort des fonctions
		if braceLevel == 0 {
			inFunction = false
		}
		
		// Vérifier les déclarations globales
		if !inFunction && braceLevel == 0 {
			if globalVarRegex.MatchString(line) && !strings.Contains(line, "const") {
				violations = append(violations, Violation{
					Rule:     r.Name(),
					Message:  "Déclaration globale non const",
					Line:     i + 1,
					Severity: "major",
					Description: "Les variables globales doivent être const",
				})
			}
		}
	}
	
	return violations
}

// FunctionParametersRule vérifie le nombre de paramètres (max 4)
type FunctionParametersRule struct{}

func (r *FunctionParametersRule) Name() string        { return "C-F4" }
func (r *FunctionParametersRule) Description() string { return "Maximum 4 paramètres par fonction" }
func (r *FunctionParametersRule) Level() int          { return 2 }

func (r *FunctionParametersRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\(([^)]*)\)`)
	
	for i, line := range ctx.Lines {
		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) > 2 {
			funcName := matches[1]
			params := strings.TrimSpace(matches[2])
			
			// Ignorer les fonctions vides ou avec void
			if params == "" || params == "void" {
				continue
			}
			
			// Compter les paramètres en séparant par les virgules
			paramCount := strings.Count(params, ",") + 1
			
			if paramCount > 4 {
				violations = append(violations, Violation{
					Rule:     r.Name(),
					Message:  "Trop de paramètres",
					Line:     i + 1,
					Severity: "major",
					Description: fmt.Sprintf("La fonction '%s' a %d paramètres (max: 4)", funcName, paramCount),
				})
			}
		}
	}
	
	return violations
}

// LoopDeclarationRule vérifie les déclarations dans les boucles
type LoopDeclarationRule struct{}

func (r *LoopDeclarationRule) Name() string        { return "C-L5" }
func (r *LoopDeclarationRule) Description() string { return "Pas de déclaration dans les boucles for" }
func (r *LoopDeclarationRule) Level() int          { return 2 }

func (r *LoopDeclarationRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	// Regex pour for avec déclaration (ex: for (int i = 0; ...))
	forDeclRegex := regexp.MustCompile(`for\s*\(\s*(int|char|float|double|long|short|unsigned)\s+\w+`)
	
	for i, line := range ctx.Lines {
		if forDeclRegex.MatchString(line) {
			violations = append(violations, Violation{
				Rule:     r.Name(),
				Message:  "Déclaration dans une boucle for",
				Line:     i + 1,
				Severity: "major",
				Description: "Les variables doivent être déclarées avant la boucle",
			})
		}
	}
	
	return violations
}

// FileMaxFunctionsRule vérifie le nombre maximum de fonctions par fichier
type FileMaxFunctionsRule struct{}

func (r *FileMaxFunctionsRule) Name() string        { return "C-O2" }
func (r *FileMaxFunctionsRule) Description() string { return "Maximum 3 fonctions par fichier (hors main)" }
func (r *FileMaxFunctionsRule) Level() int          { return 1 }

func (r *FileMaxFunctionsRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\([^)]*\)\s*$`)
	functionCount := 0
	
	for _, line := range ctx.Lines {
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName := matches[1]
			if funcName != "main" {
				functionCount++
			}
		}
	}
	
	if functionCount > 3 {
		violations = append(violations, Violation{
			Rule:     r.Name(),
			Message:  "Trop de fonctions dans le fichier",
			Line:     1,
			Severity: "major",
			Description: fmt.Sprintf("Le fichier contient %d fonctions (max: 3, hors main)", functionCount),
		})
	}
	
	return violations
}

// VariableDeclarationLocationRule vérifie l'emplacement des déclarations
type VariableDeclarationLocationRule struct{}

func (r *VariableDeclarationLocationRule) Name() string { return "C-V1" }
func (r *VariableDeclarationLocationRule) Description() string {
	return "Déclarations de variables uniquement en début de fonction"
}
func (r *VariableDeclarationLocationRule) Level() int { return 1 }

func (r *VariableDeclarationLocationRule) Check(ctx *FileContext) []Violation {
	var violations []Violation
	
	funcRegex := regexp.MustCompile(`^\s*\w+\s+(\w+)\s*\([^)]*\)\s*$`)
	varDeclRegex := regexp.MustCompile(`^\s*(int|char|float|double|long|short|unsigned)\s+\w+`)
	
	inFunction := false
	funcName := ""
	braceCount := 0
	hasNonDeclStatement := false
	
	for i, line := range ctx.Lines {
		// Début de fonction
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName = matches[1]
			inFunction = true
			braceCount = 0
			hasNonDeclStatement = false
			continue
		}
		
		if inFunction {
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			
			trimmedLine := strings.TrimSpace(line)
			
			// Ignorer les lignes vides et commentaires
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "/*") || strings.HasPrefix(trimmedLine, "//") {
				continue
			}
			
			// Si on trouve une déclaration de variable
			if varDeclRegex.MatchString(line) {
				// Si on a déjà eu des statements non-déclaratifs, c'est interdit
				if hasNonDeclStatement {
					violations = append(violations, Violation{
						Rule:     r.Name(),
						Message:  "Déclaration de variable après du code exécutable",
						Line:     i + 1,
						Severity: "major",
						Description: "Dans la fonction '" + funcName + "', les déclarations doivent être en début",
					})
				}
			} else if trimmedLine != "{" && trimmedLine != "}" {
				// C'est du code exécutable
				hasNonDeclStatement = true
			}
			
			// Fin de fonction
			if braceCount == 0 && strings.Contains(line, "}") {
				inFunction = false
			}
		}
	}
	
	return violations
}