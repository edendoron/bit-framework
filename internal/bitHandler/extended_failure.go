package bithandler

import (
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ExtendedFailure struct {
	Failure       bit.Failure
	Time          time.Time
	failureCount  uint64
	reportsCount  uint32
	startReportId float64
	endReportId   float64
}

func (e *ExtendedFailure) ExtendedFailureToBitStatusReportedFailure() bit.BitStatus_RportedFailure {
	return bit.BitStatus_RportedFailure{
		FailureData: e.Failure.Description,
		Timestamp:   timestamppb.New(e.Time),
		Count:       e.failureCount,
	}
}

func FailuresSliceToExtendedFailuresSlice(failures []bit.Failure) []ExtendedFailure {
	var extendedFailures []ExtendedFailure
	for _, fail := range failures {
		extendedFailures = append(extendedFailures, FailureToExtendedFailure(fail))
	}
	return extendedFailures
}

func FailureToExtendedFailure(failure bit.Failure) ExtendedFailure {
	return ExtendedFailure{
		Failure:       failure,
		failureCount:  0,
		reportsCount:  0,
		startReportId: 0,
		endReportId:   0,
	}
}

// function to support protobuf

func (e *ExtendedFailure) Reset() { *e = ExtendedFailure{} }

func (e *ExtendedFailure) String() string { return proto.CompactTextString(e) }

func (e *ExtendedFailure) ProtoMessage() {}
