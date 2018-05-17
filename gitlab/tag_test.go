package gitlab

import (
	"testing"
)

func Test_GetLatestTag(t *testing.T) {
}

func Test_NewTag(t *testing.T) {
	// tag, err := NewTag(TestValidParams)
	// if err != nil {
	// 	t.Error("fail at NewTag: ", err)
	// }
	// latestTag, err := GetLatestTag(TestValidParams)
	// if err != nil {
	// 	t.Error("fail at NewTag: ", err)
	// }
	// if latestTag != tag {
	// 	t.Error("fail at NewTag")
	// }
}

func Test_NewTag_Invalid(t *testing.T) {
	_, err := NewTag(TestInvalidParamsUser)
	if err == nil || err.Error() != "username or password is wrong" {
		t.Error("fail at NewTag")
	}
}
