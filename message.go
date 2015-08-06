package basin

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Message interface {
	Field(string) (interface{}, bool)
	Bytes() []byte
}

type Logplex1 struct {
	Pri       int
	Version   int
	Timestamp time.Time
	Hostname  string
	AppName   string
	ProcId    string
	MsgId     string
	Message   string
}

func (p Logplex1) Facility() int {
	return p.Pri / 8
}

func (p Logplex1) Severity() int {
	return p.Pri % 8
}

func (p Logplex1) Field(f string) (interface{}, bool) {
	switch f {
	case "Pri":
		return p.Pri, true
	case "Version":
		return p.Version, true
	case "Timestamp":
		return p.Timestamp, true
	case "Hostname":
		return p.Hostname, true
	case "AppName":
		return p.AppName, true
	case "ProcId":
		return p.ProcId, true
	case "MsgId":
		return p.MsgId, true
	case "Message":
		return p.Message, true
	default:
		return "", false
	}
}

func (p Logplex1) Bytes() []byte {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("<%d>%d %s %s %s %s %s %s",
		p.Pri, p.Version, p.Timestamp.Format(time.RFC3339),
		maybeNil(p.Hostname, "-"), maybeNil(p.AppName, "-"),
		maybeNil(p.ProcId, "-"), maybeNil(p.MsgId, "-"),
		p.Message))
	return b.Bytes()
}

type RFC5424Params map[string]string

type RFC5424 struct {
	Pri       int
	Version   int
	Timestamp time.Time
	Hostname  string
	AppName   string
	ProcId    string
	MsgId     string
	Elements  map[string]RFC5424Params
	Message   string
}

func (p RFC5424) Facility() int {
	return p.Pri / 8
}

func (p RFC5424) Severity() int {
	return p.Pri % 8
}

func (p RFC5424) Field(f string) (interface{}, bool) {
	switch f {
	case "Pri":
		return p.Pri, true
	case "Version":
		return p.Version, true
	case "Timestamp":
		return p.Timestamp, true
	case "Hostname":
		return p.Hostname, true
	case "AppName":
		return p.AppName, true
	case "ProcId":
		return p.ProcId, true
	case "MsgId":
		return p.MsgId, true
	case "Message":
		return p.Message, true
	default:
		data, exists := p.Elements[f]
		return data, exists
	}
}

func (p RFC5424) Bytes() []byte {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("<%d>%d %s %s %s %s %s",
		p.Pri, p.Version, p.Timestamp.Format(time.RFC3339),
		maybeNil(p.Hostname, "-"), maybeNil(p.AppName, "-"),
		maybeNil(p.ProcId, "-"), maybeNil(p.MsgId, "-")))

	for id, params := range p.Elements {
		b.WriteString("[" + id)
		for name, value := range params {
			b.WriteString(fmt.Sprintf(" %s=\"%s\"", name, escape(value, []string{"\\", "=", "]"})))
		}
	}
	if len(p.Elements) > 0 {

	}

	return b.Bytes()
}

func maybeNil(val string, def string) string {
	if val == "" {
		return def
	}
	return val
}

func escape(val string, es []string) string {
	for _, e := range es {
		val = strings.Replace(val, e, "\\"+e, -1)
	}
	return val
}
