package jwthelper

import (
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

//func GenerateJWT(data domain.Account, verified bool) (tokenString string, err errs.Error) {
//	// expirationTime := time.Now().Add(1 * time.Hour)
//	claims := &JWTClaim{
//		Id: data.ID,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
//		},
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	completedString, error := token.SignedString(jwtKey)
//	if error != nil {
//		return "", errs.Wrap(err)
//	}
//	// tokenString, err = token.SignedString(jwtKey)
//	return completedString, nil
//}
//func ValidateToken(signedToken string) (message string, err error) {
//	token, err := jwt.ParseWithClaims(
//		signedToken,
//		&JWTClaim{},
//		func(token *jwt.Token) (interface{}, error) {
//			return []byte(jwtKey), nil
//		},
//	)
//
//	if err != nil {
//		err = errors.New("access token not valid")
//		return "access token not valid", err
//	}
//	claims, ok := token.Claims.(*JWTClaim)
//	if !ok {
//		err = errors.New("access token not valid")
//		return "access token not valid", err
//	}
//	if claims.ExpiresAt.Before(time.Now()) {
//		err = errors.New("token expired")
//		return "Token expired, please login again", err
//	}
//	return
//}
//
//func TokenRead(signedToken string) (data *JWTClaim, err error) {
//	token, err := jwt.ParseWithClaims(
//		signedToken,
//		&JWTClaim{},
//		func(token *jwt.Token) (interface{}, error) {
//			return []byte(jwtKey), nil
//		},
//	)
//	if err != nil {
//		err = errors.New("access token not valid")
//		return nil, err
//	}
//	claims, ok := token.Claims.(*JWTClaim)
//	if !ok {
//		err = errors.New("couldn't parse claims")
//		return nil, err
//	}
//	return claims, nil
//}
