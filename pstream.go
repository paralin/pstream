// Package pstream provides the protobuf stream protocol.
package pstream

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	inet "github.com/libp2p/go-libp2p-net"
)

// Session wraps a stream in a session.
type Session struct {
	SessionDummy
	inet.Stream
	ctx       context.Context
	ctxCancel context.CancelFunc
	sendMtx   sync.Mutex
	readMtx   sync.Mutex
	compress  bool
}

// NewSession builds a new session.
func NewSession(ctx context.Context, ctxCancel context.CancelFunc, stream inet.Stream) *Session {
	return &Session{Stream: stream, ctx: ctx, ctxCancel: ctxCancel}
}

// NewSessionWithCompression builds a new session with compression enabled.
func NewSessionWithCompression(
	ctx context.Context,
	ctxCancel context.CancelFunc,
	stream inet.Stream,
) *Session {
	s := NewSession(ctx, ctxCancel, stream)
	s.compress = true
	return s
}

// Context returns the context.
func (s *Session) Context() context.Context {
	return s.ctx
}

// SendMsg tries to send a message on the wire.
func (s *Session) SendMsg(msgInter interface{}) error {
	msg := msgInter.(proto.Message)
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if s.compress {
		data = snappy.Encode(nil, data)
	}

	headerData, err := proto.Marshal(&Header{MessageLen: uint32(len(data))})
	if err != nil {
		return err
	}

	s.sendMtx.Lock()
	defer s.sendMtx.Unlock()

	if _, err := s.Stream.Write(headerData); err != nil {
		return err
	}

	if _, err := s.Stream.Write(data); err != nil {
		return err
	}

	return nil
}

// RecvMsg tries to receive a message on the wire.
func (s *Session) RecvMsg(msgInter interface{}) error {
	msg := msgInter.(proto.Message)
	header := &Header{MessageLen: 1}
	data := make([]byte, proto.Size(header))

	s.readMtx.Lock()
	defer s.readMtx.Unlock()

	if _, err := s.Read(data); err != nil {
		return err
	}

	if err := proto.Unmarshal(data, header); err != nil {
		return err
	}

	data = make([]byte, int(header.GetMessageLen()))
	if _, err := s.Read(data); err != nil {
		return err
	}

	if s.compress {
		var err error
		dataUncmp, err := snappy.Decode(nil, data)
		if err != nil {
			return err
		}
		data = dataUncmp
	}

	return proto.Unmarshal(data, msg)
}

func (s *Session) Close() error {
	if s.ctxCancel != nil {
		s.ctxCancel()
	}
	return s.Stream.Close()
}

func (s *Session) CloseSend() error {
	return s.Close()
}
