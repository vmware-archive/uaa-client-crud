#!/bin/bash
# set -x

# Get the credentials for the UAA account we created with
# cloud_controller.admin_read_only authority. 
URL=https://uaa.sys.fillmore.cf-app.com
CF_ENDPOINT=https://api.sys.fillmore.cf-app.com
CF_CLIENT=$CF_CLIENT_CREDENTIALS_USERNAME
CF_PASS=$CF_CLIENT_CREDENTIALS_PASSWORD

# Get the token that will allow us to make calls to the CF api. 
TOKEN="$(curl -s ${URL}/oauth/token  \
  -u ${CF_CLIENT}:${CF_PASS} \
  -d 'grant_type=client_credentials' -k | jq -r '.access_token')"


curl -H "Authorization: bearer ${TOKEN}" ${CF_ENDPOINT}/v2/organizations -k
