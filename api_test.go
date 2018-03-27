package vertec

import (
	"testing"
	"net/http"
	"reflect"
	"encoding/xml"
	"fmt"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestSmokeQuery(t *testing.T) {
	var client http.Client
	value := "nonsmoking"
	client.Transport = &MockRoundTripper{ content: value }

	var s Settings
	s.URL = "dummy"
	s.Password = "no-password"
	s.Username = "no-username"
	s.Connection = client

	result, _ := Query("s", s)
	assertEqual(t, value, result)
}

func TestListProjects(t *testing.T) {

	input := `<Envelope><Body>
		<QueryResponse>
			<Projekt>
				<objid>1</objid>
				<beschrieb>Phase1</beschrieb>
				<code>P1</code>
			</Projekt>
			<Projekt>
				<objid>99</objid>
				<beschrieb>Phase2</beschrieb>
				<code>P2</code>
			</Projekt>
		</QueryResponse>
	</Body></Envelope>`

	// same as above
	expected := Projects {xml.Name{"", "Envelope"}, []Project { {1, "P1", "Phase1"}, { 99, "P2", "Phase2"} }}

	buf, _:= xml.MarshalIndent(expected, "", "  ")
	fmt.Println("expected ", string(buf))

	result, _ := ListProjects("123", forResponse(input))
	assertEqual(t, len(expected.Elements), len(result.Elements))
	assertEqual(t, true, reflect.DeepEqual(expected, result))
}

func forResponse(input string) Settings {
	var client http.Client
	client.Transport = &MockRoundTripper{ content: input }
	var s Settings
	s.URL = "dummy"
	s.Password = "no-password"
	s.Username = "no-username"
	s.Connection = client
	return s
}