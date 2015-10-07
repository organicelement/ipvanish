package commands

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const SORTFLAG = "sort"
const CAPACITY = "capacity"
const LATENCY = "latency"
const DISTANCE = "distance"

var IpvanishCmd = &cobra.Command{
	Use:   "ipvanish",
	Short: "IPVanish command line utilities",
	Long: `ipvanish is the main command.

IPVanish is a VPN provider, this command lists the servers and
sorts/displays by utilization or distance.

Complete documentation is available at http://ipvanish.com/.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
	},
	PersistentPreRun: begin,
	PersistentPostRun: finish,
}

var ipvanishCmdV *cobra.Command

//Flags that are to be added to commands.
var Sort string
var Results uint

//Execute adds all child commands to the root command.
func Execute() {
	AddCommands()
	IpvanishCmd.Execute()
//	utils.StopOnErr(IpvanishCmd.Execute())
}

//AddCommands adds child commands to the root command.
func AddCommands() {
	IpvanishCmd.AddCommand(listCmd)
	IpvanishCmd.AddCommand(pingCmd)
}

//Initializes flags
func init() {
	// Persistent == available to sub commands
	IpvanishCmd.PersistentFlags().StringVarP(&Sort, SORTFLAG, "s", "capacity", "sort order for hosts (default is distance")
	IpvanishCmd.PersistentFlags().UintVarP(&Results, "results", "n", 20, "Total number of results to display and filter")

	ipvanishCmdV = IpvanishCmd

	// for Bash autocomplete
	validSortFlags := []string{"distance", "capacity", "latency"}
	IpvanishCmd.PersistentFlags().SetAnnotation("sort", cobra.BashCompOneRequiredFlag, validSortFlags)

	// This message will be shown to Windows users if the command is opened from explorer.exe
	cobra.MousetrapHelpText = `

  IPVanish is a command line tool

  You need to open cmd.exe and run it from there.`
}

func LoadDefaultSettings() {
	viper.SetDefault("Watch", false)
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
//	viper.SetConfigFile(CfgFile)

	viper.RegisterAlias("indexes", "taxonomies")

	LoadDefaultSettings()

	if ipvanishCmdV.PersistentFlags().Lookup("sort").Changed {
		viper.Set("sort", Sort)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

}

func begin(cmd *cobra.Command, args []string) {
}

func finish(cmd *cobra.Command, args []string) {
}
