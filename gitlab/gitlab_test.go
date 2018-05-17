package gitlab

var TestValidUser = User{
	Username: "lisiurday",
	Password: "Yihe210210.",
}

var TestInvalidUser = User{
	Username: "username",
	Password: "password",
}

var TestValidProject = Project{
	Maintainer: "wecloud-counselor",
	Name:       "wec-counselor-worklog-apps",
}

var TestInvalidProject = Project{
	Maintainer: "",
	Name:       "",
}

var TestValidParams = Params{
	User:    TestValidUser,
	Project: TestValidProject,
}

var TestInvalidParamsUser = Params{
	User:    TestInvalidUser,
	Project: TestValidProject,
}

var TestInvalidParamsProject = Params{
	User:    TestValidUser,
	Project: TestInvalidProject,
}
