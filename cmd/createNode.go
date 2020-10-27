/*
Copyright © 2020 NAME HERE niels@nuls.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"math/big"

	"github.com/spf13/cobra"
)

var packingAddress string
var rewardAddress string

// createNodeCmd represents the createNode command
var createNodeCmd = &cobra.Command{
	Use:   "createNode",
	Short: i18n.GetText("0004"),
	Long:  i18n.GetText("0021"),
	Run: func(cmd *cobra.Command, args []string) {
		if amount < 200000 || amount > 100000000 {
			fmt.Println(i18n.GetText("0022"))
			return
		}
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if packingAddress == msAccount.Address {
			fmt.Println(i18n.GetText("0023"))
			return
		}
		if packingAddress == rewardAddress {
			fmt.Println(i18n.GetText("0024"))
			return
		}
		tx := utils.AssembleTransferTx(m, pks, cfg.MainChainId, cfg.MainAssetsId, amount, "", msAccount.Address, 0, cfg.POCLockValue, nil, false)
		if tx == nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		tx.TxType = txprotocal.TX_TYPE_REGISTER_AGENT
		value := big.NewFloat(amount)
		value = value.Mul(value, big.NewFloat(100000000))
		x, _ := value.Uint64()

		if "" == rewardAddress {
			rewardAddress = msAccount.Address
		}

		rAddress, _ := sdk.GetBytesAddress(rewardAddress)
		pAddress, _ := sdk.GetBytesAddress(packingAddress)

		node := txdata.CreateNode{
			Amount:         big.NewInt(int64(x)),
			AgentAddress:   msAccount.AddressBytes,
			RewardAddress:  rAddress,
			PackingAddress: pAddress,
		}
		tx.Extend, err = node.Serialize()

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
	rootCmd.AddCommand(createNodeCmd)

	createNodeCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	createNodeCmd.MarkFlagRequired("m")
	createNodeCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	createNodeCmd.MarkFlagRequired("publickeys")

	createNodeCmd.Flags().StringVarP(&packingAddress, "packingAddress", "k", "", i18n.GetText("0025"))
	createNodeCmd.MarkFlagRequired("packingAddress")

	createNodeCmd.Flags().StringVarP(&rewardAddress, "rewardAddress", "r", "", i18n.GetText("0026"))

	createNodeCmd.Flags().Float64VarP(&amount, "amount", "a", 0, i18n.GetText("0019"))
	createNodeCmd.MarkFlagRequired("amount")
}
