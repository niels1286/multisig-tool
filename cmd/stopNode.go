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
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

var nodeHash string

// stopNodeCmd represents the stopNode command
var stopNodeCmd = &cobra.Command{
	Use:   "stopNode",
	Short: "注销节点",
	Long:  `注销节点，不再参与网络维护，不再得到奖励`,
	Run: func(cmd *cobra.Command, args []string) {
		if "" == strings.TrimSpace(nodeHash) {
			fmt.Println("agent hash 不能为空")
			return
		}
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		node, err := sdk.GetNode(cfg.PsUrl, nodeHash)
		if node == nil || err != nil {
			fmt.Println("网络超时导致操作失败，请重试")
			return
		}

		hash, _ := hex.DecodeString(nodeHash)

		tx := utils.AssembleTransferTxForReduce(m, pks, "")
		if tx == nil {
			fmt.Println("Failed!")
			return
		}
		tx.TxType = txprotocal.TX_TYPE_STOP_AGENT

		nonceList := utils.GetReduceNonceList(nodeHash, big.NewInt(1000000000000000000))

		totalFrom := big.NewInt(0)
		froms := []txprotocal.CoinFrom{}
		for _, item := range nonceList {
			totalFrom = totalFrom.Add(totalFrom, item.Amount)
			froms = append(froms, txprotocal.CoinFrom{
				Coin: txprotocal.Coin{
					Address:       msAccount.AddressBytes,
					AssetsChainId: cfg.MainChainId,
					AssetsId:      cfg.MainAssetsId,
					Amount:        item.Amount,
				},
				Nonce:  item.Nonce,
				Locked: 255,
			})
		}

		feeAmount := big.NewInt(int64(100000*int(len(froms)/7) + 100000))

		totalFrom.Sub(totalFrom, feeAmount)

		tos := []txprotocal.CoinTo{
			{
				Coin: txprotocal.Coin{
					Address:       msAccount.AddressBytes,
					AssetsChainId: cfg.MainChainId,
					AssetsId:      cfg.MainAssetsId,
					Amount:        totalFrom,
				},
				LockValue: uint64(tx.Time + uint32(15*24*3600)),
			},
		}

		coinData := &txprotocal.CoinData{
			Froms: froms,
			Tos:   tos,
		}
		tx.CoinData, _ = coinData.Serialize()

		txData := txdata.StopNode{
			Address:   msAccount.AddressBytes,
			AgentHash: txprotocal.NewNulsHash(hash),
		}
		tx.Extend, err = txData.Serialize()

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
	rootCmd.AddCommand(stopNodeCmd)
	stopNodeCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	stopNodeCmd.MarkFlagRequired("m")
	stopNodeCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	stopNodeCmd.MarkFlagRequired("publickeys")
	stopNodeCmd.Flags().StringVarP(&nodeHash, "nodeHash", "n", "", "节点hash")
	stopNodeCmd.MarkFlagRequired("nodeHash")

}
