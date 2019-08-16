#!/usr/bin/env bash
set -ex

cd_executed=""
if [ -e "uaa-crud-release" ]; then
   pushd "uaa-crud-release"
   cd_executed="yes"
fi

bosh upload-release ../uaa-crud-release.tgz

bosh upload-stemcell --sha1 712632e687388f335578956fceff27f0836646ae \
  https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-xenial-go_agent?v=456.3

bosh deploy -d uaa_crud ../manifests/sample-manifest.yaml -l ../manifests/values.yaml --recreate

if [ "$cd_executed" = "yes" ]; then
   popd
fi