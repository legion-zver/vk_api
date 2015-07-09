package main

import (
	"log"

	"github.com/yanple/vk_api"
)

func main() {
	// Login/pass auth
	var api = &vk_api.Api{}
    // ToDo need check captcha with  http://captchabot.com

    // Temp get token
    // Go To http://oauth.vk.com/oauth/authorize?redirect_uri=http://oauth.vk.com/blank.html&response_type=token&client_id=3087104&&v=5.0&scope=wall,offline&display=wap
    // Login and copy token from url.
    // And put token in the api.AccessToken

    api.AccessToken = "put token here"

//	err := api.LoginAuth(
//		"email/pass",
//		"pass",
//		"3087104",      // client id
//		"wall,offline", // scope (permissions)
//	)
//	if err != nil {
//		panic(err)
//	}

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
