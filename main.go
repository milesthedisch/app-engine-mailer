package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
)

type ReqBody struct {
	Properties struct {
		Firstname struct {
			Value string `json:"value"`
		} `json:"firstname"`
		Email struct {
			Value string `json:"value"`
		} `json:"email"`
	} `json:"properties"`
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Your mail server is operational")
	case "POST":
		sendMail(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var reqBody ReqBody
	ctx := appengine.NewContext(r)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	firstName, email := reqBody.Properties.Firstname.Value, reqBody.Properties.Email.Value

	if firstName == "" || email == "" {
		http.Error(w, "BAD REQUEST NO CORRECT DATA", 400)
		return
	}

	subject := "Terms and Conditions - Acumen Finance"
	body := fmt.Sprintf("<p> Dear %s,<br>,<p>Thank you for your recent loan submission and engagement via www.acumenfinance.com.au/apply,<br>,<p>This email is to confirm that your application has been received and we will contact you A.S.A.P to progress the transaction further.,<br>,Also please find attached some further information on Acumen Finance and its services and also our standard terms of engagement for your records. ", firstName)

	dat, err := ioutil.ReadFile("/tmp/dat")

	if err != nil {
		fmt.Fprintf(w, err.Error(), 500)
	}

	msg := &mail.Message{
		Sender:   "cl@acumenfinance.com.au Commercial Loans <cl@acumenfinance.com.au>",
		To:       []string{email},
		Subject:  subject,
		HTMLBody: body,
	}

	if err := mail.Send(ctx, msg); err != nil {
		fmt.Fprintf(w, "Coudn't send email: %v", err.Error())
	}
}
