package cmd

import (
	"fmt"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
	"strings"
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
		pkArray := strings.Split(pks, ",")
		if len(pkArray) < m {
			fmt.Println("Incorrect public keys")
			return
		}
		address := utils.CreateAddress(m, pkArray)
		fmt.Println("Operation Successed.\naddress:", address)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	createCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	createCmd.MarkFlagRequired("m")
	createCmd.MarkFlagRequired("publickeys")
}
