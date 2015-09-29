package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
//	"io/ioutil"
	"github.com/tatsushid/go-fastping"
	"net"
	"fmt"
	"os"
	"time"
	"os/signal"
	"syscall"
	log "github.com/Sirupsen/logrus"
)

//func init() {
//	listCmd.AddCommand(listDraftsCmd)
//}

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}


var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "List top 5 servers (by distance)",
	Long:  `List all of the drafts in your content directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)


		// Start of Ping
		p := fastping.NewPinger()
		p.Network("udp")
		netProto := "ip4:icmp"
		ra, err := net.ResolveIPAddr(netProto, "www.organicelement.com")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		results := make(map[string]*response)
		results[ra.String()] = nil
		p.AddIPAddr(ra)

		onRecv, onIdle := make(chan *response), make(chan bool)
		p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
			onRecv <- &response{addr: addr, rtt: t}
		}
		p.OnIdle = func() {
			onIdle <- true
		}

		p.MaxRTT = time.Second
		p.RunLoop()
		timeout := time.NewTimer(time.Second * 5)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)

		loop:
		for {
			select {
			case <-c:
				fmt.Println("get interrupted")
				break loop
			case res := <-onRecv:
				if _, ok := results[res.addr.String()]; ok {
					results[res.addr.String()] = res
				}
			case <-onIdle:
				for host, r := range results {
					if r == nil {
						fmt.Printf("%s : unreachable %v\n", host, time.Now())
					} else {
						fmt.Printf("%s : %v %v\n", host, r.rtt, time.Now())
					}
					results[host] = nil
				}
			case <-p.Done():
				if err = p.Err(); err != nil {
					fmt.Println("Ping failed:", err)
				}
				break loop
			case <-timeout.C:
				log.Info("Timed out")
				break loop
			}
		}
		log.Infof("Stop 1")
		timeout.Stop()
		log.Infof("Stop 2")
		signal.Stop(c)
		log.Infof("Stop 3")
		p.Stop()
		log.Info("End")
	},
}
