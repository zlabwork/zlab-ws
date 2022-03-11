package app

import (
	"github.com/bwmarrin/snowflake"
	"os"
	"strconv"
)

var Libs *libs

type libs struct {
	Snow *snowflake.Node
}

func NewLibs() *libs {
	i, _ := strconv.ParseInt(os.Getenv("APP_NODE"), 10, 64)
	snowflake.Epoch = 1498612200000
	node, _ := snowflake.NewNode(i)
	return &libs{
		Snow: node,
	}
}
