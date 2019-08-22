package pathlib

import (
	"io"
	"os"
)

// FileRO represents the opened file for reading.
type FileRO interface {
	io.ReadCloser
	io.Seeker
	io.ReaderAt
}

// FileRW represents the opened file for reading and writing.
type FileRW interface {
	FileRO
	io.Writer
}

type Path interface {
	// Path operations
	Absolute() (Path, error)
	Cwd() (Path, error)
	Parent() (Path, error)
	JoinPath(elem ...string) Path
	String() string

	// File and directory operations
	Touch() error
	Unlink() error
	RmDir() error
	MkDir(mode os.FileMode, parents bool) error
	Open() (FileRO, error)
	OpenRW(flag int, mode os.FileMode) (FileRW, error)
	Chmod(mode os.FileMode) error
	Rename(target Path) error
	Exists() bool

	// Path testing
	IsDir() bool
	IsFile() bool
	IsAbs() bool

	// File operations
	ReadBytes() ([]byte, error)
	ReadText() (string, error)
	WriteBytes(data []byte) error
	WriteText(text string) error
}
