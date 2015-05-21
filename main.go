package main

import (
	"encoding/json"
	"fmt"
	"github.com/kaneshin/giita/giita/client"
	"github.com/kaneshin/giita/giita/request"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	stdout      = os.Stdout
	stderr      = os.Stderr
	stdin       = os.Stdin
	programName = os.Args[0]
	team        string
	token       string
	conf        = ".giita"
)

func init() {
	filename := path.Join(os.Getenv("HOME"), conf)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Fprintf(stderr, "%s doesn't exist.\n", filename)
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}
	if val, ok := data["team"].(string); ok {
		team = val
	}
	if val, ok := data["token"].(string); ok {
		token = val
	}
}

func main() {
	page := 1
	limit := 100
	var users = make(map[string]int)
	for {
		fmt.Printf("Request: page %d, per_page %d\n", page, limit)
		req := request.NewItemRequestWithPageAndLimit(team, page, limit)
		cli := client.NewClient(token)
		body, err := cli.Dispatch(req)
		if err != nil {
			fmt.Fprintln(stderr, err)
		}
		var data []map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Fprintln(stderr, err)
		}
		for _, d := range data {
			if user, ok := d["user"].(map[string]interface{}); ok {
				if id, ok := user["id"].(string); ok {
					users[id]++
				}
			}
		}
		if len(data) < limit {
			break
		}
		page++
		time.Sleep(10 * time.Second)
	}
	for id, val := range users {
		fmt.Printf("%s, %d\n", id, val)
	}
}
