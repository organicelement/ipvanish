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
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/organicelement/ipvanish/ipv"
	tm "github.com/buger/goterm"
	"github.com/kellydunn/golang-geo"
	"strconv"
	"sort"
	"fmt"
)

//func init() {
//	listCmd.AddCommand(listSortCmd)
//}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List top 20 servers (by distance)",
	Long:  `List all of the drafts in your content directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)

		servers := ipv.GetServers()
		geoip, _ := ipv.GetLocation()

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

		results, _ := strconv.Atoi(cmd.Flags().Lookup("results").Value.String())
		for _, server := range servers[:results] {
			fmt.Printf("Host %v has %v utilization and is %v miles away\n",
				tm.Color(server.Properties.Hostname, tm.CYAN),
				tm.Color(strconv.Itoa(server.Properties.Capacity) + "%", capacityColor(server.Properties.Capacity)),
				tm.Color(strconv.FormatFloat(server.Distance, 'f', 0, 64), tm.RED),
			)
		}
		fmt.Println("Args: ", cmd.Flags().Lookup("sort").Value)
//		cmd.DebugFlags()
	},
}

func capacityColor(capacity int) int {
	if capacity >=75 {
		return tm.RED
	} else if capacity >=25 && capacity < 75 {
		return tm.YELLOW
	} else if capacity < 25 {
		return tm.GREEN
	}
	return tm.WHITE
}
