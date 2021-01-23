package initialization

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/huyuan1999/hi-devops-agent/apiserver/error_type"
	"github.com/pkg/errors"
	"os"
)

func (c *certificate) subject() pkix.Name {
	return pkix.Name{
		Country:            []string{c.Country},
		Locality:           []string{c.Locality},
		Province:           []string{c.Province},
		Organization:       []string{c.Organization},
		OrganizationalUnit: []string{c.OrganizationalUnit},
		CommonName:         c.CommonName,
	}
}

func (c *certificate) template() x509.CertificateRequest {
	subject := c.subject()
	return x509.CertificateRequest{
		Subject: subject,
	}
}

func (c *certificate) Generate(keyPath string) ([]byte, error) {
	template := c.template()
	// 生成秘钥信息
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Wrap(err, error_type.CertGenerateKeyError)
	}

	// 生成 csr 信息
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return nil, errors.Wrap(err, error_type.CertGenerateCsrError)
	}

	key, err := os.Create(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, error_type.OSError)
	}

	defer func() { _ = key.Close() }()

	if err := pem.Encode(key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, errors.Wrap(err, error_type.IOError)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}), nil
}
