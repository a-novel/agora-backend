package mailer

import (
	"github.com/rs/zerolog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

// Mailer is a quick configuration to send dynamic templates through Sendgrid.
//
// This interface is made for important emails directly related to in-configfiles experience (account validation, password
// reset...). It should not be used for email campaigns.
type Mailer interface {
	Send(recipient *mail.Email, templateID string, data map[string]interface{}) error
}

// NewMailer creates a new Mailer instance.
//
// The first argument is the Sendgrid API key. The second argument is the email address that will be used as sender.
// The third argument is a boolean that enables sandbox mode. In sandbox mode, emails are not sent, but printed in
// the console.
func NewMailer(apiKey string, sender *mail.Email, sandbox bool, logger zerolog.Logger) Mailer {
	return &mailerImpl{
		apiKey:  apiKey,
		sandbox: sandbox,
		from:    sender,
		logger:  logger,
	}
}

type mailerImpl struct {
	apiKey  string
	sandbox bool
	from    *mail.Email
	logger  zerolog.Logger
}

// Send emails a target recipient. It sends the content of the template identified by templateID, and fills template
// values with the content of the third argument.
func (mailer *mailerImpl) Send(recipient *mail.Email, templateID string, data map[string]interface{}) error {
	message := mail.NewV3Mail()
	personalization := mail.NewPersonalization()

	// Set senders and recipients.
	message.SetFrom(mailer.from)
	personalization.AddTos(recipient)

	// NOTE: for now, in-configfiles does not send any email campaign. The only email sent are for account management (such
	// as validating email or resetting password), so delivery is an absolute priority. It is thus acceptable to
	// bypass all inbox firewalls, especially since this method is only able to send one email at a time (which is
	// greatly inconvenient for mailing campaigns).
	// This setting SHOULD NOT be set if we eventually decide to send management campaigns from application.
	message.MailSettings = &mail.MailSettings{
		BypassListManagement: mail.NewSetting(true),
	}
	if mailer.sandbox {
		message.MailSettings.SandboxMode = &mail.Setting{Enable: &mailer.sandbox}
	}

	// Set email content.
	message.SetTemplateID(templateID)
	for k, v := range data {
		personalization.SetDynamicTemplateData(k, v)
	}
	message.AddPersonalizations(personalization)

	// Prepare request.
	rawMessage := mail.GetRequestBody(message)
	request := sendgrid.GetRequest(mailer.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = http.MethodPost
	request.Body = rawMessage

	// NOTE: emails are not sent in sandbox mode, so we are safe to use it for development (no risk to accidentally
	// ping a client).
	response, err := sendgrid.API(request)
	// Print email in sandbox mode, for testing and debugging.
	if mailer.sandbox {
		mailer.logger.
			Debug().
			Int("status", response.StatusCode).
			Str("recipient", recipient.Address).
			Str("alias", recipient.Name).
			Msgf(string(rawMessage))
	}

	return err
}
