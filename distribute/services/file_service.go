package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/distribute/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	Md5Error     = "md5 值错误"
	FileIsExists = "文件已存在"
)

type FileService struct {
}

func (f *FileService) Upload(ctx context.Context, request *UploadReq) (*Response, error) {
	if request.Subsection {
		return f.subsection(request)
	}
	return f.complete(request)
}

func (f *FileService) subsection(request *UploadReq) (*Response, error) {
	resp := &Response{}
	base := filepath.Base(request.Name)
	dir := filepath.Dir(request.Name)
	name := fmt.Sprintf("%s/.%s.tmp", dir, base)

	if utils.IsFile(request.Name) {
		if !request.Replace {
			resp.Msg = FileIsExists
			if utils.IsFile(name) {
				_ = os.Remove(name)
			}
			return resp, nil
		}
	}

	switch {
	case request.Start:
		if err := f.subsectionStart(name, request); err != nil {
			return nil, err
		}
		resp.Success = true
		return resp, nil
	case request.End:
		err := f.subsectionEnd(name, request)
		if err != nil && err.Error() != Md5Error {
			return nil, err
		} else if err != nil && err.Error() == Md5Error {
			resp.Msg = Md5Error
			return resp, nil
		}
		resp.Success = true
		return resp, nil
	default:
		fd, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		defer func() { _ = fd.Close() }()
		if _, err := fd.Write(request.Body); err != nil {
			return nil, err
		}
		resp.Success = true
		return resp, nil
	}
}

func (f *FileService) complete(request *UploadReq) (*Response, error) {
	resp := &Response{}
	if utils.IsFile(request.Name) {
		if !request.Replace {
			resp.Msg = FileIsExists
			return resp, nil
		}
	}

	fd, err := create(request.Name, request.Permission)
	if err != nil {
		return nil, err
	}

	defer func() { _ = fd.Close() }()

	if strings.ToUpper(utils.Md5sum(request.Body)) != strings.ToUpper(request.FileMd5Sum) {
		resp.Msg = Md5Error
		return resp, nil
	}

	if _, err := fd.Write(request.Body); err != nil {
		return nil, err
	}

	resp.Success = true
	return resp, nil
}

func (f *FileService) subsectionStart(name string, request *UploadReq) error {
	fd, err := create(name, request.Permission)
	if err != nil {
		return err
	}
	defer func() { _ = fd.Close() }()
	if _, err := fd.Write(request.Body); err != nil {
		return err
	}
	return nil
}

func (f *FileService) subsectionEnd(name string, request *UploadReq) error {
	fd, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer func() { _ = fd.Sync(); _ = fd.Close() }()

	if _, err := fd.Write(request.Body); err != nil {
		return err
	}

	if _, err := fd.Seek(0, 0); err != nil {
		return err
	}
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	if strings.ToUpper(utils.Md5sum(content)) != strings.ToUpper(request.FileMd5Sum) {
		return errors.New(Md5Error)
	}

	if utils.IsFile(request.Name) {
		_ = os.Remove(request.Name)
	}

	return os.Rename(name, request.Name)
}

func create(path string, perm uint32) (*os.File, error) {
	dir := filepath.Dir(path)
	if utils.IsDir(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(perm))
	if err != nil {
		return nil, err
	}

	// 解决如果文件已存在, OpenFile 清空文件的时候不会重置权限问题
	_ = os.Chmod(path, os.FileMode(perm))

	return fd, nil
}
