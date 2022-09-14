package changestreams

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
)

type monitor struct {
	Name     string // Change Stream Name
	Response chan DataChangeResponse
	Errors   chan error
}

func Subscribe(ctx context.Context, sc *spanner.Client, name string) (chan DataChangeResponse, chan error) {
	monitor := monitor{
		Name:     name,
		Response: make(chan DataChangeResponse),
		Errors:   make(chan error),
	}
	partitionsSeen := map[string]struct{}{}
	go monitor.getPartition(ctx, sc, "", time.Now().Add(time.Hour*-5), partitionsSeen)
	return monitor.Response, monitor.Errors
}

func (m *monitor) getPartition(ctx context.Context, sc *spanner.Client, partition string, startTime time.Time, seen map[string]struct{}) {
	query := Query{
		Name:           m.Name,
		HeartbeatMS:    1000,
		StartTime:      startTime,
		PartitionToken: partition,
	}
	response, errors := GetResponse(ctx, sc, query)
	go func() {
		for {
			err := <-errors
			m.Errors <- err
		}
	}()
	for sr := range response {
		for _, r := range sr.ChangeRecord {
			//fmt.Println("data changes", len(r.DataChanges))
			for _, dc := range r.DataChanges {
				m.Response <- *dc
			}
			//wg := sync.WaitGroup{}
			for _, cp := range r.ChildPartiions {
				for _, ct := range cp.ChildPartions {
					//wg.Add(1)
					token := ct.Token
					startTime := cp.StartTimeStamp.Add(time.Second)
					if _, ok := seen[token]; ok {
						continue
					}
					seen[token] = struct{}{}
					//go func() {
					go m.getPartition(ctx, sc, token, startTime, seen)
					//}()
				}
			}
			//wg.Wait()
		}
	}
	fmt.Println("end partition:", partition)
}
