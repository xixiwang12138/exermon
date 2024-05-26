package xlog

import (
	"encoding/json"
	"strings"
)

func replaceJsonHolder(format string, args []interface{}) (string, []interface{}) {
	nodes := findPercentIndices(format)
	deletedArgsIndex := make(map[int]bool)

	sb := strings.Builder{}
	sb.Grow(2 * len(format))

	beforeJsonHolderIdx := 0
	for i, node := range nodes {
		if node.percentFormat == 'j' || node.percentFormat == 'J' {
			sb.Write([]byte(format[beforeJsonHolderIdx:node.percentByteIndex]))
			jsonStr, err := json.Marshal(args[i])
			if err != nil {
				continue
			}
			sb.Write(jsonStr)
			beforeJsonHolderIdx = node.percentByteIndex + 2

			deletedArgsIndex[i] = true
		}
	}
	if len(deletedArgsIndex) == 0 {
		return format, args
	}
	sb.Write([]byte(format[beforeJsonHolderIdx:]))
	newArgs := make([]interface{}, 0, len(args)-len(deletedArgsIndex))
	for i, arg := range args {
		if _, ok := deletedArgsIndex[i]; !ok {
			newArgs = append(newArgs, arg)
		}
	}
	return sb.String(), newArgs
}

type percentNode struct {
	percentByteIndex int  // %的位置
	percentFormat    byte // %后面的格式化字符
}

// findPercentIndices 查找格式化字符串中所有%的位置
func findPercentIndices(format string) []percentNode {
	var indices []percentNode

	// 初始化起始位置
	startIndex := 0

	for {
		// 查找下一个%的位置
		index := strings.Index(format[startIndex:], "%")
		if index == -1 {
			// 没有找到更多%
			break
		}

		// 将相对位置转换为绝对位置
		index += startIndex

		// %在最后一个位置，不需要处理
		if index == len(format)-1 {
			break
		}
		// 添加找到的%的位置到切片中
		indices = append(indices, percentNode{
			percentByteIndex: index,
			percentFormat:    format[index+1],
		})

		// 更新起始位置，以便下一次查找
		startIndex = index + 1
	}

	return indices
}
