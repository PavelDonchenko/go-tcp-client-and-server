package main

import (
	"./currency"
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TLC server !")
			return
		}

		currTime := currency.GetCurrency().Cube[0].Date
		rates := currency.GetCurrency().Cube[0].Rates
		clientReq := strings.TrimSpace(netData)
		fmt.Println(currTime, rates)

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + " 1 EURO = " + rate.Rate + strings.ToUpper(netData)))
			}
		}
	}
}
