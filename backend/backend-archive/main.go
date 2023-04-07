package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/krishpranav/Mailtrix/cmd"
	sendmail "github.com/krishpranav/Mailtrix/sendmail/cmd"
)

func main() {
	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if normalize(filepath.Base(exec)) == normalize(filepath.Base(os.Args[0])) {
		cmd.Execute()
	} else {
		sendmail.Run()
	}
}

func normalize(s string) string {
	s = strings.ToLower(s)

	return strings.TrimSuffix(s, filepath.Ext(s))
}
