package server

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type TimeSpanner struct {
	SpanName              string
	ExternalCorrelationId uuid.UUID
	SpanCorrelationId     uuid.UUID
	spans                 []Span
}

type Span struct {
	SpanName          string
	SpanCorrelationId uuid.UUID
	timer             SpanTimer
	subSpans          []Span
}

type SpanTimer struct {
	SpanStartTime time.Time
	SpanDuration  time.Duration
}

// NOTE - spanName is REQUIRED (ideally, unique)
// ExternalCorrelationId is OPTIONAL (will be assigned if not provided
func NewTimeSpanner(spanName string, externalCorrelationId string) TimeSpanner {
	var err error
	var extUUID uuid.UUID
	if strings.Compare(externalCorrelationId,"") == 0 {
		extUUID = uuid.Must(uuid.New(), err)
	} else {
		extUUID, err = uuid.Parse(externalCorrelationId)
	}
	if err != nil {
		extUUID = uuid.Must(uuid.New(), err)
	}
	return TimeSpanner{spanName, extUUID, uuid.Must(uuid.New(), err), make([]Span, 2)}
}

func (h TimeSpanner) addTimeSpan(spanName string) Span {
	var err error
	span := Span{spanName, uuid.Must(uuid.New(), err), SpanTimer{time.Now(), time.Duration(0)}, make([]Span, 1)}
	h.spans = append(h.spans, span)
	return span
}

func (h Span) addSubspan(spanName string) Span {
	var err error
	span := Span{spanName, uuid.Must(uuid.New(), err), SpanTimer{time.Now(), time.Duration(0)}, make([]Span, 1)}
	h.subSpans = append(h.subSpans, span)
	return span
}
