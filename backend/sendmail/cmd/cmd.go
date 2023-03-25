package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"

	"github.com/krishpranav/Mailtrix/utils/logger"
	flag "github.com/spf13/pflag"
)

func Run() {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	username := "nobody"
	user, err := user.Current()
	if err == nil && user != nil && len(user.Username) > 0 {
		username = user.Username
	}

	fromAddr := username + "@" + host
	smtpAddr := "localhost:1025"
	var recip []string

	if len(os.Getenv("MP_SENDMAIL_SMTP_ADDR")) > 0 {
		smtpAddr = os.Getenv("MP_SENDMAIL_SMTP_ADDR")
	}
	if len(os.Getenv("MP_SENDMAIL_FROM")) > 0 {
		fromAddr = os.Getenv("MP_SENDMAIL_FROM")
	}

	var verbose bool

	flag.StringVar(&smtpAddr, "smtp-addr", smtpAddr, "SMTP server address")
	flag.StringVarP(&fromAddr, "from", "f", fromAddr, "SMTP sender")
	flag.BoolP("long-i", "i", true, "Ignored. This flag exists for sendmail compatibility.")
	flag.BoolP("long-t", "t", true, "Ignored. This flag exists for sendmail compatibility.")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose mode (sends debug output to stderr)")
	flag.Parse()

	recip = flag.Args()

	if verbose {
		fmt.Fprintln(os.Stderr, smtpAddr, fromAddr)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin")
		os.Exit(11)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("error parsing message body: %s", err))
		os.Exit(11)
	}

	if len(recip) == 0 {
		recip = append(recip, msg.Header.Get("To"))
	}

	err = smtp.SendMail(smtpAddr, nil, fromAddr, recip, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error sending mail")
		logger.Log().Fatal(err)
	}
}
