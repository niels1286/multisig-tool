package cmd

import (
	"fmt"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
)

var m int
var pks string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a multi address",
	Long:  `create a multi address`,
	Run: func(cmd *cobra.Command, args []string) {
		if m < 1 || m > 15 {
			fmt.Println("m value valid")
			return
		}

		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.MultiAccountSDK.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Operation Successed.\naddress:", msAccount.Address)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	createCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	createCmd.MarkFlagRequired("m")
	createCmd.MarkFlagRequired("publickeys")
}
