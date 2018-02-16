package vertec

/**
According to https://www.vertec.com/de/support/kb/technik-und-datenmodell/vertecservice/xml/xmlschnittstelle/
 */

import (
	"strings"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	log "github.com/inconshreveable/log15"
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
	// StringBuffer / Writer?
	var post string = `<Envelope>
  <Header>
   <BasicAuth>
     ${auth}
   </BasicAuth>
  </Header>
  <Body>${query}</Body>
</Envelope>`

	if len(settings.Token) > 0 {
		post = strings.Replace(post, "${auth}", "<Token>" + settings.Token + "</Token>", 1)
	} else {
		log.Debug("using lecacy auth")
		post = strings.Replace(post, "${auth}", "<Name>${username}</Name><Password>${password}</Password>", 1)
		// replace patterns by values from config
		post = strings.Replace(post, "${username}", settings.Username, 1)
		post = strings.Replace(post, "${password}", settings.Password, 1)
	}

	// insert query into <Body/> section
	post = strings.Replace(post, "${query}", query, 1)

	return httppost(settings, post)
}

func httppost(settings Settings, xmlQuery string) (string, error) {

	// no authentication used. username and password are submitted as cleartext in the POST section :scream:
	// req.SetBasicAuth(`username`, `password`)

	response, err := settings.Connection.Post(settings.URL, "application/xml", strings.NewReader(xmlQuery));
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Login(settings Settings, username string, password string) error {
	// Vertec specific auth url rewrite
	authurl := strings.Replace(settings.URL, "/xml", "/auth/xml", 1)

	form := url.Values{}
	form.Add("vertec_username", username)
	form.Add("password", password)

	fmt.Printf("accessing %s with form %s\n", authurl, form.Encode())

	response, err := settings.Connection.Post(authurl, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to authenticate. Status code %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	settings.Token = string(body)

	return nil
}
