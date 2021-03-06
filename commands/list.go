/**
 * Copyright 2015 Organic Element LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	tm "github.com/buger/goterm"
	"github.com/kellydunn/golang-geo"
	"github.com/organicelement/ipvanish/ipv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List top 20 servers (by distance)",
	Long:  `List all of the drafts in your content directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)

		servers := ipv.GetServers()
		geoip, err := ipv.GetLocation()

		if err != nil {
			logrus.Errorf("Could not determine this position: %s", err.Error())
		}

		p := geo.NewPoint(geoip.Latitude, geoip.Longitude)

		for i, server := range servers {
			lat := server.Properties.Latitude
			long := server.Properties.Longitude
			p2 := geo.NewPoint(lat, long)
			dist := p.GreatCircleDistance(p2)

			// Convert km to miles
			servers[i].Distance = dist * 0.621371
		}

		switch viper.GetString(SORTFLAG) {
		case DISTANCE:
			ipv.OrderedBy(ipv.Distance).Sort(servers)
		case LATENCY:
			ipv.OrderedBy(ipv.Distance, ipv.Latency).Sort(servers)
		case CAPACITY:
			ipv.OrderedBy(ipv.Distance, ipv.Capacity).Sort(servers)
		}

		results, _ := strconv.Atoi(cmd.Flags().Lookup("results").Value.String())
		for _, server := range servers[:results] {
			fmt.Printf("Host %v has %v utilization and is %v miles away\n",
				tm.Color(server.Properties.Hostname, tm.CYAN),
				tm.Color(strconv.Itoa(server.Properties.Capacity)+"%", capacityColor(server.Properties.Capacity)),
				tm.Color(strconv.FormatFloat(server.Distance, 'f', 0, 64), tm.RED),
			)
		}
	},
}

func capacityColor(capacity int) int {
	if capacity >= 75 {
		return tm.RED
	} else if capacity >= 25 && capacity < 75 {
		return tm.YELLOW
	} else if capacity < 25 {
		return tm.GREEN
	}
	return tm.WHITE
}
