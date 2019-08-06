#!/usr/bin/env bash
set -ex

bosh add-blob ../../uaa-crud.linux uaa-crud.linux

bosh create-release --name=uaa-crud --force --tarball=../uaa-crud-release.tgz

bosh upload-release ../uaa-crud-release.tgz

bosh upload-stemcell --sha1 712632e687388f335578956fceff27f0836646ae \
  https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-xenial-go_agent?v=456.3

bosh deploy -d uaa-crud ../manifests/sample-manifest.yaml -l ../manifests/values.yaml --recreate