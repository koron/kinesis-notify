package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/nieksand/gokinesis/src/kinesis"
)

type consumer struct {
	cmd     string
	args    []string
	workers *workers
	shardID string
}

var logger = log.New(os.Stderr, "", log.LstdFlags)

func (c *consumer) Init(shardID string) error {
	c.shardID = shardID
	return nil
}

func (c *consumer) ProcessRecords(records []*kinesis.KclRecord, cp *kinesis.Checkpointer) error {
	for _, r := range records {
		cmd, err := c.toCmd(r)
		if err != nil {
			logger.Printf("failed to make command: %s", err)
			continue
		}
		c.workers.Run(workerJob{Cmd: cmd})
	}
	c.workers.Wait()
	cp.CheckpointAll()
	return nil
}

func (c *consumer) Shutdown(shutdownType kinesis.ShutdownType, cp *kinesis.Checkpointer) error {
	c.workers.Wait()
	return nil
}

func (c *consumer) toCmd(r *kinesis.KclRecord) (*exec.Cmd, error) {
	b, err := base64.StdEncoding.DecodeString(r.DataB64)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(c.cmd, c.args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		stdin.Write(b)
		stdin.Close()
	}()
	return cmd, nil
}

func main() {
	numWorkers := flag.Int("worker", runtime.NumCPU(), "num of workers")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, `USAGE: %s [OPTIONS] {command args...}

OPTIONS
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	c := &consumer{
		cmd:     args[0],
		args:    args[1:],
		workers: newWorkers(*numWorkers),
	}
	kinesis.Run(c)
}
