// Copyright (c) 2013 Andrea Pavoni, Jason McVetta.  This is Free Software,
// released under an MIT license.  See http://opensource.org/licenses/MIT for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

// txstatus-cli is a command line client for Txtatus.com.
package main

import (
	"flag"
	"fmt"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/restclient"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const endpoint = "http://txtatus.com/api/status"

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [-v] \"spent %%2.hours doing #dev on @myproject\"\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	token := env.StringDefault("TXTATUS_TOKEN", "")
	if token == "" {
		fmt.Println("The TXTATUS_TOKEN environment variable must be set.")
		return
	}
	t := Txtatus{
		Token: token,
	}
	// status := flag.String("p", "", "Push a new status")
	flag.Usage = usage
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()
	args := flag.Args()
	if flag.NArg() != 1 {
		fmt.Println("Must supply a quoted status message")
		usage()
		return
	}
	res, err := t.Send(args[0])
	if err != nil {
		log.Fatal(err)
	}
	dur, err := time.ParseDuration(strconv.Itoa(res.Time) + "s")
	var timeStr string
	if err == nil {
		timeStr = dur.String()
	} else {
		timeStr = strconv.Itoa(res.Time)
	}
	fmt.Println("--- Txtatus Status Posted ---")
	fmt.Println(" Status ID:", res.Id)
	if *verbose {
		fmt.Println("   User ID:", res.UserId)
		fmt.Println("Project ID:", res.ProjectId)
		fmt.Println("Created At:", res.CreatedAt)
		fmt.Println("      Date:", res.Date)

	}
	fmt.Println("      Time:", timeStr)
	fmt.Println("      Body:", res.Body)
	fmt.Println("      Tags:", strings.Join(res.Tags, ", "))
}

// A Txtatus is a client for the txtatus.com API.
type Txtatus struct {
	Token string
}

type statusReq struct {
	Status string `json:"status"`
}

type statusResp struct {
	Id        string     `json:"_id"`
	UserId    string     `json:"user_id"`
	ProjectId string     `json:"project_id"`
	CreatedAt *time.Time `json:"created_at"`
	Date      string
	Time      int
	Body      string
	Tags      []string
}

// Send posts a status to txstatus.com.
func (t *Txtatus) Send(status string) (*statusResp, error) {

	values := make(url.Values)
	values.Set("status", status)
	payload := statusReq{
		Status: status,
	}
	res := statusResp{}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Token token=%s", t.Token))
	rr := restclient.RequestResponse{
		Url:            endpoint,
		Method:         "POST",
		Header:         &h,
		Data:           payload,
		Result:         &res,
		ExpectedStatus: 201,
	}
	_, err := restclient.Do(&rr)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
