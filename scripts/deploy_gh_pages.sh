#!/bin/bash

source ./scripts/build_wasm.sh

GIT_REPO_URL=$(git config --get remote.origin.url)

cd ./tools/ast_explorer/static/
git init .
git remote add github $GIT_REPO_URL
git checkout -b gh-pages
git add .
git commit -am "gh-pages deploy"
git push github gh-pages --force
cd ../..
