package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/niels1286/nerve-go-sdk/account"
	txprotocal "github.com/niels1286/nerve-go-sdk/tx/protocal"
	"github.com/niels1286/nerve-go-sdk/tx/txdata"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

var agentHash string

// depositCmd represents the deposit command
var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "委托",
	Long:  `委托一笔资产到节点上`,
	Run: func(cmd *cobra.Command, args []string) {
		if amount < 2000 || amount > 500000 {
			fmt.Println("委托金额不正确")
		}
		pkArray := strings.Split(pks, ",")
		if len(pkArray) < m {
			fmt.Println("Incorrect public keys")
			return
		}
		address := utils.CreateAddress(m, pkArray)
		tx := utils.AssembleTransferTx(m, pkArray, amount, "", address, 0, cfg.POCLockValue, nil)
		if tx == nil {
			return
		}
		tx.TxType = txprotocal.TX_TYPE_DEPOSIT
		value := big.NewFloat(amount)
		value = value.Mul(value, big.NewFloat(100000000))
		x, _ := value.Uint64()

		hashBytes, err := hex.DecodeString(agentHash)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ahash := txprotocal.NewNulsHash(hashBytes)

		acc, err := account.ParseAccount(address)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		depositData := txdata.Deposit{
			Amount:    big.NewInt(int64(x)),
			AgentHash: ahash,
			Address:   acc.AddressBytes,
		}
		tx.Extend, err = depositData.Serialize()
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
	rootCmd.AddCommand(depositCmd)
	depositCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	depositCmd.MarkFlagRequired("m")
	depositCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	depositCmd.MarkFlagRequired("publickeys")
	depositCmd.Flags().StringVarP(&agentHash, "agenthash", "h", "", "节点hash")
	depositCmd.MarkFlagRequired("agentHash")
	depositCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "委托金额")
	depositCmd.MarkFlagRequired("amount")
}
