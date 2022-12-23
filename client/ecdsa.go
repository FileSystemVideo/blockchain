package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"strings"
)

type EciesKey struct {
	ecdsaPrivKey *ecdsa.PrivateKey
	ecdsaPubKey  *ecdsa.PublicKey
	privKey      *ecies.PrivateKey
	pubKey       *ecies.PublicKey
	seed         string
}

//ecies
func ImportEciesPublicKey(pubKey string) (*ecies.PublicKey, error) {
	ecdsaPubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	pukInterface, err := x509.ParsePKIXPublicKey(ecdsaPubKeyBytes)
	ecdsaPubKey := pukInterface.(*ecdsa.PublicKey)
	return ecies.ImportECDSAPublic(ecdsaPubKey), nil
}

//ecies
func ImportEciesPirvKey(privKey string) (*ecies.PrivateKey, error) {
	ecdsaPrivKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	ecdasPrivKey, err := x509.ParseECPrivateKey(ecdsaPrivKeyBytes)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSA(ecdasPrivKey), nil
}

//,ecies
func ImportEciesKey(seed string) (key *EciesKey, err error) {
	var ecdsaPriv *ecdsa.PrivateKey
	ecdsaPriv, err = ecdsa.GenerateKey(elliptic.P256(), strings.NewReader(seed))
	if err != nil {
		return key, errors.New("ecdsa.GenerateKey() failed:" + err.Error())
	}
	eciesPriv := ecies.ImportECDSA(ecdsaPriv)
	ecdsaPubKey := &eciesPriv.PublicKey
	key = &EciesKey{
		ecdsaPrivKey: ecdsaPriv,
		ecdsaPubKey:  ecdsaPubKey.ExportECDSA(),
		privKey:      eciesPriv,
		pubKey:       ecdsaPubKey,
		seed:         seed,
	}
	return key, nil
}


func (this *EciesKey) Export() (ecdsaPubkey string, ecdsaPrivkey string, seed string, err error) {
	ecdsaPubKeyBytes, err := x509.MarshalPKIXPublicKey(this.pubKey.ExportECDSA())
	if err != nil {
		return "", "", "", errors.New("x509.MarshalPKIXPublicKey() failed:" + err.Error())
	}

	ecdsaPubkey = base64.StdEncoding.EncodeToString(ecdsaPubKeyBytes)
	ecdsaPrivKeyBytes, err := x509.MarshalECPrivateKey(this.privKey.ExportECDSA())
	if err != nil {
		return "", "", "", errors.New("x509.MarshalECPrivateKey() failed:" + err.Error())
	}
	//fmt.Println("priv:::::::::",ecdsaPrivKeyBytes)
	//fmt.Println("priv:::::::::",len(ecdsaPrivKeyBytes))
	ecdsaPrivkey = base64.StdEncoding.EncodeToString(ecdsaPrivKeyBytes)
	return ecdsaPubkey, ecdsaPrivkey, this.seed, nil
}


func (this *EciesKey) CheckSign(digest, signTxt []byte) bool {
	return ecdsa.VerifyASN1(this.ecdsaPubKey, digest, signTxt)
}


func (this *EciesKey) Sign(digest []byte) (signTxt []byte, err error) {
	return ecdsa.SignASN1(rand.Reader, this.ecdsaPrivKey, digest)
}


func EciesDecrypt(privKey *ecies.PrivateKey, ciphertext []byte) (original []byte, err error) {
	return privKey.Decrypt(ciphertext, nil, nil)
}


func EciesEncrypt(pubKey *ecies.PublicKey, original []byte) (ciphertext []byte, err error) {
	return ecies.Encrypt(rand.Reader, pubKey, original, nil, nil)
}
