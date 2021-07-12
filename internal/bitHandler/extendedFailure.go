package bitHandler

import (
	. "../../configs/rafael.com/bina/bit"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ExtendedFailure struct {
	Failure       Failure
	time          time.Time
	failureCount  uint64
	reportsCount  uint32
	startReportId float64
	endReportId   float64
}

func (e *ExtendedFailure) Reset() { *e = ExtendedFailure{} }

func (e *ExtendedFailure) String() string { return proto.CompactTextString(e) }

func (e *ExtendedFailure) ProtoMessage() {}

func (e *ExtendedFailure) ExtendedFailureToBitStatusReportedFailure() BitStatus_RportedFailure {
	return BitStatus_RportedFailure{
		FailureData: e.Failure.Description,
		Timestamp:   timestamppb.New(e.time),
		Count:       e.failureCount,
	}
}

func FailuresSliceToExtendedFailuresSlice(failures []Failure) []ExtendedFailure {
	var extendedFailures []ExtendedFailure
	for _, fail := range failures {
		extendedFailures = append(extendedFailures, FailureToExtendedFailure(fail))
	}
	return extendedFailures
}

func FailureToExtendedFailure(failure Failure) ExtendedFailure {
	return ExtendedFailure{
		Failure:       failure,
		failureCount:  0,
		reportsCount:  0,
		startReportId: 0,
		endReportId:   0,
	}
}
