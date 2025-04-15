package lib

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func GetNodeID() int {
	conf := GetConfig()
	if conf.SnowflakeNodeId != "" {
		id, err := strconv.Atoi(conf.SnowflakeNodeId)
		if err == nil && id >= 0 && id < 1024 {
			log.Println("Using environment variable for Snowflake Node ID:", id)
			return id
		}
	}

	// fallback: IP hash
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip := ipnet.IP.To4()
			r := int(binary.BigEndian.Uint32(ip)) % 1024
			log.Println("Using fallback IP hash as Snowflake Node ID:", r)
			return r
		}
	}

	panic("No valid node ID could be determined")
}

func InitSnowflakeNode() {
	n, err := snowflake.NewNode(int64(GetNodeID()))
	if err != nil {
		panic("Unable to initialize snowflake generator")
	}
	node = n
}

func GenerateSnowflakeID() snowflake.ID {
	if node == nil {
		InitSnowflakeNode()
	}
	return node.Generate()
}
