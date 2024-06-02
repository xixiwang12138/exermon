package id

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var (
	snode *snowflake.Node
)

func init() {
	node := os.Getenv("NODE_ID")
	if node == "" {
		node = fmt.Sprintf("%d", rand.Intn(1024))
	}
	nid, err := strconv.ParseInt(node, 10, 64)
	if err != nil {
		panic(err)
	}
	snode, err = snowflake.NewNode(nid)
	if err != nil {
		panic(err)
	}
}

//func Generate() uint {
//	return uint(snode.Generate())
//}
