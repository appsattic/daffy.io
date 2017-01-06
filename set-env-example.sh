## --------------------------------------------------------------------------------------------------------------------
#
# These are the environment variables you need to pass to your server for configuration. All values are required except
# the social keys/secrets which are optional. The values here are fake, don't try and use them! :)
#
## --------------------------------------------------------------------------------------------------------------------

 # --- General ---
 export DAFFY_PORT=8080
 export DAFFY_BASE_URL=http://localhost:8080

 # --- OAuth ---

 # Twitter:
 # * See - https://apps.twitter.com/
 export DAFFY_TWITTER_CONSUMER_KEY=
 export DAFFY_TWITTER_CONSUMER_SECRET=

 # GitHub:
 # * See - https://github.com/settings/developers
 # * See - https://github.com/organizations/<your-organization>/settings/applications
 export DAFFY_GITHUB_CLIENT_ID=
 export DAFFY_GITHUB_CLIENT_SECRET=

 # --- Sessions ---

 # generated with `pwgen -s 32 1`
 export DAFFY_SESSION_AUTH_KEY_V2=IppHcUPiJnwy8Hk5sv9m6qVFS3TP0gWq
 export DAFFY_SESSION_ENC_KEY_V2=PtwVsDVy168bBwRLSuV432s25E3ivEdy
 export DAFFY_SESSION_AUTH_KEY_V1=8svrXLgygeaG0nQ8GM6EAj3EzNMCfP6H
 export DAFFY_SESSION_ENC_KEY_V1=CXZnNBDa4LP82dA2iXkivS50EUvMfweA

## --------------------------------------------------------------------------------------------------------------------
