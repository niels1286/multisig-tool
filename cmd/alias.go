package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"github.com/spf13/cobra"
)

var alias string

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "set alias",
	Long:  `Set alias for quick transfer and node name display`,
	Run: func(cmd *cobra.Command, args []string) {
		//todo 验证别名格式及是否重复

		amount = 1
		to = cfg.BlackHoleAddress
		tx := utils.AssembleTransferTx(m, pks, cfg.MainChainId, cfg.MainAssetsId, amount, "", to, 0, 0, nil)
		if tx == nil {
			return
		}
		tx.TxType = txprotocal.TX_TYPE_ACCOUNT_ALIAS
		sdk := utils.GetOfficalSdk()
		msAccount, _ := sdk.MultiAccountSDK.CreateMultiAccount(m, pks)
		aliasData := txdata.Alias{
			Address: msAccount.AddressBytes,
			Alias:   alias,
		}
		var err error
		tx.Extend, err = aliasData.Serialize()
		if err != nil {
			fmt.Println(err.Error())
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
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	aliasCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	aliasCmd.MarkFlagRequired("m")
	aliasCmd.MarkFlagRequired("publickeys")
	aliasCmd.Flags().StringVarP(&alias, "alias", "a", "", "别名，只允许小写字母和下划线")
	aliasCmd.MarkFlagRequired("alias")
}
