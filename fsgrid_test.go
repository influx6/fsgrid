package fsgrid

import (
	"github.com/influx6/grids"
	"os"
	"testing"
)

func TestReadDirectory(t *testing.T) {
	dir := CreateFSDir()
	file := grids.GridMap{"file": "."}
	packet := grids.CreateGridPacket(file)

	ev := dir.Out("res")
	ev.Or(func(i interface{}) {
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
	packet := grids.CreateGridPacket(grids.GridMap{"file": "./fsgrid.go"})

	ev := file.Out("res")
	ev.Or(func(i interface{}) {
		res, ok := i.(*grids.GridPacket)

		if !ok {
			t.Fatalf("value is not a gridpacket", i, res, file)
		}

		data := res.Obj()

		if len(data.([]interface{})) <= 0 {
			t.Fatalf("buffer is empty", data, res, file)
		}
	})

	file.InSend("file", packet)
}
