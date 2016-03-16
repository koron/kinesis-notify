# kinesis-notify (WIP)

Execute a command for each records from AWS Kinesis.
Working with MultiLangDaemon in [amazon-kinesis-client](https://github.com/awslabs/amazon-kinesis-client).

## How to use

    $ cd $GOPATH
    $ go get -u github.com/koron/kinesis-notify
    $ cd github.com/koron/kinesis-notify
    $ cp sample.properties src/main/resources
    # edit src/main/resources/sample.properties
    $ go install
    $ gradle run -P=sample.properties

## LICENSE

MIT license.  See LICENSE.
