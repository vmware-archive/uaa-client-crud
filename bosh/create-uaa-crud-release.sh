#!/usr/bin/env bash
set -ex

cd_executed=""
if [ -e "uaa-crud-release" ]; then
   pushd "uaa-crud-release"
   cd_executed="yes"
fi

bosh add-blob ../../uaa-crud.linux uaa-crud.linux

bosh create-release --name=uaa_crud --force --tarball=../uaa-crud-release.tgz

if [ "$cd_executed" = "yes" ]; then
   popd
fi

