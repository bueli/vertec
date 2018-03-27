package vertec

import (
	"encoding/xml"
	"strings"
	"fmt"
	"io"
)

type Project struct {
	Objid uint64 `xml:"objid" json:"id"`
	Code string `xml:"code" json:"code"`
	Beschrieb string `xml:"beschrieb" json:"description"`
}


type Projects  struct {
	XMLName xml.Name
	Elements []Project `xml:"Body>QueryResponse>Projekt"`
}

func ListProjects(user string, settings Settings) (Projects, error) {

	query := `<Selection>
		<objref>[USER]</objref>
		<ocl>bearbProjekte-&gt;select(aktiv)</ocl>
	</Selection>
	<Resultdef>
		<member>code</member>
		<member>beschrieb</member>
	</Resultdef>`

	var q2 = strings.Replace(query, "[USER]", user, 1)
	response := Projects { XMLName: xml.Name{"", "Envelope"}, }
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
	fmt.Println("parsing ", response)
	d := xml.NewDecoder(strings.NewReader(response))

	return d.Decode(v) // works,  but involves correct prefix

	for {
		token, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			// propagate error
			return tokenErr
		}

		println("on token ", token)
		switch node := token.(type) {

		case xml.StartElement:
			if node.Name.Local == "QueryResponse" {
				fmt.Println("decoding ", node.Name.Local)
				// decode whole phase according to xml: annotations
				if err := d.DecodeElement(&v, &node); err != nil {
					// propagate error
					return err
				}
				buf, _ := xml.Marshal(v)
				fmt.Println("decoded ", string(buf))
			} else {
				println("decoding element ", node.Name.Local)
			}
		default:
			println("decoding unknown ", node)
		}
	}

	return nil

}