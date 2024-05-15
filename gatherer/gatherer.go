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
		testLoop(ctx, &wg, nodesPool, blacklist, n.RemoteIP)

	}

	wg.Wait()
	fmt.Println()
	log.Printf("\nTotal saved peers:%v\nOriginal node peer count: %v\nBlacklisted nodes(not reachable):n %v\n", len(nodesPool), len(node.Peers), len(blacklist))

	mu.Lock()
	for _, node := range nodesPool {
		log.Println(node)
	}
	mu.Unlock()

}

func testLoop(ctx context.Context, wg *sync.WaitGroup, pool, blacklist map[string]string, ip string) {
	wg.Add(1)
	defer wg.Done()
	log.Println("running testloop: ", ip)
	if _, exist := blacklist[ip]; exist {
		log.Printf("BLACKLISTED: %v", ip)
		return
	}
	if _, exist := pool[ip]; exist {
		log.Printf("ALREADY EXIST: %v", ip)
		return
	} else {
		ni, err := parser.GetNetInfoFromInterx(ctx, ip)
		if err != nil {
			log.Printf("%v", err.Error())
			blacklist[ip] = ip
			return
		}
		mu.Lock()
		log.Printf("adding %v", ip)
		pool[ip] = ip
		mu.Unlock()
		for _, p := range ni.Peers {
			wg.Add(1)
			go func() {
				testLoop(ctx, wg, pool, blacklist, p.RemoteIP)
				defer wg.Done()
			}()
		}
	}
}
