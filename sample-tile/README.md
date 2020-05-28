# UAA Client Crud Sample Tile

This sample-tile directory demonstrates the use of uaa-client-crud
in the context of a (Tanzu Ops Manager)[https://docs.pivotal.io/platform/2-9/installing/pcf-docs.html] tile. 

## Prerequisites
(Tile Generator)[https://docs.pivotal.io/tiledev/2-6/tile-generator.html] is used
to build the tile artifacts. 

## Building the tile
```bash
make build
```

## Installing the tile
1) Upload to OpsManager by using the "IMPORT A PRODUCT" button. The tile is
   a .pivotal file in the sample-tile/product directory of this repository. 
1) Stage and configure the tile. 
1) Apply Changes


## Verifying that the uaa account was created
1) Look up the UAA Admin Client Credentials either via the OpsMan UI or by api
   at /api/v0/deployed/products/<cf-deployment>/credentials/.uaa.admin_client_credentials
1) ssh into the OpsManager VM. 
1) Get a current token for uaac by using the credential you fetched in step 1 above. 
    ```bash
    uaac token client get admin -s <credential>
    ```
1) List the clients and look for "uaa-client-identity" (the target_client_identity property
   passed from the tile to the bosh uaa_create job)
    ```bash
    uaa clients
    ```
   
   You should see a client similar to: 
    ```
        uaa-client-identity
            scope: credhub.write openid oauth.approvals credhub.read
            resource_ids: none
            authorized_grant_types: refresh_token client_credentials
            autoapprove:
            access_token_validity: 10000
            authorities: credhub.write credhub.read oauth.login
            lastmodified: 1590680871000
        ```

    Upon deleting the tile, assuming the uaa_crud_delete errand was turned on, 
    this client will be removed. 

## Notes on scope and authorities
In this example, we've given the UAA client credhub scopes and authorities. 

(This document)[https://docs.cloudfoundry.org/concepts/architecture/uaa.html#scopes] is helpful 
when determining what scopes/authorities to provide for other purposes. For example, 
cloud_controller.admin_read_only, is good for reading the cloud controller api. 
