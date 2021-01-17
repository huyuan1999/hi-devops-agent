package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"os"
)

func NewCertificate() *certificate {
	return &certificate{}
}

type certificate struct {
}

func (c *certificate) subject(subj *ResMsg) pkix.Name {
	return pkix.Name{
		Country:            []string{subj.Country},
		Locality:           []string{subj.Locality},
		Province:           []string{subj.Province},
		Organization:       []string{subj.Organization},
		OrganizationalUnit: []string{subj.OrganizationalUnit},
		CommonName:         subj.CommonName,
	}
}

func (c *certificate) template(subj *ResMsg) x509.CertificateRequest {
	subject := c.subject(subj)
	return x509.CertificateRequest{
		Subject: subject,
	}
}

func (c *certificate) Generate(subj *ResMsg) error {
	template := c.template(subj)
	// 创建秘钥文件
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// 创建 csr 文件
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return err
	}

	csr, err := os.Create(config.CertCsr)
	if err != nil {
		return err
	}

	defer func() { _ = csr.Close() }()

	if err := pem.Encode(csr, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
		return err
	}

	key, err := os.Create(config.KeyFile)
	if err != nil {
		return err
	}

	defer func() { _ = key.Close() }()

	if err := pem.Encode(key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}
	return nil
}
