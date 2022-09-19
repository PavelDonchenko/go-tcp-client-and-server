package gracfulshurdown

import (
	"../currency"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

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
			continue
		}

		wg.Add(1)

		go func() {
			wg.Done()
			handleConnection5(c)
		}()
	}
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
		currTime := currency.GetCurrency().Cube[0].Date
		rates := currency.GetCurrency().Cube[0].Rates
		clientReq := strings.TrimSpace(netData)

		for _, rate := range rates {
			if rate.Currency == strings.ToUpper(clientReq) {
				c.Write([]byte(currTime + " 1 EURO = " + rate.Rate + strings.ToUpper(netData)))
			}
		}
	}
	c.Close()

}
