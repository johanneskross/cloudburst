package scoreboard

import ()

type ScoreboardStatistics struct {
	TargetId          int                         `json:"target_id"`
	RunDuration       int64                       `json:"run_duration"`
	StartTime         int64                       `json:"start_time"`
	EndTime           int64                       `json:"end_time"`
	TotalDropoffs     int64                       `json:"totap_dropoffs"`
	FinalScorecard    ScorecardStatistics         `json:"final_scorecard"`
	WaitTimeSummaries WaitTimeSummariesStatistics `json:"waits_stats"`
}

type ScorecardStatistics struct {
	RunDuration       int64                        `json:"run_duration"`
	IntervalDuration  int64                        `json:"interval_duration"`
	TotalOpsInitiated int64                        `json:"total_ops_initiated"`
	TotalOpsLate      int64                        `json:"total_ops_late"`
	OfferedLoadOps    float64                      `json:"offered_load_ops"`
	Summary           OperationSummaryStatistics   `json:"summary"`
	Operational       ScorecardOperationStatistics `json:"operational"`
}

type ScorecardOperationStatistics struct {
	Operations map[string]OperationSummaryStatistics `json:"operations"`
}

type OperationSummaryStatistics struct {
	OpsSuccessful     int64 `json:"ops_successful"`
	OpsFailed         int64 `json:"ops_failed"`
	OpsSeen           int64 `json:"ops_seen"`
	ActionsSuccessful int64 `json:"actions_successful"`
	Ops               int64 `json:"ops"`

	EffectiveLoadOps float64 `json:"effective_load_ops"`
	EffectiveLoadReq float64 `json:"effective_load_req"`

	RtimeTotal     int64   `json:"rtime_total"`
	RtimeThrFailed int64   `json:"rtime_thr_failed"`
	RtimeAverage   float64 `json:"rtime_average"`
	RtimeMax       int64   `json:"rtime_max"`
	RtimeMin       int64   `json:"rtime_min"`

	SamplerSamplesCollected int     `json:"sampler_samples_collected`
	SamplerSamplesSeen      int     `json:"sampler_samples_seen"`
	SamplerRtime50th        float64 `json:"sampler_rtime_50th"`
	SamplerRtime90th        float64 `json:"sampler_rtime_90th"`
	SamplerRtime95th        float64 `json:"sampler_rtime_95th"`
	SamplerRtime99th        float64 `json:"sampler_rtime_99th"`
	SamplerRtimeMean        float64 `json:"sampler_rtime_mean"`
	SamplerRtimeStdev       float64 `json:"sampler_rtime_stdev"`
	SamplerRtimeTvalue      float64 `json:"sampler_rtime_tvalue"`
}

type WaitTimeSummariesStatistics struct {
	WaitTimeSummaries []WaitTimeSummaryStatistics `json:"waits"`
}

type WaitTimeSummaryStatistics struct {
	OperationName           string  `json:"operation_name"`
	AverageWaitTime         float64 `json:"average_wait_time"`
	TotalWaitTime           int64   `json:"total_wait_time"`
	MinWaitTime             int64   `json:"min_wait_time"`
	MaxWaitTime             int64   `json:"max_wait_time"`
	PercentileWaitTime90th  int64   `json:"90th_percentile_wait_time"`
	PercentileWaitTime99th  int64   `json:"99th_percentile_wait_time"`
	SamplesCollected        int     `json:"samples_collected"`
	SamplesSeen             int     `json:"samples_seen"`
	SampleMean              float64 `json:"sample_mean"`
	SampleStandardDeviation float64 `json:"sample_standard_deviation"`
	TvalueAverageWaitTime   float64 `json:"tvalue_average_wait_time"`
}
