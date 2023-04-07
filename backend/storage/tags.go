// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package storage

import (
	"context"
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/krishpranav/Mailtrix/config"
	"github.com/krishpranav/Mailtrix/utils/logger"
	"github.com/leporo/sqlf"
)

func SetTags(id string, tags []string) error {
	applyTags := []string{}
	reg := regexp.MustCompile(`\s+`)
	for _, t := range tags {
		t = strings.TrimSpace(reg.ReplaceAllString(t, " "))

		if t != "" && config.TagRegexp.MatchString(t) && !inArray(t, applyTags) {
			applyTags = append(applyTags, t)
		}
	}

	tagJSON, err := json.Marshal(applyTags)
	if err != nil {
		logger.Log().Errorf("[db] setting tags for message %s", id)
		return err
	}

	_, err = sqlf.Update("mailbox").
		Set("Tags", string(tagJSON)).
		Where("ID = ?", id).
		ExecAndClose(context.Background(), db)

	if err == nil {
		logger.Log().Debugf("[db] set tags %s for message %s", string(tagJSON), id)
	}

	return err
}

func findTags(message *[]byte) []string {
	tags := []string{}
	if len(config.SMTPTags) == 0 {
		return tags
	}

	str := strings.ToLower(string(*message))
	for _, t := range config.SMTPTags {
		if !inArray(t.Tag, tags) && strings.Contains(str, t.Match) {
			tags = append(tags, t.Tag)
		}
	}

	sort.Strings(tags)

	return tags
}

func getMessageTags(id string) []string {
	tags := []string{}
	var data string

	q := sqlf.From("mailbox").
		Select(`Tags`).To(&data).
		Where(`ID = ?`, id)

	err := q.QueryRowAndClose(context.Background(), db)
	if err != nil {
		logger.Log().Error(err)
		return tags
	}

	if err := json.Unmarshal([]byte(data), &tags); err != nil {
		logger.Log().Error(err)
		return tags
	}

	return tags
}
