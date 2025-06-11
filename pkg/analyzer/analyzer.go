package analyzer

import (
	"path/filepath"
	"strings"

	"github.com/RaphRoss/EpicStyle/pkg/rules"
)

// FileResult contient les résultats d'analyse d'un fichier
type FileResult struct {
	Filename   string            `json:"filename"`
	Violations []rules.Violation `json:"violations"`
	Score      float64           `json:"score"`
	LineCount  int               `json:"line_count"`
}

// Analyzer est le moteur principal d'analyse
type Analyzer struct {
	ruleSet *rules.RuleSet
}

// New crée un nouvel analyseur avec toutes les règles
func New() *Analyzer {
	ruleSet := rules.NewRuleSet()
	
	// Règles de base (niveau 1)
	ruleSet.Add(&rules.LineLengthRule{})
	ruleSet.Add(&rules.EmptyLinesRule{})
	ruleSet.Add(&rules.IndentationRule{})
	ruleSet.Add(&rules.VariableDeclarationRule{})
	ruleSet.Add(&rules.VariableDeclarationLocationRule{})
	ruleSet.Add(&rules.FilenameRule{})
	ruleSet.Add(&rules.FunctionNamingRule{})
	ruleSet.Add(&rules.MacroNamingRule{})
	ruleSet.Add(&rules.FunctionLengthRule{})
	ruleSet.Add(&rules.FileMaxFunctionsRule{})
	
	// Règles avancées (niveau 2)
	ruleSet.Add(&rules.CommentFormatRule{})
	ruleSet.Add(&rules.FunctionCommentRule{})
	ruleSet.Add(&rules.GlobalVariableRule{})
	ruleSet.Add(&rules.FunctionParametersRule{})
	ruleSet.Add(&rules.LoopDeclarationRule{})
	
	return &Analyzer{
		ruleSet: ruleSet,
	}
}

// AnalyzeFile analyse un fichier et retourne les résultats
func (a *Analyzer) AnalyzeFile(filename string, level int) (*FileResult, error) {
	// Lire le fichier
	content, lines, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	// Créer le contexte
	ctx := &rules.FileContext{
		Filename: filename,
		Lines:    lines,
		Content:  content,
		IsHeader: strings.HasSuffix(filename, ".h"),
	}
	
	// Exécuter toutes les règles du niveau spécifié
	violations := a.ruleSet.CheckAll(ctx, level)
	
	// Calculer le score
	score := a.calculateScore(len(lines), len(violations))
	
	return &FileResult{
		Filename:   filename,
		Violations: violations,
		Score:      score,
		LineCount:  len(lines),
	}, nil
}

// calculateScore calcule un score de qualité basé sur le nombre de violations
func (a *Analyzer) calculateScore(lineCount, violationCount int) float64 {
	if lineCount == 0 {
		return 100.0
	}
	
	// Score basé sur le ratio violations/lignes
	// Plus il y a de violations par ligne, plus le score baisse
	violationRatio := float64(violationCount) / float64(lineCount)
	
	// Score de base : 100%
	// Chaque violation fait perdre des points selon sa gravité
	score := 100.0 - (violationRatio * 100.0)
	
	// Minimum 0, maximum 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	return score
}

// GetRulesList retourne la liste de toutes les règles disponibles
func (a *Analyzer) GetRulesList(level int) []rules.Rule {
	var result []rules.Rule
	for _, rule := range a.ruleSet.GetRules() {
		if rule.Level() <= level {
			result = append(result, rule)
		}
	}
	return result
}

// AnalyzeResults contient les résultats globaux d'analyse
type AnalyzeResults struct {
	Files        []*FileResult `json:"files"`
	TotalScore   float64       `json:"total_score"`
	TotalFiles   int           `json:"total_files"`
	TotalLines   int           `json:"total_lines"`
	Violations   int           `json:"total_violations"`
	CleanFiles   int           `json:"clean_files"`
}

// CalculateGlobalResults calcule les statistiques globales
func CalculateGlobalResults(results []*FileResult) *AnalyzeResults {
	if len(results) == 0 {
		return &AnalyzeResults{}
	}
	
	var totalScore float64
	var totalLines int
	var totalViolations int
	var cleanFiles int
	
	for _, result := range results {
		totalScore += result.Score
		totalLines += result.LineCount
		totalViolations += len(result.Violations)
		
		if len(result.Violations) == 0 {
			cleanFiles++
		}
	}
	
	avgScore := totalScore / float64(len(results))
	
	return &AnalyzeResults{
		Files:        results,
		TotalScore:   avgScore,
		TotalFiles:   len(results),
		TotalLines:   totalLines,
		Violations:   totalViolations,
		CleanFiles:   cleanFiles,
	}
}