package fsgrid

import (
	"errors"
	// "fmt"
	"github.com/influx6/grids"
	"os"
	fpath "path/filepath"
	"strings"
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

type FileReader struct {
	Path string
	*FSFile
}

type DirReader struct {
	Path string
	*FSDir
}

type FileWriter struct {
	Path string
	*FSFile
}

type DirWriter struct {
	Path string
	*FSDir
}

func CreateFSDir() *FSDir {
	dir := &FSDir{grids.NewGrid("fs.Dir")}

	dir.NewIn("read")
	dir.NewIn("write")

	dir.NewOut("err")
	dir.NewOut("res")

	root, _ := os.Getwd()

	dir.In("read").Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Body["file"].(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		if mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			list, err := rb.Readdir(-1)

			if err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			pack := grids.CreateGridPacket(fp.Body)
			pack.Body["absolutefile"] = basefile

			for _, val := range list {
				pack.Push(val)
			}

			pack.Freeze()

			dir.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a directory: " + basefile)
		fp.Body["err"] = de
		dir.OutSend("err", fp)

		return
	})

	dir.In("write").Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Body["file"].(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		if mod.Mode().IsDir() {
			fp.Offload(func(i interface{}) {
				nm, ok := i.(string)

				if !ok {
					return
				}

				dirpath := fpath.Join(basefile, nm)
				err := os.Mkdir(dirpath, 0777)

				if err != nil {
					fp.Body["err"] = err
					dir.OutSend("err", fp)
					return
				}
			})

			return
		}

		de := errors.New("path is not a directory: " + basefile)
		fp.Body["err"] = de
		dir.OutSend("err", fp)

		return
	})

	return dir
}

func CreateFSFile() *FSFile {
	var dir = &FSFile{grids.NewGrid("fs.File")}

	dir.NewIn("read")
	dir.NewIn("write")

	dir.NewOut("err")
	dir.NewOut("res")

	root, _ := os.Getwd()

	dir.In("read").Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Body["file"].(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		if !mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			data := make([]byte, mod.Size())
			n, err := rb.Read(data)

			if err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			fm := fp.Body
			fm["absoluteFile"] = basefile
			fm["read"] = n
			pack := grids.CreateGridPacket(fm)

			for _, val := range data {
				pack.Push(val)
			}

			pack.Freeze()

			dir.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a file: " + basefile)
		fp.Body["err"] = de
		dir.OutSend("err", fp)

		return
	})

	dir.In("write").Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Body["file"].(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		if !mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			data := make([]byte, 0)

			fp.Offload(func(f interface{}) {
				b, ok := f.(byte)
				if !ok {
					return
				}

				data = append(data, b)
			})

			if _, err := rb.Write(data); err != nil {
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}

			return
		}

		de := errors.New("path is not a file: " + basefile)
		fp.Body["err"] = de
		dir.OutSend("err", fp)

		return
	})

	return dir
}

func ReadFile(file string) *FileReader {
	var f = &FileReader{file, CreateFSFile()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Body["file"] = f.Path
		f.InSend("read", g)
	})

	return f
}

func ReadDir(file string) *DirReader {
	var f = &DirReader{file, CreateFSDir()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Body["file"] = f.Path
		f.InSend("read", g)
	})

	return f
}

func WriteFile(file string) *FileWriter {
	var f = &FileWriter{file, CreateFSFile()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Body["file"] = f.Path
		f.InSend("write", g)
	})

	return f
}

func WriteDir(file string) *DirWriter {
	var f = &DirWriter{file, CreateFSDir()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Body["file"] = f.Path
		f.InSend("write", g)
	})

	return f
}

func CreateFSControl(base string) (*FSControl, error) {
	var dir = &FSControl{grids.NewGrid("fs.Control")}

	dir.NewIn("file")
	dir.NewOut("err")
	dir.NewOut("res")

	root, _ := os.Getwd()
	basefile := fpath.Join(root, base)

	in := dir.In("file")

	in.Or(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Body["file"].(string)
		var bfile string

		if !fpath.IsAbs(file) {
			bfile = fpath.Join(basefile, file)
		} else {
			bfile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Body["err"] = err
			dir.OutSend("err", fp)
			return
		}

		hp := strings.HasPrefix(bfile, basefile)

		if !hp {
			dir.OutSend("err", fp)
			return
		}

		_, err := os.Stat(bfile)

		if err != nil {
			if os.IsNotExist(err) {
				err := errors.New("Invalid path" + bfile)
				fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}
		}

		fp.Body["oldfile"] = file
		fp.Body["file"] = bfile

		dir.OutSend("res", fp)
	})

	return dir, nil
}
