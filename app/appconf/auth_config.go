package appconf

import "os"

type AuthConfig struct {
	JwtSecretAccessToken []byte `validate:"required" name:"JWT_SECRET_ACCESS_TOKEN"`
}

func AuthConfigInit() *AuthConfig {
	return &AuthConfig{
		JwtSecretAccessToken: []byte(os.Getenv("JWT_SECRET_ACCESS_TOKEN")),
	}
}
