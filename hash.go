package main

import "golang.org/x/crypto/bcrypt"

func Hash(information string) (string, error) {
	str := []byte(information)
	harhStr, err := bcrypt.GenerateFromPassword(str, bcrypt.MinCost)
	return string(harhStr), err
}

func CompareHash(hashStr string, information string) bool {
	byteHash := []byte(hashStr)
	byteInformation := []byte(information)
	err := bcrypt.CompareHashAndPassword(byteHash, byteInformation)
	if err != nil {
		return false
	}
	return true
}
