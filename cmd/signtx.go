package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/niels1286/nerve-go-sdk/acc"
	cryptoutils "github.com/niels1286/nerve-go-sdk/crypto/utils"
	"github.com/niels1286/nerve-go-sdk/nerve"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/utils/seria"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

var password string
var keystore string
var prikeyHex string

// signtxCmd represents the signtx command
var signtxCmd = &cobra.Command{
	Use:   "signtx",
	Short: i18n.GetText("0008"),
	Long:  i18n.GetText("0036"),
	Run: func(cmd *cobra.Command, args []string) {
		if "" == txHex {
			fmt.Println(i18n.GetText("0032"))
			return
		}
		if "" == prikeyHex && (keystore == "" || "" == password) {
			fmt.Println(i18n.GetText("0037"))
			return
		}
		sdk := utils.GetOfficalSdk()

		nulsAccount, err := getAccount(sdk)

		if err != nil {
			fmt.Println(i18n.GetText("0038"))
			return
		}
		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			fmt.Println(i18n.GetText("0032"))
			return
		}
		tx := txprotocal.ParseTransactionByReader(seria.NewByteBufReader(txBytes, 0))

		// 判断账户是否正确
		//将签名组装到交易中
		txSign := txprotocal.MultiAddressesSignData{}
		txSign.Parse(seria.NewByteBufReader(tx.SignData, 0))
		ok := false
		for _, pk := range txSign.PubkeyList {
			address := sdk.GetAddressByPubBytes(pk, 1)
			if reflect.DeepEqual(address, nulsAccount.GetAddrBytes()) {
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println(i18n.GetText("0039"))
			return
		}
		//签名
		hash, err := tx.GetHash().Serialize()
		if err != nil {
			fmt.Println(i18n.GetText("0032"))
			return
		}
		signData, err := sdk.AccountSDK.Sign(nulsAccount, hash)
		if err != nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		sign := txprotocal.P2PHKSignature{
			SignValue: signData,
			PublicKey: nulsAccount.GetPubKey(),
		}
		txSign.Signatures = append(txSign.Signatures, sign)
		tx.SignData, err = txSign.Serialize()
		if err != nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		resultBytes, err := tx.Serialize()
		if err != nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		////判断是否需要广播
		if byte(len(txSign.Signatures)) >= txSign.M {
			sdk := utils.GetOfficalSdk()
			txHex := hex.EncodeToString(resultBytes)
			hash, err := sdk.Broadcast(txHex)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("Failed!\nNewTxHex:", txHex)
				return
			} else {
				fmt.Println(i18n.GetText("10000") + "!\ntx hash : " + hash)
				return
			}

			fmt.Println("Success!\nNewTxHex:", txHex)
			return
		}
		resultHex := hex.EncodeToString(resultBytes)
		fmt.Println("Success!\nNewTxHex:", resultHex)
	},
}

func getAccount(sdk *nerve.NerveSDK) (acc.Account, error) {
	var prikey []byte
	if "" == prikeyHex {
		encryptedPrivateKey := viper.GetString("encryptedPrivateKey")
		data, err := hex.DecodeString(encryptedPrivateKey)
		if err != nil {
			return nil, err
		}
		err = errors.New("password may be wrong!")
		pwd := cryptoutils.Sha256h([]byte(password))
		prikey = cryptoutils.AESDecrypt(data, pwd)

	} else {
		prikey, _ = hex.DecodeString(prikeyHex)
	}
	nulsAccount, err := sdk.ImportAccount(prikey)
	return nulsAccount, err
}

func init() {
	rootCmd.AddCommand(signtxCmd)
	signtxCmd.Flags().StringVarP(&txHex, "txhex", "t", "", i18n.GetText("0033"))
	signtxCmd.MarkFlagRequired("txhex")

	signtxCmd.Flags().StringVarP(&prikeyHex, "prikey", "p", "", i18n.GetText("0048"))

	signtxCmd.PersistentFlags().StringVarP(&keystore, "keystore", "k", "", i18n.GetText("0049"))

	signtxCmd.Flags().StringVarP(&password, "password", "w", "", i18n.GetText("0050"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if keystore == "" {
		return
	}
	// Use config file from the flag.
	viper.SetConfigType("json")
	viper.SetConfigFile(keystore)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println(i18n.GetText("0051")+":", viper.ConfigFileUsed())
	} else {
		fmt.Println(err.Error())
	}
}
