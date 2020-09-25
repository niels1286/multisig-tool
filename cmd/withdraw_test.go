// @Title
// @Description
// @Author  Niels
package cmd

import (
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/nerve-go-sdk/account"
	"testing"
)

func TestDWithdraw(t *testing.T) {

	ap := "3e73f764492e95362cf325bd7168d145110a75e447510c927612586c06b23e91"
	bp := "6d10f3aa23018de6bc7d1ee52badd696f0db56082c62826ba822978fdf3a59fa"
	cp := "f7bb391ab82ba9ec7a552955b2fe50d79eea085d7571e5e2480d1777bc171f5e"

	a, _ := account.GetAccountFromPrkey(ap, cfg.DefaultChainId, cfg.DefaultAddressPrefix)
	b, _ := account.GetAccountFromPrkey(bp, cfg.DefaultChainId, cfg.DefaultAddressPrefix)
	c, _ := account.GetAccountFromPrkey(cp, cfg.DefaultChainId, cfg.DefaultAddressPrefix)

	m = 2
	pks = a.GetPubKeyHex(true) + "," + b.GetPubKeyHex(true) + "," + c.GetPubKeyHex(true)
	depositTxHash = "3d95009d2b62b7ba9460c0f6ab44a487907f69db077617dc11ee0c9df4f8f029"
	withdrawCmd.Run(nil, nil)
}
