# README

## Project Name

This is a simple tool designed to emit events. It currently supports stdout and AWS SQS as output destinations.

## Prerequisites

- Go 1.16 or higher
- AWS CLI (if using AWS SQS)

## Installation

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/yourrepository.git
cd yourrepository
```

## Usage

### Emitting to stdout

To emit events to stdout, use the `-out` flag set to `stdout`:

```bash
go run main.go -out stdout
```

### Emitting to AWS SQS

To emit events to AWS SQS, you need to set the `-out` flag to `sqs` and provide a queue name with the `-queue` flag:

```bash
go run main.go -out sqs -queue your_queue_name
```

Please ensure that you have configured your AWS CLI with the appropriate credentials and region.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)