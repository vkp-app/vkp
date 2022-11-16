#!/usr/bin/env bash

for dir in ./*/; do
    dir=${dir%*/}
    dir="${dir##*/}"
    cd "$dir" || exit
    echo "Building harbor.dcas.dev/kube-glass/addons/$dir:$(cat version.txt)"
    imgpkg push --file-exclusion ".git, version.txt" -i "harbor.dcas.dev/kube-glass/addons/$dir:$(cat version.txt)" -f .
    cd ../
done
