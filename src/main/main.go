package main

import (
	"autodeploy/src/gitlab"
	"autodeploy/src/marathon"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var (
	config     Config
	gitlabCfg  gitlab.Config
	projectCfg marathon.Config
	projects   []marathon.Config
	configPath string
)

var showLogDetail = false

// Config .
type Config struct {
	Gitlab   gitlab.Config
	Project  marathon.Config
	Projects []marathon.Config
}

func init() {
	gitlabCfg.Origin = "http://172.16.7.53:9090"
	gitlabCfg.LoginAction = "/users/sign_in"
	gitlabCfg.Username = ""
	gitlabCfg.Password = ""

	projectCfg.Maintainer = ""
	projectCfg.Name = ""

	projects = []marathon.Config{}

	config = Config{
		Gitlab:   gitlabCfg,
		Project:  projectCfg,
		Projects: projects,
	}
	configPath = "./config.json"

	if !pathExist(configPath) {
		err := writeConfig()
		if err != nil {
			log.Fatalln(err)
		}

	}

	err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func readConfig() error {
	b, err := ioutil.ReadFile(configPath)
	err = json.Unmarshal(b, &config)
	return err
}

func writeConfig() error {
	b, err := json.MarshalIndent(config, "", "    ")
	err = ioutil.WriteFile(configPath, b, 0644)
	return err
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func findIndex(s []marathon.Config, e marathon.Config) int {
	for i, a := range s {
		if a.Name == e.Name && a.Maintainer == e.Maintainer {
			return i
		}
	}
	return -1
}

func start() {
	var err error

	log.Println("login gitlab...")
	_, err = gitlab.Init(config.Gitlab)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("creating new tag...")
	tag, err := gitlab.NewTag(config.Project)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("created new tag:", tag)

	log.Println("build log:")
	ok, _, image, err := gitlab.WatchBuildLog(config.Project, tag, showLogDetail)
	if err != nil || !ok {
		log.Fatalln(err)
	}
	log.Println("image pushed succeed:", image)

	log.Println("deploying, please wait...")
	ok, err = marathon.Deploy(config.Project.Name, image)
	if err != nil || !ok {
		log.Fatalln(err)
	}
	log.Println("deploy succeed")
}

func main() {
	app := cli.NewApp()
	app.Name = "autodeploy"
	app.Usage = "auto deploy app!"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Usage: "username for the gitlab",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "password for the gitlab",
		},
		cli.StringFlag{
			Name:  "maintainer, m",
			Usage: "maintainer of the project",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "name of the project",
		},
		cli.StringFlag{
			Name:  "short, s",
			Usage: "auto deploy using short of the project",
		},
		cli.IntFlag{
			Name:  "index, i",
			Usage: "auto deploy using index of the project",
		},
		cli.BoolFlag{
			Name:  "log, l",
			Usage: "show building log",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "set",
			Usage: "set username/password",
			Action: func(c *cli.Context) error {
				key := c.Args().Get(0)
				val := c.Args().Get(1)
				if key == "username" {
					config.Gitlab.Username = val
				}
				if key == "password" {
					config.Gitlab.Password = val
				}
				err := writeConfig()
				return err
			},
		},
		{
			Name:  "add",
			Usage: "add projects",
			Action: func(c *cli.Context) (err error) {
				maintainer := c.Args().Get(0)
				name := c.Args().Get(1)
				short := c.Args().Get(2)
				if maintainer != "" && name != "" {
					if short == "" {
						short = name
					}
					cfg := marathon.Config{
						Maintainer: maintainer,
						Name:       name,
						Short:      short,
					}

					foundIndex := findIndex(config.Projects, cfg)
					if foundIndex != -1 { // found
						config.Projects[foundIndex] = cfg
					} else { // not found
						config.Projects = append(config.Projects, cfg)
					}
					err = writeConfig()
				} else {
					err = errors.New("please input maintainer and project name")
				}
				return err
			},
		},
		{
			Name:  "list",
			Usage: "list projects",
			Action: func(c *cli.Context) {
				for index, project := range config.Projects {
					fmt.Println(index+1, project.Short, "->", project.Maintainer+"/"+project.Name)
				}
			},
		},
		{
			Name:  "short",
			Usage: "auto deploy using short of the project",
			Action: func(c *cli.Context) (err error) {
				showLogDetail = c.Bool("log")

				short := c.Args().Get(0)
				if short == "" {
					err = errors.New("please input short")
					return
				}

				foundIndex := -1
				for i, pj := range config.Projects {
					if pj.Short == short {
						foundIndex = i
						break
					}
				}
				if foundIndex == -1 {
					err = errors.New("not found")
					return
				}
				writeConfig()

				start()
				return
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "log, l",
					Usage: "show building log",
				},
			},
		},
		{
			Name:  "index",
			Usage: "auto deploy using index of the project",
			Action: func(c *cli.Context) (err error) {
				showLogDetail = c.Bool("log")

				indexStr := c.Args().Get(0)
				if indexStr == "" {
					err = errors.New("please input index")
				}
				index, err := strconv.Atoi(indexStr)
				if index <= 0 {
					err = errors.New("please input index (>0)")
				}
				if err != nil {
					return
				}
				project := config.Projects[index-1]
				config.Project = project
				writeConfig()

				start()
				return
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "log, l",
					Usage: "show building log",
				},
			},
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		if c.String("username") != "" {
			config.Gitlab.Username = c.String("username")
		}
		if c.String("password") != "" {
			config.Gitlab.Password = c.String("password")
		}
		if c.String("maintainer") != "" {
			config.Project.Maintainer = c.String("maintainer")
		}
		if c.String("name") != "" {
			config.Project.Name = c.String("name")
		}
		if c.Int("index") != 0 {
			if c.Int("index") > len(config.Projects) {
				err = errors.New("index out of range")
			}
			project := config.Projects[c.Int("index")-1]
			config.Project = project
		}
		if c.String("short") != "" {
			if c.Int("index") > len(config.Projects) {
				err = errors.New("index out of range")
			}
			project := config.Projects[c.Int("index")-1]
			config.Project = project
		}
		showLogDetail = c.Bool("log")

		if config.Gitlab.Username == "" {
			return errors.New("please input gitlab username")
		}
		if config.Gitlab.Password == "" {
			return errors.New("please input gitlab password")
		}
		if config.Project.Maintainer == "" {
			return errors.New("please input project maintainer")
		}
		if config.Project.Name == "" {
			return errors.New("please input project name")
		}
		start()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
