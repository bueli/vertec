package vertec

import (
	"encoding/xml"
	"bytes"
	"strings"
	"io"
)

type Project struct {
	Objid string `xml:"objid" json:"id"`
	Code string `xml:"code" json:"code"`
	Beschrieb string `xml:"beschrieb" json:"description"`
}

// this hack structure might be created at runtime ...
type Projects  struct {
	Elements []Project `xml:"QueryResponse>Projekt"`
}

func ListProjects(user string, settings Settings) (Projects, error) {

	query := `<Query><Selection>
		<objref>${USER}</objref>
		<ocl>bearbProjekte->select(aktiv)</ocl>
	</Selection>
	<Resultdef>
		<member>code</member>
		<member>beschrieb</member>
	</Resultdef></Query>`

	var q2 = strings.Replace(query, "${USER}", user, 1)
	response := Projects {}
	err := queryList(q2, &response, settings)
	if err != nil {
		return Projects{}, err
	}
	return response, nil
}

/*
    <QueryResponse>
      <OffeneLeistung>
        <objid>3740257</objid>
        <aot/>
        <minutenInt>30</minutenInt>
        <minutenIntBis>420</minutenIntBis>
        <minutenIntVon>390</minutenIntVon>
        <phase>
          <objref>2459562</objref>
        </phase>
        <projekt>
          <objref>601026</objref>
        </projekt>
        <text/>
      </OffeneLeistung>
      <OffeneLeistung>
*/

type Leistung struct {
	Objid string `xml:"objid" json:"id"`
	Minuten string `xml:"minutenInt" json:"dauer"`
	MinutenVon string `xml:"minutenIntBis" json:"start"`
	Kommentar string `xml:"text" json:"comment"`
	Phase string `xml:"phase>objref" json:"phase"`
	Projekt string `xml:"projekt>objref" json:"projekt"`
}

// this hack structure might be created at runtime ...
type Leistungen  struct {
	Elements []Leistung `xml:"QueryResponse>OffeneLeistung"`
}

func ListUserOffeneLeistungen(user string, settings Settings) (Leistungen, error) {

	ocl := bytes.NewBufferString("")
	
	// whatever phase.code='KOMP' means ... booked?
	xml.Escape(ocl, []byte("offeneleistungen->select((datum >= '01.03.2018'.strToDate) and (datum < '26.3.2018'.strToDate))->reject(phase.code='KOMP')"))
		
	query := `<Query><Selection>
		<objref>${USER}</objref>
		<ocl>
			${OCL}
		</ocl>
	</Selection>
	<Resultdef>
		<member>minutenInt</member>
		<member>minutenIntBis</member>
		<member>minutenIntVon</member>
		<member>phase</member>
		<expression><alias>code</alias><ocl>phase.code</ocl></expression>
		<member>projekt</member>
		<member>text</member>
	</Resultdef></Query>`

	var q2 = strings.Replace(query, "${USER}", user, 1)
	q2 = strings.Replace(q2, "${OCL}", ocl.String(), 1)
	response := Leistungen {}
	err := queryList(q2, &response, settings)
	if err != nil {
		return Leistungen{}, err
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