# Packet Stream

> Simple protocol for streaming packets with a length prefix.

## Usage

```go
type Session struct {
	*pstream.Session
}

// RemoteEvent is defined as a protobuf message.

// readPump reads messages.
func (s *Session) readPump(ctx context.Context) {
	for {
		var msg RemoteEvent
		if err := s.Session.RecvMsg(&msg); err != nil {
			if err != io.EOF {
				le.WithError(err).Error("peer connection closed")
			}

			return
		}
        
        // handle msg
	}
}

subCtx, subCtxCancel := context.WithCancel(r.ctx)
pSess := pstream.NewSessionWithCompression(subCtx, subCtxCancel, stream)
sess := &Session{Session: pSess}
go sess.readPump(subCtx)
```
