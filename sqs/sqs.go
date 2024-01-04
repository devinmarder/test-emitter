package sqs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog"
)

const (
	batchSize = 10
)

type Client struct {
	*sqs.Client
	log zerolog.Logger
}

func New(cfg aws.Config, log zerolog.Logger) *Client {
	return &Client{
		Client: sqs.NewFromConfig(cfg),
		log:    log.With().Str("module", "sqs").Logger(),
	}
}

func (c *Client) NewPublisher(ctx context.Context, queue string, msgs <-chan string) error {
	batch := make([]*string, 0, batchSize)
	for msg := range msgs {
		msg := msg
		batch = append(batch, &msg)
		if len(batch) == batchSize {
			if err := c.publishBatch(ctx, queue, batch); err != nil {
				return fmt.Errorf("failed to publish batch: %w", err)
			}
			batch = make([]*string, 0, batchSize)
		}
	}
	if len(batch) > 0 {
		if err := c.publishBatch(ctx, queue, batch); err != nil {
			return fmt.Errorf("failed to publish batch: %w", err)
		}
	}

	return nil
}

func (c *Client) publishBatch(ctx context.Context, queue string, batch []*string) error {
	log := c.log.With().Str("queue", queue).Logger()

	log.Debug().Int("batch_size", len(batch)).Msg("publishing batch")
	result, err := c.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queue),
		Entries:  toEntries(batch),
	})
	if err != nil {
		return err
	}

	for _, f := range result.Failed {
		log.Error().Str("id", *f.Id).Str("code", *f.Code).Str("message", *f.Message).Msg("failed to publish message")
	}
	log.Debug().Msg("published batch")
	return nil
}

func toEntries(batch []*string) []types.SendMessageBatchRequestEntry {
	entries := make([]types.SendMessageBatchRequestEntry, len(batch))
	for i, m := range batch {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(i)),
			MessageBody: m,
		}
	}
	return entries
}
