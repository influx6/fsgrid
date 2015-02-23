package fsgrid

import (
	"errors"
	// "fmt"
	"github.com/influx6/grids"
	"os"
	fpath "path/filepath"
)

type FSDir struct {
	*grids.Grid
}

type FSFile struct {
	*grids.Grid
}

type FSControl struct {
	*grids.Grid
}

func CreateFSDir() *FSDir {
	dir := &FSDir{grids.NewGrid("fs.Dir")}

	dir.NewIn("file")
	dir.NewOut("err")
	dir.NewOut("res")

	in := dir.In("file")
	root, _ := os.Getwd()

	in.Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			err := grids.CreateGridPacket(grids.GridMap{"err": errors.New("Incorrect type")})
			dir.OutSend("err", err)
			return
		}

		file, ok := fp.Body["file"].(string)
		basefile := fpath.Join(root, file)

		if !ok {
			err := grids.CreateGridPacket(grids.GridMap{"err": errors.New("Invalid packet map, no 'file' included")})
			dir.OutSend("err", err)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			err := grids.CreateGridPacket(grids.GridMap{"err": err})
			dir.OutSend("err", err)
			return
		}

		if mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				err := grids.CreateGridPacket(grids.GridMap{"err": err})
				dir.OutSend("err", err)
				return
			}

			list, err := rb.Readdir(-1)

			if err != nil {
				err := grids.CreateGridPacket(grids.GridMap{"err": err})
				dir.OutSend("err", err)
				return
			}

			pack := grids.CreateGridPacket(grids.GridMap{"file": file, "absoluteFile": basefile})

			for _, val := range list {
				pack.Push(val)
			}

			pack.Freeze()

			dir.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a directory: " + basefile)
		fe := grids.CreateGridPacket(grids.GridMap{"err": de})
		dir.OutSend("err", fe)

		return
	})

	return dir
}

func CreateFSFile() *FSFile {
	var dir = &FSFile{grids.NewGrid("fs.File")}

	dir.NewIn("file")
	dir.NewOut("err")
	dir.NewOut("res")

	in := dir.In("file")
	root, _ := os.Getwd()

	in.Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			err := grids.CreateGridPacket(grids.GridMap{"err": errors.New("Incorrect type")})
			dir.OutSend("err", err)
			return
		}

		file, ok := fp.Body["file"].(string)
		basefile := fpath.Join(root, file)

		if !ok {
			err := grids.CreateGridPacket(grids.GridMap{"err": errors.New("Invalid packet map, no 'file' included")})
			dir.OutSend("err", err)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			err := grids.CreateGridPacket(grids.GridMap{"err": err})
			dir.OutSend("err", err)
			return
		}

		if !mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				err := grids.CreateGridPacket(grids.GridMap{"err": err})
				dir.OutSend("err", err)
				return
			}

			data := make([]byte, mod.Size())
			n, err := rb.Read(data)

			if err != nil {
				err := grids.CreateGridPacket(grids.GridMap{"err": err})
				dir.OutSend("err", err)
				return
			}

			pack := grids.CreateGridPacket(grids.GridMap{"file": file, "absoluteFile": basefile, "read": n})

			for _, val := range data {
				pack.Push(val)
			}

			pack.Freeze()

			dir.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a file: " + basefile)
		fe := grids.CreateGridPacket(grids.GridMap{"err": de})
		dir.OutSend("err", fe)

		return
	})

	return dir
}
