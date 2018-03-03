package main

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"os"
)

func GpgVerify(gpgPubKey string, gpgSigPath string, assetPath string) error {
	pubKeyReader, err := os.Open(gpgPubKey)
	defer pubKeyReader.Close()
	if err != nil {
		return fmt.Errorf("Failed to open %s for reading: %s", gpgPubKey, err)
	}

	sigReader, err := os.Open(gpgSigPath)
	defer sigReader.Close()
	if err != nil {
		return fmt.Errorf("Failed to open %s for reading: %s", gpgSigPath, err)
	}

	assetReader, err := os.Open(assetPath)
	defer assetReader.Close()
	if err != nil {
		return fmt.Errorf("Failed to open %s for reading: %s", assetPath, err)
	}

	keyring, err := openpgp.ReadArmoredKeyRing(pubKeyReader)
	if err != nil {
		return fmt.Errorf("Failed to create armored key ring (pub key: %s): %s", gpgPubKey, err)
	}

	gpg, err := openpgp.CheckArmoredDetachedSignature(keyring, assetReader, sigReader)
	if err != nil {
		return fmt.Errorf("GPG verification failed (asset: %s, pub key: %s): %s", assetPath, gpgPubKey, err)
	}

	for id := range gpg.Identities {
		fmt.Printf("GPG ID for asset %s: %s\n", assetPath, id)
	}
	return nil
}
