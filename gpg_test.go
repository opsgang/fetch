package main

import (
	"fmt"
	"testing"
)

func TestGpgVerify(t *testing.T) {
	t.Parallel()

	fixDir := "test-fixtures/gpg"
	ascDir := fmt.Sprintf("%s/asc", fixDir)
	assetsDir := fmt.Sprintf("%s/assets", fixDir)
	usersDir := fmt.Sprintf("%s/users", fixDir)

	cases := []struct {
		user  string // used to locate the gpg public key
		asset string // file to verify
		asc   string // signature to use
		res   bool   // should succeed?
	}{
		{"test-1-do-not-trust", "test-1", "test-1-by-test-1-do-not-trust.asc", true},
		{"test-2-do-not-trust", "test-2", "test-2-by-test-2-do-not-trust.asc", true},
		{"test-2-do-not-trust", "test-1", "test-2-by-test-2-do-not-trust.asc", false}, // wrong sig
		{"test-1-do-not-trust", "test-2", "test-2-by-test-2-do-not-trust.asc", false}, // wrong key
		{"test-2-do-not-trust", "test-1", "test-2-by-test-2-do-not-trust.asc", false}, // invalid sig
		{"test-1-do-not-trust", "test-1", "malformed-signature.asc", false},           // invalid sig
		{"test-1-do-not-trust", "test-1", "malformed-signature2.asc", false},          // invalid sig
		{"test-1-do-not-trust", "test-1", "does-not-exist.asc", false},                // missing sig
		{"test-1-do-not-trust", "test-1", "empty-signature.asc", false},               // empty sig
	}

	for _, tc := range cases {
		gpgPubKey := fmt.Sprintf("%s/%s/gpg.pub", usersDir, tc.user)
		gpgSigPath := fmt.Sprintf("%s/%s", ascDir, tc.asc)
		assetPath := fmt.Sprintf("%s/%s", assetsDir, tc.asset)
		err := GpgVerify(gpgPubKey, gpgSigPath, assetPath)
		if tc.res && err != nil {
			msg := "Failed to successfully verify %s with key %s and sig %s: %s"
			t.Fatalf(msg, assetPath, gpgPubKey, gpgSigPath, err)
		}

		if !tc.res && err == nil {
			msg := "Unexpectedly verified %s with key %s and sig %s!"
			t.Fatalf(msg, assetPath, gpgPubKey, gpgSigPath)
		}
	}
}
