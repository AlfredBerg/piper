package insertcmd

import (
	"bufio"
	"os"
	"sync"

	"github.com/AlfredBerg/piper/internal/queue"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type insertFlags struct {
	queues []string
	input  string
}

func NewCmdInsert() *cobra.Command {
	f := insertFlags{}
	var diffCmd = &cobra.Command{
		Use:   "insert",
		Short: "Insert items to one or more queues",

		Run: func(cmd *cobra.Command, args []string) {
			insert(f)
		},
	}

	diffCmd.Flags().StringSliceVarP(&f.queues, "queue", "q", nil, "The queue to insert to, can be specified multiple times to insert to multiple queues")
	diffCmd.Flags().StringVarP(&f.input, "input", "i", "", "Input file, if empty stdin")
	return diffCmd
}

func insert(f insertFlags) {
	redisUrl := viper.GetString("redis_Url")

	q := queue.NewQueue(redisUrl)
	defer q.Close()

	wg := sync.WaitGroup{}
	insertC := make(chan string)
	wg.Add(1)
	go q.Insert(&wg, f.queues, insertC)

	var sc *bufio.Scanner
	if f.input == "" {
		sc = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(f.input)
		if err != nil {
			panic(err)
		}
		sc = bufio.NewScanner(f)
	}
	for sc.Scan() {
		insertC <- sc.Text()
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}

	close(insertC)
	wg.Wait()
}
