package umka

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	ru_nalog "github.com/temoto/ru-nalog-go"
)

type responseDoc struct {
	Protocol      int    `json:"protocol"` // 1=JSON 3=XML
	Version       string `json:"version"`  // "1.0"
	CashboxStatus Status `json:"cashboxStatus"`
	Document      struct {
		Result  int      `json:"result"` // 0=OK
		Message struct { // when result!=0
			Description string `json:"resultDescription"`
		} `json:"message"`
		Data docdata `json:"data"`
	} `json:"document"`
}

type docdata struct {
	DocNumber uint32           `json:"docNumber"`
	DocType   ru_nalog.DocType `json:"docType"`
	Name      string           `json:"name"`
	Props     []prop           `json:"fiscprops"`
}
type prop struct {
	Caption   string       `json:"caption"`
	Printable string       `json:"printable"`
	Props     []prop       `json:"fiscprops"`
	Tag       ru_nalog.Tag `json:"tag"`
	Value     interface{}  `json:"value"`
}

func ParseResponseDoc(b []byte) (*ru_nalog.Doc, error) {
	var rd responseDoc
	if err := json.Unmarshal(b, &rd); err != nil {
		err = errors.Trace(err)
		return nil, err
	}

	if !(rd.Protocol == 1 && rd.Version == "1.0") {
		return nil, errors.Errorf("unknown protocol=%d version='%s'", rd.Protocol, rd.Version)
	}

	if rd.Document.Result != 0 {
		return nil, errors.Errorf("result=%d message=%s", rd.Document.Result, rd.Document.Message.Description)
	}

	return rd.Document.Data.ToDoc()
}

func (d *docdata) ToDoc() (*ru_nalog.Doc, error) {
	fd := ru_nalog.NewDoc(d.DocNumber, d.DocType)
	errs := make([]error, 0, 8)
	for _, p := range d.Props {
		t, err := p.ToTLV()
		if err == nil {
			fd.Props.Append(t)
		} else {
			errs = append(errs, err)
		}
	}
	return fd, foldErrors(errs)
}

func (p *prop) ToTLV() (*ru_nalog.TLV, error) {
	// log.Printf("prop=%#v", p)
	t := ru_nalog.NewTLV(p.Tag)
	if t == nil {
		return nil, errors.Errorf("prop=%#v invalid tag", p)
	}
	switch t.Kind {
	case ru_nalog.DataKindSTLV:
		for _, child := range p.Props {
			subt, err := child.ToTLV()
			if err != nil {
				err = errors.Annotatef(err, "prop=%#v", p)
				return nil, err
			}
			t.Append(subt)
		}
	case ru_nalog.DataKindTime:
		crude := p.Value.(string)
		tim, err := time.Parse(TimeLayout, crude)
		if err != nil {
			err = errors.Annotatef(err, "prop=%#v", p)
			return nil, err
		}
		t.SetValue(tim)
	case ru_nalog.DataKindUint:
		switch p.Value.(type) {
		case bool:
			if p.Value.(bool) {
				t.SetValue(uint32(1))
			} else {
				t.SetValue(uint32(0))
			}
		case float64: // encoding/json approach to unmarshal JSON integer token into interface{}
			crude := p.Value.(float64)
			t.SetValue(uint32(crude))
		default:
			panic("TODO")
		}
	case ru_nalog.DataKindBytes:
		switch p.Value.(type) {
		case string:
			t.SetValue(p.Value)
		case float64: // umka joke on byte[fixed] (1077), forced print format
			crude := p.Value.(float64)
			// TODO FIXME
			t.SetValue(uint64(crude))
		default:
			panic("TODO")
		}
	case ru_nalog.DataKindVLN:
		switch p.Value.(type) {
		case float64: // encoding/json approach to unmarshal JSON integer token into interface{}
			crude := p.Value.(float64)
			t.SetValue(uint64(crude))
		default:
			panic("TODO")
		}
	case ru_nalog.DataKindFVLN:
		switch p.Value.(type) {
		case string: // umka joke on FVLN (1023), forced print format
			crude := p.Value.(string)
			// "1 333,500" -> 1333500
			crude = strings.Replace(crude, " ", "", -1)
			crude = strings.Replace(crude, ",", ".", 1)
			t.SetValue(crude)
			if err := t.Err(); err != nil {
				err = errors.Annotatef(err, "prop=%#v", p)
				return nil, err
			}
		default:
			panic("TODO")
		}
	default:
		t.SetValue(p.Value)
	}
	t.Caption = p.Caption
	t.Printable = p.Printable
	return t, nil
}

func foldErrors(errs []error) error {
	// common fast path
	if len(errs) == 0 {
		return nil
	}

	ss := make([]string, 0, 1+len(errs))
	for _, e := range errs {
		if e != nil {
			// ss = append(ss, e.Error())
			ss = append(ss, errors.ErrorStack(e))
			// ss = append(ss, errors.Details(e))
		}
	}
	switch len(ss) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf(ss[0])
	default:
		ss = append(ss, "")
		copy(ss[1:], ss[0:])
		ss[0] = "multiple errors:"
		return fmt.Errorf(strings.Join(ss, "\n- "))
	}
}
