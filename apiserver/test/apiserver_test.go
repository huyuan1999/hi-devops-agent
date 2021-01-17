package test

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huyuan1999/hi-devops-agent/apiserver/test/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"math/big"
	rand2 "math/rand"
	"net/http"
	"testing"
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

type CertRequest struct {
	Csr       string `json:"csr"`
	OnlyId    string `json:"only_id"`
	CsrMd5sum string `json:"csr_md5sum"`
}

func Md5Sum(s string) string {
	ret := md5.Sum([]byte(s))
	return hex.EncodeToString(ret[:])
}

func register(c *gin.Context) {
	msg := &ResMsg{
		Success:            true,
		OnlyId:             "8888",
		Msg:                "",
		Country:            "CN",
		Province:           "Beijing",
		Locality:           "Beijing",
		Organization:       "test",
		OrganizationalUnit: "test",
		CommonName:         "test.com",
		TimeZone:           "Asia/Shanghai",
		NtpServer:          "ntp1.aliyun.com",
	}

	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println("agent 请求加入集群, agent 信息如下: ", string(body))

	//c.JSON(http.StatusOK, gin.H{
	//	"success": false,
	//	"msg":     "你他娘的想都别想",
	//})

	c.JSON(http.StatusOK, msg)
}

func sign(c *gin.Context) {
	var p CertRequest
	_ = c.ShouldBind(&p)
	d, _ := json.Marshal(&p)
	log.Println("咦, 有人来请求证书了, 看看请求信息吧: ", string(d))

	if Md5Sum(p.Csr) != p.CsrMd5sum {
		log.Println("MD5 码验证失败, csr 无效")
		return
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"success": false,
	//	"msg":     "我就是不给你证书, 你有本事咬我啊",
	//})

	csrBlock, _ := pem.Decode([]byte(p.Csr))
	csrInfo, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "哦, 买噶的我不知道你的 csr 里面是啥, 就不给你签发证书了, 哈哈哈~~~: " + err.Error(),
		})
		return
	}

	log.Println("看看 csr 内容在决定给不给你颁发证书: ", csrInfo.Subject)

	serverPem, err := awardCert([]byte(p.Csr))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "哦, 买噶的颁发证书的时候意外的出错了: " + err.Error(),
		})
		return
	}

	caPem, err := ioutil.ReadFile("./tls/ca.pem")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "哦哦哦, 最重要的一步, 读取 ca 证书出错了: " + err.Error(),
		})
		return
	}

	data := PemResp{
		Success: true,
		CertPem: serverPem,
		Msg:     "",
		CAPem:   caPem,
	}

	c.JSON(http.StatusOK, data)
	log.Println("好吧, 给你证书. agent 加入集群这个流程你算是走完了")
}

func loadCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	// 加载 ca.pem
	caFile, err := ioutil.ReadFile("./tls/ca.pem")
	if err != nil {
		return nil, nil, err
	}
	caBlock, _ := pem.Decode(caFile)

	cert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// 加载 ca.key
	keyFile, err := ioutil.ReadFile("./tls/ca.key")
	if err != nil {
		return nil, nil, err
	}
	keyBlock, _ := pem.Decode(keyFile)
	praKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return cert, praKey, nil
}

// 通过 csr 生成证书
func awardCert(csr []byte) ([]byte, error) {
	certDERBlock, _ := pem.Decode(csr)
	if certDERBlock == nil {
		return nil, fmt.Errorf("嗯? 解析 csr 信息错误")
	}
	csrParse, err := x509.ParseCertificateRequest(certDERBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("嗯? 解析 csr 信息错误 2: %s", err.Error())
	}
	t := &x509.Certificate{
		SerialNumber:          big.NewInt(rand2.Int63()),
		Subject:               csrParse.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
	}

	rootCa, rootKey, err := loadCA()
	if err != nil {
		return nil, fmt.Errorf("加载 CA 错误: %s", err.Error())
	}

	server, err := x509.CreateCertificate(rand.Reader, t, rootCa, csrParse.PublicKey, rootKey)
	if err != nil {
		return nil, fmt.Errorf("签署证书发生错误: %s", err.Error())
	}

	serverKey := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: server,
	}

	return pem.EncodeToMemory(serverKey), nil
}

func TestGRPC(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("./tls/client.pem", "./tls/client.key")
	if err != nil {
		log.Println("LoadX509KeyPair", err.Error())
		return
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./tls/ca.pem")
	if err != nil {
		t.Fatal(err)
	}
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "test.com",
		RootCAs:      certPool,
	})
	conn, err := grpc.Dial("192.168.3.10:8088", grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	client := services.NewClientClient(conn)
	resp, err := client.GetNIC(context.Background(), &services.RequestClient{})
	if err != nil {
		log.Println("GRPC GetNIC Error: ", err.Error())
		return
	}
	data, err := json.Marshal(resp.Nic)
	if err != nil {
		log.Println("序列化 GRPC 调用结果错误: ", err.Error())
		return
	}
	t.Log("GRPC 调用结果: ", string(data))
}

func TestApiServer(t *testing.T) {
	r := gin.Default()
	r.POST("/api/v1/install/agent/register/", register)
	r.POST("/api/v1/tls/ca/sign/", sign)

	r.Run(":8888")
}
