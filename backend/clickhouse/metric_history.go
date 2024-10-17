package clickhouse

import (
	"context"
	"fmt"
	"time"

	modelInputs "github.com/highlight-run/highlight/backend/private-graph/graph/model"

	"github.com/huandu/go-sqlbuilder"
)

const MetricHistoryTable = "metric_history"
const MetricHistoryMaxLookback = time.Hour

type SavedMetricState struct {
	MetricId        string
	BlockNumberInfo []BlockNumberInfo
}

type BlockNumberInfo struct {
	Partition       string
	LastBlockNumber uint64
}

func (client *Client) GetBlockNumbers(ctx context.Context, metricId string, startDate time.Time, endDate time.Time) ([]BlockNumberInfo, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("toString(toDate(date_trunc('day', Timestamp))) as Partition",
		"maxMerge(MaxBlockNumberState) as LastBlockNumber")
	sb.From(MetricHistoryTable)
	sb.Where(sb.Equal("MetricId", metricId))
	sb.Where(sb.GreaterEqualThan("Timestamp", startDate))
	sb.Where(sb.LessThan("Timestamp", endDate))
	sb.GroupBy("1")

	sql, args := sb.BuildWithFlavor(sqlbuilder.ClickHouse)
	rows, err := client.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var results []BlockNumberInfo
	for rows.Next() {
		var result BlockNumberInfo
		if err := rows.ScanStruct(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

type AggregatedMetricResults struct {
	GroupByKey string
	BucketId   uint64
	Value      float64
}

func (client *Client) AggregateMetricStates(ctx context.Context, metricId string, endDate time.Time, interval time.Duration, aggregator modelInputs.MetricAggregator, windowSeconds *int) ([]*modelInputs.MetricBucket, error) {
	startDate := endDate.Add(-interval)
	sb := sqlbuilder.NewSelectBuilder()

	selectCols := []string{"GroupByKey"}
	groupCols := []string{"1"}

	if windowSeconds != nil {
		selectCols = append(selectCols, fmt.Sprintf("date_diff('second', %s, Timestamp) / %d as BucketId", startDate.Format("2006-01-02 15:04:05"), windowSeconds))
		groupCols = append(groupCols, "2")
	}

	switch aggregator {
	case modelInputs.MetricAggregatorCount:
		selectCols = append(selectCols, "toFloat64(countMerge(CountState)) as Value")
	case modelInputs.MetricAggregatorCountDistinct, modelInputs.MetricAggregatorCountDistinctKey:
		selectCols = append(selectCols, "toFloat64(uniqMerge(UniqState)) as Value")
	case modelInputs.MetricAggregatorMin:
		selectCols = append(selectCols, "minMerge(MinState) as Value")
	case modelInputs.MetricAggregatorAvg:
		selectCols = append(selectCols, "avgMerge(AvgState) as Value")
	case modelInputs.MetricAggregatorMax:
		selectCols = append(selectCols, "maxMerge(MaxState) as Value")
	case modelInputs.MetricAggregatorSum:
		selectCols = append(selectCols, "sumMerge(SumState) as Value")
	case modelInputs.MetricAggregatorP50:
		selectCols = append(selectCols, "quantileMerge(.5)(P50State) as Value")
	case modelInputs.MetricAggregatorP90:
		selectCols = append(selectCols, "quantileMerge(.9)(P90State) as Value")
	case modelInputs.MetricAggregatorP95:
		selectCols = append(selectCols, "quantileMerge(.95)(P95State) as Value")
	case modelInputs.MetricAggregatorP99:
		selectCols = append(selectCols, "quantileMerge(.99)(P99State) as Value")
	}
	sb.Select(selectCols...)
	sb.From(MetricHistoryTable)
	sb.Where(sb.Equal("MetricId", metricId))
	sb.Where(sb.GreaterEqualThan("Timestamp", startDate))
	sb.Where(sb.LessThan("Timestamp", endDate))
	sb.GroupBy(groupCols...)
	sb.OrderBy(groupCols...)

	sql, args := sb.BuildWithFlavor(sqlbuilder.ClickHouse)
	rows, err := client.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := []*modelInputs.MetricBucket{}
	var result AggregatedMetricResults
	for rows.Next() {
		if err := rows.ScanStruct(&result); err != nil {
			return nil, err
		}
		results = append(results, &modelInputs.MetricBucket{
			BucketID:    result.BucketId,
			Group:       []string{result.GroupByKey},
			MetricValue: &result.Value,
			YhatLower:   nil,
			YhatUpper:   nil,
		})
	}

	return results, nil
}
