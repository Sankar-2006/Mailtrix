// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package storage

import (
	"regexp"
	"strings"

	"github.com/leporo/sqlf"
)

func searchParser(args []string, start, limit int) *sqlf.Stmt {
	if limit == 0 {
		limit = 50
	}

	q := sqlf.From("mailbox").
		Select(`ID, Data, Tags, Read,
			json_extract(Data, '$.To') as ToJSON,
			json_extract(Data, '$.From') as FromJSON,
			IFNULL(json_extract(Data, '$.Cc'), '{}') as CcJSON,
			IFNULL(json_extract(Data, '$.Bcc'), '{}') as BccJSON,
			json_extract(Data, '$.Subject') as Subject,
			json_extract(Data, '$.Attachments') as Attachments
		`).
		OrderBy("Sort DESC").
		Limit(limit).
		Offset(start)

	if limit > 0 {
		q = q.Limit(limit)
	}

	for _, w := range args {
		if cleanString(w) == "" {
			continue
		}

		exclude := false
		if len(w) > 1 && (strings.HasPrefix(w, "-") || strings.HasPrefix(w, "!")) {
			exclude = true
			w = w[1:]
		}

		re := regexp.MustCompile(`[a-zA-Z0-9]+`)
		if !re.MatchString(w) {
			continue
		}

		if strings.HasPrefix(w, "to:") {
			w = cleanString(w[3:])
			if w != "" {
				if exclude {
					q.Where("ToJSON NOT LIKE ?", "%"+escPercentChar(w)+"%")
				} else {
					q.Where("ToJSON LIKE ?", "%"+escPercentChar(w)+"%")
				}
			}
		} else if strings.HasPrefix(w, "from:") {
			w = cleanString(w[5:])
			if w != "" {
				if exclude {
					q.Where("FromJSON NOT LIKE ?", "%"+escPercentChar(w)+"%")
				} else {
					q.Where("FromJSON LIKE ?", "%"+escPercentChar(w)+"%")
				}
			}
		} else if strings.HasPrefix(w, "cc:") {
			w = cleanString(w[3:])
			if w != "" {
				if exclude {
					q.Where("CcJSON NOT LIKE ?", "%"+escPercentChar(w)+"%")
				} else {
					q.Where("CcJSON LIKE ?", "%"+escPercentChar(w)+"%")
				}
			}
		} else if strings.HasPrefix(w, "bcc:") {
			w = cleanString(w[4:])
			if w != "" {
				if exclude {
					q.Where("BccJSON NOT LIKE ?", "%"+escPercentChar(w)+"%")
				} else {
					q.Where("BccJSON LIKE ?", "%"+escPercentChar(w)+"%")
				}
			}
		} else if strings.HasPrefix(w, "subject:") {
			w = cleanString(w[8:])
			if w != "" {
				if exclude {
					q.Where("Subject NOT LIKE ?", "%"+escPercentChar(w)+"%")
				} else {
					q.Where("Subject LIKE ?", "%"+escPercentChar(w)+"%")
				}
			}
		} else if strings.HasPrefix(w, "tag:") {
			w = cleanString(w[4:])
			if w != "" {
				if exclude {
					q.Where("Tags NOT LIKE ?", "%\""+escPercentChar(w)+"\"%")
				} else {
					q.Where("Tags LIKE ?", "%\""+escPercentChar(w)+"\"%")
				}
			}
		} else if w == "is:read" {
			if exclude {
				q.Where("Read = 0")
			} else {
				q.Where("Read = 1")
			}
		} else if w == "is:unread" {
			if exclude {
				q.Where("Read = 1")
			} else {
				q.Where("Read = 0")
			}
		} else if w == "has:attachment" || w == "has:attachments" {
			if exclude {
				q.Where("Attachments = 0")
			} else {
				q.Where("Attachments > 0")
			}
		} else {
			if exclude {
				q.Where("search NOT LIKE ?", "%"+cleanString(escPercentChar(w))+"%")
			} else {
				q.Where("search LIKE ?", "%"+cleanString(escPercentChar(w))+"%")
			}
		}
	}

	return q
}
