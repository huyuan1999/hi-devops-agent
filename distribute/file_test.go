package distribute

import (
	"context"
	"github.com/huyuan1999/hi-devops-agent/distribute/services"
	"github.com/huyuan1999/hi-devops-agent/distribute/utils"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestUploadComplete(t *testing.T) {
	conn, err := grpc.Dial("192.168.3.10:8089", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = conn.Close() }()

	body, err := ioutil.ReadFile("distribute.go")
	if err != nil {
		t.Fatal(err)
	}

	upload := services.UploadReq{
		FileMd5Sum: utils.Md5sum(body),
		Name:       "/tmp/test_upload_complete_1.txt",
		Permission: 0644,
		Subsection: false,
		Start:      false,
		End:        false,
		Replace:    false,
		Body:       body,
	}

	client := services.NewDistributeClient(conn)
	resp, err := client.Upload(context.Background(), &upload)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("服务器返回的值: \n", resp.Success, resp.Msg)
}

// 大于 10M 的文件应该使用分段上传的方式
func TestUploadSubsection(t *testing.T) {
	conn, err := grpc.Dial("192.168.3.10:8089", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = conn.Close() }()

	fd, err := os.Open("test.log")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = fd.Close() }()

	// 每次读取 10MB 大小
	bufSize := 1024 * 1024 * 10
	buffer := make([]byte, bufSize)

	upload := services.UploadReq{
		FileMd5Sum: "c33646d5e93c75841261f719115036a6",
		Name:       "/tmp/test_1.log",
		Permission: 0644,
		Subsection: true,
		Start:      false,
		End:        false,
		Body:       nil,
		Replace:    true,
	}

	isStart := 0
	isEnd := false
	client := services.NewDistributeClient(conn)

	for {
		if isEnd {
			break
		}
		if isStart == 0 {
			upload.Start = true
			isStart++
		} else {
			upload.Start = false
		}

		bytesRead, err := fd.Read(buffer)
		if err != nil {
			if err == io.EOF {
				isEnd = true

				upload.End = true
				upload.Body = buffer[:bytesRead]
				resp, err := client.Upload(context.Background(), &upload)
				if err != nil {
					t.Fatal(err)
				}
				t.Log("文件上传完毕: ", resp.Success, resp.Msg)
				break
			}
			t.Fatal(err)
		}

		upload.Body = buffer[:bytesRead]
		resp, err := client.Upload(context.Background(), &upload)
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Success {
			t.Fatal(resp.Success, resp.Msg)
		}
		t.Log(resp.Success, resp.Msg)
	}
}
