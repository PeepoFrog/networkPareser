package main

import (
	"context"

	v3 "github.com/PeepoFrog/networkTreeParser/gatherer/v3"
)

func main() {
	ctx := context.Background()
	// gatherer.GetAllNodesV2(ctx, "148.251.69.56")
	// gatherer.GetAllNodesV2(ctx,"")
	// gatherer.GetAllNodesV3(ctx,"")
	v3.GetAllNodesV3(ctx, "148.251.69.56", 1, true)
	// log.Printf("total: %v, good: %v, bad: %v", len(node.Peers), goodPeers, errorPeers)

}
