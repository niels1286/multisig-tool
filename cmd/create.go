package cmd

import (
	"fmt"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
)

var m int
var pks string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.GetText("0003"),
	Long:  i18n.GetText("0003"),
	Run: func(cmd *cobra.Command, args []string) {
		if m < 1 || m > 15 {
			fmt.Println(i18n.GetText("0020"))
			return
		}

		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.MultiAccountSDK.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(i18n.GetText("10000")+".\naddress:", msAccount.Address)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	createCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	createCmd.MarkFlagRequired("m")
	createCmd.MarkFlagRequired("publickeys")
}
