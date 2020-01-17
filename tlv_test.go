package ru_nalog

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	t.Parallel()

	for try := 1; ; try++ {
		freeTag := Tag(rand.Int()) // #nosec G404
		if FindTag(freeTag) == nil {
			require.Nil(t, NewTLV(freeTag))
			break
		}
		if try >= 15 {
			t.Fatal("need free Tag")
		}
	}

	for _, desc := range builtinTags {
		t.Run(fmt.Sprintf("%d-%s", desc.Tag, desc.Kind.String()), func(t *testing.T) {
			check := func(f interface{}) {
				t.Helper()
				require.NoErrorf(t, quick.Check(f, nil), "tag=%d kind=%s", desc.Tag, desc.Kind.String())
			}
			tlv := NewTLV(desc.Tag)
			require.NotNil(t, tlv, "tag is not supported")
			switch desc.Kind {
			case DataKindBool:
				check(func(n bool) bool { tlv.SetValue(n); return tlv.Bool() == n })
			case DataKindBytes:
				check(func(n []byte) bool { tlv.SetValue(n); return bytes.Equal(tlv.Bytes(), n) })
				check(func(n string) bool { tlv.SetValue(n); return string(tlv.Bytes()) == n })
			case DataKindFVLN:
				check(func(n float64) bool { tlv.SetValue(n); return tlv.Float64() == n })
			case DataKindString:
				check(func(n string) bool { tlv.SetValue(n); return tlv.String() == n })
			case DataKindTime:
				check(func(r int64) bool { n := time.Unix(r, 0); tlv.SetValue(n); return n.Equal(tlv.Time()) })
			case DataKindUint:
				check(func(n uint32) bool { tlv.SetValue(n); return tlv.Uint64() == uint64(n) })
				check(func(n uint64) bool { tlv.SetValue(n); return tlv.Uint64() == uint64(n) })
				check(func(n uint) bool { tlv.SetValue(n); return tlv.Uint64() == uint64(n) })
			case DataKindVLN:
				if desc.Length <= 6 {
					check(func(n uint32) bool { tlv.SetValue(n); return tlv.Uint32() == n })
				} else {
					check(func(n uint32) bool { tlv.SetValue(n); return tlv.Uint64() == uint64(n) })
					check(func(n uint64) bool { tlv.SetValue(n); return tlv.Uint64() == n })
				}
			default:
				t.Skipf("not implemented tag=%d kind=%s", desc.Tag, desc.Kind.String())
				return
			}
		})
	}
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
