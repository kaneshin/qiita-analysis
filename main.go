package main

import (
	"encoding/json"
	"fmt"
	"github.com/kaneshin/giita/giita/client"
	"github.com/kaneshin/giita/giita/request"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	stdout      = os.Stdout
	stderr      = os.Stderr
	stdin       = os.Stdin
	programName = os.Args[0]
	team        string
	token       string
)

const (
	conf        = ".giita"
	waitSeconds = 1
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
	var users = make(map[string]map[string]interface{})
	for {
		req := request.NewUserRequestWithPageAndLimit(team, page, limit)
		cli := client.NewClient(token)
		body, err := cli.Dispatch(req)
		if err != nil {
			panic(err)
		}
		var data []map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}
		for _, user := range data {
			if id, ok := user["id"].(string); ok {
				info := make(map[string]interface{})
				if val, ok := user["name"].(string); ok {
					if len(val) > 0 {
						info["name"] = val
					} else {
						info["name"] = id
					}
				} else {
					info["name"] = id
				}
				if val, ok := user["profile_image_url"].(string); ok {
					info["profile_image_url"] = val
				} else {
					info["profile_image_url"] = ""
				}
				if val, ok := user["items_count"].(float64); ok {
					info["items_count"] = (int)(val)
				} else {
					info["items_count"] = 0
				}
				if val, ok := user["description"].(string); ok {
					info["description"] = val
				} else {
					info["description"] = ""
				}
				info["team_items_count"] = 0
				users[id] = info
			}
		}
		if len(data) < limit {
			break
		}
		page++
		time.Sleep(waitSeconds * time.Second)
	}
	page = 1
	for {
		req := request.NewItemRequestWithPageAndLimit(team, page, limit)
		cli := client.NewClient(token)
		body, err := cli.Dispatch(req)
		if err != nil {
			panic(err)
		}
		var data []map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}
		for _, d := range data {
			if user, ok := d["user"].(map[string]interface{}); ok {
				if id, ok := user["id"].(string); ok {
					if val, ok := users[id]["team_items_count"].(int); ok {
						users[id]["team_items_count"] = val + 1
					}
				}
			}
		}
		if len(data) < limit {
			break
		}
		page++
		time.Sleep(waitSeconds * time.Second)
	}
	var formatString = func(v map[string]interface{}, name string) string {
		if val, ok := v[name].(string); ok {
			return strings.Replace(val, "\n", " ", -1)
		}
		return ""
	}
	fmt.Println("| User Image | Name | Description | Qiita:Team posts | Qiita posts |")
	fmt.Println("| :--------: | :--- | :---------- | :--------------- | :---------- |")
	for id, val := range users {
		imgURL := formatString(val, "profile_image_url")
		name := formatString(val, "name")
		desc := formatString(val, "description")
		titems := val["team_items_count"].(int)
		items := val["items_count"].(int)
		fmt.Printf(`| <img src="%s" height=60 /> | <a href="https://qiita.com/%s">%s</a> | %s | %d | %d |
`, imgURL, id, name, desc, titems, items)
	}
}

func simple() {
	page := 1
	limit := 100
	var users = make(map[string]int)
	for {
		req := request.NewItemRequestWithPageAndLimit(team, page, limit)
		cli := client.NewClient(token)
		body, err := cli.Dispatch(req)
		if err != nil {
			panic(err)
		}
		var data []map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
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
		time.Sleep(waitSeconds * time.Second)
	}
	for id, val := range users {
		fmt.Printf("%s, %d\n", id, val)
	}
}
