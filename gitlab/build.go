package gitlab

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// WatchBuildLog watch build log.
func WatchBuildLog(params Params, tag string, logChan chan *Logger, wg *sync.WaitGroup) {
	logger := &(Logger{})
	// write error info to logChan
	handleErr := func(err error) {
		logger.Status = -1
		logger.Message = err.Error()
		logChan <- logger
		return
	}

	// get build log id
	buildLogID, err := getBuildLogID(params, tag)
	if err != nil {
		handleErr(err)
		return
	}

	// get build log url
	buildLogURL, err := getBuildLogURL(params, buildLogID)
	if err != nil {
		handleErr(err)
		return
	}

	// log.Println("watching build log, waiting...")
	time.Sleep(time.Duration(60) * time.Second)

	for {
		ok, image, log, err := judgeIsFinish(params, buildLogURL)
		if err != nil {
			handleErr(err)
			break
		}

		s := string([]rune(log))
		splices := strings.Split(s, "\n")

		if ok {
			wg.Done()
			logger.Status = 1
			logger.Image = image
			logger.Log = splices
			logger.Message = "success"
			logChan <- logger
			break
		} else {
			logger.Status = 0
			logger.Image = image
			logger.Log = splices
			logger.Message = "building"
			logChan <- logger
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
	return
}

func getBuildLogID(params Params, tag string) (id string, err error) {
	session, err := login(params.User)
	if err != nil {
		return
	}

	path := gOrigin + "/" + params.Project.Maintainer + "/" + params.Project.Name
	piplinesTagsURL := path + "/pipelines?scope=tags"
	res, err := session.Get(piplinesTagsURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return
	}

	id = ""
	doc.Find(".commit").Each(func(i int, selection *goquery.Selection) {
		curTag := selection.Find(".monospace.branch-name").First().Text()
		if curTag == tag {
			href, ok := selection.Find(".commit-link").First().Find("a").First().Attr("href")
			if !ok {
				err = errors.New("not found tag pipline id")
				return
			}
			splices := strings.Split(href, "/")
			id = splices[len(splices)-1]
		}
	})
	if id == "" {
		err = errors.New("not found tag pipline id")
		return
	}
	return
}

func getBuildLogURL(params Params, buildID string) (buildLogURL string, err error) {
	session, err := login(params.User)
	if err != nil {
		return
	}

	path := gOrigin + "/" + params.Project.Maintainer + "/" + params.Project.Name
	piplinesURL := path + "/pipelines/" + buildID
	res, err := session.Get(piplinesURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return
	}

	buildLogURL, ok := doc.Find(".pipeline-graph .stage-column:last-child .build-content a").First().Attr("href")
	if !ok {
		err = errors.New("not found pipline url")
		return
	}
	buildLogURL = gOrigin + buildLogURL

	return
}

func judgeIsFinish(params Params, buildLogURL string) (ok bool, image string, log string, err error) {
	session, err := login(params.User)
	if err != nil {
		return
	}

	res, err := session.Get(buildLogURL + "/raw")
	if err != nil {
		ok = false
		return
	}

	contentBytes, err := ioutil.ReadAll(res.Body)
	log = string(contentBytes)

	matches := regexp.MustCompile(`building and pushing (.*?)\s`).FindStringSubmatch(log)
	if len(matches) > 1 {
		image = matches[1]
	}
	if strings.Contains(log, "[ERROR] ") {
		ok = false
		err = errors.New("Build failed")
		return
	}
	ok = strings.Contains(log, "Build succeeded")

	return
}
