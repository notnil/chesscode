package main

import (
	"syscall/js"

	"github.com/notnil/chess"
	"github.com/notnil/chesscode"
)

func main() {
	done := make(chan struct{}, 0)
	global := js.Global()
	global.Set("encode", js.FuncOf(encode))
	global.Set("decode", js.FuncOf(decode))
	<-done
}

func encode(this js.Value, args []js.Value) interface{} {
	s := args[0].String()
	b, err := chesscode.Encode(s)
	if err != nil {
		return chess.NewBoard(nil).String()
	}
	return b.String()
}

func decode(this js.Value, args []js.Value) interface{} {
	s := args[0].String()
	b := &chess.Board{}
	if err := b.UnmarshalText([]byte(s)); err != nil {
		return ""
	}
	out, err := chesscode.Decode(b)
	if err != nil {
		return ""
	}
	return out
}
