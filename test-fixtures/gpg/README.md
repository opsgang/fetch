# CREATING FIXTURES

cd to this dir, and run make\_keys.sh

Should work whether you are in a container or not.

**Important**

> If you are already in a container, this gpg dir
> must be with in a volume mount, so the gpg
> container can _see and touch_ the contents.

## REF: Creating fresh master and subkey (sign)

```bash
gdir=$HOME/.gnupg
u=test-1-do-not-trust
rm -rf $gdir; mkdir -p $gdir ; chmod 0700 $gdir
cat <<EOF > $gdir/gpg.conf
personal-digest-preferences SHA256
cert-digest-algo SHA256
default-preference-list SHA512 SHA384 SHA256 SHA224 AES256 AES192 AES CAST5 ZLIB BZIP2 ZIP Uncompressed
EOF
chmod 0600 $gdir/gpg.conf

echo "
INFO: See https://www.gnupg.org/documentation/manuals/gnupg/Unattended-GPG-key-generation.html
INFO: * To set passphrase, remove %no-protection line, and add Passphrase: <string>
INFO: * Valid Key-Usage (or Subkey) (csv) encrypt,sign,auth
"

cat <<EOF > /var/tmp/user-input
%echo from /var/tmp/user-input
%no-protection
Key-Type: RSA
Key-Usage: sign
Key-Length: 4096
Subkey-Type: RSA
Subkey-Usage: sign
Subkey-Length: 4096
Expire-Date: 0
Name-Real: $u
Name-Email: $u@opsgang.io
EOF

gpg --full-generate-key --batch /var/tmp/user-input

gpg --export -a > $u-gpg.pub
subkey_id=$(gpg -K --with-keygrip | grep 'Keygrip =' | tail -n 1 | awk '{print $NF}')
gpg --export-secret-keys -a $subkey_id > $u-gpg-signing.key
```
