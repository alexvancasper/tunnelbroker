package executor

import (
	"text/template"

	"github.com/sirupsen/logrus"
)

type TunnelConfig struct {
	P2P        string
	PD         string
	IPv4Remote string
	IPv4Local  string
	TunnelName string
}

type Executor struct {
	Log *logrus.Logger
}

var TunnelMaker Executor

func (t TunnelConfig) ConfigureTunnel() error {
	l := TunnelMaker.Log.WithFields(logrus.Fields{
		"function": "getIPv6Prefix",
	})

	t.IPv4Local = "185.60.45.135"
	t.TunnelName = "TEST-TUNNEL"

	l.Infof("Start to provision the tunnel %+v", t)

	tmpl, err := template.ParseFiles("./internal/executor/templates/linux.tmpl")
	if err != nil {
		l.Errorf("template error %v", err)
		return err
	}
	err = tmpl.Execute(TunnelMaker.Log.Out, t)
	if err != nil {
		l.Errorf("template execute error %v", err)
		return nil
	}

	return nil
}
