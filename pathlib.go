package pathlib

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
)

// OsPath
type OsPath struct {
	Path string
}

// New Returns a new path.
func New(path string) Path {
	p := new(OsPath)
	p.Path = path
	return p
}

// Absolute Returns an absolute representation of path.
func (p *OsPath) Absolute() (Path, error) {
	pth, err := filepath.Abs(p.Path)
	if err != nil {
		return nil, errors.Wrap(err, "get absolute failed")
	}
	newP := New(pth)
	return newP, nil
}

// Cwd Return a new path pointing to the current working directory.
func (p *OsPath) Cwd() (Path, error) {
	pth, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get cwd failed")
	}
	newP := New(pth)
	return newP, nil
}

// Parent Return a new path for current path parent.
func (p *OsPath) Parent() (Path, error) {
	pth, err := p.Absolute()
	if err != nil {
		return nil, errors.Wrap(err, "get parent failed")
	}
	dir := filepath.Dir(pth.String())
	newP := New(dir)
	return newP, nil
}

// Touch Create creates the named file with mode 0666 (before umask), regardless of whether it exists.
func (p *OsPath) Touch() error {
	f, err := os.Create(p.Path)
	if err != nil {
		return err
	}
	return f.Close()
}

// Unlink Remove this file or link.
func (p *OsPath) Unlink() error {
	err := syscall.Unlink(p.Path)
	return err
}

// RmDir Remove this directory. The directory must be empty.
func (p *OsPath) RmDir() error {
	err := os.Remove(p.Path)
	return err
}

// MkDir Create a new directory at this given path.
func (p *OsPath) MkDir(mode os.FileMode, parents bool) (err error) {
	if parents {
		err = os.MkdirAll(p.Path, mode)
	} else {
		err = os.Mkdir(p.Path, mode)
	}
	return
}

// Open opens the named file for reading.
func (p *OsPath) Open() (FileRO, error) {
	f, err := os.Open(p.Path)
	if err != nil {
		return nil, err
	}
	return f, err
}

// Open opens the named file for reading and writing.
func (p *OsPath) OpenRW(flag int, mode os.FileMode) (FileRW, error) {
	f, err := os.OpenFile(p.Path, os.O_RDWR|flag, mode)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Chmod changes the mode of the named file to mode.
func (p *OsPath) Chmod(mode os.FileMode) error {
	return os.Chmod(p.Path, mode)
}

// JoinPath Returns a new path, Combine current path with one or several arguments
func (p *OsPath) JoinPath(elem ...string) Path {
	temp := []string{p.Path}
	elem = append(temp, elem[0:]...)
	newP := New(path.Join(elem...))
	return newP
}

// String returns the file path represented by string.
func (p *OsPath) String() string {
	return p.Path
}

// Rename renames (moves) the file or directory to target.
func (p *OsPath) Rename(target Path) error {
	return os.Rename(p.Path, target.String())
}

// Exists reports current path parent exists.
func (p *OsPath) Exists() bool {
	_, err := os.Stat(p.Path)
	return err == nil || os.IsExist(err)
}

// IsDir reports Whether this path is a directory.
func (p *OsPath) IsDir() bool {
	f, err := os.Stat(p.Path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// IsFile reports Whether this path is a regular file.
func (p *OsPath) IsFile() bool {
	f, e := os.Stat(p.Path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsAbs reports whether the path is absolute.
func (p *OsPath) IsAbs() bool {
	return filepath.IsAbs(p.Path)
}

// ReadBytes reads the file named by filename and returns the contents.
func (p *OsPath) ReadBytes() ([]byte, error) {
	buf, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadText reads the file and returns the contents.
func (p *OsPath) ReadText() (string, error) {
	// TODO: transform encoding
	b, err := p.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(b), err
}

// WriteBytes writes a byte slice to the file.
func (p *OsPath) WriteBytes(data []byte) error {
	return ioutil.WriteFile(p.Path, data, os.ModePerm)
}

// WriteText writes a text to the file.
func (p *OsPath) WriteText(text string) error {
	// TODO: transform encoding
	return p.WriteBytes([]byte(text))
}
