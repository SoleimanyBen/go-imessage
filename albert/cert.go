package albert

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"howett.net/plist"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	B64FairplayKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUUMzQktyTFBJQmFiaHByKzRTdnVRSG5iRjBzc3FSSVE2Ny8xYlRmQXJWdVVGNnA5c2RjdjcwTityOHlGeGVzRG1wVG1LaXRMUDA2c3pLTkFPMWs1SlZrOS9QMWVqejA4Qk1lOWVBYjRqdUFoVldkZkFJeWFKN3NHRmplU0wwMTVtQXZyeFRGY09NMTBGL3FTbEFSQmljY3hIalBYdHVXVnIwZkxHcmhNKy9BTVFJREFRQUJBb0dBQ0dXM2JISFBOZGI5Y1Z6dC9wNFBmMDNTakoxNXVqTVkwWFk5d1VtL2gxczZyTE84Ky8xME1ETUVHTWxFZGNtSGlXUmt3T1ZpalJIeHpOUnhFQU1JODdBcnVvZmhqZGRiTlZMdDZwcFcybkxDSzdjRURRSkZhaFRXOUdRRnpwVlJRWFhmeHI0Y3MxWDNrdXRsQjZ1WTJWR2x0eFFGWXNqNWRqdjdEK0E3MkEwQ1FRRFpqMVJHZHhiZU9vNFh6eGZBNm40MkdwWmF2VGxNM1F6R0ZvQkpnQ3FxVnUxSlFPem9vQU1SVCtOUGZnb0U4K3VzSVZWQjRJbzBiQ1VUV0xwa0V5dFRBa0VBMTFyeklwR0loRmtQdE5jLzMzZnZCRmd3VWJzalRzMVY1RzZ6NWx5L1huRzlFTmZMYmxnRW9iTG1TbXozaXJ2QlJXQURpd1V4NXpZNkZOL0RtdGk1NndKQWRpU2Nha3VmY255dnp3UVo3UndwLzYxK2VyWUpHTkZ0YjJDbXQ4Tk82QU9laGNvcEhNWlFCQ1d5MWVjbS83dUovb1ozYXZmSmRXQkkzZkd2L2twZW13SkFHTVh5b0RCanB1M2oyNmJEUno2eHRTczc2N3IrVmN0VExTTDYrTzRFYWFYbDNQRW1DcngvVSthVGpVNDVyN0RuaThaK3dkaElKRlBkbkpjZEZrd0dId0pBUFErd1ZxUmpjNGgzSHd1OEk2bGxrOXdocEs5TzcwRkxvMUZNVmRheXRFbE15cXpRMi8wNWZNYjdGNnlhV2h1K1EyR0dYdmRsVVJpQTN0WTBDc2ZNMHc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="
	B64CertChain   = "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"
)

type CertRequestBody struct {
	ActivationInfoComplete bool   `plist:"ActivationInfoComplete"`
	ActivationInfoXML      []byte `plist:"ActivationInfoXML"`
	FairPlayCertChain      []byte `plist:"FairPlayCertChain"`
	FairPlaySignature      []byte `plist:"FairPlaySignature"`
}

type ActivationInformation struct {
	ActivationRandomness string `plist:"ActivationRandomness"`
	ActivationState      string `plist:"ActivationState"`
	BuildVersion         string `plist:"BuildVersion"`
	DeviceCertRequest    []byte `plist:"DeviceCertRequest"`
	DeviceClass          string `plist:"DeviceClass"`
	ProductType          string `plist:"ProductType"`
	ProductVersion       string `plist:"ProductVersion"`
	SerialNumber         string `plist:"SerialNumber"`
	UniqueDeviceID       string `plist:"UniqueDeviceID"`
}

func generateActivationInformation(csr []byte) ActivationInformation {
	return ActivationInformation{
		ActivationRandomness: uuid.NewString(),
		ActivationState:      "Unactivated",
		BuildVersion:         "10.6.4",
		DeviceCertRequest:    csr,
		DeviceClass:          "Windows",
		ProductType:          "windows1,1",
		ProductVersion:       "10.6.4",
		SerialNumber:         "WindowSerial",
		UniqueDeviceID:       uuid.NewString(),
	}
}

func GeneratePushCert() (*rsa.PrivateKey, *pem.Block, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	csr, err := generateCSR(key)
	if err != nil {
		return nil, nil, err
	}

	buf, err := plist.Marshal(generateActivationInformation(csr), plist.XMLFormat)
	if err != nil {
		return nil, nil, err
	}

	fairplayKey, err := getFairplayKey()
	if err != nil {
		return nil, nil, err
	}

	hashed := sha1.Sum(buf)
	signature, err := rsa.SignPKCS1v15(rand.Reader, fairplayKey, crypto.SHA1, hashed[:])
	if err != nil {
		return nil, nil, err
	}

	certChain, err := base64.StdEncoding.DecodeString(B64CertChain)
	if err != nil {
		return nil, nil, err
	}

	body := CertRequestBody{
		ActivationInfoComplete: true,
		ActivationInfoXML:      buf,
		FairPlayCertChain:      certChain,
		FairPlaySignature:      signature,
	}

	bodyBuf, err := plist.Marshal(&body, plist.XMLFormat)
	if err != nil {
		return nil, nil, err
	}

	query := fmt.Sprintf("activation-info=%s", url.QueryEscape(string(bodyBuf)))
	r := strings.NewReader(query)

	resp, err := http.Post("https://albert.apple.com/WebObjects/ALUnbrick.woa/wa/deviceActivation?device=Windows", "application/x-www-form-urlencoded", r)
	if err != nil {
		return nil, nil, err
	}

	resBuf := bytes.Buffer{}
	if _, err := io.Copy(&resBuf, resp.Body); err != nil {
		return nil, nil, err
	}

	pushCert, err := parsePushCerts(resBuf.Bytes())
	if err != nil {
		return nil, nil, err
	}

	return key, pushCert, nil
}

type pushCertResponse struct {
	DeviceActivation deviceActivation `plist:"device-activation"`
}

type deviceActivation struct {
	AckReceived      bool             `plist:"ack-received"`
	ActivationRecord activationRecord `plist:"activation-record"`
	ShowSettings     bool             `plist:"show-settings"`
}

type activationRecord struct {
	DeviceCertificate []byte `plist:"DeviceCertificate"`
}

var certQuery = regexp.MustCompile("<Protocol>(.*)</Protocol>")

func parsePushCerts(data []byte) (*pem.Block, error) {
	res := certQuery.FindSubmatch(data)
	if len(res) != 2 {
		return nil, errors.New("unable to find plist data within bytes passed")
	}

	var certRes pushCertResponse
	if _, err := plist.Unmarshal(res[1], &certRes); err != nil {
		return nil, err
	}

	block, _ := pem.Decode(certRes.DeviceActivation.ActivationRecord.DeviceCertificate)
	if block == nil {
		return nil, errors.New("unable to decode PEM block from data")
	}

	return block, nil
}

func generateCSR(key *rsa.PrivateKey) ([]byte, error) {
	req := x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Cupertino"},
			Organization:       []string{"Apple Inc."},
			OrganizationalUnit: []string{"iPhone"},
			CommonName:         uuid.New().String(),
		},
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &req, key)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getFairplayKey() (*rsa.PrivateKey, error) {
	fairplayKey, err := base64.StdEncoding.DecodeString(B64FairplayKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(fairplayKey)
	if block == nil {
		return nil, errors.New("unable to find PEM block")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
