package gitlab

import (
	"testing"
)

func Test_GetAuthenticityToken(t *testing.T) {
	token, cookie, err := getAuthenticityToken()
	if token == "" || cookie == "" || err != nil {
		t.Error("fail at getAuthenticityToken")
	}
}

func Test_GetCookie(t *testing.T) {
	cookie, err := getCookie(TestValidUser)
	if cookie == "" || err != nil {
		t.Error("fail at getCookie")
	}
}

func Test_Login(t *testing.T) {
	gSession = nil
	session, err := login(TestValidUser)
	if session == nil || err != nil {
		t.Error("fail at login")
	}
}

// multi login should get same session
func Test_Login_MultiLogin(t *testing.T) {
	gSession = nil
	session1, err := login(TestValidUser)
	if session1 == nil || err != nil {
		t.Error("fail at login")
	}
	session2, err := login(TestValidUser)
	if err != nil {
		t.Error("fail at login")
	}
	if session1 != session2 {
		t.Error("multi login should get same session")
	}
}

func Test_Login_InvalidAccess(t *testing.T) {
	gSession = nil
	_, err := login(TestInvalidUser)
	if err == nil || err.Error() != "username or password is wrong" {
		t.Error("fail at login invalid access")
	}
}
