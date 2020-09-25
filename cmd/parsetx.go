package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/nerve-go-sdk/account"
	txprotocal "github.com/niels1286/nerve-go-sdk/tx/protocal"
	"github.com/niels1286/nerve-go-sdk/tx/txdata"
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
	5:                                 "委托交易",
	txprotocal.TX_TYPE_ACCOUNT_ALIAS:  "设置别名",
	txprotocal.TX_TYPE_CANCEL_DEPOSIT: "退出委托",
	txprotocal.TX_TYPE_STOP_AGENT:     "停止节点",
	txprotocal.TX_TYPE_REGISTER_AGENT: "创建节点",
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

func getTxInfo(tx *txprotocal.Transaction) *TxInfo {
	typeStr := TypeMap[tx.TxType]
	txData := map[string]string{}
	switch tx.TxType {
	case txprotocal.TX_TYPE_DEPOSIT:
		deposit := &txdata.Deposit{}
		deposit.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = account.GetStringAddress(deposit.Address, cfg.DefaultAddressPrefix)
		txData["agentHash"] = deposit.AgentHash.String()
		txData["amount"] = fmt.Sprintf("%d", deposit.Amount.Uint64()/100000000)
	case txprotocal.TX_TYPE_REGISTER_AGENT:
		agent := &txdata.Agent{}
		agent.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["agentAddress"] = account.GetStringAddress(agent.AgentAddress, cfg.DefaultAddressPrefix)
		txData["packingAddress"] = account.GetStringAddress(agent.PackingAddress, cfg.DefaultAddressPrefix)
		txData["rewardAddress"] = account.GetStringAddress(agent.RewardAddress, cfg.DefaultAddressPrefix)
		txData["amount"] = fmt.Sprintf("%d", agent.Amount.Uint64()/100000000)
		txData["commissionRate"] = fmt.Sprintf("%d", agent.CommissionRate)
	case txprotocal.TX_TYPE_STOP_AGENT:
		info := &txdata.StopAgent{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["agentHash"] = info.AgentHash.String()
	case txprotocal.TX_TYPE_CANCEL_DEPOSIT:
		info := txdata.Withdraw{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["depositTxHash"] = info.DepositTxHash.String()
	case txprotocal.TX_TYPE_ACCOUNT_ALIAS:
		info := &txdata.Alias{}
		info.Parse(seria.NewByteBufReader(tx.Extend, 0))
		txData["address"] = account.GetStringAddress(info.Address, cfg.DefaultAddressPrefix)
		txData["alias"] = info.Alias
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
		message += "\t" + account.GetStringAddress(from.Address, cfg.DefaultAddressPrefix) + "(" + fmt.Sprintf("%d", from.AssetsChainId) + "-" + fmt.Sprintf("%d", from.AssetsId) + ") (" + nonce + "):: " + mathutils.GetStringAmount(from.Amount, 8) + "\n"
	}
	message += "To:\n"
	for _, to := range cd.Tos {
		lock := fmt.Sprintf("%d", to.LockValue)
		if to.LockValue == uint64(18446744073709551615) {
			lock = "-1"
		}
		message += "\t" + account.GetStringAddress(to.Address, cfg.DefaultAddressPrefix) + "(" + fmt.Sprintf("%d", to.AssetsChainId) + "-" + fmt.Sprintf("%d", to.AssetsId) + ") :: " + mathutils.GetStringAmount(to.Amount, 8) + " (lock:" + lock + ")\n"
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
