// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/tg123/go-htpasswd"
)

/**
 * @brief: Default Config Values For This Whole Project.
 */
var (
	SMTPListen = "0.0.0.0:1025"

	HTTPListen = "0.0.0.0:8025"

	DataFile string

	MaxMessages = 500

	UseMessageDates bool

	VerboseLogging = false

	QuietLogging = false

	NoLogging = false

	UITLSCert string

	UITLSKey string

	UIAuthFile string

	UIAuth *htpasswd.File

	Webroot = "/"

	SMTPTLSCert string

	SMTPTLSKey string

	SMTPAuthFile string

	SMTPAuth *htpasswd.File

	SMTPAuthAllowInsecure bool

	SMTPAuthAcceptAny bool

	SMTPCLITags string

	TagRegexp = regexp.MustCompile(`^([a-zA-Z0-9\-\ \_]){3,}$`)

	SMTPTags []Tag

	ContentSecurityPolicy = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; frame-src 'self'; img-src * data: blob:; font-src 'self' data:; media-src 'self'; connect-src 'self' ws: wss:; object-src 'none'; base-uri 'self';"

	Version = "dev"

	Repo = "krishpranav/Mailtrix"

	RepoBinaryName = "mailtrix"
)

/**
 * @brief: Tag Struct
 */
type Tag struct {
	Tag   string
	Match string
}

/**
 * @brief: Verify maildb Struct
 */
func VerifyConfig() error {
	if DataFile != "" && isDir(DataFile) {
		DataFile = filepath.Join(DataFile, "mailtrix.db")
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9\.\-]{3,}:\d{2,}$`)
	if !re.MatchString(SMTPListen) {
		return errors.New("SMTP bind should be in the format of <ip>:<port>")
	}
	if !re.MatchString(HTTPListen) {
		return errors.New("HTTP bind should be in the format of <ip>:<port>")
	}

	if UIAuthFile != "" {
		if !isFile(UIAuthFile) {
			return fmt.Errorf("HTTP password file not found: %s", UIAuthFile)
		}

		a, err := htpasswd.New(UIAuthFile, htpasswd.DefaultSystems, nil)
		if err != nil {
			return err
		}
		UIAuth = a
	}

	if UITLSCert != "" && UITLSKey == "" || UITLSCert == "" && UITLSKey != "" {
		return errors.New("You must provide both a UI TLS certificate and a key")
	}

	if UITLSCert != "" {
		if !isFile(UITLSCert) {
			return fmt.Errorf("TLS certificate not found: %s", UITLSCert)
		}

		if !isFile(UITLSKey) {
			return fmt.Errorf("TLS key not found: %s", UITLSKey)
		}
	}

	if SMTPTLSCert != "" && SMTPTLSKey == "" || SMTPTLSCert == "" && SMTPTLSKey != "" {
		return errors.New("You must provide both an SMTP TLS certificate and a key")
	}

	if SMTPTLSCert != "" {
		if !isFile(SMTPTLSCert) {
			return fmt.Errorf("SMTP TLS certificate not found: %s", SMTPTLSCert)
		}

		if !isFile(SMTPTLSKey) {
			return fmt.Errorf("SMTP TLS key not found: %s", SMTPTLSKey)
		}
	}

	if SMTPAuthFile != "" {
		if !isFile(SMTPAuthFile) {
			return fmt.Errorf("SMTP password file not found: %s", SMTPAuthFile)
		}

		if SMTPAuthAcceptAny {
			return errors.New("SMTP authentication can either use --smtp-auth-file or --smtp-auth-accept-any")
		}

		a, err := htpasswd.New(SMTPAuthFile, htpasswd.DefaultSystems, nil)
		if err != nil {
			return err
		}
		SMTPAuth = a
	}

	if SMTPTLSCert == "" && (SMTPAuthFile != "" || SMTPAuthAcceptAny) && !SMTPAuthAllowInsecure {
		return errors.New("SMTP authentication requires TLS encryption, run with `--smtp-auth-allow-insecure` to allow insecure authentication")
	}

	validWebrootRe := regexp.MustCompile(`[^0-9a-zA-Z\/\-\_\.]`)
	if validWebrootRe.MatchString(Webroot) {
		return fmt.Errorf("Invalid characters in Webroot (%s). Valid chars include: [a-z A-Z 0-9 _ . - /]", Webroot)
	}

	s := strings.TrimRight(path.Join("/", Webroot, "/"), "/") + "/"
	Webroot = s

	SMTPTags = []Tag{}

	p := shellwords.NewParser()

	if SMTPCLITags != "" {
		args, err := p.Parse(SMTPCLITags)
		if err != nil {
			return fmt.Errorf("Error parsing tags (%s)", err)
		}

		for _, a := range args {
			t := strings.Split(a, "=")
			if len(t) > 1 {
				tag := strings.TrimSpace(t[0])
				if !TagRegexp.MatchString(tag) || len(tag) == 0 {
					return fmt.Errorf("Invalid tag (%s) - can only contain spaces, letters, numbers, - & _", tag)
				}
				match := strings.TrimSpace(strings.ToLower(strings.Join(t[1:], "=")))
				if len(match) == 0 {
					return fmt.Errorf("Invalid tag match (%s) - no search detected", tag)
				}
				SMTPTags = append(SMTPTags, Tag{Tag: tag, Match: match})
			} else {
				return fmt.Errorf("Error parsing tags (%s)", a)
			}
		}

	}

	return nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.Mode().IsRegular() {
		return false
	}

	return true
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return false
	}

	return true
}
