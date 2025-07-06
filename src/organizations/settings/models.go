package settings

import (
	"donation-mgmt/src/dal"
	"time"
)

type EmailProvider string

var (
	SMTPEmailProvider EmailProvider = "SMTP"
)

type EncryptedEmailProviderSettings struct {
	Provider      EmailProvider `json:"provider"`
	EncryptedSMTP *string       `json:"smtp"`
}

type EmailProviderSettings struct {
	Provider EmailProvider `json:"provider"`
	SMTP     *SMTPSettings `json:"smtp"`
}

type SMTPSettings struct {
	Host     string  `json:"host"`
	Port     uint16  `json:"port"`
	Username *string `json:"username"`
	Password *string `json:"password"`

	SenderEmail string `json:"senderEmail"`
}

type OrganizationSettings struct {
	OrganizationID int64
	Environment    dal.Environment
	Timezone       string
	IsValid        bool
	UpdatedAt      time.Time
}
