package main

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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

	_, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go handleConnection(c)
		}
	}()

	wg.Wait()
	fmt.Println("Main done")
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
		clientReq := strings.TrimSpace(netData)

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + " 1 EURO = " + rate.Rate + strings.ToUpper(netData)))
			}
		}
	}
	c.Close()

}
