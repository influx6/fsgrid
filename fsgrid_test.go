package fsgrid

import (
	// "fmt"
	"github.com/influx6/grids"
	"os"
	"testing"
)

func TestReadDirectory(t *testing.T) {
	dir := CreateFSDir()
	packet := grids.NewPacket()
	packet.Set("file", ".")

	ev := dir.Out("res")
	ev.Receive(func(i interface{}) {
		res, ok := i.(*grids.GridPacket)

		if !ok {
			t.Fatalf("value is not a gridpacket", i, res, dir)
		}

		res.Offload(func(v interface{}) {
			k, ok := v.(os.FileInfo)
			if !ok {
				t.Fatalf("path value is not a fileInfo", k, v, dir)
			}
		})
	})

	dir.InSend("file", packet)
}

func TestReadFile(t *testing.T) {
	file := CreateFSFile()
	packet := grids.NewPacket()
	packet.Set("file", "./fsgrid.go")

	ev := file.Out("res")
	ev.Receive(func(i interface{}) {
		res, ok := i.(*grids.GridPacket)

		if !ok {
			t.Fatalf("value is not a gridpacket", i, res, file)
		}

		data := res.Obj()

		if len(data.([]interface{})) <= 0 {
			t.Fatalf("buffer is empty", data, res, file)
		}
	})

	file.InSend("read", packet)
}

func TestFileControl(t *testing.T) {
	file, _ := CreateFSControl("./")
	packet := grids.NewPacket()
	packet.Set("file", "./fsgrid.go")

	epacket := grids.NewPacket()
	epacket.Set("file", "./reflowj.go")

	rev := file.Out("res")
	re := file.Out("err")

	re.Receive(func(i interface{}) {
		res, ok := i.(*grids.GridPacket)

		if !ok {
			t.Fatalf("value is not a gridpacket", i, res, file)
		}

		if !res.Has("err") {
			t.Fatalf("should have an error attr", res, file, epacket)
		}
	})

	rev.Receive(func(i interface{}) {
		res, ok := i.(*grids.GridPacket)

		if !ok {
			t.Fatalf("value is not a gridpacket", i, res, file)
		}

		if res != packet {
			t.Fatalf("value is a different object", res, file)
		}
	})

	file.InSend("file", packet)
	file.InSend("file", epacket)
}
