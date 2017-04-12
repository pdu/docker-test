package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
	"strings"

	pb "github.com/pdu/docker-test/proto"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

var size = flag.Int("size", 4096, "message size")
var ips = flag.String("ips", "", "ip list")
var limit = flag.Int("limit", 10000, "The times of running")
var sleep = flag.Int("sleep", 1, "The sleep gap in ms")

func main() {
	flag.Parse()
    ipList := strings.Split(*ips, ",")
	var wg sync.WaitGroup
	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string, port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			for i := 0; i < *limit; i++ {
				select {
				case <-time.After(time.Millisecond * time.Duration(*sleep)):
					start := time.Now()
					r, err := c.SayHello(context.Background(), &pb.HelloRequest{Size: int64(*size)})
					elapsed := time.Since(start) / time.Millisecond
					if err != nil {
						log.Fatalf("could not greet: %v", err)
					}
					if elapsed > 10 {
						log.Printf("port:%d seq: %d reply len: %d, took %dms", port, i, len(r.Message), elapsed)
					}
				}
			}
		}(ip, 10000)
	}
	wg.Wait()
}
