package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var to string
var amount float64
var remark string
var assets string

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: i18n.GetText("0011"),
	Long:  i18n.GetText("0058"),
	Run: func(cmd *cobra.Command, args []string) {
		cId := cfg.MainChainId
		aId := cfg.MainAssetsId

		if "" != assets {
			arr := strings.Split(assets, "-")
			val, err := strconv.Atoi(arr[0])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			cId = uint16(val)
			val2, err := strconv.Atoi(arr[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			aId = uint16(val2)
		}
		tx := utils.AssembleTransferTx(m, pks, cId, aId, amount, remark, to, 0, 0, nil, false)
		if tx == nil {
			return
		}
		txBytes, err := tx.Serialize()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		txHex := hex.EncodeToString(txBytes)

		fmt.Println("Successed:\ntxHex : " + txHex)
		fmt.Println("txHash : " + tx.GetHash().String())
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	transferCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	transferCmd.MarkFlagRequired("m")
	transferCmd.MarkFlagRequired("publickeys")

	transferCmd.Flags().StringVarP(&to, "to", "t", "", i18n.GetText("0059"))
	transferCmd.MarkFlagRequired("to")
	transferCmd.Flags().Float64VarP(&amount, "amount", "a", 0, i18n.GetText("0060"))
	transferCmd.MarkFlagRequired("amount")
	transferCmd.Flags().StringVarP(&remark, "remark", "r", "", i18n.GetText("0061"))

	transferCmd.Flags().StringVarP(&assets, "assets", "s", "", i18n.GetText("0062"))

}
