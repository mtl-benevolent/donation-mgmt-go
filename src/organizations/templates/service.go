package templates

import (
	"context"
	"donation-mgmt/src/dal"
	"fmt"
)

type EmailTemplateType string

var (
	NoMailingAddrEmailTemplate         EmailTemplateType = "no_mailing_addr_email"
	NoMailingAddrReminderEmailTemplate EmailTemplateType = "no_mailing_addr_reminder_email"
	ReceiptEmailTemplate               EmailTemplateType = "receipt_email"
)

type TemplateType string

var (
	ReceiptPDFTemplate TemplateType = "receipt_pdf"
)

type UpdateEmailTemplateParams struct {
	OrgID        int64
	Environment  dal.Environment
	TemplateType EmailTemplateType

	TemplateTitle   *string
	TemplateContent *string
}

type UpdateTemplateParams struct {
	OrgID       int64
	Environment dal.Environment

	TemplateType    TemplateType
	TemplateContent *string
}

type OrgTemplatesService struct {
}

func NewOrgTemplatesService() *OrgTemplatesService {
	return &OrgTemplatesService{}
}

func (s *OrgTemplatesService) UpdateEmailTemplate(ctx context.Context, querier dal.Querier, params UpdateEmailTemplateParams) (dal.OrganizationTemplate, error) {
	switch params.TemplateType {
	case NoMailingAddrEmailTemplate:
		return querier.UpsertNoMailingAddrEmailTemplate(ctx, dal.UpsertNoMailingAddrEmailTemplateParams{
			OrganizationID:          params.OrgID,
			Environment:             params.Environment,
			NoMailingAddrEmailTitle: params.TemplateTitle,
			NoMailingAddrEmail:      params.TemplateContent,
		})
	case NoMailingAddrReminderEmailTemplate:
		return querier.UpsertNoMailingAddrReminderEmailTemplate(ctx, dal.UpsertNoMailingAddrReminderEmailTemplateParams{
			OrganizationID:          params.OrgID,
			Environment:             params.Environment,
			NoMailingAddrReminderEmailTitle: params.TemplateTitle,
			NoMailingAddrReminderEmail:      params.TemplateContent,
		})
	case ReceiptEmailTemplate:
		return querier.UpsertReceiptEmailTemplate(ctx, querier, dal.UpsertReceiptEmailTemplateParams{
			OrganizationID: params.OrgID,
			Environment: params.Environment
			ReceiptEmailTitle: params.TemplateTitle,
			ReceiptEmail: params.TemplateContent,
		})
	default:
		return dal.OrganizationTemplate{}, fmt.Errorf("invalid template type: %s", params.TemplateType)
	}
}

func (s *OrgTemplatesService) UpdateTemplate(ctx context.Context, querier dal.Querier, params UpdateTemplateParams) (dal.OrganizationTemplate, error) {
	switch params.TemplateType {
	case ReceiptPDFTemplate:
		return querier.UpsertReceiptPDFTemplate(ctx, dal.UpsertReceiptPDFTemplateParams{
			OrganizationID: params.OrgID,
			Environment: params.Environment,
			ReceiptPdf: params.TemplateContent,
		})
	default:
		return dal.OrganizationTemplate{}, fmt.Errorf("invalid template type: %s", params.TemplateType)
	}
}
