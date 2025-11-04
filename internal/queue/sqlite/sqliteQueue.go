package sqliteQueue

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"maragu.dev/goqite"
)

type sQueue struct {
	db *sql.DB
}

func NewQueue(sqliteFile string) sQueue {
	q := sQueue{}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal=WAL&_timeout=5000&_fk=true", sqliteFile))
	if err != nil {
		panic(err)
	}

	// check if schema exists
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?;")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var tblName string
	err = stmt.QueryRowContext(ctx, "goqite").Scan(&tblName)
	if err == sql.ErrNoRows {
		err = goqite.Setup(context.Background(), db)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	q.db = db

	return q
}

func (q *sQueue) Close() {
	q.db.Close()
}

func (s *sQueue) Insert(wg *sync.WaitGroup, queues []string, c <-chan string) {
	var q []*goqite.Queue
	for _, qName := range queues {
		q = append(q, goqite.New(goqite.NewOpts{
			DB:   s.db,
			Name: qName,
		}))
	}
	for item := range c {
		for _, queue := range q {
			err := queue.Send(context.Background(), goqite.Message{Body: []byte(item)})
			if err != nil {
				panic(err)
			}
		}
	}
	wg.Done()
}

func (s *sQueue) Stream(queue string) {
	q := goqite.New(goqite.NewOpts{
		DB:   s.db,
		Name: queue,
	})

	for {
		res, err := q.ReceiveAndWait(context.Background(), time.Millisecond*100)
		if err != nil {
			log.Println("receive error: ", err)
			continue
		}
		fmt.Println(string(res.Body))
		if err := q.Delete(context.Background(), res.ID); err != nil {
			log.Println("error deleting message from queue", "error", err)
			return
		}
	}

}
