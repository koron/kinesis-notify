package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var (
	region = flag.String("r", "ap-northeast-1", "AWS region")
	stream = flag.String("s", "", "stream name")
	pk     = flag.String("pk", "default", "partition key")
)

func main() {
	flag.Parse()
	if *stream == "" {
		log.Fatal("require stream name with '-s'")
	}
	ss := session.New(
		&aws.Config{
			Region:      aws.String(*region),
			Credentials: credentials.NewSharedCredentials("", ""),
		})
	svc := kinesis.New(ss)

	resp, err := svc.PutRecord(&kinesis.PutRecordInput{
		StreamName:   aws.String(*stream),
		PartitionKey: aws.String(*pk),
		Data:         []byte("foobar"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
