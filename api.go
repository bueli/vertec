package vertec

/**
According to https://www.vertec.com/de/support/kb/technik-und-datenmodell/vertecservice/xml/xmlschnittstelle/
 */

import (
	"strings"
	"net/http"
	neturl "net/url"
	"io/ioutil"
)

type Settings struct {
	URL         string
	Username    string
	Password    string
	// reuse connection HTTP 1.1
	Connection	http.Client
	// Authentication Token
	Token		string
}

func Version() string {
	return "0.0.2"
}

func Query(query string, settings Settings) (string, error) {
	// not very sophisticated: just replace markers inside the fixed structure 
	var post string = `<Envelope>
  <Header>
   <BasicAuth>
     <Name>${username}</Name>
     <Password>${password}</Password>
   </BasicAuth>
  </Header>
  <Body>${query}</Body>
</Envelope>`

	// replace patterns by values from config
	post = strings.Replace(post, "${username}", settings.Username, 1)
	post = strings.Replace(post, "${password}", settings.Password, 1)

	// insert query into <Body/> section
	post = strings.Replace(post, "${query}", query, 1)

	return httppost(settings, post)
}

func httppost(settings Settings, xmlQuery string) (string, error) {

	// no authentication used. username and password are submitted as cleartext in the POST section :scream:
	//req.SetBasicAuth(`username`, `password`)

	response, err := settings.Connection.Post(settings.URL, "application/xml", strings.NewReader(xmlQuery));
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body), nil
}

func Login(settings Settings, username string, password string) error {
	// Vertec specific auth url rewrite
	url := strings.Replace(settings.URL, "/xml", "/auth/xml", 1)

	parameters := "?vertec_username=" +
		neturl.QueryEscape(username) + "&password=" + neturl.QueryEscape(password)

	response, err := settings.Connection.Post(url, "application/x-www-form-urlencoded", strings.NewReader(parameters))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	settings.Token = string(body)

	return nil
}