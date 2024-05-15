package main

import (
	"context"

	"github.com/PeepoFrog/networkTreeParser/gatherer"
)

func main() {
	ctx := context.Background()
	gatherer.GetAllNodes(ctx, "148.251.69.56")
	// log.Printf("total: %v, good: %v, bad: %v", len(node.Peers), goodPeers, errorPeers)
}
