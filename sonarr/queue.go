package sonarr

import (
	"context"
	"fmt"
	"golift.io/starr"
	"net/url"
)

// GetQueue returns a single page from the Sonarr Queue (processing, but not yet imported).
// WARNING: 12/30/2021 - this method changed.
// If you need control over the page, use sonarr.GetQueuePage().
// This function simply returns the number of queue records desired,
// up to the number of records present in the application.
// It grabs records in (paginated) batches of perPage, and concatenates
// them into one list.  Passing zero for records will return all of them.
func (s *Sonarr) GetQueue(records, perPage int) (*Queue, error) {
	return s.GetQueueContext(context.Background(), records, perPage)
}

func (s *Sonarr) GetQueueContext(ctx context.Context, records, perPage int) (*Queue, error) {
	queue := &Queue{Records: []*QueueRecord{}}
	perPage = starr.SetPerPage(records, perPage)

	for page := 1; ; page++ {
		curr, err := s.GetQueuePageContext(ctx, &starr.Req{PageSize: perPage, Page: page})
		if err != nil {
			return nil, err
		}

		queue.Records = append(queue.Records, curr.Records...)

		if len(queue.Records) >= curr.TotalRecords ||
			(len(queue.Records) >= records && records != 0) ||
			len(curr.Records) == 0 {
			queue.PageSize = curr.TotalRecords
			queue.TotalRecords = curr.TotalRecords
			queue.SortDirection = curr.SortDirection
			queue.SortKey = curr.SortKey

			break
		}

		perPage = starr.AdjustPerPage(records, curr.TotalRecords, len(queue.Records), perPage)
	}

	return queue, nil
}

// GetQueuePage returns a single page from the Sonarr Queue.
// The page size and number is configurable with the input request parameters.
func (s *Sonarr) GetQueuePage(params *starr.Req) (*Queue, error) {
	return s.GetQueuePageContext(context.Background(), params)
}

func (s *Sonarr) GetQueuePageContext(ctx context.Context, params *starr.Req) (*Queue, error) {
	var queue Queue

	params.CheckSet("sortKey", "timeleft")
	params.CheckSet("includeUnknownSeriesItems", "true")

	_, err := s.GetInto(ctx, "v3/queue", params.Params(), &queue)
	if err != nil {
		return nil, fmt.Errorf("api.Get(queue): %w", err)
	}

	return &queue, nil
}

type DeleteQueueRecordParam struct {
	Blacklist bool
}

// https://github.com/Sonarr/Sonarr/wiki/Queue
func (s *Sonarr) DeleteQueueRecord(ctx context.Context, record *QueueRecord, p *DeleteQueueRecordParam) error {
	params := make(url.Values)
	params.Set("blacklist", fmt.Sprintf("%t", p.Blacklist))

	_, err := s.Delete(ctx, fmt.Sprintf("v3/queue/%d", record.ID), params)
	if err != nil {
		return fmt.Errorf("api.Delete(queue): %w", err)
	}
	return nil
}
