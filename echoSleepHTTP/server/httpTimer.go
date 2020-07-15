package server

import (
	"github.com/google/uuid"
	"log"

	//"log"
	"strings"
	"time"
)

type TimeSpanner struct {
	SpanName              string
	ExternalCorrelationId uuid.UUID
	SpanCorrelationId     uuid.UUID
	Spans                 []Span
}

type Span struct {
	SpanName          string
	SpanCorrelationId uuid.UUID
	Timer             SpanTimer
	SubSpans          []Span
}

type SpanTimer struct {
	SpanStartTime time.Time
	SpanStopTime time.Time
	SpanDuration  time.Duration
}

// NOTE - spanName is REQUIRED (ideally, unique)
// ExternalCorrelationId is OPTIONAL (will be assigned if not provided
func NewTimeSpanner(spanName string, externalCorrelationId string) *TimeSpanner {
	var err error
	var extUUID uuid.UUID
	if strings.Compare(externalCorrelationId, "") == 0 {
		extUUID = uuid.Must(uuid.New(), err)
	} else {
		extUUID, err = uuid.Parse(externalCorrelationId)
	}
	if err != nil {
		extUUID = uuid.Must(uuid.New(), err)
	}
	return &TimeSpanner{spanName, extUUID, uuid.Must(uuid.New(), err), make([]Span, 0, 5)}
}

// Add a new time span to the slice
func (h *TimeSpanner) addTimeSpan(spanName string) (span *Span, position int) {
	var err error
	span = &Span{spanName, uuid.Must(uuid.New(), err), SpanTimer{time.Now(), time.Time{}, 0}, make([]Span, 0, 5)}
	h.Spans = append(h.Spans,*span)
  for i,v := range (h.Spans) {
  	 if v.SpanName == "" {
  	 	break
	 }
  	 span = &h.Spans[i]
  	 position = i
	}
	log.Println("position:" +string(position))
  	log.Println(span )
  log.Println("======")
	return
}

// Start the timer by setting start time to Now
func (s *Span) StartTimer() {
	s.Timer.SpanStartTime = time.Now()
}

// Retrieve the recorded start time
func (s *Span) GetStartTime() string {
	return s.Timer.SpanStartTime.String()
}

// Retrieve the saved duration
func (s *Span) GetDuration() string {
	s.Timer.SpanDuration = s.Timer.SpanStopTime.Sub(s.Timer.SpanStartTime)
//	log.Println("my duration time is " + s.Timer.SpanDuration.String())
	return s.Timer.SpanDuration.String()
}

// Stop the timer and record duration
func (s *Span) StopTimer() {
	s.Timer.SpanStopTime = time.Now()
}

func (s *Span) addSubspan(spanName string) (span *Span, position int) {
	var err error
	span = &Span{SpanName: spanName, SpanCorrelationId: uuid.Must(uuid.New(), err), Timer: SpanTimer{SpanStartTime: time.Now() }, SubSpans: make([]Span, 1)}
	s.SubSpans = append(s.SubSpans,*span)
	for i,v := range (s.SubSpans) {
		if v.SpanName == "" {
			break
		}
		span = &s.SubSpans[i]
		position = i
	}

	return
}
