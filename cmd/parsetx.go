package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"github.com/niels1286/nerve-go-sdk/utils/mathutils"
	"github.com/niels1286/nerve-go-sdk/utils/seria"
	"github.com/spf13/cobra"
)

var txHex string

type TxInfo struct {
	Hash     string
	TxType   string
	TxData   map[string]string
	CoinData string
	Remark   string
}

var TypeMap = map[uint16]string{
	1:                                 "共识奖励",
	2:                                 "转账交易",
	5:                                 "质押交易",
	txprotocal.TX_TYPE_ACCOUNT_ALIAS:  "设置别名",
	txprotocal.TX_TYPE_CANCEL_DEPOSIT: "退出Staking",
	txprotocal.TX_TYPE_STOP_AGENT:     "停止节点",
	txprotocal.TX_TYPE_REGISTER_AGENT: "创建节点",
	txprotocal.APPEND_AGENT_DEPOSIT:   "追加保证金",
	txprotocal.REDUCE_AGENT_DEPOSIT:   "退出保证金",
}

func (ti *TxInfo) String() string {
	bus := "TxExtend:\n"
	for key, val := range ti.TxData {
		bus += "\t" + key + " : " + val + "\n"
	}
	value := fmt.Sprintf("===========tx info============\nhash:%s\ntype:%s\n%s%s\nRemark : %s", ti.Hash, ti.TxType, bus, ti.CoinData, ti.Remark)
	return value
}

// parsetxCmd represents the parsetx command
var parsetxCmd = &cobra.Command{
	Use:   "parsetx",
	Short: "Deserialize transactions to readable content",
	Long:  `Deserialize the transaction into readable content. Mainly focus on transaction type, coindata content and txdata content.`,
	Run: func(cmd *cobra.Command, args []string) {
		if "" == txHex {
			fmt.Println("txHex is valid.")
			return
		}
		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tx := txprotocal.ParseTransactionByReader(seria.NewByteBufReader(txBytes, 0))
		tx.CalcHash()
		info := getTxInfo(tx)
		fmt.Println(info.String())
	},
}
var timeStrArr = []string{"3月", "半年", "一年", "两年", "三年", "五年", "十年"}

func getTxInfo(tx *txprotocal.Transaction) *TxInfo {
	typeStr := TypeMap[tx.TxType]
	txData := map[string]string{}
	sdk := utils.GetOfficalSdk()
	switch tx.TxType {
	case txprotocal.TX_TYPE_DEPOSIT:
		deposit := &txdata.Staking{}
		deposit.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(deposit.Address)
		txData["amount"] = fmt.Sprintf("%d", deposit.Amount.Uint64()/100000000)
		timeStr := "活期"
		if deposit.DepositType != 0 {
			timeStr = timeStrArr[deposit.TimeType]
		}
		txData["time"] = timeStr
	case txprotocal.TX_TYPE_REGISTER_AGENT:
		agent := &txdata.CreateNode{}
		agent.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["createAddress"] = sdk.GetStringAddress(agent.AgentAddress)
		txData["packingAddress"] = sdk.GetStringAddress(agent.PackingAddress)
		txData["rewardAddress"] = sdk.GetStringAddress(agent.RewardAddress)
		txData["amount"] = fmt.Sprintf("%d", agent.Amount.Uint64()/100000000)
	case txprotocal.TX_TYPE_STOP_AGENT:
		info := &txdata.StopNode{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(info.Address)
		txData["nodeHash"] = info.AgentHash.String()
	case txprotocal.TX_TYPE_CANCEL_DEPOSIT:
		info := txdata.Withdraw{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(info.Address)
		txData["stakingTxHash"] = info.StakingTxHash.String()
	case txprotocal.TX_TYPE_ACCOUNT_ALIAS:
		info := &txdata.Alias{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(info.Address)
		txData["alias"] = info.Alias
	case txprotocal.APPEND_AGENT_DEPOSIT:
		info := &txdata.ChangeNodeDeposit{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(info.Address)
		txData["nodeHash"] = info.NodeHash.String()
		txData["amount"] = fmt.Sprintf("%d", info.Amount.Uint64()/100000000)
	case txprotocal.REDUCE_AGENT_DEPOSIT:
		info := &txdata.ChangeNodeDeposit{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = sdk.GetStringAddress(info.Address)
		txData["nodeHash"] = info.NodeHash.String()
		txData["amount"] = fmt.Sprintf("%d", info.Amount.Uint64()/100000000)
	default:
		if tx.Extend != nil {
			txData["hex"] = hex.EncodeToString(tx.Extend)
		}
	}
	cd := &txprotocal.CoinData{}
	cd.Parse(seria.NewByteBufReader(tx.CoinData, 0))
	var message = "From:\n"
	for _, from := range cd.Froms {
		nonce := hex.EncodeToString(from.Nonce)
		message += "\t" + sdk.GetStringAddress(from.Address) + "(" + fmt.Sprintf("%d", from.AssetsChainId) + "-" + fmt.Sprintf("%d", from.AssetsId) + ") (" + nonce + "):: " + mathutils.GetStringAmount(from.Amount, 8) + "\n"
	}
	message += "To:\n"
	for _, to := range cd.Tos {
		lock := fmt.Sprintf("%d", to.LockValue)
		if to.LockValue == uint64(18446744073709551615) {
			lock = "-1"
		}
		message += "\t" + sdk.GetStringAddress(to.Address) + "(" + fmt.Sprintf("%d", to.AssetsChainId) + "-" + fmt.Sprintf("%d", to.AssetsId) + ") :: " + mathutils.GetStringAmount(to.Amount, 8) + " (lock:" + lock + ")\n"
	}

	return &TxInfo{
		Hash:     tx.GetHash().String(),
		TxType:   typeStr,
		TxData:   txData,
		CoinData: message,
		Remark:   string(tx.Remark),
	}
}

func init() {
	rootCmd.AddCommand(parsetxCmd)
	parsetxCmd.Flags().StringVarP(&txHex, "txhex", "t", "", "Transaction serialization data in hexadecimal string format")
	parsetxCmd.MarkFlagRequired("txhex")
}
