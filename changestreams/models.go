package changestreams

import (
	"time"

	"cloud.google.com/go/spanner"
)

type StreamResponse struct {
	ChangeRecord []*ChangeRecord `spanner:"ChangeRecord"`
}

type ChangeRecord struct {
	DataChanges    []*DataChangeResponse    `spanner:"data_change_record"`
	Heartbeats     []*HeatbeatResponse      `spanner:"heartbeat_record"`
	ChildPartiions []*ChildPartionsResponse `spanner:"child_partitions_record"`
}

type HeatbeatResponse struct {
	Timestamp time.Time
}

type DataChangeResponse struct {
	CommitTimestamp                      time.Time     `spanner:"commit_timestamp"`
	RecordSequence                       string        `spanner:"record_sequence"`
	ServerTransactionID                  string        `spanner:"server_transaction_id"`
	IsLastRecordInTransactionInPartition bool          `spanner:"is_last_record_in_transaction_in_partition"`
	TableName                            string        `spanner:"table_name"`
	ValueCaptureType                     string        `spanner:"value_capture_type"`
	ColumnTypes                          []*ColumnType `spanner:"column_types"`
	Mods                                 []*Mod        `spanner:"mods"`
	ModType                              string        `spanner:"mod_type"`
	NumberOfRecordsInTransaction         int64         `spanner:"number_of_records_in_transaction"`
	NumberOfPartitionsInTransaction      int64         `spanner:"number_of_partitions_in_transaction"`
	TransactionTag                       string        `spanner:"transaction_tag"`
	IsSystemTransaction                  bool          `spanner:"is_system_transaction"`
}

type ColumnType struct {
	Name            string           `spanner:"name"`
	Type            spanner.NullJSON `spanner:"type"`
	IsPrimaryKey    bool             `spanner:"is_primary_key"`
	OrdinalPosition int64            `spanner:"ordinal_position"`
}

type Mod struct {
	Keys      spanner.NullJSON `spanner:"keys"`
	NewValues spanner.NullJSON `spanner:"new_values"`
	OldValues spanner.NullJSON `spanner:"old_values"`
}

type ChildPartionsResponse struct {
	StartTimeStamp time.Time         `spanner:"start_timestamp"`
	RecordSequence string            `spanner:"record_sequence"`
	ChildPartions  []*ChildPartition `spanner:"child_partitions"`
}

type ChildPartition struct {
	Token               string   `spanner:"token"`
	ParentPartionTokens []string `spanner:"parent_partition_tokens"`
}

type Query struct {
	Name           string
	StartTime      time.Time
	EndTime        time.Time
	PartitionToken string
	HeartbeatMS    int64
}

// ModType possible values INSERT/UPDATE/DELETE
type ModType string

const (
	// ModTypeInsert is raised when an insert occurs to subscribed tables
	ModTypeInsert ModType = "INSERT"
	// ModTypeUpdate is raised when an updated occurs on the subscribed table
	ModTypeUpdate ModType = "UPDATE"
	// ModTypeDelete is raised when a delete occurs on the subscribed table
	ModTypeDelete ModType = "DELETE"
)
