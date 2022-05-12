package app

import (
	"github.com/bwmarrin/snowflake"
)

var Libs *libs

type libs struct {
	Snow *snowflake.Node
}

func NewLibs(id int32) *libs {
	snowflake.Epoch = 1498612200000
	node, _ := snowflake.NewNode(int64(id))
	return &libs{
		Snow: node,
	}
}
