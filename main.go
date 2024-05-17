package main

import (
	"context"
	"log"

	v3 "github.com/PeepoFrog/networkTreeParser/gatherer/v3"
)

func main() {
	ctx := context.Background()
	nodes, err := v3.GetAllNodesV3(ctx, "148.251.69.56", 1, true)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	for _, n := range nodes {
		log.Println(n)
	}

}
