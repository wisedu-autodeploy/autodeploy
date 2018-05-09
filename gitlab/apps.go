package gitlab

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup

// App .
type App struct {
	Maintainer string
	Name       string
}

// GetAllApps .
func GetAllApps() (apps []App, err error) {
	res, err := session.Get(projectsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// 获取最后一页的页码
	var maxPage int
	lastURL, exists := doc.Find(".pagination .last a").First().Attr("href")
	if !exists {
		// err = errors.New("not found lastURL")
		maxPage = 1
	} else {
		getPageRegexp := regexp.MustCompile(`.*?page=(\d*)`)
		maxPageStr := getPageRegexp.FindStringSubmatch(lastURL)[1]
		maxPage, _ = strconv.Atoi(maxPageStr)
	}

	// 多线程请求每个页面的数据
	wg.Add(maxPage)
	appsChan := make(chan []App)

	for i := 1; i <= maxPage; i++ {
		go getOnePageApps(i, appsChan)
	}

	appsChan <- apps

	wg.Wait()

	// 将通道里的值取出
	apps = <-appsChan
	close(appsChan)

	return apps, err
}

func getOnePageApps(page int, appsChan chan []App) {
	pageStr := strconv.Itoa(page)
	res, err := session.Get(projectsURL + "/?page=" + pageStr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tmpApps := []App{}
	doc.Find(".projects-list .project-row a.project").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			log.Fatalln(false)
			return
		}
		splices := strings.Split(href, "/")
		tmpApps = append(tmpApps, App{
			Maintainer: splices[1],
			Name:       splices[2],
		})
	})
	wg.Done()
	appsChan <- append(<-appsChan, tmpApps...)
	return
}
