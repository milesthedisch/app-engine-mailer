package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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

type Attachment struct {
	// Name must be set to a valid file name.
	Name      string
	Data      []byte
	ContentID string
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path != "/" {
		log.Errorf(ctx, "Not found")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		_, err := fmt.Fprintf(w, "Operational")

		if err != nil {
			log.Errorf(ctx, "Error GET %v", err)
			panic(err)
		}

	case "POST":
		sendMail(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		log.Debugf(ctx, "Sorry, only GET and POST methods are supported.")
	}
}

func serverError(ctx context.Context, m string, w http.ResponseWriter, err error) {
	if err != nil {
		log.Errorf(ctx, "%s %v", m, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	return
}

func sendMail(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	dir, err := os.Getwd()
	serverError(ctx, "DIRECTORY ERROR", w, err)

	var reqBody ReqBody

	buf, _ := ioutil.ReadAll(r.Body)

	log.Debugf(ctx, "%s", buf)

	var result map[string]interface{}
	json.Unmarshal(buf, &result)

	log.Debugf(ctx, "Result: %+v", result)

	props := result["properties"].(map[string]interface{})

	log.Debugf(ctx, "Properties %+v", result)

	for v, k := range props {
		log.Infof(ctx, "\n Key: %+v \n Value: %+v", v, k.(map[string]interface{})["value"])
	}

	// // somehow this makes the buffer into json
	// rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	decoder := json.NewDecoder(rdr2)
	decodeErr := decoder.Decode(&reqBody)
	serverError(ctx, "DECODER ERROR", w, decodeErr)

	firstName, email := reqBody.Properties.Firstname.Value, reqBody.Properties.Email.Value

	log.Infof(ctx, "reqBody: %+v", reqBody)

	if firstName == "" || email == "" {
		log.Errorf(ctx, "Bad request not correct data, data: %v %v", firstName, email)
		http.Error(w, "BAD REQUEST NO CORRECT DATA", 400)
		return
	}

	log.Debugf(ctx, "Sending mail to %s", email)

	subject := "Terms and Conditions - Acumen Finance"
	body := fmt.Sprintf("<p> Dear %s,<br><p>Thank you for your recent loan submission and engagement via www.acumenfinance.com.au/apply,<br><p>This email is to confirm that your application has been received and we will contact you A.S.A.P to progress the transaction further.<br>Also please find attached some further information on Acumen Finance and its services and also our standard terms of engagement for your records. ", firstName)

	file := filepath.Join(dir, "./terms_and_conditions.pdf")

	data, err := ioutil.ReadFile(file)
	serverError(ctx, "READ FILE ERROR", w, err)

	msg := &mail.Message{
		Sender:   "cl@acumenfinance.com.au Commercial Loans <cl@acumenfinance.com.au>",
		Cc:       []string{"nd@acumenfinance.com.au"},
		To:       []string{email},
		Subject:  subject,
		HTMLBody: body,
		Attachments: []mail.Attachment{
			{
				Name:      "terms-and-conditions.html",
				Data:      data,
				ContentID: "<fieldid>",
			},
		},
	}

	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "%v", err.Error())
		fmt.Fprintf(w, "Coudn't send email: %v", err.Error())
	}
}
