# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.999; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

go build -ldflags "-X maahsome/tool-notes/cmd.semVer=${SEMVER} -X maahsome/tool-notes/cmd.buildDate=${BUILD_DATE} -X maahsome/tool-notes/cmd.gitCommit=${GIT_COMMIT} -X maahsome/tool-notes/cmd.gitRef=/refs/tags/${SEMVER}" && \
./tool-notes version | jq .
```

