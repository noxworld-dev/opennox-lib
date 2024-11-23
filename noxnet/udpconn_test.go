package noxnet

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func TestSeqBefore(t *testing.T) {
	cases := []struct {
		name   string
		v, cur byte
		exp    bool
	}{
		{"zero", 0, 0, true},
		{"equal", 100, 100, true},
		{"future", 10, 5, false},
		{"max left", maxAckMsgs / 4, maxAckMsgs / 2, true},
		{"max right", 0xff - maxAckMsgs/2, 0xff - maxAckMsgs/4, true},
		{"overflow", 0xff - maxAckMsgs/4, maxAckMsgs / 4, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			must.EqOp(t, c.exp, seqBefore(c.v, c.cur))
		})
	}
}

func TestStream(t *testing.T) {
	const maxMessages = 300

	newTest := func(t testing.TB, fast, debug bool, drop func(b []byte) bool) (*Conn, <-chan Message) {
		srvC, cliC := NewPipe(slog.Default(), maxAckMsgs)
		t.Cleanup(func() {
			_ = cliC.Close()
			_ = srvC.Close()
		})
		cliC.Drop = drop
		cliC.Debug = debug
		srvC.Debug = debug
		srvAddr := srvC.Addr

		srvRecv := make(chan Message, maxMessages)
		srv := srvC.Port
		srv.OnMessage(func(conn *Conn, sid StreamID, m Message) bool {
			if debug {
				t.Logf("server recv: %#v", m)
			}
			if fast {
				_ = conn.Ack()
			}
			srvRecv <- m
			return true
		})
		t.Cleanup(srv.Close)

		cli := cliC.Port
		cli.OnMessage(func(conn *Conn, sid StreamID, m Message) bool {
			if debug {
				t.Logf("client recv: %#v", m)
			}
			return true
		})
		t.Cleanup(cli.Close)

		srv.Start()
		cli.Start()

		return cli.Conn(srvAddr), srvRecv
	}

	// Test that general ACK mechanism works for a long sequence of messages.
	t.Run("sequential", func(t *testing.T) {
		cliConn, srvRecv := newTest(t, true, true, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 5*resendTick)
		defer cancel()

		timer := time.NewTimer(resendTick)

		for i := 0; i < maxMessages; i++ {
			exp := &MsgDiscover{Token: uint32(i + 1)}
			err := cliConn.SendReliableMsg(ctx, 0, exp)
			if err != nil {
				t.Fatal(err)
			}
			timer.Reset(resendTick)
			select {
			case m := <-srvRecv:
				must.Eq[Message](t, exp, m)
			case <-timer.C:
				t.Fatal("expected a message")
			}
		}
	})

	// Test that large queue works.
	t.Run("long queue", func(t *testing.T) {
		cliConn, srvRecv := newTest(t, false, false, nil)

		ctx := context.Background()

		var expected []Message
		for i := 0; i < maxMessages; i++ {
			exp := &MsgDiscover{Token: uint32(i + 1)}
			_, err := cliConn.QueueReliableMsg(ctx, 0, []Message{exp}, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			expected = append(expected, exp)
		}
		for _, exp := range expected {
			select {
			case m := <-srvRecv:
				must.Eq[Message](t, exp, m)
			case <-ctx.Done():
				t.Fatal("expected a message")
			}
		}
	})

	// Test redeliveries.
	t.Run("redelivery", func(t *testing.T) {
		dropped := 0
		cliConn, srvRecv := newTest(t, false, true, func(data []byte) bool {
			if dropped < resendRetries-1 {
				dropped++
				return true
			}
			return false
		})

		ctx := context.Background()

		exp := &MsgDiscover{Token: 0x123}
		err := cliConn.SendReliableMsg(ctx, 0, exp)
		if err != nil {
			t.Fatal(err)
		}
		select {
		case m := <-srvRecv:
			must.Eq[Message](t, exp, m)
		default:
			t.Fatal("expected a message")
		}
	})
}
