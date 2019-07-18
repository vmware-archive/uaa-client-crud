# uaa-client-crud

### Background

Pivotal's ISV partners build BOSH Releases and tiles that are deployed to Ops Manager.
Often those partners products need to interact with the Cloud Foundry API or CredHub.
Cloud Foundry's API and CredHub are secured by accounts in
User Account and Authentication (UAA)

UAA accounts have [scopes](https://docs.cloudfoundry.org/concepts/architecture/uaa.html#cc-scopes) that control 
what the account is capable of accessing.
Depending on what the partner product is doing, they may need only "cloud_controller.admin_read_only".
Additional scopes may be necessary for things like storing bind credentials in CredHub. 

This repo is a CLI and BOSH Release that can create and delete a UAA client for partner products to use.

### Developing

Clone the project. No need to put it into the GOPATH, as this is using go modules.

Run `make bootstrap` to acquire dependencies for the project. 

Install counterfeiter with
```bash
go get -u github.com/maxbrunsfeld/counterfeiter
```

Enable gomod integration in Goland (Preferences->Go->Go Modules).

You may need to run `GO111MODULE=on go mod vendor` if getting errors with `make`

### Using

#### Running the BOSH Release
1. Update the [sample manifest](/bosh/manifests/sample-manifest.yaml).
1. Update the [sample values](/bosh/manifests/values.yaml) with your product necessary scopes and identity to create.
1. Run `make` to build the `uaa-crud.linux` and `uaa-crud.darwin` binaries.
1. `cd bosh/uaa-crud-release`
1. Create the release and deploy: 
```
bosh add-blob ../../uaa-crud.linux uaa-crud.linux
bosh create-release --name=uaa_crud --force --tarball=../uaa-crud-release.tgz
bosh upload-release ../uaa-crud-release.tgz
bosh upload-stemcell --sha1 712632e687388f335578956fceff27f0836646ae \
  "https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-xenial-go_agent?v=456.3"
bosh deploy -d uaa_crud ../manifests/sample-manifest.yaml -l ../manifests/values.yaml --recreate
```
1. To delete the newly created UAA account run: `bosh run-errand uaa_delete -d uaa-crud`

#### CLI Only
1. Run `make` to build the `uaa-crud.linux` and `uaa-crud.darwin` binaries.

1. Create Client: 
```
./uaa-crud.darwin create --uaa-endpoint <uaa-url> --admin-identity <uaa-admin-username> \
--admin-pwd <uaa-admin-secret> --auth-grant-types <list, of, grant, types> --authorities <list, of, authorities> \
--scopes <list, of, uaa, scopes> --target-client-identity <identity-to-create> --target-client-pwd <secret-to-create> \
--token-validity <validity in seconds>
```

* Optional flags to create CredHub permissions: 
```
--credential-path <credhub-credential-path> --credhub-endpoint <credhub-url> \
--credhub-identity <credhub-admin-identity> --credhub-secret <credhub-admin-secret> --credhub-permissions <list, of, permissions>`

```
1. Delete Client:
```
./uaa-crud.darwin delete --uaa-endpoint <uaa-url> --admin-identity <uaa-admin-username> \
--admin-pwd <uaa-admin-secret> --target-client-identity <identity-to-delete>

```

* Optional flags to delete CredHub permissions: 
```
--credential-path <credhub-credential-path> --credhub-endpoint <credhub-url> \
--credhub-identity <credhub-admin-identity> --credhub-secret <credhub-admin-secret>`

```

#### Enhancements
- [ ] Don't always skip SSL validation
