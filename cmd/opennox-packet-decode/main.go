package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/noxworld-dev/opennox-lib/noxnet"
)

var (
	fIn  = flag.String("i", "network.jsonl", "input file with packet capture")
	fOut = flag.String("o", "network-dec.jsonl", "output file for decoded packets")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	f, err := os.Open(*fIn)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	w, err := os.Create(*fOut)
	if err != nil {
		return err
	}
	defer w.Close()
	enc := json.NewEncoder(w)

	for {
		var r RecordIn
		err := dec.Decode(&r)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		r2 := r.Decode()
		if err = enc.Encode(r2); err != nil {
			return err
		}
	}
	return w.Close()
}

type RecordIn struct {
	SrcID uint32 `json:"src_id"`
	DstID uint32 `json:"dst_id"`
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Data  string `json:"data"`
}

func (r RecordIn) Decode() RecordOut {
	o := RecordOut{
		SrcID: r.SrcID,
		DstID: r.DstID,
		Src:   r.Src,
		Dst:   r.Dst,
		Data:  r.Data,
	}
	raw, err := hex.DecodeString(r.Data)
	if err != nil {
		return o
	}
	o.Len = len(raw)
	if len(raw) < 2 {
		return o
	}
	hdr, data := raw[:2], raw[2:]
	o.Hdr = hex.EncodeToString(hdr)
	reliable := hdr[0]&0x80 != 0
	o.SID = hdr[0] &^ 0x80
	if seq := hdr[1]; reliable {
		o.Syn = &seq
	} else {
		o.Ack = &seq
	}
	if len(data) == 1 {
		op := noxnet.Op(data[0])
		if _, _, err := noxnet.DecodeAnyPacket(data); err != nil {
			s := op.String()
			o.Op = &s
			return o
		}
	}
	allSplit := true
	for len(data) != 0 {
		op := noxnet.Op(data[0])
		sz := len(data)
		var v any
		lenOK := false
		if n := op.Len(); n >= 0 && n <= len(data) {
			sz = n + 1
			lenOK = true
		}
		if m, n, err := noxnet.DecodeAnyPacket(data); err == nil && n > 0 {
			sz = n
			v = m
			lenOK = true
		} else if err != nil {
			slog.Error("cannot decode message", "op", op, "err", err)
		}
		if !lenOK {
			allSplit = false
		}
		msg := data[:sz]
		data = data[sz:]
		m := Msg{
			Op:     op.String(),
			Len:    len(msg),
			Data:   hex.EncodeToString(msg),
			Fields: v,
		}
		o.Msgs = append(o.Msgs, m)
	}
	if allSplit {
		o.Data = ""
	}
	return o
}

type RecordOut struct {
	SrcID uint32  `json:"src_id"`
	DstID uint32  `json:"dst_id"`
	Src   string  `json:"src"`
	Dst   string  `json:"dst"`
	Hdr   string  `json:"hdr"`
	SID   byte    `json:"sid"`
	Syn   *byte   `json:"syn,omitempty"`
	Ack   *byte   `json:"ack,omitempty"`
	Len   int     `json:"len"`
	Op    *string `json:"op,omitempty"`
	Msgs  []Msg   `json:"msgs,omitempty"`
	Data  string  `json:"data,omitempty"`
}

type Msg struct {
	Op     string `json:"op,omitempty"`
	Fields any    `json:"fields,omitempty"`
	Len    int    `json:"len"`
	Data   string `json:"data"`
}
