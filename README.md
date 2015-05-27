Api client for VKontakte with login/pass authorization (hack) on Go (golang).
==========
###Plus: masking client_id to the iPhone, Android, iPad, Windows Phone clients.

go (golang) api client for vk.com

###Get
```Bash
    go get github.com/yanple/vk_api
    // and dependence
    go get github.com/PuerkitoBio/goquery
```

###Import
```Go
    @import "github.com/yanple/vk_api"
```

##How to use

###Login/pass auth

```Go
	var api vk_api.Api
	err := api.LoginAuth(
		"email/phone",
		"pass",
		"3087104", // client id
		"wall,offline", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}
```

###OAuth (click "allow" on special vk page)
[See martini based example](https://github.com/yanple/vk_api/blob/master/examples/martini.go)

###Make query to API
```Go
	params := make(map[string]string)
	params["domain"] = "yanple"
	params["count"] = "1"

	strResp, err := api.Request("wall.get", params)
    if err != nil {
        panic(err)
    }
    log.Println(strResp)
```
[See example](https://github.com/yanple/vk_api/blob/master/examples/example.go)

All api methods on https://vk.com/dev/methods

###Client ids (Masking only for login/pass auth)
```Go
    // client_id = "28909846" # Vk application ID (Android) doesn't work.
	// client_id = "3502561"  # Vk application ID (Windows Phone)
	// client_id = "3087106"  # Vk application ID (iPhone)
	// client_id = "3682744"  # Vk application ID (iPad)
```

### License
Vk_api by Yanple is [BSD licensed](./LICENSE)
