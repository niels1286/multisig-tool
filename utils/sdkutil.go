// @Title
// @Description
// @Author  Niels
package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/nerve-go-sdk/multisig"
	"github.com/niels1286/nerve-go-sdk/nerve"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/utils/mathutils"
	"github.com/niels1286/nerve-go-sdk/utils/rpc"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

func GetOfficalSdk() *nerve.NerveSDK {
	return nerve.GetSDK(cfg.ApiUrl, cfg.MainChainId, cfg.AddressPrefix)
}

func AssembleTransferTx(m int, pkArrayHex string, assetsChainId uint16, assetsId uint16, amount float64, remark string, to string, fromLocked byte, toLockValue uint64, nonce []byte, needFeeNonce bool) *txprotocal.Transaction {
	tx := txprotocal.Transaction{
		TxType:   txprotocal.TX_TYPE_TRANSFER,
		Time:     uint32(time.Now().Unix()),
		Remark:   []byte(remark),
		Extend:   nil,
		CoinData: nil,
		SignData: nil,
	}

	if m < 1 || m > 15 {
		fmt.Println(i18n.GetText("0020"))
		return nil
	}
	sdk := GetOfficalSdk()
	msAccount, err1 := sdk.MultiAccountSDK.CreateMultiAccount(m, pkArrayHex)
	if err1 != nil {
		fmt.Println(err1.Error())
		return nil
	}
	if nil == msAccount || "" == msAccount.Address {
		fmt.Println("")
		return nil
	}

	tx.CoinData = fillCoinData(sdk, msAccount, fromLocked, to, toLockValue, amount, assetsChainId, assetsId, nonce, needFeeNonce)

	pkArray := strings.Split(pkArrayHex, ",")
	publicKeys := [][]byte{}
	for _, pk := range pkArray {
		bytes, err := hex.DecodeString(pk)
		if err != nil {
			fmt.Println(i18n.GetText("0068"))
			return nil
		}
		publicKeys = append(publicKeys, bytes)
	}

	txSign := txprotocal.MultiAddressesSignData{
		M:              uint8(m),
		PubkeyList:     publicKeys,
		CommonSignData: txprotocal.CommonSignData{},
	}
	var err error
	tx.SignData, err = txSign.Serialize()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &tx
}

func AssembleTransferTxForReduce(m int, pkArrayHex string, remark string) *txprotocal.Transaction {
	tx := txprotocal.Transaction{
		TxType:   txprotocal.TX_TYPE_TRANSFER,
		Time:     uint32(time.Now().Unix()),
		Remark:   []byte(remark),
		Extend:   nil,
		CoinData: nil,
		SignData: nil,
	}

	if m < 1 || m > 15 {
		fmt.Println(i18n.GetText("0020"))
		return nil
	}
	sdk := GetOfficalSdk()
	msAccount, err1 := sdk.MultiAccountSDK.CreateMultiAccount(m, pkArrayHex)
	if err1 != nil {
		fmt.Println(err1.Error())
		return nil
	}
	if nil == msAccount || "" == msAccount.Address {
		fmt.Println("")
		return nil
	}

	pkArray := strings.Split(pkArrayHex, ",")
	publicKeys := [][]byte{}
	for _, pk := range pkArray {
		bytes, err := hex.DecodeString(pk)
		if err != nil {
			fmt.Println(i18n.GetText("0068"))
			return nil
		}
		publicKeys = append(publicKeys, bytes)
	}

	txSign := txprotocal.MultiAddressesSignData{
		M:              uint8(m),
		PubkeyList:     publicKeys,
		CommonSignData: txprotocal.CommonSignData{},
	}
	var err error
	tx.SignData, err = txSign.Serialize()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &tx
}
func fillCoinData(sdk *nerve.NerveSDK, msAccount *multisig.MultiAccount, fromLocked byte, to string, toLockValue uint64, amount float64, assetsChainId uint16, assetsId uint16, nonce []byte, feeNonce bool) []byte {
	value := decimal.NewFromFloat(amount)
	scale := cfg.AssetsMap[fmt.Sprintf("%d", assetsChainId)+"-"+fmt.Sprintf("%d", assetsId)]
	if scale == 0 {
		fmt.Println(i18n.GetText("0069"))
		return nil
	}
	value = value.Mul(decimal.NewFromFloat(math.Pow10(scale)))
	val := value.BigInt()

	coinData := &txprotocal.CoinData{
		Froms: nil,
		Tos:   nil,
	}

	if cfg.MainChainId == assetsChainId && assetsId == cfg.MainAssetsId {

		if nonce == nil {
			nonce = GetNonce(msAccount.Address, assetsChainId, assetsId)
		}

		fromVal := big.NewInt(100000)
		fromVal.Add(fromVal, val)
		from1 := txprotocal.CoinFrom{
			Coin: txprotocal.Coin{
				Address:       msAccount.AddressBytes,
				AssetsChainId: assetsChainId,
				AssetsId:      assetsId,
				Amount:        fromVal,
			},
			Nonce:  nonce,
			Locked: fromLocked,
		}
		coinData.Froms = []txprotocal.CoinFrom{from1}
	} else {

		var mainNonce []byte

		if nonce == nil {
			nonce = GetNonce(msAccount.Address, assetsChainId, assetsId)
			mainNonce = GetNonce(msAccount.Address, cfg.MainChainId, cfg.MainAssetsId)
		} else if feeNonce {
			mainNonce = GetNonce(msAccount.Address, cfg.MainChainId, cfg.MainAssetsId)
		} else {
			mainNonce = nonce
		}
		x := int64(100000)

		if msAccount.M > 7 {
			x = 2 * x
		}

		fromVal := big.NewInt(x)

		from1 := txprotocal.CoinFrom{
			Coin: txprotocal.Coin{
				Address:       msAccount.AddressBytes,
				AssetsChainId: assetsChainId,
				AssetsId:      assetsId,
				Amount:        val,
			},
			Nonce:  nonce,
			Locked: fromLocked,
		}
		from2 := txprotocal.CoinFrom{
			Coin: txprotocal.Coin{
				Address:       msAccount.AddressBytes,
				AssetsChainId: cfg.MainChainId,
				AssetsId:      cfg.MainAssetsId,
				Amount:        fromVal,
			},
			Nonce:  mainNonce,
			Locked: 0,
		}
		coinData.Froms = []txprotocal.CoinFrom{from1, from2}
	}
	toAddress, err2 := sdk.AccountSDK.GetBytesAddress(to)
	if nil != err2 {
		fmt.Println(err2.Error())
		return nil
	}
	to1 := txprotocal.CoinTo{
		Coin: txprotocal.Coin{
			Address:       toAddress,
			AssetsChainId: assetsChainId,
			AssetsId:      assetsId,
			Amount:        val,
		},
		LockValue: toLockValue,
	}
	coinData.Tos = []txprotocal.CoinTo{to1}

	result, _ := coinData.Serialize()

	return result
}

func GetNonce(address string, chainId uint16, assetsId uint16) []byte {
	sdk := GetOfficalSdk()
	status, err := sdk.ApiSDK.GetBalance(address, chainId, assetsId)
	if err != nil {
		return nil
	}
	if status == nil {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return status.Nonce
}

type ReduceNonce struct {
	Nonce  []byte
	Amount *big.Int
}

func GetReduceNonceList(nodeHash string, reduceAmount *big.Int) []*ReduceNonce {
	sdk := GetOfficalSdk()
	client := rpc.GetJsonRPCClient(sdk.GetApiUrl())
	param := client.NewRequestParam(rand.Intn(10000), "getReduceNonceList", []interface{}{cfg.MainChainId, nodeHash, 0, 1})
	result, err := client.Request(param)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if nil == result || nil == result.Result {
		return nil
	}
	list := []*ReduceNonce{}
	resultMap := result.Result.([]interface{})
	totalAmount := big.NewInt(0)
	for _, item := range resultMap {
		nonceStr := item.(map[string]interface{})["nonce"].(string)
		amountStr := item.(map[string]interface{})["deposit"].(string)

		nonceByte, _ := hex.DecodeString(nonceStr)
		amount, _ := mathutils.StringToBigInt(amountStr)
		list = append(list, &ReduceNonce{
			Nonce:  nonceByte,
			Amount: amount,
		})
		totalAmount = totalAmount.Add(totalAmount, amount)
		if totalAmount.Cmp(reduceAmount) >= 0 {
			break
		}
	}
	return list
}
