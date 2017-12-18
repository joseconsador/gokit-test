package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/benmanns/goworker"
	"github.com/subosito/gotenv"
)

type (
	// TicketComment gago
	TicketComment struct {
		Body string `json:"body"`
	}

	// Ticket struct
	Ticket struct {
		Subject string        `json:"subject"`
		Status  string        `json:"status"`
		Comment TicketComment `json:"comment"`
		URL     string        `json:"url"`
	}

	// WeirdTicketWrapper struct
	WeirdTicketWrapper struct {
		Ticket Ticket `json:"ticket"`
	}

	// SlackMessage gago eh
	SlackMessage struct {
		Text        string                   `json:"text"`
		Attachments []SlackMessageAttachment `json:"attachments"`
	}

	// SlackMessageAttachment gago e
	SlackMessageAttachment struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	SlackCreds struct {
		Webhook     string
		AccessToken string
	}

	ZDCreds struct {
		Subdomain string
		Username  string
		Token     string
	}
)

func myFunc(queue string, args ...interface{}) error {
	fmt.Printf("From %s, %v\n", queue, args)

	go func() {
		ticket := getTicket(args[1].(string))
		postToSlack(ticket)
	}()

	return nil
}

func postToSlack(ticket Ticket) {
	someMessage := SlackMessageAttachment{Title: "GAGO KA", Text: ticket.Subject}
	attachments := []SlackMessageAttachment{someMessage}

	creds := SlackCreds{Webhook: os.Getenv("SLACK_WEBHOOK"), AccessToken: os.Getenv("SLACK_ACCESS_TOKEN")}

	message := SlackMessage{Text: fmt.Sprintf("A ticket has been %s", ticket.Status), Attachments: attachments}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", creds.Webhook, b)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", creds.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	fmt.Println(message)
	fmt.Print(b)
	fmt.Println(resp)
}

func getTicket(ticketID string) Ticket {
	creds := ZDCreds{Subdomain: os.Getenv("ZENDESK_SUBDOMAIN"), Username: os.Getenv("ZENDESK_USER"), Token: os.Getenv("ZENDESK_TOKEN")}

	client := &http.Client{}

	var url = fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets/%s.json", creds.Subdomain, ticketID)

	req, _ := http.NewRequest("GET", url, new(bytes.Buffer))

	data := []byte(fmt.Sprintf("%s/token:%s", creds.Username, creds.Token))
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(data))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	ticket := WeirdTicketWrapper{}
	if err := json.Unmarshal(bodyBytes, &ticket); err != nil {
		panic(err)
	}

	return ticket.Ticket
}

func init() {
	gotenv.Load()
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{"myqueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       2,
	}
	goworker.SetSettings(settings)
	goworker.Register("MyClass", myFunc)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
