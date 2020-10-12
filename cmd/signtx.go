package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
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
	Short: "sign a transaction",
	Long:  `对多签交易进行签名，当签名数量足够时，自动将交易广播到网络中`,
	Run: func(cmd *cobra.Command, args []string) {
		if "" == txHex {
			fmt.Println("txHex is valid.")
			return
		}
		if "" == prikeyHex && (keystore == "" || "" == password) {
			fmt.Println("need prikey")
			return
		}
		sdk := utils.GetOfficalSdk()

		nulsAccount, err := getAccount(sdk)

		if err != nil {
			fmt.Println("account wrong.")
			return
		}
		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			fmt.Println("txhex wrong.")
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
			fmt.Println("The address is not necessary")
			return
		}
		//签名
		hash, err := tx.GetHash().Serialize()
		if err != nil {
			fmt.Println("txhex wrong.")
			return
		}
		signData, err := sdk.AccountSDK.Sign(nulsAccount, hash)
		if err != nil {
			fmt.Println("sign failed.")
			return
		}
		sign := txprotocal.P2PHKSignature{
			SignValue: signData,
			PublicKey: nulsAccount.GetPubKey(),
		}
		txSign.Signatures = append(txSign.Signatures, sign)
		tx.SignData, err = txSign.Serialize()
		if err != nil {
			fmt.Println("sign failed.")
			return
		}
		resultBytes, err := tx.Serialize()
		if err != nil {
			fmt.Println("sign failed.")
			return
		}
		////判断是否需要广播
		if byte(len(txSign.Signatures)) >= txSign.M {
			sdk := utils.GetOfficalSdk()
			txHex := hex.EncodeToString(resultBytes)
			hash, err := sdk.Broadcast(txHex)
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				fmt.Println("Broadcast Success!\ntx hash : " + hash)
				return
			}
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
	signtxCmd.Flags().StringVarP(&txHex, "txhex", "t", "", "Transaction serialization data in hexadecimal string format")
	signtxCmd.MarkFlagRequired("txhex")

	signtxCmd.Flags().StringVarP(&prikeyHex, "prikey", "p", "", "签名使用的私钥，程序将自动验证其是否属于多签成员")

	signtxCmd.PersistentFlags().StringVarP(&keystore, "keystore", "k", "", "当不是用prikey时，可以指定同目录的keystore文件名")

	signtxCmd.Flags().StringVarP(&password, "password", "w", "", "使用keystore时，需要使用密码")

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
		fmt.Println("Using Account from:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err.Error())
	}
}
