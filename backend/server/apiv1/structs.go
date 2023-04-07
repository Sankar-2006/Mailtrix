package apiv1

import "github.com/krishpranav/Mailtrix/storage"

type MessageSummary = storage.MessageSummary

type MessagesSummary struct {
	Total    int              `json:"total"`
	Unread   int              `json:"unread"`
	Count    int              `json:"count"`
	Start    int              `json:"start"`
	Tags     []string         `json:"tags"`
	Messages []MessageSummary `json:"messages"`
}

type Message = storage.Message

type Attachment = storage.Attachment
