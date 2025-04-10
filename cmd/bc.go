package cmd

import (
	"github.com/0xJeti/bbscope/internal/utils"
	"github.com/0xJeti/bbscope/pkg/platforms/bugcrowd"
	"github.com/0xJeti/bbscope/pkg/whttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bcCmd represents the bc command
var bcCmd = &cobra.Command{
	Use:   "bc",
	Short: "Bugcrowd",
	Long:  "Gathers data from Bugcrowd (https://bugcrowd.com/)",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		token, _ := cmd.Flags().GetString("token")
		categories, _ := cmd.Flags().GetString("categories")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		outputFlags, _ := rootCmd.PersistentFlags().GetString("output")
		delimiterCharacter, _ := rootCmd.PersistentFlags().GetString("delimiter")
		includeOOS, _ := rootCmd.PersistentFlags().GetBool("oos")

		proxy, _ := rootCmd.PersistentFlags().GetString("proxy")
		bbpOnly, _ := rootCmd.Flags().GetBool("bbpOnly")
		pvtOnly, _ := rootCmd.Flags().GetBool("pvtOnly")

		email := viper.GetViper().GetString("bugcrowd-email")
		password := viper.GetViper().GetString("bugcrowd-password")

		if proxy != "" {
			whttp.SetupProxy(proxy)
		}

		if email != "" && password != "" && token == "" {
			token, err = bugcrowd.Login(email, password, proxy)
			if err != nil {
				utils.Log.Fatal("[bc] ", err)
			}
		}

		_, err = bugcrowd.GetAllProgramsScope(token, bbpOnly, pvtOnly, categories, outputFlags, concurrency, delimiterCharacter, includeOOS, true, nil)

		if err != nil {
			utils.Log.Fatal("[bc] ", err)
		}

		utils.Log.Info("bbscope run successfully")
	},
}

func init() {
	rootCmd.AddCommand(bcCmd)
	bcCmd.Flags().StringP("token", "t", "", "Bugcrowd session token (_bugcrowd_session cookie)")
	bcCmd.Flags().StringP("categories", "c", "all", "Scope categories, comma separated (Available: all, url, api, mobile, android, apple, other, hardware)")

	// Useless as of now since we're forcing 1 http request per second due to Bugcrowd's WAF
	bcCmd.Flags().IntP("concurrency", "", 1, "Concurrency threshold")

	bcCmd.Flags().StringP("email", "E", "", "Login email")
	viper.BindPFlag("bugcrowd-email", bcCmd.Flags().Lookup("email"))

	bcCmd.Flags().StringP("password", "P", "", "Login password")
	viper.BindPFlag("bugcrowd-password", bcCmd.Flags().Lookup("password"))

}
