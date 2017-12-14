package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/benmanns/goworker"
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

	message := SlackMessage{Text: fmt.Sprintf("A ticket has been %s", ticket.Status), Attachments: attachments}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://hooks.slack.com/services/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", b)

	req.Header.Add("Authorization", "Bearer <SLACK API TOKEN>")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	fmt.Println(message)
	fmt.Print(b)
	fmt.Println(resp)
}

func getTicket(ticketID string) Ticket {
	client := &http.Client{}

	var url = fmt.Sprintf("https://<subdomain>.zendesk.com/api/v2/tickets/%s.json", ticketID)

	req, _ := http.NewRequest("GET", url, new(bytes.Buffer))

	data := []byte("username/token:<TOKEN>")
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
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{"myqueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       0.01,
	}
	goworker.SetSettings(settings)
	goworker.Register("MyClass", myFunc)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
