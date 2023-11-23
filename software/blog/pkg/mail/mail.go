package mail

import (
	"blog/global"
	"blog/templates"
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
)

// Options database option
type Options struct {
	Host     string `toml:"host" json:"host" yaml:"host" env:"SMTP_HOST"`
	Port     int    `toml:"port" json:"port" yaml:"port" env:"SMTP_PORT"`
	Username string `toml:"username" json:"username" yaml:"username" env:"SMTP_USERNAME"`
	Password string `toml:"password" json:"password" yaml:"password" env:"SMTP_PASSWORD"`
}

type Client struct {
	Username string
	*gomail.Dialer
}

// NewOptions for redis
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("email", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}

	logger.Info("load email options success", zap.Any("email options", o))
	return o, err
}

// New redis pool conn
func New(o *Options) (*Client, error) {
	d := gomail.NewDialer(
		o.Host,
		o.Port,
		o.Username,
		o.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Client{Username: o.Username, Dialer: d}, nil
}

func (c Client) sendMail(to, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Artemis <%s>", c.Username))
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	return c.DialAndSend(m)
}

func (c Client) sendTemplateMail(email, title string, templateFS embed.FS, templatePath string, params map[string]string) error {
	var content bytes.Buffer
	t, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		return errors.Wrap(err, "parse template file")
	}
	if err := t.Execute(&content, params); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return c.sendMail(email, title, content.String())
}

func (c Client) SendNewAnswerMail(email, domain string, articleID int64, article, comment string) error {
	params := map[string]string{
		"link":    fmt.Sprintf("%s/%s/%d", global.WEBSITE, domain, articleID),
		"article": article,
		"comment": comment,
	}
	return c.sendTemplateMail(email, "【Artemis】您的提问有了回复", templates.FS, "mail/new-comment.html", params)
}
