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
	cmd      string
	args     []string

	workers  *workers
	retryMax int
	cpFirst  bool
	logname  string

	shardID  string
	logFile  *os.File
}

var logger = log.New(os.Stderr, "", log.LstdFlags)

func (c *consumer) Init(shardID string) error {
	c.shardID = shardID
	// Replace logger with file
	if c.logname != "" {
		name := fmt.Sprintf("%s-%s.log", c.logname, c.shardID)
		f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		c.logFile = f
		logger = log.New(c.logFile, "", log.LstdFlags)
	}
	return nil
}

func (c *consumer) ProcessRecords(records []*kinesis.KclRecord, cp *kinesis.Checkpointer) error {
	if c.cpFirst {
		cp.CheckpointAll()
	} else {
		defer cp.CheckpointAll()
	}
	for _, r := range records {
		c.run(r, 0)
	}
	c.workers.Wait()
	return nil
}

func (c *consumer) Shutdown(shutdownType kinesis.ShutdownType, cp *kinesis.Checkpointer) error {
	c.workers.Wait()
	c.logFile.Close()
	return nil
}

func (c *consumer) run(r *kinesis.KclRecord, failedCount int) {
	cmd, err := c.toCmd(r)
	if err != nil {
		logger.Printf("failed to make command: %s", err)
		return
	}
	c.workers.Run(workerJob{
		Cmd:    cmd,
		Finish: c.newFin(r, failedCount),
	})
}

func (c *consumer) newFin(rec *kinesis.KclRecord, failedCount int) func(workerResult) {
	return func(res workerResult) {
		if res.Success() {
			return
		}
		if failedCount += 1; failedCount > c.retryMax {
			logger.Printf("gave up retry, last error: %s", res.Error)
			return
		}
		c.run(rec, failedCount)
	}
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
	numRetry := flag.Int("retry", 0, "retry count")
	cpFirst := flag.Bool("checkpointfirst", false, "update check point at first of ProcessRecords")
	logname := flag.String("logname", "", "core name for log files")
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
		cmd:      args[0],
		args:     args[1:],
		workers:  newWorkers(*numWorkers),
		retryMax: *numRetry,
		cpFirst:  *cpFirst,
		logname:  *logname,
	}
	kinesis.Run(c)
}
