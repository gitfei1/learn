package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

//大华RSA加密获取token
func main() {
	x := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCrb+5TI4x/uMnH+DdVFI4pmCh8OaIBLHEZCuiTDe16BewUHWi4CR5/sO1YX47aoi5+2gppy5cIURI7QR1g4ezC0pUTUEQVEmOZTiNVW+c/cpCru+USHdPwOzcrv1bJCNZ0IV37SywUBdgKf/9h/CeqLaFwK05qzOyw4TYA/48rzwIDAQAB"
	y := []byte(x)
	d := make([]byte, base64.StdEncoding.DecodedLen(len(y)))
	n, err := base64.StdEncoding.Decode(d, y)
	if err != nil {
		fmt.Println("encode err", err)
		return
	}
	d = d[:n]
	derPKix, err := x509.ParsePKIXPublicKey(d)
	if err != nil {
		fmt.Println("x509 err")
		return
	}
	pub := derPKix.(*rsa.PublicKey)
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte("zysj1234"))
	if err != nil {
		fmt.Println("jiami")
	}
	data := base64.StdEncoding.EncodeToString(result)
	fmt.Println("data:", data)
}
