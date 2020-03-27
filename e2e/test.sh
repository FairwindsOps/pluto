#!/bin/bash

set -e


printf "\n\n"
echo "***************************"
echo "** Install and Run Venom **"
echo "***************************"
printf "\n\n"

curl -LO https://github.com/ovh/venom/releases/download/v0.27.0/venom.linux-amd64
mv venom.linux-amd64 /usr/local/bin/venom
chmod +x /usr/local/bin/venom

cd /pluto/e2e
mkdir -p /tmp/test-results
venom run tests/* --log debug --output-dir=/tmp/test-results --strict
exit $?
