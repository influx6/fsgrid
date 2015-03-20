package fsgrid

import (
	"errors"
	// "fmt"
	"os"
	fpath "path/filepath"
	"strings"

	"github.com/influx6/grids"
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

	dir.In("read").Receive(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Get("file").(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			// fp.Body["err"] = err
			fp.Set("err", err)
			dir.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			// fp.Body["err"] = err
			fp.Set("err", err)
			dir.OutSend("err", fp)
			return
		}

		if mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				// fp.Body["err"] = err
				fp.Set("err", err)
				dir.OutSend("err", fp)
				return
			}

			list, err := rb.Readdir(-1)

			if err != nil {
				// fp.Body["err"] = err
				fp.Set("err", err)
				dir.OutSend("err", fp)
				return
			}

			pack := grids.NewPacket()
			pack.Copy(fp.ToMap())
			pack.Set("absolutefile", basefile)

			for _, val := range list {
				pack.Push(val)
			}

			pack.Freeze()

			dir.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a directory: " + basefile)
		// fp.Body["err"] = de
		fp.Set("err", de)
		dir.OutSend("err", fp)

		return
	})

	dir.In("write").Receive(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Get("file").(string)

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Set("err", err)
			dir.OutSend("err", fp)
			return
		}

		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			os.Mkdir(basefile, 0777)
		} else if !mod.Mode().IsDir() {
			// fp.Body["err"] = err
			fp.Set("err", err)
			dir.OutSend("err", fp)
			return
		}

		fp.Offload(func(i interface{}) {
			nm, ok := i.(string)

			if !ok {
				return
			}

			dirpath := fpath.Join(basefile, nm)
			err := os.Mkdir(dirpath, 0777)

			if err != nil {
				fp.Set("err", err)
				dir.OutSend("err", fp)
				return
			}
		})

		return

	})

	return dir
}

func CreateFSFile() *FSFile {
	var thefile = &FSFile{grids.NewGrid("fs.File")}

	thefile.NewIn("read")
	thefile.NewIn("write")

	thefile.NewOut("err")
	thefile.NewOut("res")

	root, _ := os.Getwd()

	thefile.In("read").Receive(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		// file, ok := fp.Body["file"].(string)
		file, ok := fp.Get("file").(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			// fp.Body["err"] = err
			fp.Set("err", err)
			thefile.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			// fp.Body["err"] = err
			fp.Set("err", err)
			thefile.OutSend("err", fp)
			return
		}

		if !mod.Mode().IsDir() {
			rb, err := os.Open(basefile)

			if err != nil {
				fp.Set("err", err)
				// fp.Body["err"] = err
				thefile.OutSend("err", fp)
				return
			}

			data := make([]byte, mod.Size())
			n, err := rb.Read(data)

			if err != nil {
				// fp.Body["err"] = err
				fp.Set("err", err)
				thefile.OutSend("err", fp)
				return
			}

			fp.Set("absoluteFile", basefile)
			fp.Set("read", n)

			pack := grids.NewPacket()
			pack.Clone(fp.Map)

			for _, val := range data {
				pack.Push(val)
			}

			pack.Freeze()

			thefile.OutSend("res", pack)

			return
		}

		de := errors.New("path is not a file: " + basefile)
		// fp.Body["err"] = de
		fp.Set("err", de)
		thefile.OutSend("err", fp)

		return
	})

	thefile.In("write").Receive(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		// file, ok := fp.Body["file"].(string)
		file, ok := fp.Get("file").(string)
		var basefile string

		if !fpath.IsAbs(file) {
			basefile = fpath.Join(root, file)
		} else {
			basefile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			// fp.Body["err"] = err
			fp.Set("err", err)
			thefile.OutSend("err", fp)
			return
		}

		mod, err := os.Stat(basefile)

		if err != nil {
			// fp.Body["err"] = err
			fp.Set("err", err)
			thefile.OutSend("err", fp)
			return
		}

		if !mod.Mode().IsDir() {
			rb, err := os.Open(basefile)
			if err != nil {
				// fp.Body["err"] = err
				fp.Set("err", err)
				thefile.OutSend("err", fp)
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
				// fp.Body["err"] = err
				fp.Set("err", err)
				thefile.OutSend("err", fp)
				return
			}

			return
		}

		de := errors.New("path is not a file: " + basefile)
		fp.Set("err", de)
		// fp.Body["err"] = de
		thefile.OutSend("err", fp)

		return
	})

	return thefile
}

func ReadFile(file string) *FileReader {
	var f = &FileReader{file, CreateFSFile()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Set("file", f.Path)
		// g.Body["file"] = f.Path
		f.InSend("read", g)
	})

	return f
}

func ReadDir(file string) *DirReader {
	var f = &DirReader{file, CreateFSDir()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Set("file", f.Path)
		f.InSend("read", g)
	})

	return f
}

func WriteFile(file string) *FileWriter {
	var f = &FileWriter{file, CreateFSFile()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Set("file", f.Path)
		f.InSend("write", g)
	})

	return f
}

func WriteDir(file string) *DirWriter {
	var f = &DirWriter{file, CreateFSDir()}
	f.NewIn("file")

	f.OrIn("file", func(g *grids.GridPacket) {
		g.Set("file", f.Path)
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

	in.Receive(func(i interface{}) {
		fp, ok := i.(*grids.GridPacket)

		if !ok {
			return
		}

		file, ok := fp.Get("file").(string)
		var bfile string

		if !fpath.IsAbs(file) {
			bfile = fpath.Join(basefile, file)
		} else {
			bfile = file
		}

		if !ok {
			err := errors.New("Invalid packet map, no 'file' included")
			fp.Set("err", err)
			// fp.Body["err"] = err
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
				fp.Set("err", err)
				// fp.Body["err"] = err
				dir.OutSend("err", fp)
				return
			}
		}

		fp.Set("oldfile", file)
		fp.Set("file", bfile)

		dir.OutSend("res", fp)
	})

	return dir, nil
}
