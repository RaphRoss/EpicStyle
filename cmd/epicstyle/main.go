package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RaphRoss/EpicStyle/pkg/analyzer"
	"github.com/RaphRoss/EpicStyle/pkg/reporter"
)

type Config struct {
	Path     string
	Verbose  bool
	JSON     bool
	Silent   bool
	Level    int
}

func main() {
	config := parseFlags()
	
	if config.Path == "" {
		fmt.Println("Usage: epicstyle [options] <file_or_directory>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	analyzer := analyzer.New()
	
	// Analyse du fichier ou dossier
	results, err := analyzeTarget(analyzer, config.Path, config.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}

	// Génération du rapport
	rep := reporter.New(config.JSON, config.Verbose, config.Silent)
	rep.Generate(results)
	
	// Code de sortie basé sur le nombre de violations
	if hasViolations(results) {
		os.Exit(1)
	}
}

func parseFlags() Config {
	var config Config
	
	flag.StringVar(&config.Path, "path", "", "Chemin du fichier ou dossier à analyser")
	flag.BoolVar(&config.Verbose, "verbose", false, "Sortie détaillée")
	flag.BoolVar(&config.JSON, "json", false, "Sortie au format JSON")
	flag.BoolVar(&config.Silent, "silent", false, "Sortie silencieuse (code de retour uniquement)")
	flag.IntVar(&config.Level, "level", 1, "Niveau de vérification (1=base, 2=avancé)")
	
	flag.Parse()
	
	// Si un argument positionnel est fourni, l'utiliser comme path
	if len(flag.Args()) > 0 {
		config.Path = flag.Args()[0]
	}
	
	return config
}

func analyzeTarget(analyzer *analyzer.Analyzer, path string, level int) ([]*analyzer.FileResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var files []string
	
	if info.IsDir() {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".c" || filepath.Ext(path) == ".h" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if filepath.Ext(path) == ".c" || filepath.Ext(path) == ".h" {
			files = append(files, path)
		} else {
			return nil, fmt.Errorf("le fichier doit avoir une extension .c ou .h")
		}
	}

	var results []*analyzer.FileResult
	for _, file := range files {
		result, err := analyzer.AnalyzeFile(file, level)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func hasViolations(results []*analyzer.FileResult) bool {
	for _, result := range results {
		if len(result.Violations) > 0 {
			return true
		}
	}
	return false
}