package vertec

import (
	"strings"
	"net/http"
	"io/ioutil"
)

type Settings struct {
	URL         string
	Username    string
	Password    string
}

func Version() string {
	return "0.0.1"
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

	return httppost(settings.URL, post)
}

func httppost(url string, xmlQuery string) (string, error) {

	req, err := http.NewRequest("POST", url, strings.NewReader(xmlQuery))
	if err != nil {
		return "", err
	}

	// no authentication used. username and password are submitted as cleartext in the POST section :scream:
	//req.SetBasicAuth(`username`, `password`)

	// TODO verify/improve client handling
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body), nil
}
