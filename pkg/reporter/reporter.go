package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/epicstyle/pkg/analyzer"
	"github.com/your-username/epicstyle/pkg/rules"
)

// Reporter gère l'affichage des résultats
type Reporter struct {
	jsonOutput bool
	verbose    bool
	silent     bool
}

// New crée un nouveau reporter
func New(jsonOutput, verbose, silent bool) *Reporter {
	return &Reporter{
		jsonOutput: jsonOutput,
		verbose:    verbose,
		silent:     silent,
	}
}

// Generate génère et affiche le rapport
func (r *Reporter) Generate(results []*analyzer.FileResult) {
	if r.silent {
		return
	}

	globalResults := analyzer.CalculateGlobalResults(results)

	if r.jsonOutput {
		r.generateJSONReport(globalResults)
	} else {
		r.generateTextReport(globalResults)
	}
}

// generateJSONReport génère un rapport JSON
func (r *Reporter) generateJSONReport(results *analyzer.AnalyzeResults) {
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la génération JSON: %v\n", err)
		return
	}
	
	fmt.Println(string(output))
}

// generateTextReport génère un rapport texte
func (r *Reporter) generateTextReport(results *analyzer.AnalyzeResults) {
	// En-tête
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                            EPICSTYLE - RAPPORT D'ANALYSE                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Résumé global
	r.printSummary(results)
	fmt.Println()

	// Détails par fichier
	for _, fileResult := range results.Files {
		r.printFileResult(fileResult)
	}

	// Score final
	r.printFinalScore(results)
}

// printSummary affiche le résumé global
func (r *Reporter) printSummary(results *analyzer.AnalyzeResults) {
	fmt.Printf("📊 RÉSUMÉ GLOBAL\n")
	fmt.Printf("   • Fichiers analysés: %d\n", results.TotalFiles)
	fmt.Printf("   • Lignes de code: %d\n", results.TotalLines)
	fmt.Printf("   • Violations totales: %d\n", results.Violations)
	fmt.Printf("   • Fichiers propres: %d/%d\n", results.CleanFiles, results.TotalFiles)
	
	// Barre de progression visuelle
	cleanPercentage := float64(results.CleanFiles) / float64(results.TotalFiles) * 100
	fmt.Printf("   • Propreté: %.1f%% ", cleanPercentage)
	r.printProgressBar(cleanPercentage)
	fmt.Println()
}

// printFileResult affiche les résultats d'un fichier
func (r *Reporter) printFileResult(result *analyzer.FileResult) {
	filename := filepath.Base(result.Filename)
	
	if len(result.Violations) == 0 {
		fmt.Printf("✅ %s (%.1f%% - %d lignes)\n", filename, result.Score, result.LineCount)
		return
	}

	// Fichier avec violations
	fmt.Printf("❌ %s (%.1f%% - %d lignes - %d violations)\n", 
		filename, result.Score, result.LineCount, len(result.Violations))

	if r.verbose {
		// Grouper les violations par règle
		violationsByRule := make(map[string][]rules.Violation)
		for _, violation := range result.Violations {
			violationsByRule[violation.Rule] = append(violationsByRule[violation.Rule], violation)
		}

		for rule, violations := range violationsByRule {
			fmt.Printf("   🔸 %s (%d violations)\n", rule, len(violations))
			
			for _, violation := range violations {
				severity := r.getSeverityIcon(violation.Severity)
				fmt.Printf("      %s Ligne %d: %s\n", severity, violation.Line, violation.Message)
				
				if violation.Description != "" {
					fmt.Printf("         💡 %s\n", violation.Description)
				}
			}
		}
		fmt.Println()
	}
}

// printFinalScore affiche le score final avec style
func (r *Reporter) printFinalScore(results *analyzer.AnalyzeResults) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║                           SCORE GLOBAL: %.1f%%", results.TotalScore)
	
	// Padding pour centrer
	padding := 79 - len(fmt.Sprintf("SCORE GLOBAL: %.1f%%", results.TotalScore)) - 27
	fmt.Print(strings.Repeat(" ", padding))
	fmt.Println("║")
	
	// Barre de score visuelle
	fmt.Print("║ ")
	r.printProgressBar(results.TotalScore)
	fmt.Print(" ║")
	fmt.Println()
	
	// Message selon le score
	message := r.getScoreMessage(results.TotalScore)
	messagePadding := (78 - len(message)) / 2
	fmt.Printf("║%s%s%s║\n", 
		strings.Repeat(" ", messagePadding), 
		message, 
		strings.Repeat(" ", 78-len(message)-messagePadding))
	
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
}

// printProgressBar affiche une barre de progression
func (r *Reporter) printProgressBar(percentage float64) {
	barLength := 50
	filled := int(percentage * float64(barLength) / 100)
	
	fmt.Print("[")
	for i := 0; i < barLength; i++ {
		if i < filled {
			if percentage >= 80 {
				fmt.Print("█") // Vert
			} else if percentage >= 60 {
				fmt.Print("▓") // Orange
			} else {
				fmt.Print("▒") // Rouge
			}
		} else {
			fmt.Print("░")
		}
	}
	fmt.Printf("] %.1f%%", percentage)
}

// getSeverityIcon retourne l'icône selon la gravité
func (r *Reporter) getSeverityIcon(severity string) string {
	switch severity {
	case "major":
		return "🚨"
	case "minor":
		return "⚠️ "
	case "info":
		return "ℹ️ "
	default:
		return "❓"
	}
}

// getScoreMessage retourne un message selon le score
func (r *Reporter) getScoreMessage(score float64) string {
	switch {
	case score >= 95:
		return "🏆 EXCELLENT! Code parfaitement conforme!"
	case score >= 85:
		return "🎉 TRÈS BIEN! Quelques petits détails à corriger."
	case score >= 70:
		return "👍 BIEN! Bon travail, continuez les améliorations."
	case score >= 50:
		return "⚠️  MOYEN. Plusieurs points à améliorer."
	default:
		return "❌ INSUFFISANT. Révision majeure nécessaire."
	}
}