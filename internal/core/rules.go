package core

func DependencyRules() map[string][]string {
	return map[string][]string{
		"tests":   {"api", "base", "feature", "mock", "app"},
		"app":     {"api", "base", "feature", "mock"},
		"mock":    {"api", "base"},
		"feature": {"api", "base"},
		"base":    {"api", "base"},
		"api":     {},
	}
}
