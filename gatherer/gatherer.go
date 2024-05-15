package gatherer

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

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

	node, err := parser.GetNetInfoFromInterx(ctx, firstNode)
	if err != nil {
		panic(err)
	}

	for _, n := range node.Peers {
		go func() {

			n1, err := parser.GetNetInfoFromInterx(ctx, n.RemoteIP)
			if err != nil {
				return
			}
			

		}()
	}

}

func GetAllNodes(ctx context.Context, firstNode string) map[string]Node {

	var finalTree = make(map[string]Node)

	var poolOfNodes = make(map[string]Node)

	node, err := parser.GetNetInfoFromInterx(ctx, firstNode)
	if err != nil {
		panic(err)
	}
	errorPeers := 0
	goodPeers := 0
	var wg sync.WaitGroup

	poolChan := make(chan Node)

	// go finalAdder(&finalTree, poolChan)
	// wg.Add(1)
	go func() {
		for n := range poolChan {
			currentKey := fmt.Sprintf("%v@%v", n.IP, n.ID)
			if _, exist := finalTree[currentKey]; !exist {
				// mu.Lock()
				finalTree[currentKey] = n
				// mu.Unlock()
				// log.Println("range")
			}
		}
		defer log.Printf("DONE")
		defer wg.Done()
	}()
	for _, peer := range node.Peers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			peerCtx, cancel := context.WithTimeout(ctx, time.Second*1)
			defer cancel()

			_, err = parser.GetNetInfoFromInterx(peerCtx, peer.RemoteIP)
			mu.Lock()
			if err != nil {
				errorPeers++
			} else {
				goodPeers++
				fmt.Println(peer.RemoteIP, peer.NodeInfo.ID)
			}
			mu.Unlock()
			u, err := url.Parse(peer.NodeInfo.ListenAddr)
			if err != nil {
				log.Println("Error parsing URL:", err)
				return
			}
			host := u.Host
			parts := strings.Split(host, ":")
			if len(parts) != 2 {
				log.Println("Invalid host format")
				return
			}

			port, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Println("WrongPort: ", err.Error())
			}
			currentKey := fmt.Sprintf("%v@%v", peer.RemoteIP, peer.NodeInfo.ID)
			mu.Lock()
			if _, exist := poolOfNodes[currentKey]; !exist {
				poolOfNodes[currentKey] = Node{
					IP:       peer.RemoteIP,
					P2P_port: port,
					ID:       peer.NodeInfo.ID,
				}
			}

			mu.Unlock()
			poolChan <- poolOfNodes[currentKey]
		}()

	}

	wg.Wait()
	for _, n := range poolOfNodes {
		fmt.Println(n)
	}
	for _, n := range finalTree {
		fmt.Println(n)
	}
	fmt.Println(len(finalTree), len(poolOfNodes))
	return nil
}

// func finalAdder(final *map[string]Node, c <-chan Node) {
// 	for n := range c {
// 		mu.Lock()
// 		// &final[fmt.Sprintf("%v@%v", n.IP, n.ID)] <- n

// 		mu.Unlock()
// 	}

// }
