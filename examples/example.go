package main

import (
	"log"

	"github.com/yanple/vk_api"
)

func main() {
	// Login/pass auth
	var api = &vk_api.Api{}
	err := api.LoginAuth(
		"email/phone",
		"pass",
		"3087104",      // client id
		"wall,offline", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}

	// Make query
	params := make(map[string]string)
	params["domain"] = "yanple"
	params["count"] = "1"

	strResp, err := api.Request("wall.get", params)
    if err != nil {
        panic(err)
    }
	if strResp != "" {
		log.Println(strResp)
	}
}
