package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ExtendedFailure struct {
	failure       Failure
	time          time.Time
	failureCount  uint64
	reportsCount  uint64
	startReportId float64
	endReportId   float64
}

func (e *ExtendedFailure) Reset() { *e = ExtendedFailure{} }

func (e *ExtendedFailure) String() string { return proto.CompactTextString(e) }

func (e *ExtendedFailure) ProtoMessage() {}

func (e *ExtendedFailure) extendedFailureToBitStatusReportedFailure() BitStatus_RportedFailure {
	return BitStatus_RportedFailure{
		FailureData: e.failure.Description,
		Timestamp:   timestamppb.New(e.time),
		Count:       e.failureCount,
	}
}

func failureToExtendedFailure(failure Failure, timestamp time.Time, countFailed uint64) ExtendedFailure {
	return ExtendedFailure{
		failure:      failure,
		time:         timestamp,
		failureCount: countFailed,
	}
}
