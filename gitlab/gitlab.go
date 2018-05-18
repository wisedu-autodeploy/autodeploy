package gitlab

import (
	"github.com/wisedu-autodeploy/autodeploy/client"
)

var (
	gSession     client.Sessioner
	gOrigin      = "http://172.16.7.53:9090"
	gLoginURL    = gOrigin + "/users/sign_in"
	gProjectsURL = gOrigin + "/dashboard/projects"
)

var mode = ""

// User is gitlab login params.
type User struct {
	Username string
	Password string
}

// Project is gitlab project info.
type Project struct {
	Maintainer string
	Name       string
}

// Params is gitlab login params and target project info.
type Params struct {
	User
	Project
}

// Logger .
type Logger struct {
	Log     []string
	Image   string
	Status  int // -1 => fail; 0 => init or running; 1 => success
	Message string
}

func Debugger() {
	mode = "debug"
}
