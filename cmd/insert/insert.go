package insertcmd

import (
	"bufio"
	"os"
	"sync"

	"github.com/AlfredBerg/piper/internal/queue"
	"github.com/spf13/cobra"
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
	queue := queue.NewQueue()

	defer queue.Close()

	wg := sync.WaitGroup{}
	insertC := make(chan string)
	wg.Add(1)
	go queue.Insert(&wg, f.queues, insertC)

	var sc *bufio.Scanner
	if f.input == "" {
		sc = bufio.NewScanner(os.Stdin)
	} else {
		fh, err := os.Open(f.input)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		sc = bufio.NewScanner(fh)
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
