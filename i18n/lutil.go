// @Title
// @Description
// @Author  Niels  2020/10/22
package i18n

import "time"

var currentMap map[string]string

var cnMap = map[string]string{
	"0001": "设置别名",
	"0002": "追加保证金",
	"0003": "生成多签地址",
	"0004": "创建节点",
	"0005": "帮助",
	"0006": "解析交易",
	"0007": "减少保证金",
	"0008": "签名交易",
	"0009": "质押",
	"0010": "注销节点",
	"0011": "转账",
	"0012": "退出质押",
	"0013": "设置别名，用于标识账户地址，可以用于快捷转账和节点名称",
	"0014": "发起交易的最小签名个数",
	"0015": "多签地址的成员公钥，以','分隔不同的公钥",
	"0016": "别名，只允许小写字母和下划线",
	"0017": "节点hash不能为空",
	"0018": "节点hash",
	"0019": "委托金额",
	"0020": "m值不正确",
	"0021": `创建共识节点，参与网络维护`,
	"0022": "保证金金额不正确",
	"0023": "打包地址不能和创建地址一致",
	"0024": "打包地址不能和奖励地址一致",
	"0025": "节点打包地址，该地址必须放在节点钱包中",
	"0026": "奖励地址，不填则默认为创建地址",
	"0027": "共识奖励",
	"0028": "转账交易",
	"0029": "质押交易",
	"0030": "反序列化交易为可读内容",
	"0031": `反序列化交易为可读内容,聚焦交易类型、Coin数据、业务数据、备注内容`,
	"0032": "交易Hex不正确",
	"0033": "十六进制字符串格式的事务序列化数据",
	"0034": "NULS多签工具",
	"0035": "NULS多签工具，包含账户创建、转账、质押、节点管理、签名功能",
	"0036": "对多签交易进行签名，当签名数量足够时，自动将交易广播到网络中",
	"0037": "私钥不能为空",
	"0038": "账户错误",
	"0039": "这个私钥不是必须的",
	"0040": "3月", "0041": "半年", "0042": "一年", "0043": "两年", "0044": "三年", "0045": "五年", "0046": "十年",
	"0047":  "活期",
	"0048":  "签名使用的私钥，程序将自动验证其是否属于多签成员",
	"0049":  "当不使用prikey时，可以指定同目录的keystore文件名",
	"0050":  "使用keystore时，需要使用密码",
	"0051":  "账户",
	"0052":  `质押资产获取收益`,
	"0053":  "资产标识,格式为chainId-assetsId，NVT:9-1,NULS:1-1",
	"0054":  "质押时间类型：0-活期，1-3月，2-6月，3-1年，4-2年，5-3年，6-5年，7-10年",
	"0055":  "注销节点，不再参与网络维护，不再得到奖励",
	"0056":  "节点hash 不能为空",
	"0057":  "网络超时导致操作失败，请重试",
	"0058":  "根据参数组装一个转账交易，并返回交易hex",
	"0059":  "目标地址",
	"0060":  "金额，到账数量，以NULS为单位",
	"0061":  "交易备注，可以为空",
	"0062":  "资产标识,格式为chainId-assetsId，NVT:9-1,NULS:1-1",
	"0063":  "退出指定一笔委托，立即解锁对应的资产",
	"0064":  "找不到质押交易",
	"0065":  "查询的质押交易解析失败",
	"0066":  "Hex格式不正确",
	"0067":  "委托交易的交易hash",
	"0068":  "公钥不正确",
	"0069":  "未能获取到该资产的小数位数",
	"10000": "成功",
	"10001": "失败",
}
var enMap = map[string]string{
	"0001": "set alias",
	"0002": "Additional margin for node",
	"0003": "Generate multi signature address",
	"0004": "create consensus node",
	"0005": "Help",
	"0006": "Deserialize transactions to readable content",
	"0007": "Reduce node margin",
	"0008": "sign a transaction",
	"0009": "Income from staking",
	"0010": "Cancel consensus node",
	"0011": "Assemble a transfer transaction",
	"0012": "Exit staking",
	"0013": "Set alias for quick transfer and node name display",
	"0014": "minimum number of signatures to initiate a transaction",
	"0015": "public keys of members with multiple signature addresses, separating different public keys with ','",
	"0016": "alias, only lowercase letters and underscores are allowed",
	"0017": "node hash cannot be empty",
	"0018": "node hash",
	"0019": "entrusted amount",
	"0020": "m value valid",
	"0021": "creating consensus nodes and participating in network maintenance",
	"0022": "the amount of deposit is incorrect",
	"0023": "the package address cannot be the same as the created address,",
	"0024": "the package address cannot be the same as the reward address.",
	"0025": "node package address, which must be put in the node wallet",
	"0026": "award address. If it is not filled in, it will be the created address by default.",
	"0027": "consensus award",
	"0028": "transfer transactions",
	"0029": "staking transaction",
	"0030": "Deserialize transactions to readable content",
	"0031": `Deserialize the transaction into readable content. Mainly focus on transaction type, coindata content and txdata content.`,
	"0032": "txHex is valid.",
	"0033": "Transaction serialization data in hexadecimal string format",
	"0034": "NULS MultiAddress Tools",
	"0035": "nuls multi signature tool basic fund, including account creation and transfer signature functions",
	"0036": "sign multiple transactions. When the number of signatures is sufficient, the transactions will be broadcast to the network automatically.",
	"0037": "PrikeyHex can not be nil",
	"0038": "account wrong.",
	"0039": "The address is not necessary",
	"0040": "March", "0041": "half a year", "0042": "one year", "0043": "two years", "0044": "three years", "0045": "five years", "0046": "ten years",
	"0047":  "current",
	"0048":  "the private key used in the signature, the program will automatically verify whether it belongs to a multi signer member",
	"0049":  "when prikey is not used, you can specify the keystore file name of the same directory.",
	"0050":  "when using keystore, you need to use a password",
	"0051":  "Using Account from",
	"0052":  "income from stakingd assets",
	"0053":  "asset ID, in the form of chainid assetsid, NVT:9-1 , NULS:1-1 ",
	"0054":  "staking time type: 0-current, 1-3 months, 2-6 months, 3-1 years, 4-2 years, 5-3 years, 6-5 years, 7-10 years",
	"0055":  "log off a node, no longer participate in network maintenance, no longer get rewards.",
	"0056":  "agent hash cannot be empty",
	"0057":  "the operation failed due to network timeout, please try again",
	"0058":  "assemble a transfer transaction according to the parameters and return the transaction hex",
	"0059":  "target address",
	"0060":  "amount, received quantity, in nuls unit",
	"0061":  "transaction remarks, which can be blank",
	"0062":  "asset ID, in the form of chainid assetsid, NVT:9-1 , NULS:1-1 ",
	"0063":  "exit the specified delegation and unlock the corresponding asset immediately",
	"0064":  "Can't find the deposit transaction.",
	"0065":  "query staking transaction resolution failed",
	"0066":  "Failed to parse the deposit transaction.",
	"0067":  "transaction hash of entrusted transaction",
	"0068":  "public key not right.",
	"0069":  "the decimal digits of the asset could not be obtained,",
	"10000": "Success",
	"10001": "Failed",
}

func initLang(key string) {
	switch key {
	case "cn":
		currentMap = cnMap
		return
	case "en":
		currentMap = enMap
		return
	}
}

func GetText(code string) string {
	if len(currentMap) == 0 {
		initLang(getLangType())
	}

	return currentMap[code]
}

func getLangType() string {
	zone, offset := time.Now().Local().Zone()
	if zone == "CST" && offset == 28800 {
		return "cn"
	}
	return "en"
}
