package main

import (
  "net/http"
  "strings"
  "net/url"
  "fmt"
  "os"
  "flag"
)

const (
  /* Url = "http://localhost:5000/api/status"*/
)

func main() {
  status := flag.String("p", "", "Push a new status")
  flag.Parse()

  PostStatus(*status)
}

func PostStatus(status string) (err error) {
  token := os.Getenv("TXTATUS_TOKEN")
  endpoint := "http://txtatus.com/api/status"

  values := make(url.Values)
  values.Set("status", status)

  req, err := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
  if (err != nil) { return err }

  req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", token))

  resp, err := http.DefaultClient.Do(req)
  if (err != nil) { return err }

  defer resp.Body.Close()
  fmt.Println("RESP: ", resp)

  return nil
}
