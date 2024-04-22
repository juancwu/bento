package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	// generate rsa key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating key pair: %v\n", err)
	}

	publicKey := &privateKey.PublicKey

	// save private key as PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	err = os.WriteFile("private_key.pem", privateKeyPEM, 0600)
	if err != nil {
		log.Fatalf("Error saving private key: %v\n", err)
	}

	// save public key as PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Error marshaling public key: %v\n", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})
	err = os.WriteFile("public_key.pem", publicKeyPEM, 0644)
	if err != nil {
		log.Fatalf("Error saving public key: %v\n", err)
	}

	// sign data
	message := []byte("hello world")
	hash := sha256.New()
	_, err = hash.Write(message)
	if err != nil {
		log.Fatalf("Error hashing message: %v\n", err)
	}
	hashedMsg := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedMsg)
	if err != nil {
		log.Fatalf("Error signing message: %v\n", err)
	}

	// verify signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedMsg, signature)
	if err != nil {
		log.Fatalf("Error verifying signature: %v\n", err)
	} else {
		fmt.Println("Signature verified.")
	}

	// encrypt data
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	if err != nil {
		log.Fatalf("Error encrypting message: %v\n", err)
	}

	// decrypt data
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		log.Fatalf("Error decrypting message: %v\n", err)
	}

	if string(plaintext) == string(message) {
		fmt.Println("Decrypted message matches original.")
	}
}
