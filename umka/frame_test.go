package umka

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ru_nalog "github.com/temoto/ru-nalog-go"
)

func TestParseResponseDoc(t *testing.T) {
	t.Parallel()

	type Case struct {
		name  string
		input string
		check func(t testing.TB, in string)
	}
	cases := []Case{
		{name: "response-fiscal-document", input: `{
"protocol": 1,
"version": "1.0",
"document": {
  "result": 0,
  "data": {
		"docNumber": 20,
		"docType": 21,
		"name": "Отчет о текущем состоянии расчетов",
		"fiscprops" : [
			{"tag": 1018, "value": "7725225244  "},
			{"tag": 1048, "value": "ЗАВОД ТРУДНОДОБЫВАЕМЫХ КЕРАМИЧЕСКИХ ЦВЕТОВ"},
			{"tag": 1002, "value": false},
			{"tag": 1012, "value": "18 Jun 2019 21:26:00 +0300"},
			{"tag": 1116, "value": 19},
			{"tag": 1209, "value": 2},
			{"tag": 1059, "fiscprops": [
				{"tag": 1079, "value": 3300, "printable": "33,00"},
				{"tag": 1023, "value": "22,000", "printable": "22"}
			]},
			{"tag": 1077, "value": "1359967045", "printable": "ФП\t1359967045"}
		]
  }
}}`,
			check: func(t testing.TB, in string) {
				d, err := ParseResponseDoc([]byte(in))
				require.NoError(t, err, "input=%s", in)
				assert.Equal(t, d.Number, uint32(20))
				assert.Equal(t, d.Type, ru_nalog.FDStateReport)

				findCheckEqual(t, d, 1002, false)
				findCheckEqual(t, d, 1012, time.Date(2019, 6, 18, 21, 26, 0, 0, time.FixedZone("+0300", 3*3600)))

				findCheckEqual(t, d, 1018, "7725225244")
				findCheckEqual(t, d, 1048, "ЗАВОД ТРУДНОДОБЫВАЕМЫХ КЕРАМИЧЕСКИХ ЦВЕТОВ")
				findCheckEqual(t, d, 1077, "1359967045")
				findCheckEqual(t, d, 1116, uint32(19))

				stlv := d.FindByTag(1059)
				require.NotNil(t, stlv)
				findCheckEqual(t, stlv, 1023, float64(22))
				findCheckEqual(t, stlv, 1079, uint32(3300))
			}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) { c.check(t, c.input) })
	}
}

func findCheckEqual(t testing.TB, f ru_nalog.FindByTager, tag ru_nalog.Tag, expected interface{}) {
	t.Helper()
	message := fmt.Sprintf("tag=%d", tag)
	if tlv := f.FindByTag(tag); assert.NotNil(t, tlv, message) {
		switch expected := expected.(type) {
		case bool:
			require.Equal(t, expected, tlv.Bool(), message)
		case time.Time:
			require.Equal(t, expected.UnixNano(), tlv.Time().UnixNano(), message)
		case string:
			switch tlv.TagDesc.Kind {
			case ru_nalog.DataKindString:
				require.Equal(t, expected, tlv.String(), message)
			case ru_nalog.DataKindBytes:
				require.Equal(t, expected, string(tlv.Bytes()), message)
			default:
				t.Fatalf("not implemented")
			}
		default:
			require.Equal(t, expected, tlv.Value(), message)
		}
	}
}
