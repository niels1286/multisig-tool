package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/spf13/cobra"
	"strings"
)

var to string
var amount float64
var remark string
var assets string

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Assemble a transfer transaction",
	Long:  `根据参数组装一个转账交易，并返回交易hex`,
	Run: func(cmd *cobra.Command, args []string) {
		pkArray := strings.Split(pks, ",")
		if len(pkArray) < m {
			fmt.Println("Incorrect public keys")
			return
		}

		tx := utils.AssembleTransferTx(m, pkArray, amount, remark, to, 0, 0, nil)
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

	transferCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	transferCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	transferCmd.MarkFlagRequired("m")
	transferCmd.MarkFlagRequired("publickeys")

	transferCmd.Flags().StringVarP(&to, "to", "t", "", "转入地址")
	transferCmd.MarkFlagRequired("to")
	transferCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "金额，到账数量，以NULS为单位")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.Flags().StringVarP(&remark, "remark", "r", "", "交易备注，可以为空")

	transferCmd.Flags().StringVarP(&assets, "assets", "", "9-1", "资产标识,格式为chainId-assetsId，NVT:9-1,NULS:1-1")

}
