package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

type Envelope2 struct {
	Cube []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

func GetCurrency2() Envelope2 {
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

	var env Envelope2
	err = xml.Unmarshal(xmlCurrenciesData, &env)

	if err != nil {
		fmt.Println(err)
	}
	currency := env
	return currency
}
func main() {

	//SERVER
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	port := ":" + arguments[1]
	l, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
func handleConnection(c net.Conn) {
	fmt.Print(".")
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		currTime := GetCurrency2().Cube[0].Date
		rates := GetCurrency2().Cube[0].Rates
		clientReq := strings.TrimSpace(string(netData))

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + "  "))
				c.Write([]byte("1 EURO = " + rate.Rate))
			}
		}
		c.Write([]byte(strings.ToUpper(netData)))
	}
	c.Close()
}
