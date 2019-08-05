# uaa-client-crud

### Background

Pivotal's ISV partners build tiles that are deployed to Ops Manager.
Often those partners products need to interact with the cloud foundry api.
Cloud foundry's api is secured by accounts in
User Account and Authentication (UAA)

UAA accounts have scopes: https://docs.cloudfoundry.org/concepts/architecture/uaa.html#cc-scopes.
Depending on what the partner product is doing, they may need only "cloud_controller.admin_read_only".

This is a cli that encapsulates the creation and deletion of a UAA client.

Once working, wrap in a bosh release and run it from a
post-deploy errand (https://bosh.io/docs/job-lifecycle/)

### Developing

Add this project to your gopath with: 
`go get -u github.com/cf-platform-eng/uaa-client-crud` 

Enable gomod integration in Goland

You may need to run `GO111MODULE=on go mod vendor` 