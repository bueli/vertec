package vertec

import (
	"encoding/xml"
	"strings"
	"io"
)

type Project struct {
	Objid uint64 `xml:"objid" json:"id"`
	Code string `xml:"code" json:"code"`
	Beschrieb string `xml:"beschrieb" json:"description"`
}

// this hack structure might be created at runtime ...
type Projects  struct {
	Elements []Project `xml:"QueryResponse>Projekt"`
}

func ListProjects(user string, settings Settings) (Projects, error) {

	query := `<Selection>
		<objref>${USER}</objref>
		<ocl>bearbProjekte-&gt;select(aktiv)</ocl>
	</Selection>
	<Resultdef>
		<member>code</member>
		<member>beschrieb</member>
	</Resultdef>`

	var q2 = strings.Replace(query, "${USER}", user, 1)
	response := Projects {}
	err := queryList(q2, &response, settings)
	if err != nil {
		return Projects{}, err
	}
	return response, nil
}

func queryList(query string, v interface{}, settings Settings) (error) {

	response, err := Query(query, settings)
	if err != nil {
		return err
	}
	d := xml.NewDecoder(strings.NewReader(response))

	for {
		token, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				// end of document
				break
			}
			// propagate error
			return tokenErr
		}

		switch node := token.(type) {
		case xml.StartElement:
			if node.Name.Local == "Body" {
				// decode whole phase according to xml: annotations
				if err := d.DecodeElement(v, &node); err != nil {
					// propagate error
					return err
				}
			}
		}
	}

	return nil

}