#!/bin/bash
# set -x

# Get the UAA endpoint and the UAA admin client credentials
# from environment variables
# UAA_ENDPOINT is "uaa." prepended to the PAS system domain. 
# Get identity/password from .uaa.admin_client_credentials in the PAS tile. 
URL=$UAA_ENDPOINT
USER=$UAA_ADMIN_CLIENT_IDENTITY
PASS=$UAA_ADMIN_CLIENT_PASSWORD

# Get the credentials we want to use for the UAA account.
# This account is not yet created: this is what we are creating with this script. 
CF_CLIENT=$CF_CLIENT_CREDENTIALS_USERNAME
CF_PASS=$CF_CLIENT_CREDENTIALS_PASSWORD

# Get the token that will allow us to make calls to the UAA api. 
TOKEN="$(curl -s ${URL}/oauth/token  \
  -u ${USER}:${PASS} \
  -d 'grant_type=client_credentials' -k | jq -r '.access_token')"

# Exit if the new client exists in UAA. 
CLIENT=$(curl -s ${URL}/oauth/clients/${CF_CLIENT} -k -H "Authorization: Bearer ${TOKEN}")
if [ ! -z $CLIENT ]; then
    echo "Client ${CF_CLIENT} already exists"
    exit 0
fi

# Create the UAA client. 
# This persists a "client" with cloud_controller.admin_read_only authority
NEW_CLIENT=$(curl -s ${URL}/oauth/clients -k \
-H "Content-Type: application/json" -H "Accept: application/json" \
-H "Authorization: Bearer ${TOKEN}" \
-d "{\"client_id\":\"${CF_CLIENT}\", \
     \"client_secret\":\"${CF_PASS}\", \
     \"authorized_grant_types\":[\"refresh_token\",\"client_credentials\"], \
     \"authorities\":[\"cloud_controller.admin_read_only\"], \
     \"name\":\"${CF_CLIENT}\"}")

