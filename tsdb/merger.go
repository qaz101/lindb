package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/tblstore/metricsmeta"
	"github.com/lindb/lindb/tsdb/tblstore/metricsnameid"
)

const (
	invertedIndexMerger = "inverted_index_merger"
	metricNameIDsMerger = "metric_name_ids_merger"
	metricMetaMerger    = "metric_meta_merger"
	//defaultTTLDuration  = time.Hour * 24 * 30
	nopMerger = "nop_merger"
)

func init() {
	// FIXME stone1100
	kv.RegisterMerger(
		invertedIndexMerger,
		metricsnameid.NewMerger())

	kv.RegisterMerger(
		metricNameIDsMerger,
		metricsnameid.NewMerger())

	kv.RegisterMerger(
		metricMetaMerger,
		metricsmeta.NewMerger())

	kv.RegisterMerger(nopMerger, &_nopMerger{})
}

// nopMerger does nothing
type _nopMerger struct{}

func (m *_nopMerger) Merge(key uint32, value [][]byte) ([]byte, error) {
	return nil, nil
}
