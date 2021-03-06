// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/prometheus/client_golang/prometheus"
)

var (
	rangeCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "range_total",
			Help:      "Total number of ranges seen by this member.",
		})

	putCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "put_total",
			Help:      "Total number of puts seen by this member.",
		})

	deleteCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "delete_total",
			Help:      "Total number of deletes seen by this member.",
		})

	txnCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "txn_total",
			Help:      "Total number of txns seen by this member.",
		})

	keysGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "keys_total",
			Help:      "Total number of keys.",
		})

	indexCompactionPauseDurations = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "index_compaction_pause_duration_milliseconds",
			Help:      "Bucketed histogram of index compaction puase duration.",
			// 0.5ms -> 1second
			Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
		})

	dbCompactionPauseDurations = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "db_compaction_pause_duration_milliseconds",
			Help:      "Bucketed histogram of db compaction puase duration.",
			// 1ms -> 4second
			Buckets: prometheus.ExponentialBuckets(1, 2, 13),
		})

	dbCompactionTotalDurations = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "storage",
			Name:      "db_compaction_total_duration_milliseconds",
			Help:      "Bucketed histogram of db compaction total duration.",
			// 100ms -> 800second
			Buckets: prometheus.ExponentialBuckets(100, 2, 14),
		})
)

func init() {
	prometheus.MustRegister(rangeCounter)
	prometheus.MustRegister(putCounter)
	prometheus.MustRegister(deleteCounter)
	prometheus.MustRegister(txnCounter)
	prometheus.MustRegister(keysGauge)
	prometheus.MustRegister(indexCompactionPauseDurations)
	prometheus.MustRegister(dbCompactionPauseDurations)
	prometheus.MustRegister(dbCompactionTotalDurations)
}
