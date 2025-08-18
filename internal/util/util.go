package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"europm/internal/bank/model"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Hashing(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	fmt.Println("get private key: ", path)
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func SignWithRSA(data string, priv *rsa.PrivateKey) (string, error) {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func VerifySignature(data, signatureBase64 string, pub *rsa.PublicKey) error {
	sig, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, sig)
}

func GetConfig(key string) string {
	return viper.Get(key).(string)
}

func LoadConfig(path string) (model.Config, error) {
	var cfg model.Config
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	return cfg, err
}

func FindAccount(customerAcc string, cfg model.Config) (model.Account, bool) {
	for _, acc := range cfg.Accounts {
		if acc.CustomerAcc == customerAcc {
			return acc, true
		}
	}
	return model.Account{}, false
}
