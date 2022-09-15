package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Envelope3 struct {
	Cube []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

func GetCurrency3() Envelope3 {
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

	var env Envelope3

	err = xml.Unmarshal(xmlCurrenciesData, &env)
	if err != nil {
		fmt.Println(err)
	}
	currency := env
	return currency
}
func NewServer() {
	cmdAddr, _ := net.ResolveTCPAddr("tcp", "localhost:6666")
	lcmd, err := net.ListenTCP("tcp", cmdAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer lcmd.Close()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	for {
		select {
		case <-quitChan:
			lcmd.Close()
			wg.Wait()
			return
		default:
		}
		lcmd.SetDeadline(time.Now().Add(1e9))
		c, err := lcmd.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		}
		if err != nil {
			//log.WithError(err).Errorln("Listener accept")
			continue
		}
		wg.Add(1)
		go func() {
			wg.Done()
			handleConnection5(c)
		}()

	}
}
func main() {
	NewServer()
}
func handleConnection5(c net.Conn) {
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
		currTime := GetCurrency3().Cube[0].Date
		rates := GetCurrency3().Cube[0].Rates
		clientReq := strings.TrimSpace(netData)

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + " 1 EURO = " + rate.Rate + strings.ToUpper(netData)))
			}
		}
	}
	c.Close()

}
