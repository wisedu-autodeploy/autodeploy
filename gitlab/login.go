package gitlab

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wisedu-autodeploy/autodeploy/client"
)

// login return a new gitlab session.
func login(user User) (session client.Sessioner, err error) {
	if gSession != nil {
		return gSession, nil
	}

	// get cookie
	cookie, err := getCookie(user)
	if err != nil {
		return
	}

	// set session
	session = client.NewSession().SetCookie(cookie)

	gSession = session
	return session, err
}

// getCookie return a valid login cookie.
func getCookie(user User) (cookie string, err error) {
	authenticityToken, cookie, err := getAuthenticityToken()
	if err != nil {
		return
	}

	formData := url.Values{
		"utf8":               {"✓"},
		"authenticity_token": {authenticityToken},
		"user[login]":        {user.Username},
		"user[password]":     {user.Password},
		"user[remember_me]":  {"0"},
	}

	session := client.NewSession().
		SetCookie(cookie).
		AddHeader("Origin", gOrigin).
		AddHeader("Referer", gLoginURL)

	res, err := session.PostForm(gLoginURL, formData)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 302 {
		return "", errors.New("username or password is wrong")
	}

	cookies := strings.Join(res.Header["Set-Cookie"], ";")
	matches := regexp.MustCompile(`(_gitlab_session=.*?);`).FindStringSubmatch(cookies)
	cookie = matches[1]

	return
}

// getAuthenticityToken return authenticityToken and preset cookie.
func getAuthenticityToken() (authenticityToken string, cookie string, err error) {
	res, err := http.Get(gLoginURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	cookie = strings.Split(res.Header["Set-Cookie"][0], ";")[0]

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", cookie, err
	}

	authenticityToken, exists := doc.Find("input[name=authenticity_token]").First().Attr("value")
	if !exists {
		err = errors.New("not found authenticity_token")
	}
	if err != nil {
		return
	}

	return authenticityToken, cookie, err
}
