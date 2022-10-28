package cmd

// import (
// 	"fmt"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/chux0519/yeti/pkg/service/rank"
// 	"github.com/chux0519/yeti/pkg/utils"
// )

// func main() {
// 	defaultIgns := "POTAFREE"

// 	igns := strings.SplitN(defaultIgns, ",", -1)
// 	count := len(igns)
// 	ch := make(chan *rank.RankData, count)

// 	wg := sync.WaitGroup{}
// 	for _, ign := range igns {
// 		wg.Add(1)
// 		name := ign
// 		go func() {
// 			fmt.Printf("Get rank of %s...\n", name)
// 			rank, err := rank.GetGMSRank(name)
// 			if err != nil {
// 				fmt.Printf("Failed to get rank of %s: %s\n", name, err.Error())
// 			}
// 			ch <- rank
// 		}()
// 	}
// 	utils.WaitTimeout(&wg, 30*time.Second)

// 	results := []*rank.RankData{}

// 	for {
// 		res := <-ch
// 		if res != nil {
// 			results = append(results, res)
// 		}
// 		count -= 1
// 		if count <= 0 {
// 			break
// 		}
// 	}

// 	for _, res := range results {
// 		process(res)
// 	}
// }

// func process(rank *rank.RankData) {
// 	fmt.Printf("res: %v\n", *rank)
// }
