package stadtreinigung

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"strings"
)

type HouseNumber struct {
	Name string
	Url  string
}

func ParseHouseNumber(content string, bremerStadtreinigungRootUrl string) []HouseNumber {
	houseNumbers := make([]HouseNumber, 0)

	decoder := xml.NewDecoder(strings.NewReader(content))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose

	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			fmt.Println(tokenErr)
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if matchesHouseNumber(startElement) {
				var td Td
				err := decoder.DecodeElement(&td, &startElement)

				if err != nil {
					fmt.Printf("Unable to decode tag %s, Tag skipped", startElement.Name.Local)
					continue
				}

				if td.matches() {
					houseNumbers = append(houseNumbers, td.mapToHouseNumber(bremerStadtreinigungRootUrl))
				}
			}
		}
	}

	return houseNumbers
}

func (td Td) mapToHouseNumber(bremerStadtreinigungRootUrl string) HouseNumber {
	return HouseNumber{html.UnescapeString(td.A.Value), bremerStadtreinigungRootUrl + td.A.Href}
}

func matchesHouseNumber(startElement xml.StartElement) bool {
	if startElement.Name.Local == `td` {
		for _, attribute := range startElement.Attr {
			if attribute.Name.Local == `class` && strings.Contains(attribute.Value, `BAKStr`) {
				return true
			}
		}
	}

	return false
}
