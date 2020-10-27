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
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

var reduceCmd = &cobra.Command{
	Use:   "reduce",
	Short: i18n.GetText("0007"),
	Long:  i18n.GetText("0007"),
	Run: func(cmd *cobra.Command, args []string) {
		if "" == nodeHash || strings.TrimSpace(nodeHash) == "" {
			fmt.Println(i18n.GetText("0017"))
			return
		}
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		cId := cfg.MainChainId
		aId := cfg.MainAssetsId
		tx := utils.AssembleTransferTxForReduce(m, pks, "")
		if tx == nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		tx.TxType = txprotocal.REDUCE_AGENT_DEPOSIT

		dec := decimal.NewFromFloat(amount)
		dec = dec.Mul(decimal.New(1, 8))
		x := dec.BigInt()

		nonceList := utils.GetReduceNonceList(nodeHash, x)

		totalFrom := big.NewInt(0)
		froms := []txprotocal.CoinFrom{}
		for _, item := range nonceList {
			totalFrom = totalFrom.Add(totalFrom, item.Amount)
			froms = append(froms, txprotocal.CoinFrom{
				Coin: txprotocal.Coin{
					Address:       msAccount.AddressBytes,
					AssetsChainId: cId,
					AssetsId:      aId,
					Amount:        item.Amount,
				},
				Nonce:  item.Nonce,
				Locked: 255,
			})
		}

		lockAmount := totalFrom.Sub(totalFrom, x)
		feeAmount := big.NewInt(int64(100000*int(len(froms)/7) + 100000))

		timeLockAmount := big.NewInt(0).Add(big.NewInt(0), x)
		timeLockAmount.Sub(timeLockAmount, feeAmount)

		tos := []txprotocal.CoinTo{
			{
				Coin: txprotocal.Coin{
					Address:       msAccount.AddressBytes,
					AssetsChainId: cId,
					AssetsId:      aId,
					Amount:        timeLockAmount,
				},
				LockValue: uint64(tx.Time + uint32(15*24*3600)),
			},
		}
		if lockAmount.Cmp(big.NewInt(0)) > 0 {
			tos = append(tos, txprotocal.CoinTo{
				Coin: txprotocal.Coin{
					Address:       msAccount.AddressBytes,
					AssetsChainId: cId,
					AssetsId:      aId,
					Amount:        lockAmount,
				},
				LockValue: cfg.POCLockValue,
			})
		}

		coinData := &txprotocal.CoinData{
			Froms: froms,
			Tos:   tos,
		}
		tx.CoinData, _ = coinData.Serialize()
		depositData := txdata.ChangeNodeDeposit{
			Amount:   x,
			Address:  msAccount.AddressBytes,
			NodeHash: txprotocal.ImportNulsHash(nodeHash),
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
	rootCmd.AddCommand(reduceCmd)
	reduceCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	reduceCmd.MarkFlagRequired("m")
	reduceCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	reduceCmd.MarkFlagRequired("publickeys")
	reduceCmd.Flags().Float64VarP(&amount, "amount", "a", 0, i18n.GetText("0019"))
	reduceCmd.MarkFlagRequired("amount")
	reduceCmd.Flags().StringVarP(&nodeHash, "nodeHash", "n", "", i18n.GetText("0018"))
	reduceCmd.MarkFlagRequired("nodeHash")

}
