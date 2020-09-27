// @Title
// @Description
// @Author  Niels
package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/nerve-go-sdk/multisig"
	"github.com/niels1286/nerve-go-sdk/nerve"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"math/big"
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
		fmt.Println("m value valid")
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

	fillCoinData(tx, sdk, msAccount, fromLocked, to, toLockValue, amount, assetsChainId, assetsId, nonce, needFeeNonce)

	pkArray := strings.Split(pkArrayHex, ",")
	publicKeys := [][]byte{}
	for _, pk := range pkArray {
		bytes, err := hex.DecodeString(pk)
		if err != nil {
			fmt.Println("public key not right.")
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

func fillCoinData(tx txprotocal.Transaction, sdk *nerve.NerveSDK, msAccount *multisig.MultiAccount, fromLocked byte, to string, toLockValue uint64, amount float64, assetsChainId uint16, assetsId uint16, nonce []byte, feeNonce bool) {
	value := big.NewFloat(amount)
	value = value.Mul(value, big.NewFloat(100000000))
	x, _ := value.Uint64()
	val := new(big.Int)
	val.SetUint64(x)

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
		return
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

	var err error
	tx.CoinData, err = coinData.Serialize()
	if err != nil {
		fmt.Println(err.Error())
	}
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
