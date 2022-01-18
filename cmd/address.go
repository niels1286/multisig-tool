package cmd

import (
	"fmt"

	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var addressCmd = &cobra.Command{
	Use: "address",
	Run: func(cmd *cobra.Command, args []string) {
		sdk := utils.GetOfficalSdk()
		account, err := sdk.AccountSDK.CreateAccount()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(account.GetAddr())
			fmt.Println("prikey: " + account.GetPriKeyHex())
			fmt.Println("pubkey: " + account.GetPubKeyHex())
		}
		fmt.Println("done!")
	},
}

func init() {
	rootCmd.AddCommand(addressCmd)
}
