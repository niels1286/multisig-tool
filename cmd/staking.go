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
	"strconv"
	"strings"
)

var timeType byte

// depositCmd represents the deposit command
var depositCmd = &cobra.Command{
	Use:   "staking",
	Short: "质押",
	Long:  `质押资产获取收益`,
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
			fmt.Println("Failed!")
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
	depositCmd.Flags().IntVarP(&m, "m", "m", 0, "发起交易的最小签名个数")
	depositCmd.MarkFlagRequired("m")
	depositCmd.Flags().StringVarP(&pks, "publickeys", "p", "", "多签地址的成员公钥，以','分隔不同的公钥")
	depositCmd.MarkFlagRequired("publickeys")

	depositCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "委托金额")
	depositCmd.MarkFlagRequired("amount")

	depositCmd.Flags().StringVarP(&assets, "assets", "", "9-1", "资产标识,格式为chainId-assetsId，NVT:9-1,NULS:1-1")
	depositCmd.Flags().Uint8VarP(&timeType, "timeType", "", 0, "质押时间类型：0-活期，1-3月，2-6月，3-1年，4-2年，5-3年，6-5年，7-10年")
}
