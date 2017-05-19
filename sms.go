package zang

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
	"net/url"
)

// PhoneNumber represents a phone number.
type PhoneNumber string
func (p PhoneNumber) toString() string { return string(p)}

// Message represents
type Message struct {
	From PhoneNumber
	To PhoneNumber
	Body string
}

// Response is the object returned from a Sender when a message has been sent.
type Response struct {
	APIVersion string `json:"api_version"`
	Sid string `json:"sid"`
	AccountSid string `json:"account_sid"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	DateSent string `json:"date_sent"`
	To string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
	Status string `json:"status"`
	Direction string `json:"direction"`
	Price string `json:"price"`
	URI string `json:"uri"`
}

// Sender represents something that can send sms messages.
type Sender interface {
	Send(Message) (*Response, error)
}

// NewClient returns a Sender capable of sending messages.
func NewClient(sid, token string) *ZangSender {
	return &ZangSender{
		AccountSID: sid,
		AuthToken: token,
		client: http.DefaultClient,
	}
}

// ZangSender is used to send messages.
type ZangSender struct {
	AccountSID string
	AuthToken string
	client *http.Client
}

// WithClient returns a new ZangSender with the specified key.
func (z *ZangSender) WithClient(client *http.Client) *ZangSender {
	return &ZangSender{
		client: client,
		AuthToken: z.AuthToken,
		AccountSID: z.AccountSID,
	}
}

// Send sends a message.
func (z *ZangSender) Send(message Message) (*Response, error) {
	request, err := z.makeRequest(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}


	response, err := z.client.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}


	zResponse := &Response{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(zResponse); err != nil {
		return nil, errors.WithStack(err)
	}
	return zResponse, nil
}

func (z *ZangSender) makeRequest(message Message) (*http.Request, error) {
	endpoint := fmt.Sprintf("https://api.zang.io/v2/Accounts/%s/SMS/Messages.json", z.AccountSID)

	form := url.Values{}
	form.Add("From", message.From.toString())
	form.Add("To", message.To.toString())
	form.Add("Body", message.Body)

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request.SetBasicAuth(z.AccountSID, z.AuthToken)

	return request, nil
}
