package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

var (
	host       = os.Getenv("EMAIL_HOST")
	username   = os.Getenv("EMAIL_USERNAME")
	password   = os.Getenv("EMAIL_PASSWORD")
	portNumber = os.Getenv("EMAIL_PORT")
)

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To             []string
	Subject        string
	Body           string
	Attachments    io.Reader
	AttachmentName string
}

func NewEmailSender() *Sender {
	auth := smtp.PlainAuth("", username, password, host)
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.auth, username, m.To, m.WriteEmailContent())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b}
}

// WriteEmailContent write the header, body and file attachments to a bytes.Buffer and return a slice of byte
func (m *Message) WriteEmailContent() []byte {
	// creates a bytes.Buffer and read from io.Reader
	buf := &bytes.Buffer{}

	if _, err := buf.ReadFrom(m.Attachments); err != nil {
		return []byte{}
	}

	// retrieve a byte slice from bytes.Buffer
	data := buf.Bytes()

	content := bytes.NewBuffer(nil)

	// write to header
	content.WriteString(fmt.Sprintf("From: %s\n", username))
	content.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ", ")))
	content.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))

	content.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(content)
	boundary := writer.Boundary()

	content.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	content.WriteString(fmt.Sprintf("--%s\n", boundary))

	// write to body
	content.WriteString(m.Body)

	// attach csv file
	content.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
	content.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(data)))
	content.WriteString("Content-Transfer-Encoding: base64\n")
	content.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", m.AttachmentName))

	// for multipart MIME messages
	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encodedBytes, data)
	content.Write(encodedBytes)
	content.WriteString(fmt.Sprintf("\n--%s", boundary))

	return content.Bytes()
}
