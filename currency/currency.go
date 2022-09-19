package currency

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Envelope struct {
	Cube []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

func GetCurrency() Envelope {
	// get the latest exchange rate
	resp, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	xmlCurrenciesData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var env Envelope
	err = xml.Unmarshal(xmlCurrenciesData, &env)
	if err != nil {
		fmt.Println(err)
	}

	currency := env
	return currency
}
