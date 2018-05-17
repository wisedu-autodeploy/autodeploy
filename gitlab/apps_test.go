package gitlab

import (
	"testing"
)

var user = User{
	Username: "lisiurday",
	Password: "Yihe210210.",
}

func Test_GetAllApps(t *testing.T) {
	_, err := GetAllApps(TestValidUser)
	if err != nil {
		t.Error("fail at GetAllApps")
	}
}
