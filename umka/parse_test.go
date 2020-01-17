package umka

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ru_nalog "github.com/temoto/ru-nalog-go"
)

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
			require.Equal(t, expected, tlv.String(), message)
		default:
			require.Equal(t, expected, tlv.Value(), message)
		}
	}
}

func TestParseResponseDoc(t *testing.T) {
	t.Parallel()

	type Case struct {
		name  string
		input string
		check func(t testing.TB, in string)
	}
	cases := []Case{
		{name: "fisc-doc", input: `{
"protocol": 1,
"version": "1.0",
"cashboxStatus": {
  "agentFlags": 0,
  "allowGames": false,
  "allowLotteries": false,
  "allowServices": false,
  "automatMode": false,
  "cash": 11421746,
  "cashBoxNumber": 1,
  "cashier": 99,
  "cycleNumber": 17,
  "cycleOpened": "21 Jun 2019 16:02:46 +0300",
  "dt": "21 Jun 2019 22:46:07 +0300",
  "email": "email@example.com",
  "excisableGoods": false,
  "externPrinter": false,
  "fSfDfVersion": 2,
  "fdfVersion": 2,
  "flags": 75,
  "fnsSite": "www.nalog.ru",
  "fsNumber": "9999078900003063",
  "fsStatus": {
    "cycleIsOpen": 1,
    "debugMode": true,
    "fsNumber": "9999078900003063",
    "fsVersion": "fn debug v 1.32",
    "lastDocDt": "2019-06-21T20:53:00",
    "lastDocNumber": 76,
    "lifeTime": {
      "availableRegistrations": 11,
      "completedRegistrations": 1,
      "expirationDt": "2020-07-20"
    },
    "phase": 3,
    "transport": {
      "docIsReading": true,
      "firstDocDt": null,
      "firstDocNumber": 0,
      "offlineDocsCount": 0,
      "state": 0
    }
  },
  "internetOnly": false,
  "introductions": 0,
  "introductionsSum": 0,
  "ipAddresses": "192.168.1.38",
  "lastCheckNumber": 13,
  "makeBso": false,
  "mode": 1,
  "model": 200,
  "modelstr": "УМКА-01-ФА",
  "ofdInn": "2310031475",
  "ofdName": "Акционерное общество Тандер",
  "offlineMode": false,
  "paymentAddress": "г. Верхненижнесевероюгонефтегазовск, ул. Народной Гомеопатии, д. 234, корп. 3234",
  "paymentPlace": "ЭТО ТЕСТОВАЯ КАССА. СТАТУС ПО НЕЙ НЕ СВЕРЯТЬ!",
  "payouts": 0,
  "payoutsSum": 0,
  "regCashierInn": "000000000000",
  "regCashierName": "Аглая Никифоровна Джоумобетова",
  "regDate": "2019-06-19",
  "regDocNumber": 1,
  "regNumber": "1693666568053977",
  "serial": "16999987",
  "shortFlags": 0,
  "subMode": 0,
  "subver": 1,
  "taxes": 63,
  "useEncryption": false,
  "userInn": "7725225244",
  "userName": "ТОРГОВЫЕ АВТОМАТЫ ЯЙЦЕКЛАДУЩЕГО КОМБИНАТА",
  "ver": 0
}}`, check: func(t testing.TB, in string) {
			_, err := ParseResponseDoc([]byte(in))
			require.NoError(t, err, "input=%s", in)
		}},
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
			{"tag": 1077, "value": 1359967045, "printable": "ФП\t1359967045"}
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
				findCheckEqual(t, d, 1077, uint64(1359967045))
				findCheckEqual(t, d, 1116, uint32(19))

				stlv := d.FindByTag(1059)
				require.NotNil(t, stlv)
				findCheckEqual(t, stlv, 1023, uint64(22000))
				findCheckEqual(t, stlv, 1079, uint32(3300))
			}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) { c.check(t, c.input) })
	}
}
