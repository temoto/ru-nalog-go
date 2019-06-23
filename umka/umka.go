package umka

import (
	"net/http"
	"net/url"

	"github.com/juju/errors"
	ru_nalog "github.com/temoto/ru-nalog-go"
)

// RFC822Z +century,seconds or RFC1123Z -dayofweek
const TimeLayout = "02 Jan 2006 15:04:05 -0700"

type Umker interface {
	Status() (*Status, error)
	CalcReport() (*ru_nalog.Doc, error)
	XReport() (*ru_nalog.Doc, error)
	GetDoc(number uint32) (*ru_nalog.Doc, error)

	CycleOpen() (*ru_nalog.Doc, error)
	CycleClose() (*ru_nalog.Doc, error)
	NewCheck(sessionId string, d *ru_nalog.Doc) (*ru_nalog.Doc, error)
	Fiscalize(sessionId string, d *ru_nalog.Doc) (*ru_nalog.Doc, error)
	Danger_CloseFiscalStorage(sessionId string) (*ru_nalog.Doc, error)
}

type Umka struct {
	config *UmkaConfig
	h      *http.Client
	url    *url.URL
}
type UmkaConfig struct {
	BaseURL   string
	SecretFun func() (string, string)
	Transport *http.Transport
}

var _ /*type check*/ Umker = &Umka{}

func NewUmka(config *UmkaConfig) (Umker, error) {
	self := &Umka{
		config: config,
		h:      &http.Client{Transport: config.Transport},
	}
	if u, err := url.Parse(config.BaseURL); err == nil {
		self.url = u
	} else {
		err = errors.Trace(err)
		return nil, err
	}

	return self, nil
}

func (self *Umka) Status() (*Status, error) {
	return nil, errors.NotImplementedf("status")
}

func (self *Umka) CalcReport() (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("calc-report")
}

func (self *Umka) XReport() (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("x-report")
}

func (self *Umka) GetDoc(number uint32) (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("get-doc")
}

func (self *Umka) NewCheck(sessionId string, d *ru_nalog.Doc) (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("new-check")
}

func (self *Umka) Fiscalize(sessionId string, d *ru_nalog.Doc) (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("fiscalize")
}

func (self *Umka) CycleOpen() (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("cycleopen")
}

func (self *Umka) CycleClose() (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("cycleclose")
}

func (self *Umka) Danger_CloseFiscalStorage(sessionId string) (*ru_nalog.Doc, error) {
	return nil, errors.NotImplementedf("danger_closefiscalstorage")
}
