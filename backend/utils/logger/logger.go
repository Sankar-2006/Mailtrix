// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package logger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/krishpranav/Mailtrix/config"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

// Log returns the logger instance
func Log() *logrus.Logger {
	if log == nil {
		log = logrus.New()
		log.SetLevel(logrus.InfoLevel)
		if config.VerboseLogging {
			// verbose logging (debug)
			log.SetLevel(logrus.DebugLevel)
		} else if config.QuietLogging {
			// show errors only
			log.SetLevel(logrus.ErrorLevel)
		} else if config.NoLogging {
			// disable all logging (tests)
			log.SetLevel(logrus.PanicLevel)
		}

		log.Out = os.Stdout
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006/01/02 15:04:05",
			ForceColors:     true,
		})
	}

	return log
}

// PrettyPrint for debugging
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
