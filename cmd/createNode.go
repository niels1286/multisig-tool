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
	Short: "创建节点",
	Long:  `创建共识节点，参与网络维护`,
	Run: func(cmd *cobra.Command, args []string) {
		if amount < 1000 || amount > 100000000 {
			fmt.Println("staking金额不正确")
			return
		}
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if packingAddress == msAccount.Address {
			fmt.Println("打包地址不能和创建地址一致")
			return
		}
		if packingAddress == rewardAddress {
			fmt.Println("打包地址不能和奖励地址一致")
			return
		}
		tx := utils.AssembleTransferTx(m, pks, cfg.MainChainId, cfg.MainAssetsId, amount, "", msAccount.Address, 0, cfg.POCLockValue, nil, false)
		if tx == nil {
			fmt.Println("Failed!")
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

	createNodeCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	createNodeCmd.MarkFlagRequired("m")
	createNodeCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	createNodeCmd.MarkFlagRequired("publickeys")

	createNodeCmd.Flags().StringVarP(&packingAddress, "packingAddress", "k", "", "节点打包地址，改地址必须放在节点钱包中")
	createNodeCmd.MarkFlagRequired("packingAddress")

	createNodeCmd.Flags().StringVarP(&rewardAddress, "rewardAddress", "r", "", "奖励地址，不填则默认为创建地址")

	createNodeCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "委托金额")
	createNodeCmd.MarkFlagRequired("amount")
}
