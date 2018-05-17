package gitlab

import (
	"testing"
)

func Test_AddTagVersion_Patch(t *testing.T) {
	tag := "version!@#$%^&*.2.3.4"
	targetTag := "version!@#$%^&*.2.3.5"
	if targetTag != addTagVersion(tag, "patch") {
		t.Error("fail at addTagVersion for patch")
	}
}

func Test_AddTagVersion_Minor(t *testing.T) {
	tag := "version!@#$%^&*.2.3.4"
	targetTag := "version!@#$%^&*.2.4.4"
	if targetTag != addTagVersion(tag, "minor") {
		t.Error("fail at addTagVersion for minor")
	}
}

func Test_AddTagVersion_Major(t *testing.T) {
	tag := "version!@#$%^&*.3.3.4"
	targetTag := "version!@#$%^&*.4.3.4"
	if targetTag != addTagVersion(tag, "major") {
		t.Error("fail at addTagVersion for major")
	}
}
