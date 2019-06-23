package ru_nalog

import (
	"strconv"
	"testing"
)

// FindTag is called very often
func BenchmarkFindDescOnlyBuiltin(b *testing.B) {
	newCheck := func(tag Tag) func(*testing.B) {
		return func(b *testing.B) {
			b.ResetTimer()
			for i := 1; i <= b.N; i++ {
				found := FindTag(tag)
				if found == nil {
					b.Fatal("not found")
				}
			}
		}
	}
	for i := 0; i < len(builtinTags); i += (len(builtinTags) / 9) {
		t := builtinTags[i].Tag
		b.Run(strconv.Itoa(int(t)), newCheck(t))
	}
}
