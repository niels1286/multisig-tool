// @Title
// @Description
// @Author  Niels
package cmd

import (
	"encoding/hex"
	"github.com/niels1286/multisig-tool/utils"
	"testing"
)

func TestStopNode(t *testing.T) {
	aph := "031c41ae3aa0463345b8e406805d48332c100b50706e6a324969c7dd9522b090fe"
	bph := "032e3a43a3b7949dfdff663467be6ebf06419afd9ed42bfb9d266e0c9d986d4cf1"
	cph := "03cedd7b86c823b365d192816f70663ff781f808cf266e52beddc43932afe9b339"
	sdk := utils.GetOfficalSdk()
	ap, _ := hex.DecodeString(aph)
	bp, _ := hex.DecodeString(bph)
	cp, _ := hex.DecodeString(cph)
	a, _ := sdk.AccountSDK.ImportAccount(ap)
	b, _ := sdk.AccountSDK.ImportAccount(bp)
	c, _ := sdk.AccountSDK.ImportAccount(cp)

	m = 2
	pks = a.GetPubKeyHex() + "," + b.GetPubKeyHex() + "," + c.GetPubKeyHex()
	nodeHash = "6eed8478564b1d0f7deb0e20b630c10c4c1fb90873af41f653d51f0484191d0e"
	stopNodeCmd.Run(nil, nil)
}
