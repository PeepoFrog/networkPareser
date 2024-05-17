package v3

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/PeepoFrog/networkTreeParser/parser"
)

var mu sync.Mutex

type Node struct {
	IP       string
	P2P_port int
	ID       string
}

func GetAllNodesV3(ctx context.Context, firstNode string, depth int, ignoreDepth bool) (map[string]string, error) {
	nodesPool := make(map[string]string)
	blacklist := make(map[string]string)
	processed := make(map[string]string)
	node, err := parser.GetNetInfoFromInterx(ctx, firstNode)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, n := range node.Peers {
		log.Println(n.RemoteIP)
		wg.Add(1)
		go testLoop(ctx, &wg, nodesPool, blacklist, processed, n.RemoteIP, 0, depth, ignoreDepth)
	}

	wg.Wait()
	fmt.Println()
	log.Printf("\nTotal saved peers:%v\nOriginal node peer count: %v\nBlacklisted nodes(not reachable): %v\n", len(nodesPool), len(node.Peers), len(blacklist))

	return nodesPool, nil
}

func testLoop(ctx context.Context, wg *sync.WaitGroup, pool, blacklist, processed map[string]string, ip string, currentDepth, totalDepth int, ignoreDepth bool) {

	defer wg.Done()
	if !ignoreDepth {
		if currentDepth >= totalDepth {
			log.Println("DEPTH LIMIT REACHED")
			return
		}
	}

	log.Printf("Current depth: %v, IP: %v", currentDepth, ip)

	mu.Lock()
	if _, exist := blacklist[ip]; exist {
		mu.Unlock()
		log.Printf("BLACKLISTED: %v", ip)
		return
	}
	if _, exist := pool[ip]; exist {
		mu.Unlock()
		log.Printf("ALREADY EXIST: %v", ip)
		return
	}
	if _, exist := processed[ip]; exist {
		mu.Unlock()
		log.Printf("ALREADY PROCESSED: %v", ip)
		return
	} else {
		processed[ip] = ip
	}
	mu.Unlock()
	currentDepth++

	ni, err := parser.GetNetInfoFromInterx(ctx, ip)
	if err != nil {
		log.Printf("%v", err.Error())
		mu.Lock()
		blacklist[ip] = ip
		cleanValue(processed, ip)
		mu.Unlock()
		return
	}

	mu.Lock()
	cleanValue(processed, ip)
	pool[ip] = ip
	mu.Unlock()

	for _, p := range ni.Peers {
		wg.Add(1)
		go testLoop(ctx, wg, pool, blacklist, processed, p.RemoteIP, currentDepth, totalDepth, ignoreDepth)
	}

}

func cleanValue(toClean map[string]string, key string) {
	delete(toClean, key)
}
