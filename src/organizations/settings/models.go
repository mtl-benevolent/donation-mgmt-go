package settings

import (
	"donation-mgmt/src/dal"
	"time"
)

type EmailProviderSettings struct {
	SMTP *SMTPSettings `json:"smtp"`
}

type SMTPSettings struct {
	Host     string  `json:"host"`
	Port     uint16  `json:"port"`
	Username *string `json:"username"`
	Password *string `json:"password"`

	SenderEmail string `json:"senderEmail"`
}

type OrganizationSettings struct {
	OrganizationID        int64
	Environment           dal.Environment
	Timezone              string
	EmailProvider         dal.EmailProvider
	EmailProviderSettings EmailProviderSettings
	IsValid               bool
	UpdatedAt             time.Time
}
