package mailer

import (
	"context"
	"embed"
	"time"
)

const (
	fromName               = "Gopher Social"
	initialBackoff         = 200 * time.Millisecond
	maxBackoff             = 2 * time.Second
	maxSendAttempts        = 5
	UserInvitationTemplate = "user_invitation.tmpl"
)

//go:embed templates/*
var FS embed.FS

type Mailer interface {
	Send(
		ctx context.Context,
		templateFile string,
		mail Mail,
		isSandbox bool, // For development purposes, if true, the mail will not be sent
	) error
}

type Mail struct {
	To      []To
	Subject string
	Args    any
}

type To struct {
	Name  string
	Email string
}

type InvitationEmailData struct {
	To            To
	ActivationURL string
}
