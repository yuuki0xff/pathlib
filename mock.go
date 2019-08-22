package pathlib

import (
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func NewMock(fs afero.Fs, path string) Path {
	return &MockPath{
		Fs:   fs,
		Path: path,
	}
}

type MockPath struct {
	Fs   afero.Fs
	Path string
}

func (p *MockPath) Absolute() (Path, error) {
	cwd, err := p.Cwd()
	if err != nil {
		return nil, err
	}
	return p.JoinPath(cwd.String(), p.Path), nil
}
func (p *MockPath) Cwd() (Path, error) {
	return NewMock(p.Fs, "/"), nil
}
func (p *MockPath) Parent() (Path, error) {
	return NewMock(p.Fs, path.Dir(p.Path)), nil
}
func (p *MockPath) JoinPath(elem ...string) Path {
	temp := []string{p.Path}
	elem = append(temp, elem...)
	return NewMock(p.Fs, path.Join(elem...))
}
func (p *MockPath) String() string {
	return p.Path
}
func (p *MockPath) Touch() error {
	f, err := p.Fs.Create(p.Path)
	if err != nil {
		return err
	}
	return f.Close()
}
func (p *MockPath) Unlink() error {
	return p.Fs.Remove(p.Path)
}
func (p *MockPath) RmDir() error {
	return p.Unlink()
}
func (p *MockPath) MkDir(mode os.FileMode, parents bool) error {
	if parents {
		return p.Fs.MkdirAll(p.Path, mode)
	}
	return p.Fs.Mkdir(p.Path, mode)
}
func (p *MockPath) Open() (io.ReadCloser, error) {
	f, err := p.Fs.Open(p.Path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (p *MockPath) OpenRW(flag int, mode os.FileMode) (io.ReadWriteCloser, error) {
	f, err := p.Fs.OpenFile(p.Path, flag, mode)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (p *MockPath) Chmod(mode os.FileMode) error {
	return p.Fs.Chmod(p.Path, mode)
}
func (p *MockPath) Rename(target Path) error {
	return p.Fs.Rename(p.Path, target.String())
}
func (p *MockPath) Exists() bool {
	_, err := p.Fs.Stat(p.Path)
	return err == nil
}
func (p *MockPath) IsDir() bool {
	fi, err := p.Fs.Stat(p.Path)
	return err == nil && fi.IsDir()
}
func (p *MockPath) IsFile() bool {
	fi, err := p.Fs.Stat(p.Path)
	return err == nil && !fi.IsDir()
}
func (p *MockPath) IsAbs() bool {
	return filepath.IsAbs(p.Path)
}
func (p *MockPath) ReadBytes() ([]byte, error) {
	r, err := p.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
func (p *MockPath) ReadText() (string, error) {
	// TODO: transform encoding
	b, err := p.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (p *MockPath) WriteBytes(data []byte) error {
	w, err := p.OpenRW(os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return io.ErrShortWrite
	}
	return nil
}
func (p *MockPath) WriteText(text string) error {
	// TODO: transform encoding
	return p.WriteBytes([]byte(text))
}
