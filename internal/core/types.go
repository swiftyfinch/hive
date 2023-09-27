package core

import "regexp"

type Type struct {
	Name    string
	Regexps []regexp.Regexp
}

func DefaultTypes() map[string]Type {
	// Can use any regexp for different platforms
	rawTypes := map[string][]string{
		"tests":   {".*Tests$"},
		"app":     {".*Example$"},
		"mock":    {".*Mock$"},
		"feature": {},
		"base":    {},
		"api":     {".*IO$", ".*Interfaces$"},
	}

	types := map[string]Type{}
	for name, regexpStrings := range rawTypes {
		regexps := []regexp.Regexp{}
		for _, regexpString := range regexpStrings {
			regexps = append(regexps, *regexp.MustCompile(regexpString))
		}
		types[name] = Type{Name: name, Regexps: regexps}
	}
	return types
}
