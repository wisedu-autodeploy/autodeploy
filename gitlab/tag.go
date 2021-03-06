package gitlab

import (
	"errors"
	"log"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

// NewTag create a new tag.
func NewTag(params Params) (newTag string, err error) {
	// get session
	session, err := login(params.User)
	if err != nil {
		return
	}

	// get latest tag
	if mode == "debug" {
		log.Println("try to get latest tag...")
	}
	latestTag, err := GetLatestTag(params)
	if err != nil {
		return
	}
	if mode == "debug" {
		log.Println("get latest tag:", latestTag)
	}

	// calc target tag
	newTag = addTagVersion(latestTag, "patch")
	if mode == "debug" {
		log.Println("calc new tag:", newTag)
	}

	// get authenticity_token
	if mode == "debug" {
		log.Println("try to get tag authenticity token")
	}
	path := params.Project.Maintainer + "/" + params.Project.Name + "/tags"
	newTagURL := gOrigin + "/" + path + "/new"
	res, err := session.Get(newTagURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	authenticityToken, exists := doc.Find("input[name=authenticity_token]").First().Attr("value")
	if !exists {
		err = errors.New("not found new tag authenticity_token")
		return
	}
	if mode == "debug" {
		log.Println("got tag authenticity token:", authenticityToken)
	}

	// post request to create target tag.
	if mode == "debug" {
		log.Println("try to create new tag:", newTag)
	}
	postURL := gOrigin + "/" + path
	formData := url.Values{
		"utf8":                {"✓"},
		"authenticity_token":  {authenticityToken},
		"tag_name":            {newTag},
		"ref":                 {"master"},
		"message":             {"tagged with autodeploy by " + params.User.Username},
		"release_description": {""},
	}
	res, err = session.PostForm(postURL, formData)

	if res.StatusCode != 302 {
		err = errors.New("not 302 Found, create new tag failed")
	}
	if mode == "debug" {
		log.Println("created new tag:", newTag)
	}
	return
}

// GetLatestTag get lastest tag.
func GetLatestTag(params Params) (tag string, err error) {
	session, err := login(params.User)
	if err != nil {
		return
	}

	path := params.Project.Maintainer + "/" + params.Project.Name + "/tags"
	projectTagsURL := gOrigin + "/" + path

	res, err := session.Get(projectTagsURL)

	if res.StatusCode == 404 {
		err = errors.New("project not found")
	}
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}

	tagText := doc.Find("div.tags > ul").First().Find(".item-title").First().Text()
	reg, err := regexp.Compile("[^a-zA-Z0-9_.]+")
	tag = reg.ReplaceAllString(tagText, "")
	return tag, err
}
