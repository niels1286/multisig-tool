// @Title
// @Description
// @Author  Niels  2020/9/25
package cfg

const ApiUrl = "http://beta.api.nerve.network/jsonrpc/"
const PsUrl = "http://beta.public.nerve.network/"
const MainChainId = uint16(5)
const MainAssetsId = uint16(1)
const AddressPrefix = "TNVT"
const BlackHoleAddress = "TNVTdTSPGwjgRMtHqjmg8yKeMLnpBpVN5ZuuY"

//const ApiUrl = "https://api.nerve.network/jsonrpc/"
//const PsUrl = "https://public.nerve.network/"
//const MainChainId = uint16(9)
//const MainAssetsId = uint16(1)
//const AddressPrefix = "NERVE"
//const BlackHoleAddress = "NERVEepb63T1M8JgQ26jwZpZXYL8ZMLdUAK31L"

const POCLockValue = 18446744073709551615

var AssetsMap = map[string]int{"2-1": 8, "5-1": 8, "5-2": 18, "5-6": 6, "5-7": 6, "5-8": 18, "5-9": 18}

//var AssetsMap = map[string]int{"1-1": 8, "9-1": 8, "9-2": 18, "9-3": 6, "9-5": 6, "9-6": 18, "9-7": 18}
