package emailer

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
	"github.com/mailgun/mailgun-go"
)

var (
	name = "emailer"
)

// Emailer sends email
type Emailer struct {
	name                string
	app                 string
	mailgun             mailgun.Mailgun
	MailgunDomain       string
	MailgunAPIKey       string
	MailgunPublicAPIKey string
	MailgunSender       string
	MailgunReceiver     []string
	MailgunSubject      string
}

// Option for emailer
type Option struct {
	App                 string
	MailgunDomain       string
	MailgunAPIKey       string
	MailgunPublicAPIKey string
	MailgunSender       string
	MailgunReceiver     []string
	MailgunSubject      string
}

// NewEmailer returns emailer
func NewEmailer(opt *Option) (*Emailer, error) {
	m := new(Emailer)
	m.name = name
	m.app = opt.App
	m.MailgunDomain = opt.MailgunDomain
	m.MailgunAPIKey = opt.MailgunAPIKey
	m.MailgunPublicAPIKey = opt.MailgunPublicAPIKey
	m.mailgun = mailgun.NewMailgun(m.MailgunDomain, m.MailgunAPIKey, m.MailgunPublicAPIKey)
	m.MailgunSender = opt.MailgunSender
	m.MailgunReceiver = opt.MailgunReceiver
	m.MailgunSubject = opt.MailgunSubject
	return m, nil
}

// NewOption for emailer
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	sec := config.Section(fmt.Sprintf("%s_%s", app, name))
	opt.MailgunDomain = sec.Key("mailgun_domain").String()
	opt.MailgunAPIKey = sec.Key("mailgun_api_key").String()
	opt.MailgunPublicAPIKey = sec.Key("mailgun_public_api_key").String()
	opt.MailgunSender = sec.Key("mailgun_sender").String()
	opt.MailgunReceiver = make([]string, 0)
	ss := strings.Split(sec.Key("mailgun_receiver").String(), ",")
	for _, s := range ss {
		e := strings.TrimSpace(s)
		if e != "" {
			opt.MailgunReceiver = append(opt.MailgunReceiver, e)
		}
	}
	opt.MailgunSubject = sec.Key("mailgun_subject").String()
	return opt, nil
}

// Execute ...
func (m *Emailer) Execute(msg *message.Message) error {
	logrus.Debug("emailer.Execute")
	if !msg.OK {
		for _, receiver := range m.MailgunReceiver {
			mail := m.mailgun.NewMessage(m.MailgunSender, m.MailgunSubject, msg.Err.Error(), receiver)
			_, _, err := m.mailgun.Send(mail)
			if err != nil {
				return errors.Trace(err)
			}
		}
	}
	return nil
}

// Name returns plugin name
func (m *Emailer) Name() string {
	return m.name
}

// App returns app name
func (m *Emailer) App() string {
	return m.app
}
