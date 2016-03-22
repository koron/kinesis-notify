# kinesis-notify (WIP)

Execute a command for each records from AWS Kinesis.
Working with MultiLangDaemon in [amazon-kinesis-client](https://github.com/awslabs/amazon-kinesis-client).

## How to use

    $ go get -u github.com/koron/kinesis-notify
    $ cd $GOPATH/github.com/koron/kinesis-notify
    $ go install
    $ go install ./cmd/kinesis-echo
    $ cp sample.properties kinesis-notify.properties

    # edit kinesis-notify.properties
    $ vi kinesis-notify.properties

    $ ./run-kinesis-notify

## LICENSE

MIT license.  See LICENSE.
