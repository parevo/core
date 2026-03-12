package mfa

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type PquernaVerifier struct{}

func NewPquernaVerifier() *PquernaVerifier {
	return &PquernaVerifier{}
}

func (p *PquernaVerifier) Verify(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (p *PquernaVerifier) GenerateSecret() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Parevo",
		AccountName: "user",
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

func (p *PquernaVerifier) GenerateCode(secret string, t time.Time) (string, error) {
	return totp.GenerateCode(secret, t)
}
