package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/shirou/gopsutil/v3/host"
	"io/ioutil"
	"net/http"
	"time"
)

type ResMsg struct {
	Success            bool   `json:"success"`
	OnlyId             string `json:"only_id"`
	Msg                string `json:"msg"`
	Country            string `json:"country"`
	Province           string `json:"province"`
	Locality           string `json:"locality"`
	Organization       string `json:"organization"`
	OrganizationalUnit string `json:"organizational_unit"`
	CommonName         string `json:"common_name"`
	TimeZone           string `json:"time_zone"`
	NtpServer          string `json:"ntp_server"`
}

type PemResp struct {
	Success bool   `json:"success"`
	CertPem []byte `json:"cert_pem"`
	Msg     string `json:"msg"`
	CAPem   []byte `json:"ca_pem"`
}

type Host struct {
	Hostname      string   `json:"hostname"`
	HostID        string   `json:"host_id"`
	KernelVersion string   `json:"kernel_version"`
	Release       string   `json:"release"`
	IPAddress     []string `json:"ip_address"`
	Uptime        uint64   `json:"uptime"`
}

type initialization struct {
}

func NewInitialization() *initialization {
	return &initialization{}
}

func (i *initialization) signCertificate() error {
	csr, err := ioutil.ReadFile(config.CertCsr)
	if err != nil {
		return err
	}

	request := make(map[string]string)
	request["csr"] = string(csr)
	request["only_id"] = config.OnlyId
	request["csr_md5sum"] = utils.Md5Sum(string(csr))
	data, err := json.Marshal(&request)
	if err != nil {
		return err
	}

	client := i.client(time.Hour * 24)
	resp, err := client.Post(config.CfgServer+config.ServerRequestCASign, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Request.Close }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	pem := &PemResp{}
	if err := json.Unmarshal(body, pem); err != nil {
		return nil
	}

	if !pem.Success {
		return fmt.Errorf("sign certificate error: %s", pem.Msg)
	}

	// 将服务器根据 csr 签发的 pem 文件写入到本地文件
	if err := ioutil.WriteFile(config.CertFile, pem.CertPem, 0644); err != nil {
		return err
	}

	// 将服务器端的 ca.pem 写入到本地文件
	if err := ioutil.WriteFile(config.CAPem, pem.CAPem, 0644); err != nil {
		return err
	}

	return nil
}

func (i *initialization) generateCertificate(initialization *ResMsg) error {
	// 1. 生成证书请求文件
	// 2. 将证书请求文件发送给服务器端
	// 3. 获取服务器端返回的证书文件
	// 4. 获取服务器端 ca.pem 文件
	// 5. 将证书和 ca.pem 写入到对应文件中
	cert := NewCertificate()
	if err := cert.Generate(initialization); err != nil {
		return err
	}
	return i.signCertificate()
}

func (i *initialization) client(timeout time.Duration) http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

func (i *initialization) agent() ([]byte, error) {
	infoStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	d := &Host{
		Hostname:      infoStat.Hostname,
		HostID:        infoStat.HostID,
		KernelVersion: infoStat.KernelVersion,
		Release:       infoStat.Platform + " " + infoStat.PlatformVersion,
		Uptime:        infoStat.Uptime,
		IPAddress:     []string{config.CfgPublicIP},
	}
	return json.Marshal(d)
}

func (i *initialization) Register() error {
	// 1. 向服务器端发送 agent 所在服务器的基本信息, 请求将 agent 添加到集群中
	// 2. 服务器端同意 agent 加入集群之后返回 OnlyId, 生成 csr 证书请求文件的配置, ntp server 地址(可选), 时区(可选)
	// 3. 将 OnlyId 写入到 ${WORKDIR}/.init 文件中
	// 4. 将 res 对象传递给 generateCertificate 函数进行证书相关的一系列操作
	buf, err := i.agent()
	if err != nil {
		return err
	}

	client := i.client(time.Hour * 24)
	resp, err := client.Post(config.CfgServer+config.ServerRequestRegister, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	res := &ResMsg{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("server side reject: %s", res.Msg)
	}

	config.OnlyId = res.OnlyId
	config.NtpServer = res.NtpServer
	config.TimeZone = res.TimeZone

	d := &initData{
		OnlyId:    res.OnlyId,
		NtpServer: res.NtpServer,
		TimeZone:  res.TimeZone,
	}

	if config.CfgSetZone && config.CfgSetTime {
		// 设置时区
		utils.SetZone()
		// 同步时间
		utils.SyncTime()
	}

	if err := i.generateCertificate(res); err != nil {
		return err
	}

	data, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.InitializationFile, data, 0644)
}
