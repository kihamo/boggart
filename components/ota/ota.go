package ota

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/kardianos/osext"
)

const (
	ArchitectureUnknown = "unknown"
)

type Release interface {
	Version() string
	BinFile() io.ReadCloser
	Checksum() []byte
	Size() int64
	Architecture() string
}

type Repository interface {
	Releases(arch string) ([]Release, error)
	ReleaseLatest(arch string) (Release, error)
}

type Updater struct {
}

func NewUpdater() *Updater {
	return &Updater{}
}

// очистка старых не используемых релизов при запуске
func (u *Updater) AutoClean() error {
	return nil
}

func (u *Updater) UpdateTo(release Release, path string) error {
	if release.Architecture() != runtime.GOARCH {
		return errors.New("not valid architecture")
	}

	// TODO: проверка подписи к файлу
	// TODO: проверка что текущий файл не является релизным

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New(path + "is directory, not executable")
	}

	//  раскрываем симлинк
	if stat.Mode()&os.ModeSymlink != 0 {
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
	}

	// 1. создаем файл path.new в него копируем новый релиз
	newPath := path + ".new"
	newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, release.BinFile())
	if err != nil {
		return err
	}

	// windows
	// newFile.Close()

	// 2. переименовываем текущий файл в path.old
	oldPath := path + ".old"
	_ = os.Remove(oldPath)

	err = os.Rename(path, oldPath)
	if err != nil {
		return err
	}

	// 3. копируем path.new в path
	err = os.Rename(newPath, path)
	if err != nil {
		// rollback
		return os.Rename(oldPath, path)
	}

	// 4. удалить старые релизы
	_ = os.Remove(oldPath)

	return err
}

func (u *Updater) Update(release Release) error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	return u.UpdateTo(release, execName)
}

func (u *Updater) Restart() error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	execDir := filepath.Dir(execName)

	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	_, err = os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	return err
}

type goArchReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

func GoArch(reader goArchReader) string {
	data := make([]byte, 16)
	if _, err := io.ReadFull(reader, data); err != nil {
		return ArchitectureUnknown
	}
	reader.Seek(0, 0)

	if bytes.HasPrefix(data, []byte("\x7FELF")) {
		if _elf, err := elf.NewFile(reader); err == nil {
			switch _elf.Machine {
			case elf.EM_386:
				return "386"
			case elf.EM_X86_64:
				return "amd64"
			case elf.EM_ARM:
				return "arm"
			case elf.EM_AARCH64:
				return "arm64"
			case elf.EM_PPC64:
				if _elf.ByteOrder == binary.LittleEndian {
					return "ppc64le"
				}
				return "ppc64"
			case elf.EM_S390:
				return "s390x"
			}
		}
	}

	if bytes.HasPrefix(data, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(data[1:], []byte("\xFA\xED\xFE")) {
		if _macho, err := macho.NewFile(reader); err == nil {
			switch _macho.Cpu {
			case macho.Cpu386:
				return "386"
			case macho.CpuAmd64:
				return "amd64"
			case macho.CpuArm:
				return "arm"
			case macho.CpuArm64:
				return "arm64"
			case macho.CpuPpc:
				return "ppc"
			case macho.CpuPpc64:
				return "ppc64"
			}
		}
	}

	return ArchitectureUnknown
}
