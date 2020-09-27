package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"github.com/niels1286/nerve-go-sdk/utils/seria"
	"github.com/spf13/cobra"
	"math/big"
	"time"
)

var depositTxHash string

// withdrawCmd represents the withdraw command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "退出委托",
	Long:  `退出指定一笔委托，立即解锁对应的资产`,
	Run: func(cmd *cobra.Command, args []string) {
		hashBytes, err := hex.DecodeString(depositTxHash)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ahash := txprotocal.NewNulsHash(hashBytes)

		sdk := utils.GetOfficalSdk()

		txJson, err := sdk.GetTxJson(depositTxHash)
		if err != nil {
			fmt.Println("Can't find the deposit transaction.")
			return
		}
		//fmt.Println(txJson)
		txmap := map[string]interface{}{}
		json.Unmarshal([]byte(txJson), &txmap)
		txDataHex := txmap["txDataHex"].(string)
		if txDataHex == "" {
			fmt.Println("Failed to parse the deposit transaction.")
			return
		}
		txDataBytes, err := hex.DecodeString(txDataHex)
		if err != nil {
			fmt.Println("Failed to parse the deposit transaction.")
			return
		}
		depositData := txdata.Staking{}
		depositData.Parse(seria.NewByteBufReader(txDataBytes, 0))
		value := depositData.Amount.Div(depositData.Amount, big.NewInt(100000000))
		amount = float64(value.Uint64()) - 0.001

		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		toLockValue := uint64(0)
		if cfg.MainChainId == depositData.AssetsChainId && cfg.MainAssetsId == depositData.AssetsId {
			toLockValue = uint64(time.Now().Unix() + 7*24*3600)
		}
		tx := utils.AssembleTransferTx(m, pks, depositData.AssetsChainId, depositData.AssetsId, amount, "", msAccount.Address, cfg.POCLockValue, toLockValue, hashBytes[24:], true)
		if tx == nil {
			return
		}
		tx.TxType = txprotocal.TX_TYPE_CANCEL_DEPOSIT

		withdrawData := txdata.Withdraw{StakingTxHash: ahash}

		tx.Extend, err = withdrawData.Serialize()
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
	rootCmd.AddCommand(withdrawCmd)
	withdrawCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	withdrawCmd.MarkFlagRequired("m")
	withdrawCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	withdrawCmd.MarkFlagRequired("publickeys")
	withdrawCmd.Flags().StringVarP(&depositTxHash, "depositTxHash", "h", "", "委托交易的交易hash")
	withdrawCmd.MarkFlagRequired("depositTxHash")
}
