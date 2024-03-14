package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type queue struct {
	client *redis.Client
}

func NewQueue(redisUrl string) queue {
	if redisUrl == "" {
		log.Fatal("A redis url must be set, e.g. with the env \"PIPER_REDIS_URL\" following the format redis://<user>:<password>@<host>:<port>/<db>")
	}

	q := queue{}
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	q.client = redis.NewClient(opt)

	return q
}

func (q *queue) Close() {
	q.client.Close()
}

// Stream reads items from a queue. It Blocks until there are items, and gets multiple items at the same time to be efficent
func (q *queue) Stream(queue string) {
	for {
		_, res, err := q.client.BLMPop(context.Background(), time.Minute, "right", 100, queue).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Println("receive error: ", err)
			continue
		}
		for _, r := range res {
			fmt.Println(r)
		}
	}
}

// Insert buffers the data and does an insert every second or when all data has been read
func (q *queue) Insert(wg *sync.WaitGroup, queues []string, c <-chan string) {
	var buffer []string
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if len(buffer) > 0 {
				redisInsert(q.client, buffer, queues)
				buffer = nil
			}
			timer.Reset(time.Second)
		case data, ok := <-c:
			//Is the channel closed?
			if !ok {
				redisInsert(q.client, buffer, queues)
				wg.Done()
				return
			}
			buffer = append(buffer, data)
		}
	}
}

func redisInsert(client *redis.Client, buffer []string, queues []string) {
	p := client.Pipeline()
	for _, queue := range queues {
		for _, data := range buffer {
			p.LPush(context.Background(), queue, data)
		}
	}
	_, err := p.Exec(context.Background())
	if err != nil {
		log.Println("insert failed: ", err)
	}
}
