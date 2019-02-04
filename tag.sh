#!bin/sh

curl -L -k -H \"Content-Type: application/json\" -X POST -d '{ \"projectName\": \"kernel-concierge")\"'\", \"version\":\"'\"$(git describe --tags $(git rev-list --tags --max-count=1))\"'\" }' api-semver.gtw.dev.io.bb.com.br/project/version/canary
git push
git push --tags