// @Title
// @Description
// @Author  Niels
package cmd

import (
	"encoding/hex"
	"github.com/niels1286/multisig-tool/utils"
	"testing"
)

func TestDWithdraw(t *testing.T) {

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
	depositTxHash = "3d95009d2b62b7ba9460c0f6ab44a487907f69db077617dc11ee0c9df4f8f029"
	withdrawCmd.Run(nil, nil)
}
