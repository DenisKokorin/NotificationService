package models

type EmailMessage struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	TextBody    string   `json:"text_body"`
	HtmlBody    string   `json:"html_body"`
	Attachments []string `json:"attachments"`
}
