package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
)

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

	// ctx := appengine.NewContext(r)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	fn := gjson.get(body, "properties.firstname")

	fmt.Fprintf(w, string(fn))

	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }

	// ctx := appengine.NewContext(r)
	// msg := &mail.Message{
	// 		Sender: "cl@acumenfinance.com.au"
	// 		To:
	// }

	// fmt.Fprintf(w, "SENDING MAIL")
}
