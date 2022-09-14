package changestreams

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func GetResponse(ctx context.Context, sc *spanner.Client, query Query) (chan StreamResponse, chan error) {
	fmt.Println("get response, partition:", query.PartitionToken)
	response := make(chan StreamResponse)
	errors := make(chan error)
	go func() {
		stmt := prepareStatement(query)
		tx := sc.Single()
		defer tx.Close()
		iter := tx.Query(ctx, stmt)
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				errors <- err
			}
			val := StreamResponse{}
			err = row.ToStruct(&val)
			if err != nil {
				errors <- err //todo wrap for context
			}
			response <- val
		}
	}()
	return response, errors
}

func prepareStatement(query Query) spanner.Statement {
	sql := fmt.Sprintf(`SELECT ChangeRecord FROM READ_%s(
		start_timestamp => @startTime,
		end_timestamp => @endTime,
		partition_token => @partitionToken,
		heartbeat_milliseconds => @heartbeatMS
	  );`, query.Name)
	stmt := spanner.NewStatement(sql)
	stmt.Params = map[string]interface{}{
		"startTime":      query.StartTime,
		"heartbeatMS":    query.HeartbeatMS,
		"endTime":        nil,
		"partitionToken": nil,
	}
	if !query.EndTime.IsZero() {
		stmt.Params["endTime"] = query.EndTime
	}
	if query.PartitionToken != "" {
		stmt.Params["partitionToken"] = query.PartitionToken
	}

	return stmt
}
