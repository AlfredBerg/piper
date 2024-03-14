package streamcmd

import (
	"os"

	"github.com/AlfredBerg/piper/internal/pipebuffer"
	"github.com/AlfredBerg/piper/internal/queue"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type streamFlags struct {
	queue string
}

func NewCmdStream() *cobra.Command {
	f := streamFlags{}
	var streamCmd = &cobra.Command{
		Use:   "stream",
		Short: "Read items from a queue",

		Run: func(cmd *cobra.Command, args []string) {
			stream(f)
		},
	}

	streamCmd.Flags().StringVarP(&f.queue, "queue", "q", "", "Queue to read items from")
	return streamCmd
}

func stream(f streamFlags) {
	//Since linux 2.6.11 the pipe buffer is by default 16 pages.
	//We generally want a fairly small pipe buffer, so if e.g.
	//the application that consumes the output crashes the number
	//of items lost is limited. It also makes it easier to divide
	//the items between different workers
	//
	//The smallest possible pipe buffer is 1 page size (4096 bytes)
	fi, err := os.Stdout.Stat()
	if err != nil {
		panic(err)
	}

	//Is a pipe connected?
	if fi.Mode()&os.ModeNamedPipe > 0 {
		pipebuffer.Set(os.Stdout.Fd(), 4096)
	}

	redisUrl := viper.GetString("redis_Url")
	q := queue.NewQueue(redisUrl)
	defer q.Close()

	q.Stream(f.queue)
}
