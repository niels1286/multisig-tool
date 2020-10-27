package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"github.com/spf13/cobra"
	"math/big"
	"strconv"
	"strings"
)

var timeType byte

// depositCmd represents the deposit command
var depositCmd = &cobra.Command{
	Use:   "staking",
	Short: i18n.GetText("0009"),
	Long:  i18n.GetText("0052"),
	Run: func(cmd *cobra.Command, args []string) {
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

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

		tx := utils.AssembleTransferTx(m, pks, cId, aId, amount, "", msAccount.Address, 0, cfg.POCLockValue, nil, false)
		if tx == nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		tx.TxType = txprotocal.TX_TYPE_DEPOSIT
		value := big.NewFloat(amount)
		value = value.Mul(value, big.NewFloat(100000000))
		x, _ := value.Uint64()

		dType := uint8(0)
		tType := uint8(0)

		if timeType > 0 {
			dType = 1
			tType = timeType - 1
		}

		depositData := txdata.Staking{
			Amount:        big.NewInt(int64(x)),
			Address:       msAccount.AddressBytes,
			AssetsChainId: cId,
			AssetsId:      aId,
			DepositType:   dType,
			TimeType:      tType,
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
	depositCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	depositCmd.MarkFlagRequired("m")
	depositCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	depositCmd.MarkFlagRequired("publickeys")

	depositCmd.Flags().Float64VarP(&amount, "amount", "a", 0, i18n.GetText("0019"))
	depositCmd.MarkFlagRequired("amount")

	depositCmd.Flags().StringVarP(&assets, "assets", "", "9-1", i18n.GetText("0053"))
	depositCmd.Flags().Uint8VarP(&timeType, "timeType", "", 0, i18n.GetText("0054"))
}
