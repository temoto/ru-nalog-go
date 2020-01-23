package umka

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ru_nalog "github.com/temoto/ru-nalog-go"
)

func Test401(t *testing.T) {
	t.Parallel()

	const stubHeader401 = "HTTP/1.0 401 Unauthorized\r\n\r\n"
	u, err := NewUmka(&UmkaConfig{
		BaseURL: "mock",
		RT:      &mockRT{header: []byte(stubHeader401)},
	})
	require.NoError(t, err)
	st, err := u.Status()
	assert.Nil(t, st)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status=401")
}

func TestFiscalCheck(t *testing.T) {
	t.Parallel()

	checkError := func(t testing.TB, u Umker, doc1 *ru_nalog.Doc) {
		doc2, err := u.FiscalCheck("TODO_session_id", doc1)
		require.Error(t, err)
		t.Logf(err.Error())
		assert.Nil(t, doc2)
	}
	checkOK := func(t testing.TB, u Umker, doc1 *ru_nalog.Doc) {
		doc2, err := u.FiscalCheck("TODO_session_id", doc1)
		require.NoError(t, err)
		t.Logf(doc2.String())
	}
	cases := []struct {
		name     string
		response string
		check    func(t testing.TB, u Umker, doc1 *ru_nalog.Doc)
	}{
		{"result106", "{\"document\":{\"message\":{\"resultDescription\":\"\320\235\320\265\320\262\320\265\321\200\320\275\321\213\320\271 \321\202\320\270\320\277 \321\207\320\265\320\272\320\260\"},\"result\":106,\"sessionId\":\"xyz\"},\"protocol\":1,\"version\":\"1.0\"}", checkError},

		// "POST /fiscalcheck.json HTTP/1.1\r\nHost: office.armax.ru:58088\r\nContent-Length: 513\r\nAuthorization: Basic NzA6NzA=\r\nAccept-Encoding: gzip\r\n\r\n{\"protocol\":0,\"version\":\"\",\"document\":{\"message\":{\"resultDescription\":\"\"},\"data\":{\"moneyType\":1,\"sum\":0,\"type\":1,\"fiscprops\":[{\"tag\":1054,\"value\":1},{\"tag\":1055,\"value\":2},{\"tag\":1036,\"value\":\"18446744073709551615\"},{\"fiscprops\":[{\"tag\":1023,\"value\":1},{\"tag\":1030,\"value\":\"FIXME_code\"},{\"tag\":1079,\"value\":200},{\"tag\":1199,\"value\":6},{\"tag\":1212,\"value\":1},{\"tag\":1214,\"value\":1}],\"tag\":1059}]},\"sessionId\":\"vm=-1,time=2020-01-20T19:31:43.430224+05:00,code=1,name=FIXME_code2020-01-25T07:31:59+05:00\",\"print\":0}}"
		{"result10", "{\"document\":{\"message\":{\"resultDescription\":\"\320\235\320\265\320\262\320\265\321\200\320\275\320\276\320\265 \320\272\320\276\320\273\320\270\321\207\320\265\321\201\321\202\320\262\320\276\"},\"result\":10,\"sessionId\":\"xyz\"},\"protocol\":1,\"version\":\"1.0\"}", checkError},

		{"success", "{\"document\":{\"data\":{\"docNumber\":8493,\"docType\":3,\"fiscprops\":[{\"printable\":\"\320\220\320\264\321\200\320\265\321\201 \321\200\320\260\321\201\321\207\320\265\321\202\320\276\320\262\",\"tag\":1009,\"value\":\"\320\220\320\264\321\200\320\265\321\201 \321\200\320\260\321\201\321\207\320\265\321\202\320\276\320\262\"},{\"caption\":\"\320\230\320\235\320\235\",\"printable\":\"\320\230\320\235\320\235\\t7725225244\",\"tag\":1018,\"value\":\"7725225244\"},{\"caption\":\"\320\232\320\220\320\241\320\241\320\230\320\240\",\"printable\":\"\320\232\320\220\320\241\320\241\320\230\320\240\\t\320\232\320\220\320\241\320\241\320\230\320\240 70\",\"tag\":1021,\"value\":\"\320\232\320\220\320\241\320\241\320\230\320\240 70\"},{\"printable\":\"\320\220\321\200\320\274\320\260\320\272\321\201\",\"tag\":1048,\"value\":\"\320\220\321\200\320\274\320\260\320\272\321\201\"},{\"caption\":\"\320\234\320\225\320\241\320\242\320\236 \320\240\320\220\320\241\320\247\320\225\320\242\320\236\320\222\",\"printable\":\"\320\234\320\225\320\241\320\242\320\236 \320\240\320\220\320\241\320\247\320\225\320\242\320\236\320\222\\t\320\234\320\265\321\201\321\202\320\276 \321\200\320\260\321\201\321\207\320\265\321\202\320\276\320\262\",\"tag\":1187,\"value\":\"\320\234\320\265\321\201\321\202\320\276 \321\200\320\260\321\201\321\207\320\265\321\202\320\276\320\262\"},{\"printable\":\"25.01.20 06:18\",\"tag\":1012,\"value\":\"25 Jan 2020 06:18:19 +0300\"},{\"caption\":\"\320\230\320\242\320\236\320\223\",\"printable\":\"\320\230\320\242\320\236\320\223\\t2,00\",\"tag\":1020,\"value\":200},{\"caption\":\"\320\241\320\234\320\225\320\235\320\220\",\"printable\":\"\320\241\320\234\320\225\320\235\320\220\\t372\",\"tag\":1038,\"value\":372},{\"caption\":\"\320\247\320\225\320\232\",\"printable\":\"\320\247\320\225\320\232 1\",\"tag\":1042,\"value\":1},{\"printable\":\"\320\237\320\240\320\230\320\245\320\236\320\224\",\"tag\":1054,\"value\":1},{\"caption\":\"\320\241\320\235\320\236\",\"printable\":\"\320\241\320\235\320\236\\t\320\243\320\241\320\235 \320\264\320\276\321\205\320\276\320\264\",\"tag\":1055,\"value\":2},{\"fiscprops\":[{\"printable\":\"\320\237\320\240\320\225\320\224\320\236\320\237\320\233\320\220\320\242\320\220 100%\",\"tag\":1214,\"value\":1},{\"printable\":\"\\t\320\242\",\"tag\":1212,\"value\":1},{\"printable\":\"FIXME_code\",\"tag\":1030,\"value\":\"FIXME_code\"},{\"printable\":\"2,00\",\"tag\":1079,\"value\":200},{\"printable\":\"1\",\"tag\":1023,\"value\":\"1,000\"},{\"tag\":1199,\"value\":6},{\"printable\":\"2,00\",\"tag\":1043,\"value\":200}],\"tag\":1059},{\"tag\":1209,\"value\":2},{\"printable\":\"t=20200125T0618&s=2.00&fn=9999078900003063&i=8493&fp=1765583868&n=1\",\"tag\":1196,\"value\":\"t=20200125T0618&s=2.00&fn=9999078900003063&i=8493&fp=1765583868&n=1\"},{\"caption\":\"\320\235\320\220\320\233\320\230\320\247\320\235\320\253\320\234\320\230\",\"printable\":\"\320\235\320\220\320\233\320\230\320\247\320\235\320\253\320\234\320\230\\t2,00\",\"tag\":1031,\"value\":200},{\"caption\":\"\320\241\320\243\320\234\320\234\320\220 \320\221\320\225\320\227 \320\235\320\224\320\241\",\"printable\":\"\320\241\320\243\320\234\320\234\320\220 \320\221\320\225\320\227 \320\235\320\224\320\241\\t2,00\",\"tag\":1105,\"value\":200},{\"caption\":\"\320\240\320\235 \320\232\320\232\320\242\",\"printable\":\"\320\240\320\235 \320\232\320\232\320\242\\t0329868379061673\",\"tag\":1037,\"value\":\"0329868379061673\"},{\"caption\":\"\320\244\320\224\",\"printable\":\"\320\244\320\224\\t8493\",\"tag\":1040,\"value\":8493},{\"caption\":\"\320\244\320\235\",\"printable\":\"\320\244\320\235\\t9999078900003063\",\"tag\":1041,\"value\":\"9999078900003063\"},{\"caption\":\"\320\241\320\220\320\231\320\242 \320\244\320\235\320\241\",\"printable\":\"\320\241\320\220\320\231\320\242 \320\244\320\235\320\241\\tnalog.ru\",\"tag\":1060,\"value\":\"nalog.ru\"},{\"caption\":\"\320\244\320\237\",\"printable\":\"\320\244\320\237\\t1765583868\",\"tag\":1077,\"value\":\"1765583868\"},{\"caption\":\"\320\227\320\235 \320\232\320\232\320\242\",\"printable\":\"\320\227\320\235 \320\232\320\232\320\242\\t16999987\",\"tag\":1013,\"value\":\"16999987\"},{\"caption\":\"\320\244\320\244\320\224 \320\232\320\232\320\242\",\"printable\":\"\320\244\320\244\320\224 \320\232\320\232\320\242\\t 1.05\",\"tag\":1189,\"value\":2},{\"caption\":\"\320\255\320\233. \320\220\320\224\320\240. \320\236\320\242\320\237\320\240\320\220\320\222\320\230\320\242\320\225\320\233\320\257\",\"printable\":\"\320\255\320\233. \320\220\320\224\320\240. \320\236\320\242\320\237\320\240\320\220\320\222\320\230\320\242\320\225\320\233\320\257\\taaa@bbb.ru\",\"tag\":1117,\"value\":\"aaa@bbb.ru\"}],\"name\":\"\320\232\320\260\321\201\321\201\320\276\320\262\321\213\320\271 \321\207\320\265\320\272\"},\"result\":0,\"sessionId\":\"vm=-1,time=2020-01-20T19:31:43.430224+05:00,code=1,name=FIXME_code2020-01-25T08:18:19+05:00\"},\"protocol\":1,\"version\":\"1.0\"}", checkOK},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			u, err := NewUmka(&UmkaConfig{
				BaseURL: "mock",
				RT:      &mockRT{body: []byte(c.response)},
			})
			require.NoError(t, err)
			doc1 := ru_nalog.NewDoc(0, ru_nalog.FDCheck)
			doc1.AppendNew(1054, 1)
			doc1.AppendNew(1055, 2)
			doc1.AppendNew(1036, "18446744073709551615")
			row := doc1.AppendNew(1059, nil)
			row.AppendNew(1023, 1)
			row.AppendNew(1030, "FIXME_code")
			row.AppendNew(1079, 200)
			row.AppendNew(1199, 6)
			row.AppendNew(1212, 1)
			row.AppendNew(1214, 1)
			c.check(t, u, doc1)
		})
	}
}

func TestStatus(t *testing.T) {
	t.Parallel()

	const stubResponseBody = `{"cashboxStatus":{"agentFlags":127,"allowGames":false,"allowLotteries":false,"allowServices":false,"automatMode":false,"cash":2156099220,"cashBoxNumber":1,"cashier":1,"cycleNumber":369,"cycleOpened":"22 Jan 2020 19:34:54 +0300","dt":"23 Jan 2020 15:30:48 +0300","email":"aaa@bbb.ru","excisableGoods":false,"externPrinter":false,"fSfDfVersion":2,"fdfVersion":2,"flags":75,"fnsSite":"nalog.ru","fsNumber":"9999078900003063","fsStatus":{"cycleIsOpen":1,"debugMode":true,"fsNumber":"9999078900003063","fsVersion":"fn debug v 1.32","lastDocDt":"2020-01-23T15:20:00","lastDocNumber":8433,"lifeTime":{"availableRegistrations":11,"completedRegistrations":1,"expirationDt":"2020-10-01"},"phase":3,"transport":{"docIsReading":true,"firstDocDt":null,"firstDocNumber":0,"offlineDocsCount":0,"state":0}},"internetOnly":false,"introductions":0,"introductionsSum":0,"ipAddresses":"192.168.1.38","lastCheckNumber":20,"makeBso":false,"mode":1,"model":200,"modelstr":"УМКА-01-ФА","ofdInn":"2310031475","ofdName":"Акционерное общество Тандер","offlineMode":false,"paymentAddress":"Адрес расчетов","paymentPlace":"Место расчетов","payouts":0,"payoutsSum":0,"regCashierInn":"000000000000","regCashierName":"СИС. АДМИН","regDate":"2019-08-22","regDocNumber":1,"regNumber":"0329868379061673","serial":"16999987","shortFlags":0,"subMode":0,"subver":1,"taxes":63,"useEncryption":false,"userInn":"7725225244","userName":"Армакс","ver":0},"protocol":1,"version":"1.0"}`
	u, err := NewUmka(&UmkaConfig{
		BaseURL: "mock",
		RT:      &mockRT{body: []byte(stubResponseBody)},
	})
	require.NoError(t, err)
	st, err := u.Status()
	require.NoError(t, err)
	// t.Logf("st=%#v", st)
	assert.Equal(t, uint64(2156099220), st.Cash)
	assert.True(t, st.IsCycleOpen())
	assert.Equal(t, uint32(0), st.OfdOfflineCount())
	assert.Equal(t, "9999078900003063", st.FsNumber)
	assert.Equal(t, "2020-01-23T15:20:00", st.FsStatus.LastDocDt)
}

type mockRT struct {
	header []byte
	body   []byte
	err    error
}

func (m *mockRT) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	header := m.header
	if header == nil {
		header = []byte("HTTP/1.0 200 OK\r\n\r\n")
	}
	rb := make([]byte, 0, len(header)+len(m.body))
	rb = append(rb, header...)
	rb = append(rb, m.body...)
	return http.ReadResponse(bufio.NewReader(bytes.NewReader(rb)), req)
}
