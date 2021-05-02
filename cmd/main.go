package main

import (
	"context"
	"github.com/akhrorov/wallet/pkg/search"
	"log"
)

func main() {
	root := context.Background()
	ctx, cancel := context.WithCancel(root)
	files := []string{
		"1. aaa, bbb, ccc\n2. ddd, eee, fff\n3. ggg, hhh, iii\n",
		"4. jjj, kkk, lll\n5. xxx, nnn, ooo\n6. ppp, qqq, rrr\n",
		"7. sss, ttt, uuu\n8. vvv, www, xxx\n9. yyy, bbb, zzz\n",
	}
	result := <- search.All(ctx, "xxx", files)
	cancel()
	log.Print(result)


	//root := context.Background()
	//ctx, cancel := context.WithCancel(root)
	//ch := make(chan int)
	//
	//for i := 0; i < 5; i++ {
	//	go func(ctx context.Context, index int, ch chan<- int) {
	//		wait := rand.Intn(10)
	//		log.Printf("%d wait for %d", index, wait)
	//		select {
	//			case <-ctx.Done():
	//				log.Printf("%d canceled", index)
	//		case <-time.After(time.Second * time.Duration(wait)):
	//			ch <- index
	//			log.Printf("%d done with %d", index, wait)
	//		}
	//	}(ctx, i, ch)
	//}
	//
	//winner := <-ch
	//cancel()
	//log.Print(winner)
	//<-time.After(time.Second)

	//for i := 0; i <5; i++ {
	//	go func(ctx context.Context, index int) {
	//		<-ctx.Done()
	//		log.Printf("done %d", index)
	//		return
	//	}(ctx, i)
	//}
	//
	//<-ctx.Done()
	//<-time.After(time.Second)
	//log.Print("done")

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func(ctx context.Context) {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			log.Print("done")
	//			wg.Done()
	//			return
	//		case <-time.After(time.Second):
	//			log.Print("tick")
	//		}
	//	}
	//}(ctx)
	//wg.Wait()
	//log.Print("main done")

	//data := make([]int,1_000_000)
	//for i := range data{
	//	data[i] = i
	//}
	//
	//parts := 10
	//size := len(data) / parts
	//channels := make([]<-chan int, parts)
	//
	//for i := 0; i < parts; i++ {
	//	ch := make(chan int)
	//	channels[i] = ch
	//	go func(ch chan <- int, data []int) {
	//		defer close(ch)
	//		sum := 0
	//		for _, v := range data{
	//			sum += v
	//		}
	//		ch <- sum
	//	}(ch, data[i*size:(i+1)*size])
	//}
	//
	//total := 0
	//for value := range merge(channels){
	//	total += value
	//}
	//log.Print(total)
}

//func merge(channels []<- chan int) <- chan int {
//	wg := sync.WaitGroup{}
//	wg.Add(len(channels))
//	merged := make(chan int)
//
//
//	for _, ch := range channels {
//		go func(ch <- chan int) {
//			defer wg.Done()
//			for val := range ch{
//				merged <- val
//			}
//		}(ch)
//	}
//
//	go func() {
//		defer close(merged)
//		wg.Wait()
//	}()
//
//	return merged
//}

//func main() {
//	data := make([]int, 1_000_000)
//	for i := range data {
//		data[i] = i
//	}
//
//	ch := make(chan int)
//	defer close(ch)
//	parts := 10
//	size := len(data) / parts
//
//	for i := 0; i < parts; i++ {
//		go func(ch chan <- int, data []int) {
//			sum := 0
//			for _, v := range data {
//				sum += v
//			}
//			ch <- sum
//		}(ch, data[i*size:(i+1)*size])
//	}
//
//	total := 0
//	for i := 0; i < parts; i++ {
//		total += <-ch
//	}
//
//	log.Print(total)
//}

//func main()  {
//	ch := make(chan struct{})
//	go func() {
//		<-time.After(time.Second)
//		close(ch)
//	}()
//
//	val, ok := <-ch
//	if !ok {
//		log.Print("channel closed")
//		return
//	}
//	log.Print(val)
//}

//func main()  {
//	ch := tick()
//
//	for channel := range ch {
//		log.Print(channel)
//	}
//}
//
//func tick() <- chan int {
//	ch := make(chan int)
//
//	go func() {
//		for i:= 0; i < 10; i++ {
//			ch <- i
//		}
//		close(ch)
//	}()
//	return ch
//}

//func main() {
//	done := make(chan struct{})
//
//	go tick(done)
//
//	<-time.After(time.Second * 5)
//	done <- struct{}{}
//}
//
//func tick(done <- chan struct{})  {
//	for  {
//		select {
//		case <-done:
//			return
//		case <-time.After(time.Second):
//			log.Print("tick")
//		}
//	}
//}
