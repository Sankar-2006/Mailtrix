// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package storage

import (
	"context"
	"database/sql"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jhillyerd/enmime"
	"github.com/k3a/html2text"
	"github.com/krishpranav/Mailtrix/config"
	"github.com/krishpranav/Mailtrix/server/websockets"
	"github.com/krishpranav/Mailtrix/utils/logger"
	"github.com/leporo/sqlf"
)

func addressToSlice(env *enmime.Envelope, key string) []*mail.Address {
	data, err := env.AddressList(key)
	if err != nil || data == nil {
		return []*mail.Address{}
	}

	return data
}

func createSearchText(env *enmime.Envelope) string {
	var b strings.Builder

	b.WriteString(env.GetHeader("From") + " ")
	b.WriteString(env.GetHeader("Subject") + " ")
	b.WriteString(env.GetHeader("To") + " ")
	b.WriteString(env.GetHeader("Cc") + " ")
	b.WriteString(env.GetHeader("Bcc") + " ")
	h := strings.TrimSpace(
		html2text.HTML2TextWithOptions(
			env.HTML,
			html2text.WithLinksInnerText(),
		),
	)
	if h != "" {
		b.WriteString(h + " ")
	} else {
		b.WriteString(env.Text + " ")
	}
	for _, a := range env.Attachments {
		b.WriteString(a.FileName + " ")
	}

	d := cleanString(b.String())

	return d
}

func cleanString(str string) string {
	re := regexp.MustCompile(`(\r?\n|\t|>|<|"|\,|;)`)
	str = re.ReplaceAllString(str, " ")

	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(str)), " "))
}

func dbCron() {
	for {
		time.Sleep(60 * time.Second)
		start := time.Now()

		currentTime := time.Now()
		diff := currentTime.Sub(dbLastAction)
		if dbDataDeleted && diff.Minutes() > 5 {
			dbDataDeleted = false
			_, err := db.Exec("VACUUM")
			if err == nil {
				elapsed := time.Since(start)
				logger.Log().Debugf("[db] compressed idle database in %s", elapsed)
			}
			continue
		}

		if config.MaxMessages > 0 {
			q := sqlf.Select("ID").
				From("mailbox").
				OrderBy("Sort DESC").
				Limit(5000).
				Offset(config.MaxMessages)

			ids := []string{}
			if err := q.Query(nil, db, func(row *sql.Rows) {
				var id string

				if err := row.Scan(&id); err != nil {
					logger.Log().Errorf("[db] %s", err.Error())
					return
				}
				ids = append(ids, id)

			}); err != nil {
				logger.Log().Errorf("[db] %s", err.Error())
				continue
			}

			if len(ids) == 0 {
				continue
			}

			tx, err := db.BeginTx(context.Background(), nil)
			if err != nil {
				logger.Log().Errorf("[db] %s", err.Error())
				continue
			}

			args := make([]interface{}, len(ids))
			for i, id := range ids {
				args[i] = id
			}

			_, err = tx.Query(`DELETE FROM mailbox WHERE ID IN (?`+strings.Repeat(",?", len(ids)-1)+`)`, args...)
			if err != nil {
				logger.Log().Errorf("[db] %s", err.Error())
				continue
			}

			_, err = tx.Query(`DELETE FROM mailbox_data WHERE ID IN (?`+strings.Repeat(",?", len(ids)-1)+`)`, args...)
			if err != nil {
				logger.Log().Errorf("[db] %s", err.Error())
				continue
			}

			err = tx.Commit()

			if err != nil {
				logger.Log().Errorf(err.Error())
				if err := tx.Rollback(); err != nil {
					logger.Log().Errorf(err.Error())
				}
			}

			dbDataDeleted = true

			elapsed := time.Since(start)
			logger.Log().Debugf("[db] auto-pruned %d messages in %s", len(ids), elapsed)

			websockets.Broadcast("prune", nil)
		}
	}
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.Mode().IsRegular() {
		return false
	}

	return true
}

func inArray(k string, arr []string) bool {
	k = strings.ToLower(k)
	for _, v := range arr {
		if strings.ToLower(v) == k {
			return true
		}
	}

	return false
}

func escPercentChar(s string) string {
	return strings.ReplaceAll(s, "%", "%%")
}
