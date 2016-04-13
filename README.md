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

## Required runtime files

*   `kinesis-notify` in PATH
*   `run-kinesis-notify` in PATH or current dir
*   `./lib/*.jar`
*   `./kinesis-notify.properties`

## Options

*   `-worker {NUM}` Number of workers (parallel jobs).  Default is equals with
    number of logical CPU core.
*   `-retry {NUM}` Max count of retry failed command.  Default is zero - don't
    retry.
*   `-checkpointfirst` Update check point at first of receiving records (each
    calls of ProcessingRecords).  Default false.
*   `-logname` Used for core name of log file.  Filename is determined by
    format `{logname}-{YYYYMMDD}-{shardID}.log` and log files are rotated
    daily.  And it can be included path separators.
    If not specified, output to STDERR instead of file.

## LICENSE

MIT license.  See LICENSE.
