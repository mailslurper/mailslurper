package api

import "time"

type MailListItemDTO struct {
	ID              string   `json:"id"`
	DateSent        string   `json:"dateSent"`
	FromAddress     string   `json:"fromAddress"`
	ToAddresses     []string `json:"toAddresses"`
	Subject         string   `json:"subject"`
	XMailer         string   `json:"xMailer"`
	ContentType     string   `json:"contentType"`
	AttachmentCount int      `json:"attachmentCount"`
}

type MailListResponse struct {
	MailItems    []MailListItemDTO `json:"mailItems"`
	Page         int               `json:"page"`
	PageSize     int               `json:"pageSize"`
	TotalRecords int               `json:"totalRecords"`
	TotalPages   int               `json:"totalPages"`
}

type AttachmentDTO struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	SizeBytes   int    `json:"sizeBytes"`
}

type MailDetailDTO struct {
	ID          string          `json:"id"`
	DateSent    string          `json:"dateSent"`
	FromAddress string          `json:"fromAddress"`
	ToAddresses []string        `json:"toAddresses"`
	Subject     string          `json:"subject"`
	XMailer     string          `json:"xMailer"`
	ContentType string          `json:"contentType"`
	TextBody    string          `json:"textBody"`
	HTMLBody    string          `json:"htmlBody"`
	Attachments []AttachmentDTO `json:"attachments"`
}

type CountResponse struct {
	Count int `json:"count"`
}

type PruneOptionDTO struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PruneRequest struct {
	Code string `json:"code"`
}

type PruneResponse struct {
	DeletedCount int64 `json:"deletedCount"`
}

type SettingsResponse struct {
	AuthenticationScheme string `json:"authenticationScheme"`
	ServerVersion        string `json:"serverVersion"`
	PublicURL            string `json:"publicURL"`
}

type VersionResponse struct {
	Version string `json:"version"`
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
