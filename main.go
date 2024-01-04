package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/devinmarder/test-emitter/sqs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

var (
	count    = flag.Int("count", 1, "number of times to print message")
	out      = flag.String("out", "stdout", "output destination")
	queue    = flag.String("queue", "", "sqs queue name")
	logLevel = flag.String("log-level", "info", "log level")
	msg      = flag.String("msg", "message {{.ID}}", "message to print")
	file     = flag.String("file", "", "file to read message from")
)

func main() {
	flag.Parse()

	logLevel, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	log := zerolog.New(zerolog.NewConsoleWriter()).Level(logLevel)

	var tmpl *template.Template

	switch {
	case *file != "":
		tmpl, err = loadTemplate(*file)
		if err != nil {
			panic(err)
		}
	default:
		tmpl, err = template.New("msg").Parse(*msg)
		if err != nil {
			panic(err)
		}
	}

	var (
		msgs   chan string
		g, ctx = errgroup.WithContext(context.Background())
	)

	switch *out {
	case "stdout":
		g.Go(func() error {
			return newPublisher(os.Stdout, msgs)
		})
	case "sqs":
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			panic(err)
		}
		if *queue == "" {
			panic("queue name required")
		}
		g.Go(func() error {
			return sqs.New(cfg, log).NewPublisher(context.Background(), *queue, msgs)
		})
	default:
		panic("invalid output destination")
	}

	fmt.Printf("time to print %d messages: %s\n", *count, timed(func() {
		msgs = make(chan string)
		g.Go(func() error {
			return printMsg(ctx, msgs, *count, tmpl)
		})

		if err := g.Wait(); err != nil {
			log.Error().Err(err).Msg("failed to publish messages")
		}
	}))
}

func timed(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

type Params struct {
	ID        int
	Timestamp string
}

func printMsg(ctx context.Context, msgs chan<- string, count int, tmpl *template.Template) error {
	defer close(msgs)
	var b bytes.Buffer
	for i := 0; i < count; i++ {
		msg := Params{
			ID:        i,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		if err := tmpl.Execute(&b, msg); err != nil {
			panic(err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case msgs <- b.String():
		}
		b.Reset()
	}
	return nil
}

func loadTemplate(file string) (*template.Template, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return template.New("msg").Parse(string(b))
}

func newPublisher(w io.Writer, msgs chan string) error {
	for msg := range msgs {
		if _, err := fmt.Fprintln(w, msg); err != nil {
			return err
		}
	}
	return nil
}
