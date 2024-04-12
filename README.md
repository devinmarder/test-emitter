# Test Emitter

This is a simple tool designed to emit events. It currently supports stdout and
AWS SQS as output destinations.

## Requires

- Go 1.20 or higher
- AWS CLI (if using AWS SQS)

## Installation

```bash
go install github.com/devinmarder/test-emitter@latest
```

## Usage

### Emitting to stdout

To emit events to stdout, use the `-out` flag set to `stdout`:

```bash
test-emitter -out stdout
```

### Emitting to AWS SQS

To emit events to AWS SQS, you need to set the `-out` flag to `sqs` and provide
a queue name with the `-queue` flag:

```bash
test-emitter -out sqs -queue your_queue_name
```

Please ensure that you have configured your AWS CLI with the appropriate
credentials and region.

### Options

The message to be sent can be specified with the `-msg` parameter or a file
location can be specified with `-file`. The Message supports go templates
where `ID` and `Timestamp` can be used to embed an incremental id and the
current time in the message.

The number of messages sent can be specified with `-count` and the logging
level can be specified using `-log-level`

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
