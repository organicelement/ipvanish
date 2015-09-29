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
}

var ipvanishCmdV *cobra.Command

//Flags that are to be added to commands.
var BuildWatch, IgnoreCache, Draft, Future, UglyURLs, Verbose, Logging, VerboseLog, DisableRSS, DisableSitemap, PluralizeListTitles, PreserveTaxonomyNames, NoTimes bool
var Source, CacheDir, Destination, Theme, BaseURL, CfgFile, LogFile, Editor string

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
	IpvanishCmd.PersistentFlags().BoolVarP(&Draft, "buildDrafts", "D", false, "include content marked as draft")
	IpvanishCmd.PersistentFlags().BoolVar(&DisableRSS, "disableRSS", false, "Do not build RSS files")
	IpvanishCmd.PersistentFlags().StringVarP(&Source, "source", "s", "", "filesystem path to read files relative from")
	IpvanishCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")

	ipvanishCmdV = IpvanishCmd

	// for Bash autocomplete
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	IpvanishCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
	IpvanishCmd.PersistentFlags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})

	// This message will be shown to Windows users if Hugo is opened from explorer.exe
	cobra.MousetrapHelpText = `

  Hugo is a command line tool

  You need to open cmd.exe and run it from there.`
}

func LoadDefaultSettings() {
	viper.SetDefault("Watch", false)
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
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

	if ipvanishCmdV.PersistentFlags().Lookup("buildDrafts").Changed {
		viper.Set("BuildDrafts", Draft)
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
