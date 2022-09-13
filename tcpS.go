package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	GetCurrency()
	//TCP SERVER
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
	}

	port := ":" + arguments[1]
	l, err := net.Listen("tcp", port)
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

		currTime := GetCurrency().Cube[0].Date
		rates := GetCurrency().Cube[0].Rates
		clientReq := strings.TrimSpace(string(netData))
		fmt.Println(currTime, rates)

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + "  "))
				c.Write([]byte("1 EURO = " + rate.Rate))
			}

		}
		c.Write([]byte(strings.ToUpper(netData)))
	}
}
