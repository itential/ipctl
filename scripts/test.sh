# #!/bin/bash

# Copyright 2024 Itential Inc. All Rights Reserved

PKG_LIST=$(go list ./... | grep -v /vendor/)
TASK=""
RETURNCODE=0

coverage() {
    if [ -d cover ]; then
        rm -rf cover 
    fi

    mkdir cover

    echo "mode: cover" > cover/.coverage

    for package in ${PKG_LIST}; do
        go test -covermode=count -coverprofile "cover/${package##*/}.cov" "$package" ;
        rc=$?
        if [ $rc -ne 0 ]
        then
            RETURNCODE=$rc
        fi
    done

    tail -q -n +2 cover/*.cov >> cover/.coverage

    go tool cover -func=cover/.coverage
}

unittest()  {
    go fmt $(go list ./... | grep -v /vendor/)

    if ! go vet $(go list ./... | grep -v /vendor/); then
        exit 1
    fi

    go clean -testcache 

    for package in ${PKG_LIST}; do
        go test -v "$package";
        rc=$?
        if [ $rc -ne 0 ]
        then
            RETURNCODE=$rc
        fi
    done
}

debugtest()  {
    go fmt $(go list ./... | grep -v /vendor/)

    if ! go vet $(go list ./... | grep -v /vendor/); then
        exit 1
    fi

    go clean -testcache 

    for package in ${PKG_LIST}; do
        dlv test "$package";
        rc=$?
        if [ $rc -ne 0 ]
        then
            RETURNCODE=$rc
        fi
    done
}


help() {
    cat<<EOF
usage: test.sh [command]

Commands:
  unittest    - Run all unit tests locally
  debugtest   - Run all unit tests locally with debug
  coverage    - Run tests with coverage

EOF
    exit 0
}

case "$1" in
    unittest) TASK="unittest";;
    debugtest) TASK="debugtest";;
    coverage) TASK="coverage";;
esac

if [[ "$TASK" == "" ]]; then
    help 
else
    $TASK
fi

exit $RETURNCODE
