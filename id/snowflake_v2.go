package id

import (
	"os"
	"strconv"

	"github.com/yitter/idgenerator-go/idgen"
)

func init() {
	// 创建 IdGeneratorOptions 对象，可在构造函数中输入 WorkerId：
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "63"
	}
	nodeID_, _ := strconv.ParseInt(nodeID, 10, 64)
	var options = idgen.NewIdGeneratorOptions(uint16(nodeID_))
	// options.WorkerIdBitLength = 10  // 默认值6，限定 WorkerId 最大值为2^6-1，即默认最多支持64个节点。
	// options.SeqBitLength = 6; // 默认值6，限制每毫秒生成的ID个数。若生成速度超过5万个/秒，建议加大 SeqBitLength 到 10。
	// options.BaseTime = Your_Base_Time // 如果要兼容老系统的雪花算法，此处应设置为老系统的BaseTime。

	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)
}

func Generate() uint {
	return uint(idgen.NextId())
}
