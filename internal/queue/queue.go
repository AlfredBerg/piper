package queue

import (
	"sync"

	redisqueue "github.com/AlfredBerg/piper/internal/queue/redis"
	sqliteQueue "github.com/AlfredBerg/piper/internal/queue/sqlite"
	"github.com/spf13/viper"
)

type Queue interface {
	Stream(queue string)
	Insert(wg *sync.WaitGroup, queues []string, c <-chan string)
	Close()
}

func NewQueue() Queue {
	redisUrl := viper.GetString("redis_Url")
	queue := Queue(nil)
	if redisUrl == "" {
		sf := viper.GetString("sqlite_file")
		q := sqliteQueue.NewQueue(sf)
		queue = &q
	} else {
		q := redisqueue.NewQueue(redisUrl)
		queue = &q
	}
	return queue
}
