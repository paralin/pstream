package pstream

import (
	"google.golang.org/grpc/metadata"
)

// SessionDummy contains common functions we dont use.
type SessionDummy struct{}

// SetHeader sets the header metadata. It may be called multiple times.
// This is no-op
func (s *SessionDummy) SetHeader(metadata.MD) error { return nil }

// SendHeader sends the header metadata.
// This is no-op
// The provided md and headers set by SetHeader() will be sent.
// It fails if called multiple times.
func (s *SessionDummy) SendHeader(metadata.MD) error { return nil }

// SetTrailer sets the trailer metadata which will be sent with the RPC status.
// This is no-op
// When called more than once, all the provided metadata will be merged.
func (s *SessionDummy) SetTrailer(metadata.MD) {}

// Header returns the header metadata received from the server if there
// is any. It blocks if the metadata is not ready to read.
// This is no-op
func (s *SessionDummy) Header() (metadata.MD, error) { return nil, nil }

// Trailer returns the trailer metadata from the server, if there is any.
// It must only be called after stream.CloseAndRecv has returned, or
// stream.Recv has returned a non-nil error (including io.EOF).
// This is no-op
func (s *SessionDummy) Trailer() metadata.MD { return nil }
