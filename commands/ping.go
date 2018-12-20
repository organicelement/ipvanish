package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//	"io/ioutil"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/briandowns/spinner"
	tm "github.com/buger/goterm"
	"github.com/kellydunn/golang-geo"
	"github.com/organicelement/ipvanish/ipv"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//func init() {
//	listCmd.AddCommand(listDraftsCmd)
//}

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

var spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "List top 5 servers (by latency)",
	Long:  `List all of the drafts in your content directory.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		spin.Prefix = "Scanning: "
		//spin.Suffix = ""
		//spin.FinalMSG = "Finished..."
		spin.Color("green")
		spin.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		spin.Stop()
	},
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)

		servers := ipv.GetServers()
		if len(servers) <= 0 {
			fmt.Printf("IPVanish returned %v servers.\n",
				tm.Color(strconv.Itoa(len(servers)), tm.RED),
			)
			return
		}

		geoip, _ := ipv.GetLocation()

		p := geo.NewPoint(geoip.Latitude, geoip.Longitude)

		for i, _ := range servers {
			srv := &servers[i]
			lat := srv.Properties.Latitude
			long := srv.Properties.Longitude
			p2 := geo.NewPoint(lat, long)
			dist := p.GreatCircleDistance(p2)

			// Convert km to miles
			srv.Distance = dist * 0.621371
			//log.Infof("Server %s is %.0f miles away\n", srv.Properties.IP, srv.Distance)
		}

		//log.Infof("Servers with Geo Distance: %v", servers)

		// Start of Ping
		ping := fastping.NewPinger()
		ping.Network("udp")
		netProto := "ip4:icmp"

		results := make(map[string]*ipv.IPVServer)
		for i, _ := range servers {
			srv := &servers[i]
			ra, err := net.ResolveIPAddr(netProto, srv.Properties.IP)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			results[ra.String()] = srv // maybe &servers[i]?
			ping.AddIPAddr(ra)

		}

		onRecv, onIdle := make(chan *response), make(chan bool)

		ping.OnRecv = func(addr *net.IPAddr, t time.Duration) {
			onRecv <- &response{addr: addr, rtt: t}
		}
		ping.OnIdle = func() {
			onIdle <- true
		}

		ping.MaxRTT = 2 * time.Second
		ping.RunLoop()
		timeout := time.NewTimer(time.Second * 10)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)

		//log.Infof("Geo Results: %v", results)
		//iter := 0
	loop:
		for {
			select {
			case <-c:
				fmt.Println("got interrupted")
				break loop
			case res := <-onRecv:
				//iter++
				//log.Debugln(iter)
				if res.rtt > 0 {
					server := results[res.addr.String()]
					if server != nil {
						server.Latency = res.rtt
						//log.Infof("Server is %s", server.Properties.IP)
						//log.Infof("Response is %s, %s", res.addr, res.rtt)
					} else {
						log.Infof("Didn't find the server for %s", res.addr.String())
					}
				}
			case <-onIdle:
				log.Debugln("Idle Called")
				for host, r := range results {
					if Debug && r.Latency == 0 {
						log.Debugf("%s : unreachable", host)
					}
				}
			case <-ping.Done():
				log.Debugln("Done Pinging")
				if err := ping.Err(); err != nil {
					fmt.Println("Ping failed:", err)
				}
				break loop
			case <-timeout.C:
				log.Debugln("Timed out")
				break loop
			}
		}
		timeout.Stop()
		signal.Stop(c)
		ping.Stop()

		srvrs := make([]ipv.IPVServer, len(results))
		//log.Infof("%v results", len(results))
		i := 0
		for _, v := range results {
			//log.Infof("This server is %v", *v)
			if v != nil {
				srvrs[i] = *v
				//log.Infof("Added %s", v.Properties.IP)
			} else {
				//log.Info("No server in this iteration")
			}
			i++
		}

		switch viper.GetString(SORTFLAG) {
		case DISTANCE:
			ipv.OrderedBy(ipv.Distance).Sort(servers)
		case LATENCY:
			ipv.OrderedBy(ipv.Latency).Sort(servers)
		case CAPACITY:
			ipv.OrderedBy(ipv.Distance, ipv.Capacity).Sort(servers)
		}

		numresults, _ := strconv.Atoi(cmd.Flags().Lookup("results").Value.String())
		fmt.Print("\n")
		for _, server := range servers[:numresults] {
			fmt.Printf("Host %v has %v utilization, is %v miles away, with %vms latency\n",
				// TODO: Use https://github.com/gosuri/uitable
				tm.Color(server.Properties.Hostname, tm.CYAN),
				tm.Color(strconv.Itoa(server.Properties.Capacity)+"%", capacityColor(server.Properties.Capacity)),
				tm.Color(strconv.FormatFloat(server.Distance, 'f', 0, 64), tm.RED),
				tm.Color(strconv.FormatInt(int64(time.Duration(server.Latency)/time.Millisecond), 10), tm.GREEN),
			)
		}

	},
}
