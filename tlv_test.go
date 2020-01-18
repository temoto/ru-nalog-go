package ru_nalog

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
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
			tlv := NewTLV(desc.Tag)
			require.NotNil(t, tlv, "tag is not supported")

			Q := func(f interface{}) {
				t.Helper()
				require.NoErrorf(t, quick.Check(f, nil), "Tag=%d Kind=%s", tlv.Tag, tlv.Kind.String())
			}
			Eq := assert.New(t).Equal
			F := func(format string) func(interface{}) string {
				return func(n interface{}) string { return fmt.Sprintf("(#%d "+format+")", tlv.Tag, n) }
			}
			check := func(n interface{}, ideq func() bool, gostr func(interface{}) string) bool {
				tlv.SetValue(n)
				require.NoError(t, tlv.Err())
				okGS := assert.Equal(t, gostr(n), tlv.GoString())
				return ideq() && okGS
			}

			switch desc.Kind {
			case DataKindBool:
				Q(func(n bool) bool { return check(n, func() bool { return Eq(n, tlv.Bool()) }, F("%t")) })
			case DataKindBytes:
				Q(func(n []byte) bool {
					return check(n, func() bool { return assert.True(t, bytes.Equal(n, tlv.Bytes())) }, F("%x"))
				})
				Q(func(n string) bool {
					return check(n, func() bool { return assert.True(t, n == string(tlv.Bytes())) }, F("%x"))
				})
			case DataKindFVLN:
				Q(func(n float64) bool {
					ns := fmt.Sprintf("%.f", n)
					return check(n, func() bool { return Eq(n, tlv.Float64()) },
						func(interface{}) string { return F("%s")(ns) })
				})
				Q(func(r float64) bool {
					n := fmt.Sprintf("%.f", r)
					return check(n, func() bool { return Eq(r, tlv.Float64()) },
						func(interface{}) string { return F("%s")(n) })
				})
			case DataKindString:
				Q(func(n string) bool { return check(n, func() bool { return Eq(n, tlv.String()) }, F("%s")) })
			case DataKindTime:
				Q(func(r uint32) bool {
					n := time.Unix(int64(r), 0)
					ns := n.Format(time.RFC3339)
					return check(n, func() bool { return Eq(n, tlv.Time()) }, func(interface{}) string { return F("%s")(ns) })
				})
			case DataKindUint:
				Q(func(n uint) bool {
					return check(n, func() bool { return Eq(uint32(n), tlv.Uint32()) },
						func(interface{}) string { return F("%x")(uint32(n)) })
				})
				Q(func(n uint32) bool {
					return check(n, func() bool { return Eq(n, tlv.Uint32()) },
						func(interface{}) string { return F("%x")(n) })
				})
				Q(func(n uint64) bool {
					return check(n, func() bool { return Eq(uint32(n), tlv.Uint32()) },
						func(interface{}) string { return F("%x")(uint32(n)) })
				})
			case DataKindVLN:
				if desc.Length <= 6 {
					Q(func(n uint) bool {
						return check(n, func() bool { return Eq(uint32(n), uint32(tlv.Uint64())) },
							func(interface{}) string { return F("%d")(uint32(n)) })
					})
					Q(func(n uint32) bool {
						return check(n, func() bool { return Eq(n, uint32(tlv.Uint64())) },
							func(interface{}) string { return F("%d")(uint32(n)) })
					})
					Q(func(n uint64) bool {
						return check(n, func() bool { return Eq(uint32(n), uint32(tlv.Uint64())) },
							func(interface{}) string { return F("%d")(uint32(n)) })
					})
				} else {
					Q(func(n uint) bool { return check(n, func() bool { return Eq(n, uint(tlv.Uint64())) }, F("%d")) })
					Q(func(n uint32) bool { return check(n, func() bool { return Eq(n, uint32(tlv.Uint64())) }, F("%d")) })
					Q(func(n uint64) bool { return check(n, func() bool { return Eq(n, tlv.Uint64()) }, F("%d")) })
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
