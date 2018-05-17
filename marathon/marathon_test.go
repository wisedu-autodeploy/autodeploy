package marathon

import (
	"testing"
)

var TestConfig = Config{
	MarathonID:   "wec-counselor-worklog-apps-v.0.0.7",
	MarathonName: "wec-counselor-worklog-apps-v-0-0-7",
}
var image = "172.16.9.100:5000/wec-counselor-worklog-apps:test_v.0.3.37_70ac6a7ed83156c1798b0c8901182fead76abf2f"

func Test_GetApps(t *testing.T) {
	_, err := GetApps()
	if err != nil {
		t.Error("fail at GetApps:", err)
	}

}

func Test_Deploy(t *testing.T) {
	t.Skip("skip Test_Deploy")
	ok, err := Deploy(TestConfig, image)
	if err != nil {
		t.Error("fail at Deploy:", err)
	}
	if !ok {
		t.Error("fail at Deploy")
	}
}
