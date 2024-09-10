package utils

import (
	"regexp"
	"time"
)

const (
	MAX_DIMENSION_SET_SIZE     = 30
	MAX_DIMENSION_NAME_LENGTH  = 250
	MAX_DIMENSION_VALUE_LENGTH = 1024
	MAX_METRIC_NAME_LENGTH     = 1024
	MAX_NAMESPACE_LENGTH       = 256
	MAX_TIMESTAMP_PAST_AGE     = 14 * 24 * time.Hour // 2 weeks in time.Duration
	MAX_TIMESTAMP_FUTURE_AGE   = 2 * time.Hour       // 2 hours in time.Duration
	DEFAULT_NAMESPACE          = "aws-embedded-metrics"
	MAX_METRICS_PER_EVENT      = 100
	MAX_VALUES_PER_METRIC      = 100
	DEFAULT_AGENT_HOST         = "0.0.0.0"
	DEFAULT_AGENT_PORT         = 25888
)

var (
	VALID_NAMESPACE_REGEX = regexp.MustCompile(`^[a-zA-Z0-9._#:/-]+$`)
	VALID_DIMENSION_REGEX = regexp.MustCompile(`^[\x00-\x7F]+$`)
)

type StorageResolution int

const (
	High     StorageResolution = 1
	Standard StorageResolution = 60
)

var StorageResolutions = []StorageResolution{High, Standard}

type Unit string

const (
	Seconds            Unit = "Seconds"
	Microseconds       Unit = "Microseconds"
	Milliseconds       Unit = "Milliseconds"
	Bytes              Unit = "Bytes"
	Kilobytes          Unit = "Kilobytes"
	Megabytes          Unit = "Megabytes"
	Gigabytes          Unit = "Gigabytes"
	Terabytes          Unit = "Terabytes"
	Bits               Unit = "Bits"
	Kilobits           Unit = "Kilobits"
	Megabits           Unit = "Megabits"
	Gigabits           Unit = "Gigabits"
	Terabits           Unit = "Terabits"
	Percent            Unit = "Percent"
	Count              Unit = "Count"
	BytesPerSecond     Unit = "Bytes/Second"
	KilobytesPerSecond Unit = "Kilobytes/Second"
	MegabytesPerSecond Unit = "Megabytes/Second"
	GigabytesPerSecond Unit = "Gigabytes/Second"
	TerabytesPerSecond Unit = "Terabytes/Second"
	BitsPerSecond      Unit = "Bits/Second"
	KilobitsPerSecond  Unit = "Kilobits/Second"
	MegabitsPerSecond  Unit = "Megabits/Second"
	GigabitsPerSecond  Unit = "Gigabits/Second"
	TerabitsPerSecond  Unit = "Terabits/Second"
	CountPerSecond     Unit = "Count/Second"
	None               Unit = "None"
)

var Units = []Unit{Seconds, Microseconds, Milliseconds, Bytes, Kilobytes, Megabytes, Gigabytes, Terabytes, Bits, Kilobits, Megabits, Gigabits, Terabits, Percent, Count, BytesPerSecond, KilobytesPerSecond, MegabytesPerSecond, GigabytesPerSecond, TerabytesPerSecond, BitsPerSecond, KilobitsPerSecond, MegabitsPerSecond, GigabitsPerSecond, TerabitsPerSecond, CountPerSecond, None}
