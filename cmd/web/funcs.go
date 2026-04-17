package main

import (
	"strings"

	"github.com/DAlba-sudo/verb"
)

// registerFuncs registers custom template helper functions on the verb instance.
func registerFuncs(v *verb.Verb) {
	v.Func("lower", func(s string) string {
		return strings.ToLower(s)
	})

	v.Func("default", func(a, b any) any {
		if b == nil {
			return a
		}
		return b
	})

	v.Func("split", func(s string, sep string) []string {
		return strings.Split(s, sep)
	})

	v.Func("trim", func(a any) string {
		s, ok := a.(string)
		if !ok {
			return ""
		}
		return strings.TrimSpace(s)
	})
}
