package initialization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/apiserver/error_type"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"strings"
	"time"
)

const (
	register = "/api/v1/install/agent/register/"
	caSign   = "/api/v1/tls/ca/sign/"
)

type descAgent struct {
	Hostname      string `json:"hostname"`
	HostID        string `json:"host_id"`
	KernelVersion string `json:"kernel_version"`
	Release       string `json:"release"`
	Uptime        uint64 `json:"uptime"`
	IPAddress     string `json:"ip_address"`
}

type certificate struct {
	Country            string `json:"country"`
	Province           string `json:"province"`
	Locality           string `json:"locality"`
	Organization       string `json:"organization"`
	OrganizationalUnit string `json:"organizational_unit"`
	CommonName         string `json:"common_name"`
}

type reqCert struct {
	Csr       string `json:"csr"`
	OnlyId    string `json:"only_id"`
	CsrMd5sum string `json:"csr_md5sum"`
}

type resCert struct {
	Success bool   `json:"success"`
	CertPem []byte `json:"cert_pem"`
	Msg     string `json:"msg"`
	CAPem   []byte `json:"ca_pem"`
}

type initData struct {
	OnlyId    string `json:"only_id"`
	TimeZone  string `json:"time_zone"`
	NtpServer string `json:"ntp_server"`
	Server    string `json:"server"`
}

type resMsg struct {
	Success     bool   `json:"success"`
	Msg         string `json:"msg"`
	InitData    *initData
	Certificate *certificate
}

type controller struct {
	context *cli.Context
}

func (c *controller) GetDescAgent() ([]byte, error) {
	infoStat, err := host.Info()
	if err != nil {
		return nil, errors.Wrap(err, error_type.OSError)
	}

	d := &descAgent{
		Hostname:      infoStat.Hostname,
		HostID:        infoStat.HostID,
		KernelVersion: infoStat.KernelVersion,
		Release:       infoStat.Platform + " " + infoStat.PlatformVersion,
		Uptime:        infoStat.Uptime,
		IPAddress:     c.context.String("public-ip"),
	}
	data, err := json.Marshal(d)
	if err != nil {
		return nil, errors.Wrap(err, error_type.SerializationError)
	}
	return data, nil
}

func (c *controller) SaveCertificate(cert []byte, ca []byte) error {
	// 将服务器根据 csr 签发的 pem 文件写入到本地文件
	if err := ioutil.WriteFile(c.context.String("cert-pem"), cert, 0644); err != nil {
		return errors.Wrap(err, error_type.IOError)
	}

	// 将服务器端的 ca.pem 写入到本地文件
	if err := ioutil.WriteFile(c.context.String("ca"), ca, 0644); err != nil {
		return errors.Wrap(err, error_type.IOError)
	}

	return nil
}

func (c *controller) WriteInitData(data *initData) error {
	if data.Server == "" {
		data.Server = c.context.String("server")
	}

	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return errors.Wrap(err, error_type.SerializationError)
	}
	if err := ioutil.WriteFile(".init.json", buf, 0644); err != nil {
		return errors.Wrap(err, error_type.IOError)
	}
	return nil
}

func (c *controller) client(uri string, timeout time.Duration, buf []byte) ([]byte, error) {
	client := utils.Client(timeout)

	server := strings.TrimRight(c.context.String("server"), "/")
	url := fmt.Sprintf("%s%s", server, uri)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, errors.Wrap(err, error_type.HttpRequestError)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, error_type.IOError)
	}
	return body, nil
}

func (c *controller) Certificate(cert *certificate, onlyId string) (*resCert, error) {
	var res resCert
	csr, err := cert.Generate(c.context.String("cert-key"))
	if err != nil {
		return nil, err
	}

	s := string(csr)
	req := &reqCert{
		Csr:       s,
		OnlyId:    onlyId,
		CsrMd5sum: utils.Md5Sum(s),
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, error_type.SerializationError)
	}

	body, err := c.client(caSign, time.Hour*24, data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.Wrap(err, error_type.Deserialization)
	}

	if !res.Success {
		logrus.Fatalln(res.Msg)
	}

	return &res, nil
}

func (c *controller) Register(context *cli.Context) error {
	var res resMsg
	c.context = context

	buf, err := c.GetDescAgent()
	if err != nil {
		return err
	}

	body, err := c.client(register, time.Hour*24, buf)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return errors.Wrap(err, error_type.Deserialization)
	}

	// 验证服务器是否允许此 agent 加入集群
	if !res.Success {
		logrus.Fatalln(res.Msg)
	}

	// 生成证书请求文件
	cert, err := c.Certificate(res.Certificate, res.InitData.OnlyId)
	if err != nil {
		return err
	}

	// 将证书信息写入本地文件系统
	if err := c.SaveCertificate(cert.CertPem, cert.CAPem); err != nil {
		return err
	}

	return c.WriteInitData(res.InitData)
}

func initialization(context *cli.Context) error {
	c := &controller{}
	return c.Register(context)
}
