package commands

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/utils"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Main command, all other commands are attached to this one via Execute/AddCommands
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
//	PreRun: func(cmd *cobra.Command, args []string) {
//		fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
//	},
//	PostRun: func(cmd *cobra.Command, args []string) {
//		fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
//	},
}

var ipvanishCmdV *cobra.Command

//Flags that are to be added to commands.
var BuildWatch, IgnoreCache, Draft, Future, UglyURLs, Verbose, Logging, VerboseLog, DisableRSS, DisableSitemap, PluralizeListTitles, PreserveTaxonomyNames, NoTimes bool
var Source, CacheDir, Destination, Theme, Sort string
var Results uint

//Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
func Execute() {
	AddCommands()
	utils.StopOnErr(IpvanishCmd.Execute())
}

//AddCommands adds child commands to the root command HugoCmd.
func AddCommands() {
	IpvanishCmd.AddCommand(listCmd)
	IpvanishCmd.AddCommand(pingCmd)
}

//Initializes flags
func init() {
	// Persistent == available to sub commands
	IpvanishCmd.PersistentFlags().StringVarP(&Sort, "sort", "s", "capacity", "sort order for hosts (default is distance")
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
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if Source == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(Source)
	}
/*
	err := viper.ReadInConfig()
	if err != nil {
		jww.ERROR.Println("Unable to locate Config file. Perhaps you need to create a new site. Run `hugo help new` for details")
	}
*/

	viper.RegisterAlias("indexes", "taxonomies")

	LoadDefaultSettings()

	if ipvanishCmdV.PersistentFlags().Lookup("sort").Changed {
		viper.Set("sort", Sort)
	}

/*
	if !viper.GetBool("RelativeURLs") && viper.GetString("BaseURL") == "" {
		jww.ERROR.Println("No 'baseurl' set in configuration or as a flag. Features like page menus will not work without one.")
	}
*/

	//	Set defaults if empty
	if Theme != "" {
		viper.Set("theme", Theme)
	}

	if Destination != "" {
		viper.Set("PublishDir", Destination)
	}

	if Source != "" {
		viper.Set("WorkingDir", Source)
	} else {
		dir, _ := os.Getwd()
		viper.Set("WorkingDir", dir)
	}

	if viper.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}

	if VerboseLog {
		jww.SetLogThreshold(jww.LevelInfo)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

}

func begin(cmd *cobra.Command, args []string) {
}

func finish(cmd *cobra.Command, args []string) {
}
