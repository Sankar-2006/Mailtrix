// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package storage

import (
	"net/mail"
	"time"

	"github.com/jhillyerd/enmime"
)

type Message struct {
	ID          string
	Read        bool
	From        *mail.Address
	To          []*mail.Address
	Cc          []*mail.Address
	Bcc         []*mail.Address
	Subject     string
	Date        time.Time
	Tags        []string
	Text        string
	HTML        string
	Size        int
	Inline      []Attachment
	Attachments []Attachment
}

type Attachment struct {
	PartID      string
	FileName    string
	ContentType string
	ContentID   string
	Size        int
}

type MessageSummary struct {
	ID          string
	Read        bool
	From        *mail.Address
	To          []*mail.Address
	Cc          []*mail.Address
	Bcc         []*mail.Address
	Subject     string
	Created     time.Time
	Tags        []string
	Size        int
	Attachments int
}

type MailboxStats struct {
	Total  int
	Unread int
	Tags   []string
}

func AttachmentSummary(a *enmime.Part) Attachment {
	o := Attachment{}
	o.PartID = a.PartID
	o.FileName = a.FileName
	if o.FileName == "" {
		o.FileName = a.ContentID
	}
	o.ContentType = a.ContentType
	o.ContentID = a.ContentID
	o.Size = len(a.Content)

	return o
}
