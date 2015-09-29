package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
//	"io/ioutil"
	"encoding/json"
	"github.com/organicelement/ipvanish/ipv"
	"github.com/kellydunn/golang-geo"
	"strconv"
	"sort"
)

//func init() {
//	listCmd.AddCommand(listDraftsCmd)
//}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List top 5 servers (by distance)",
	Long:  `List all of the drafts in your content directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)

		var servers []ipv.IPVServer
		var geoip ipv.GeoIP

		res, err := http.Get("https://www.ipvanish.com/api/servers.geojson")

		if err != nil {
			return
		}

		if err = json.NewDecoder(res.Body).Decode(&servers); err != nil {
			return
		}

		res, err = http.Get("http://freegeoip.net/json/")

		if err != nil {
			return
		}

		if err = json.NewDecoder(res.Body).Decode(&geoip); err != nil {
			return
		}

		p := geo.NewPoint(geoip.Latitude, geoip.Longitude)


		for i, server := range servers {
			lat, err := strconv.ParseFloat(server.Properties.Latitude, 64)
			if err != nil {
				log.Warnf("Error parsing latitude : %s", err)
			}
			long, err := strconv.ParseFloat(server.Properties.Longitude, 64)
			if err != nil {
				log.Warnf("Error parsing longitude : %s", err)
			}
			p2 := geo.NewPoint(lat, long)
			dist := p.GreatCircleDistance(p2)
			// Convert km to miles
			servers[i].Distance = dist * 0.621371
		}

		sort.Sort(ipv.DistanceSorter(servers))

		for _, server := range servers {
			log.Infof("Server %v is %v miles from your location with %v%% utilization", server.Properties.Hostname, server.Distance, server.Properties.Capacity)
		}
	},
}
