# EpicStyle - Vérificateur de Style Epitech

EpicStyle est un outil en ligne de commande développé en Go pour analyser automatiquement la conformité des fichiers C (.c) et headers (.h) avec la norme de style Epitech.

## 🚀 Fonctionnalités

### Vérifications de Base (Niveau 1)
- ✅ Taille maximale d'une ligne (80 caractères)
- ✅ Aucune ligne vide en début/fin de fichier
- ✅ Aucune ligne vide consécutive
- ✅ Indentation en TAB uniquement
- ✅ Une seule variable déclarée par ligne
- ✅ Déclarations de variables en début de fonction uniquement
- ✅ Nom de fichier en snake_case
- ✅ Nom de fonction en snake_case
- ✅ Nom de macro en SCREAMING_SNAKE_CASE
- ✅ Fonction de 25 lignes maximum
- ✅ Fichier de 3 fonctions maximum (hors main)

### Vérifications Avancées (Niveau 2)
- ✅ Format de commentaires correct (/* */ uniquement)
- ✅ Commentaire de fonction obligatoire
- ✅ Pas de déclaration globale non const
- ✅ Maximum 4 paramètres par fonction
- ✅ Pas de déclaration dans les boucles for

### Fonctionnalités Complémentaires
- 📊 Rapport détaillé dans le terminal
- 🎯 Score global de conformité
- 📋 Sortie JSON pour automatisation
- 🎨 Interface colorée et intuitive

## 📦 Installation

### Prérequis
- Go 1.21 ou supérieur

### Compilation
```bash
# Cloner le projet
git clone https://github.com/your-username/epicstyle.git
cd epicstyle

# Initialiser le module Go
go mod init github.com/your-username/epicstyle

# Compiler
go build -o epicstyle cmd/epicstyle/main.go

# Ou installer globalement
go install cmd/epicstyle/main.go
```

## 🎯 Utilisation

### Syntaxe de base
```bash
epicstyle [options] <fichier_ou_dossier>
```

### Options disponibles
- `-path` : Chemin du fichier ou dossier à analyser
- `-verbose` : Affichage détaillé des violations
- `-json` : Sortie au format JSON
- `-silent` : Mode silencieux (code de retour uniquement)
- `-level` : Niveau de vérification (1=base, 2=avancé)

### Exemples d'utilisation

```bash
# Analyser un fichier
epicstyle mon_fichier.c

# Analyser un dossier avec sortie détaillée
epicstyle -verbose src/

# Générer un rapport JSON
epicstyle -json -level 2 projet/

# Mode silencieux pour scripts
epicstyle -silent fichier.c
echo $?  # 0 = succès, 1 = violations détectées
```

## 📊 Format de Sortie

### Sortie Standard
```
╔══════════════════════════════════════════════════════════════════════════════╗
║                            EPICSTYLE - RAPPORT D'ANALYSE                     ║
╚══════════════════════════════════════════════════════════════════════════════╝

📊 RÉSUMÉ GLOBAL
   • Fichiers analysés: 3
   • Lignes de code: 127
   • Violations totales: 5
   • Fichiers propres: 1/3
   • Propreté: 33.3% [████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 33.3%

✅ utils.c (95.2% - 42 lignes)
❌ main.c (78.5% - 65 lignes - 3 violations)
❌ parser.c (82.1% - 20 lignes - 2 violations)

╔══════════════════════════════════════════════════════════════════════════════╗
║                           SCORE GLOBAL: 85.3%                                ║
║ [██████████████████████████████████████████████░░░░░░░░] 85.3%               ║
║                    🎉 TRÈS BIEN! Quelques petits détails à corriger.         ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Sortie JSON
```json
{
  "files": [
    {
      "filename": "main.c",
      "violations": [
        {
          "rule": "C-L1",
          "message": "Ligne trop longue",
          "line": 15,
          "severity": "major",
          "description": "La ligne contient plus de 80 caractères"
        }
      ],
      "score": 78.5,
      "line_count": 65
    }
  ],
  "total_score": 85.3,
  "total_files": 3,
  "total_lines": 127,
  "total_violations": 5,
  "clean_files": 1
}
```

## 🏗️ Architecture du Projet

```
epicstyle/
├── cmd/epicstyle/          # Point d'entrée principal
│   └── main.go
├── pkg/                    # Packages principaux
│   ├── analyzer/           # Moteur d'analyse
│   │   ├── analyzer.go
│   │   └── file_reader.go
│   ├── rules/              # Règles de style
│   │   ├── rule_interface.go
│   │   ├── base_rules.go
│   │   └── advanced_rules.go
│   ├── reporter/           # Génération de rapports
│   │   └── reporter.go
│   └── utils/              # Utilitaires
│       └── file_utils.go
├── examples/               # Exemples de fichiers
├── go.mod
└── README.md
```

## 🧪 Tests

```bash
# Lancer les tests
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests verbeux
go test -v ./...
```

## 📋 Codes de Règles

### Règles de Base (Niveau 1)
- `C-L1` : Longueur de ligne (80 caractères max)
- `C-L2` : Lignes vides interdites
- `C-L3` : Indentation en TAB
- `C-L4` : Une variable par ligne
- `C-V1` : Déclarations en début de fonction
- `C-O1` : Nom de fichier snake_case
- `C-O2` : Maximum 3 fonctions par fichier
- `C-F1` : Nom de fonction snake_case
- `C-F2` : Nom de macro SCREAMING_SNAKE_CASE
- `C-F3` : Fonction 25 lignes max

### Règles Avancées (Niveau 2)
- `C-C1` : Format de commentaires
- `C-C2` : Commentaire de fonction obligatoire
- `C-G1` : Pas de globales non const
- `C-F4` : Maximum 4 paramètres
- `C-L5` : Pas de déclaration dans les boucles

## 🤝 Contribution

Les contributions sont les bienvenues ! Voici comment contribuer :

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📝 License

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🎯 Roadmap

- [ ] Option `--fix` pour corrections automatiques
- [ ] Support des fichiers de configuration
- [ ] Intégration CI/CD
- [ ] Plugin VSCode
- [ ] Interface web
- [ ] Métriques de complexité
- [ ] Règles personnalisables

## 🐛 Signaler un Bug

Si vous trouvez un bug, merci de créer une issue avec :
- Description du problème
- Fichier exemple qui cause le problème
- Version de Go utilisée
- Système d'exploitation

## 📞 Support

Pour toute question ou suggestion :
- Créer une issue sur GitHub
- Envoyer un email à : support@epicstyle.dev

---

Développé avec ❤️ pour la communauté Epitech
