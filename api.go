package vk_api

import (
    "errors"
    "io/ioutil"
    "log"
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strconv"
    "encoding/json"

    "github.com/PuerkitoBio/goquery"
)

const API_METHOD_URL = "https://api.vk.com/method/"
const AUTH_HOST = "https://oauth.vk.com/authorize"
const AUTH_HOST_GET_TOKEN = "https://oauth.vk.com/access_token"
const API_VERSION = "5.62"

type Api struct {
    AccessToken string
    UserId      int
    ExpiresIn   int
    debug       bool
}

func ParseResponseUrl(responseUrl string) (string, string, string, error) {
    u, err := url.Parse("?" + responseUrl)
    if err != nil {
        return "", "", "", err
    }

    q := u.Query()
    return q.Get("access_token"), q.Get("user_id"), q.Get("expires_in"), nil
}

func parse_form(doc *goquery.Document) (url.Values, string, error) {
    _origin, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=_origin]").Attr("value")
    if exists == false {
        return nil, "", errors.New("Not _origin attr in vk form")
    }

    ip_h, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=ip_h]").Attr("value")
    if exists == false {
        return nil, "", errors.New("Not ip_h attr in vk form")
    }

    to, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=to]").Attr("value")
    if exists == false {
        return nil, "", errors.New("Not 'to' attr in vk form")
    }

    // ToDo: need realese captha verify (show client and send).
    _, exists = doc.Find("#captcha").Attr("src")
    if exists == false {
        return nil, "", errors.New("Error: captcha blocked. Login not available in the time.")
    }

    formData := url.Values{}
    formData.Add("_origin", _origin)
    formData.Add("ip_h", ip_h)
    formData.Add("to", to)    

    url, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form").Attr("action")
    if exists == false {
        return nil, "", errors.New("Not action attr in vk form")
    }
    return formData, url, nil
}

func auth_user(email string, password string, client_id string, scope string, client *http.Client) (*http.Response, error) {
    var auth_url = AUTH_HOST+"?" +
    "redirect_uri=https://oauth.vk.com/blank.html&response_type=token&" +
    "client_id=" + client_id + "&v="+API_VERSION+"&scope=" + scope + "&display=wap"

    res, e := client.Get(auth_url)
    if e != nil {
        return nil, e
    }

    doc, err := goquery.NewDocumentFromResponse(res)
    if err != nil {
        return nil, err
    }

    formData, url, err := parse_form(doc)
    if err != nil {
        return nil, err
    }
    formData.Add("email", email)
    formData.Add("pass", password)    

    res, e = client.PostForm(url, formData)
    if e != nil {
        return nil, e
    }
    return res, nil
}

func get_permissions(response *http.Response, client *http.Client) (*http.Response, error) {
    doc, err := goquery.NewDocumentFromResponse(response)
    if err != nil {
        return nil, err
    }

    formData, url, err := parse_form(doc)
    if err != nil {
        return nil, err
    }

    res, err := client.PostForm(url, formData)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (vk *Api) Request(methodName string, params map[string]string) (string, error) {
    u, err := url.Parse(API_METHOD_URL + methodName)
    if err != nil {
        return "", err
    }

    q := u.Query()
    for k, v := range params {
        q.Set(k, v)
    }
    q.Set("access_token", vk.AccessToken)
    q.Set("v", API_VERSION)
    u.RawQuery = q.Encode()

    resp, err := http.Get(u.String())
    if err != nil {
        return "", err
    }

    defer resp.Body.Close()
    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(content), nil
}

func (vk *Api) LoginAuth(email string, password string, client_id string, scope string) error {
    cookieJar, _ := cookiejar.New(nil)
    client := &http.Client{
        Jar: cookieJar,
    }

    res, err := auth_user(email, password, client_id, scope, client)
    if err != nil {
        return err
    }

    if res.Request.URL.Path != "/blank.html" {
        res, err = get_permissions(res, client)
        if err != nil {
            return err
        }
        if res.Request.URL.Path != "/blank.html" {
            return errors.New("Not auth")
        }
    }
    accessToken, userId, expiresIn, err := ParseResponseUrl(res.Request.URL.Fragment)

    int_userId, err := strconv.Atoi(userId)
    if err != nil {
        return err
    }

    int_expiresIn, err := strconv.Atoi(expiresIn)
    if err != nil {
        return err
    }

    if vk.debug {
        log.Printf("Access token %s for user %s", accessToken, userId)
    }

    vk.AccessToken = accessToken
    vk.ExpiresIn = int_userId
    vk.UserId = int_expiresIn

    return nil
}

func (vk *Api) GetAuthUrl(redirect_uri string, client_id string, scope string) (string, error) {
    u, err := url.Parse(AUTH_HOST)
    if err != nil {
        return "", err
    }

    q := u.Query()
    q.Set("v", API_VERSION)
    q.Set("client_id", client_id)
    q.Set("scope", scope)
    q.Set("redirect_uri", redirect_uri)
    q.Set("response_type", "code")
    u.RawQuery = q.Encode()

    return u.String(), nil
}

type OAuthResponse struct {
    AccessToken         string  `json:"access_token"`
    ExpiresIn           int     `json:"expires_in"`
    UserId              int     `json:"user_id"`
    Error               string  `json:"error"`
    ErrorDescription    string  `json:"error_description"`
}

func parse_oauth_response(response []byte) (OAuthResponse, error) {
    var jsonResp OAuthResponse
    if err := json.Unmarshal(response, &jsonResp); err != nil {
        return jsonResp, err
    }
    return jsonResp, nil
}

func (vk *Api) OAuth(redirect_uri string, client_secret string, client_id string, code string) (error) {
    u, err := url.Parse(AUTH_HOST_GET_TOKEN)
    if err != nil {
        return err
    }
    q := u.Query()
    q.Set("v", API_VERSION)
    q.Set("redirect_uri", redirect_uri)
    q.Set("client_secret", client_secret)
    q.Set("client_id", client_id)
    q.Set("code", code)
    u.RawQuery = q.Encode()

    resp, err := http.Get(u.String())
    if err != nil {
        return err
    }

    defer resp.Body.Close()
    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    data, err := parse_oauth_response(content)
    if err != nil {
        return err
    }
    if data.Error != "" {
        return errors.New(fmt.Sprintf("%s: %s", data.Error, data.ErrorDescription))
    }
    if vk.debug {
        log.Printf("Access token %s for user %s", data.AccessToken, data.UserId)
    }

    vk.AccessToken = data.AccessToken
    vk.ExpiresIn = data.ExpiresIn
    vk.UserId = data.UserId
    return nil
}

func (vk *Api) SetDebug(s bool) {
    vk.debug = s
}
