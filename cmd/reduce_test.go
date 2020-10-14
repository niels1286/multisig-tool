// @Title
// @Description
// @Author  Niels  2020/10/13
package cmd

import (
	"encoding/hex"
	"github.com/niels1286/multisig-tool/utils"
	"testing"
)

func TestReduce(t *testing.T) {
	aph := "3e73f764492e95362cf325bd7168d145110a75e447510c927612586c06b23e91"
	bph := "6d10f3aa23018de6bc7d1ee52badd696f0db56082c62826ba822978fdf3a59fa"
	cph := "f7bb391ab82ba9ec7a552955b2fe50d79eea085d7571e5e2480d1777bc171f5e"
	sdk := utils.GetOfficalSdk()
	ap, _ := hex.DecodeString(aph)
	bp, _ := hex.DecodeString(bph)
	cp, _ := hex.DecodeString(cph)
	a, _ := sdk.AccountSDK.ImportAccount(ap)
	b, _ := sdk.AccountSDK.ImportAccount(bp)
	c, _ := sdk.AccountSDK.ImportAccount(cp)

	m = 2
	pks = a.GetPubKeyHex() + "," + b.GetPubKeyHex() + "," + c.GetPubKeyHex()
	nodeHash = "e07d6195b1f757a06fcb040e29f75e5a03149fc677d88f941e4eb724da82bae8"
	amount = 4000
	reduceCmd.Run(nil, nil)
}