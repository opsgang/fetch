#!/bin/bash
# vim: et sr sw=4 ts=4 smartindent:
#

docker_id() {
    if awk -F/ '$2 == "docker"' /proc/self/cgroup | read
    then
        cat /proc/self/cgroup | grep -Po '(?<=:memory:/docker/)([a-zA-Z0-9]+)$'
    else
        return 1
    fi
}

cat << 'EOF' >run_in_docker.sh
#!/bin/bash
gdir=$HOME/.gnupg
u=$GPG_USER
rm -rf $gdir; mkdir -p $gdir ; chmod 0700 $gdir
cat <<EOC > $gdir/gpg.conf
personal-digest-preferences SHA256
cert-digest-algo SHA256
default-preference-list SHA512 SHA384 SHA256 SHA224 AES256 AES192 AES CAST5 ZLIB BZIP2 ZIP Uncompressed
EOC
echo "INFO: Generating gpg config"
chmod 0600 $gdir/gpg.conf

echo '
INFO: See https://www.gnupg.org/documentation/manuals/gnupg/Unattended-GPG-key-generation.html
INFO: * To set passphrase, remove %no-protection line, and add Passphrase: <string>
INFO: * Valid Key-Usage (or Subkey) (csv) encrypt,sign,auth
'

echo "INFO: Generating pub/prv keys for $u"
cat <<EOC > /var/tmp/user-input
%echo from /var/tmp/user-input
%no-protection
Key-Type: RSA
Key-Usage: sign
Key-Length: 4096
Expire-Date: 0
Name-Real: $u
Name-Email: $u@opsgang.io
EOC

gpg --full-generate-key --batch /var/tmp/user-input

gpg --export -a > users/$u/gpg.pub
subkey_id=$(gpg -K --with-keygrip | grep 'Keygrip =' | tail -n 1 | awk '{print $NF}')
gpg --export-secret-keys -a > users/$u/gpg-signing.key

(
    cd asc
    for f in ../assets/test-* ; do
        echo "... generating asc signature for $f, by $u"
        gpg --armor --detach-sign --output $(basename $f)-by-$u.asc $f
    done
)

EOF

chmod a+x run_in_docker.sh


THIS_DIR=$(dirname $(realpath -- $0))
USERS="test-1-do-not-trust test-2-do-not-trust"

if id=$(docker_id)
then
    CWORKDIR=$THIS_DIR
    VOL_OPTS="--volumes-from $id"
else
    CWORKDIR=/$(basename $THIS_DIR)
    VOL_OPTS="-v $THIS_DIR:$CWORKDIR"
fi

RC=0
for u in $USERS; do
    mkdir -p $THIS_DIR/users/$u
    docker run -it --rm --name gpg-create-$(date '+%Y%m%d%H%M%S') \
        $VOL_OPTS  -w $CWORKDIR -e GPG_USER=$u \
        opsgang/devbox_aws:stable $CWORKDIR/run_in_docker.sh || RC=1
done

exit $RC
