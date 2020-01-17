package ru_nalog

import (
	"fmt"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func TestTagType(t *testing.T) {
	t.Parallel()

	for _, desc := range builtinTags {
		t.Run(fmt.Sprintf("%d-%s", desc.Tag, desc.Kind.String()), func(t *testing.T) {
			switch desc.Kind {
			case DataKindVLN:
				if desc.Length <= 6 {
					checkUint32(t, desc.Tag)
				} else {
					checkUint64(t, desc.Tag)
				}
			}
			t.Skip("not implemented")
		})
	}
}

func checkUint32(t testing.TB, tag Tag) {
	tlv := NewTLV(tag)
	require.NotNil(t, tlv, "tag is not supported")
	err := quick.Check(
		func(n uint32) bool { tlv.SetValue(n); return tlv.Uint32() == n },
		nil)
	require.NoErrorf(t, err, "tag=%d kind=%s", tag, tlv.TagDesc.Kind.String())
}

func checkUint64(t testing.TB, tag Tag) {
	tlv := NewTLV(tag)
	require.NotNil(t, tlv, "tag is not supported")
	err := quick.Check(
		func(n uint64) bool { tlv.SetValue(n); return tlv.Uint64() == n },
		nil)
	require.NoErrorf(t, err, "tag=%d kind=%s", tag, tlv.TagDesc.Kind.String())
}

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
