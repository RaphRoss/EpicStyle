package rules

// Violation représente une violation de règle
type Violation struct {
	Rule        string `json:"rule"`
	Message     string `json:"message"`
	Line        int    `json:"line"`
	Column      int    `json:"column,omitempty"`
	Severity    string `json:"severity"`
	Description string `json:"description,omitempty"`
}

// FileContext contient les informations sur le fichier analysé
type FileContext struct {
	Filename string
	Lines    []string
	Content  string
	IsHeader bool
}

// Rule interface pour toutes les règles de style
type Rule interface {
	Name() string
	Description() string
	Level() int // 1 = base, 2 = avancé
	Check(ctx *FileContext) []Violation
}

// RuleSet contient un ensemble de règles
type RuleSet struct {
	rules []Rule
}

// NewRuleSet crée un nouveau set de règles
func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make([]Rule, 0),
	}
}

// Add ajoute une règle au set
func (rs *RuleSet) Add(rule Rule) {
	rs.rules = append(rs.rules, rule)
}

// CheckAll exécute toutes les règles du niveau spécifié
func (rs *RuleSet) CheckAll(ctx *FileContext, level int) []Violation {
	var violations []Violation
	
	for _, rule := range rs.rules {
		if rule.Level() <= level {
			ruleViolations := rule.Check(ctx)
			violations = append(violations, ruleViolations...)
		}
	}
	
	return violations
}

// GetRules retourne toutes les règles du set
func (rs *RuleSet) GetRules() []Rule {
	return rs.rules
}