package vertec

import (
	"testing"
	"net/http"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestSmokeQuiery(t *testing.T) {
	var client http.Client;
	value := "nonsmoking"
	client.Transport = &MockRoundTripper{ content: value }

	var s Settings;
	s.URL = "dummy"
	s.Password = "no-password"
	s.Username = "no-username"
	s.Connection = client;

	result, _ := Query("s", s);
	assertEqual(t, value, result)
}