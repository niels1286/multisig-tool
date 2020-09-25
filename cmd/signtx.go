package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/utils"
	"github.com/niels1286/nerve-go-sdk/account"
	txprotocal "github.com/niels1286/nerve-go-sdk/tx/protocal"
	"github.com/niels1286/nerve-go-sdk/utils/seria"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
		nulsAccount, err := getAccount()

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
			address := account.GetAddressByPubBytes(pk, cfg.DefaultChainId, account.NormalAccountType, cfg.DefaultAddressPrefix)
			if reflect.DeepEqual(address, nulsAccount.AddressBytes) {
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
		signData, err := nulsAccount.Sign(hash)
		if err != nil {
			fmt.Println("sign failed.")
			return
		}
		sign := txprotocal.P2PHKSignature{
			SignValue: signData,
			PublicKey: nulsAccount.GetPubKeyBytes(true),
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

			hash, err := sdk.BroadcastTx(resultBytes)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Success!\ntx hash : " + hash)
		}
		resultHex := hex.EncodeToString(resultBytes)

		fmt.Println("Success!\nNewTxHex:", resultHex)
	},
}

func getAccount() (*account.Account, error) {
	if "" != prikeyHex {
		nulsAccount, err := account.GetAccountFromPrkey(prikeyHex, cfg.DefaultChainId, cfg.DefaultAddressPrefix)
		if err != nil {
			return nil, err
		}
		return nulsAccount, nil
	} else {
		ks := account.KeyStore{
			Address:             viper.GetString("address"),
			EncryptedPrivateKey: viper.GetString("encryptedPrivateKey"),
			Pubkey:              viper.GetString("pubkey"),
		}
		return ks.GetAccount(password, cfg.DefaultChainId, cfg.DefaultAddressPrefix)
	}
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
	if keystore != "" {
		// Use config file from the flag.
		viper.SetConfigType("json")
		viper.SetConfigFile(keystore)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".nmt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using Account from:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err.Error())
	}
}
