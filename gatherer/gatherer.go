package gatherer

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

func GetAllNodesV2(ctx context.Context, firstNode string) {
	nodesPool := make(map[string]string)
	blacklist := make(map[string]string)

	node, err := parser.GetNetInfoFromInterx(ctx, firstNode)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, n := range node.Peers {
		log.Println(n.RemoteIP)
		wg.Add(1)
		go testLoop(ctx, &wg, nodesPool, blacklist, n.RemoteIP)
	}

	wg.Wait()
	fmt.Println()
	log.Printf("\nTotal saved peers:%v\nOriginal node peer count: %v\nBlacklisted nodes(not reachable): %v\n", len(nodesPool), len(node.Peers), len(blacklist))

	mu.Lock()
	for _, node := range nodesPool {
		log.Println(node)
	}
	mu.Unlock()

	fmt.Println("Done")

}

func testLoop(ctx context.Context, wg *sync.WaitGroup, pool, blacklist map[string]string, ip string) {
	defer wg.Done()
	log.Println("running testloop: ", ip)

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
	mu.Unlock()

	ni, err := parser.GetNetInfoFromInterx(ctx, ip)
	if err != nil {
		log.Printf("%v", err.Error())
		mu.Lock()
		blacklist[ip] = ip
		mu.Unlock()
		return
	}

	mu.Lock()
	pool[ip] = ip
	mu.Unlock()

	for _, p := range ni.Peers {
		wg.Add(1)
		go testLoop(ctx, wg, pool, blacklist, p.RemoteIP)
	}
}
