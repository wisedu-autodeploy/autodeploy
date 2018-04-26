package main

import (
	"autodeploy/src/gitlab"
	"autodeploy/src/marathon"
	"log"
)

var gitlabCfg = map[string]string{
	"origin":      "http://172.16.7.53:9090",
	"loginAction": "/users/sign_in",
	"username":    "lisiurday",
	"password":    "Yihe210210.",
}

var projectCfg = map[string]string{
	"maintainer":  "wecloud-counselor",
	"projectName": "wec-counselor-collector-apps",
}

var showLogDetail = false

func main() {
	var err error

	log.Println("login gitlab...")
	_, err = gitlab.Init(gitlabCfg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("creating new tag...")
	tag, err := gitlab.NewTag(projectCfg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("created new tag:", tag)

	log.Println("build log:")
	ok, _, image, err := gitlab.WatchBuildLog(projectCfg, tag, showLogDetail)
	if err != nil || !ok {
		log.Fatalln(err)
	}
	log.Println("image pushed succeed:", image)

	log.Println("deploying, please wait...")
	ok, err = marathon.Deploy(projectCfg["projectName"], image)
	if err != nil || !ok {
		log.Fatalln(err)
	}
	log.Println("deploy succeed")
}
