package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	mrand "math/rand"
	"time"
)

const keySize = 2048
const expiryDays = 28
const certPrefix = "UNLIMITED"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mrand.Intn(len(letterRunes))]
	}
	return string(b)
}

type caInfo struct {
	cacrt, cakey string
}

func NewCaInfo(cacrt string, cakey string) caInfo {
	return caInfo{cacrt, cakey}
}

// CertWriter main handler

func (ci *caInfo) CertWriter(out io.Writer) {
	catls, err := tls.LoadX509KeyPair(ci.cacrt, ci.cakey)

	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	subjectKeyID := make([]byte, 20)
	rand.Read(subjectKeyID)

	_ = randStringRunes(25)
	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: serialNumber,

		Subject: pkix.Name{
			//CommonName: certPrefix + randStringRunes(25),
			CommonName: certPrefix,
		},
		NotBefore: time.Now().AddDate(0, 0, -7),
		NotAfter:  time.Now().AddDate(0, 0, expiryDays),

		SubjectKeyId: subjectKeyID,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, keySize)
	pub := &priv.PublicKey

	// Sign the certificate
	certB, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)

	// Write the private Key
	pem.Encode(out, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Write the public key
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certB})
}
