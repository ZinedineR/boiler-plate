package domain

const (
	MailTemplateCollectionName = "mailtemplate"
)

type MailTemplate struct {
	ID      string `bson:"_id" json:"id" validate:"required,uuid"`
	Name    string `bson:"name" json:"name" validate:"required"`
	Subject string `bson:"subject" json:"subject" validate:"required"`
	Body    string `bson:"body" json:"body" validate:"required,min=2"`
}

type MailTemplateList struct {
	ID   string `bson:"_id" json:"id" validate:"required,uuid"`
	Name string `bson:"name" json:"name" validate:"required"`
}

type HandlebarsDetail struct {
	ID         string   `bson:"_id" json:"id" validate:"required,uuid"`
	Name       string   `bson:"name" json:"name" validate:"required"`
	Handlebars []string `json:"handlebars"`
}

type UpsertMailTemplate struct {
	ID         string   `bson:"_id" json:"id"`
	MailLayout string   `json:"mail_layout"`
	Name       string   `bson:"name" json:"name"`
	Tags       []string `json:"tags"`
	Subject    string   `bson:"subject" json:"subject"`
	Body       string   `bson:"body" json:"body"`
}

func (model *MailTemplate) InitiateMailTemplate(name, subject, body string) {
	model.Name = name
	model.Subject = subject
	model.Body = body
}

func (model *MailTemplate) CollectionName() string {
	return MailTemplateCollectionName
}
