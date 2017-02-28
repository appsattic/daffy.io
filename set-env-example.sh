## --------------------------------------------------------------------------------------------------------------------
#
# These are the environment variables you need to pass to your server for configuration. All values are required except
# the social keys/secrets which are optional. The values here are fake, don't try and use them! :)
#
## --------------------------------------------------------------------------------------------------------------------

 # --- General ---
 export APP_NAME=Example
 export NAKED_DOMAIN=example.com
 export DAFFY_BASE_URL=https://example.com
 export DAFFY_PORT=8080
 export DAFFY_DB_DUMP_DIR=/var/lib/daffy/db

 # --- OAuth ---

 # Twitter:
 #
 # * See : https://apps.twitter.com/
 #
 export DAFFY_TWITTER_CONSUMER_KEY=
 export DAFFY_TWITTER_CONSUMER_SECRET=

 # Google (Note: we say Google, but Goth uses `gplus` instead):
 #
 # * See : https://developers.google.com/identity/sign-in/web/devconsole-project
 # * See : https://support.google.com/cloud/answer/6158849?hl=en
 #
 # Goth has a `gplus` provider and not a `google`, however, I wonder if that is actually deprecated and will stop
 # working eventually. Perhaps we don't care since it works the same, but just ends up a `gplus` prefix with each
 # social entity, instead of a `google` one. ¯\_(ツ)_/¯
 export DAFFY_GPLUS_CLIENT_ID=
 export DAFFY_GPLUS_CLIENT_SECRET=

 # GitHub:
 #
 # * See : https://github.com/settings/developers
 # * See : https://github.com/organizations/<your-organization>/settings/applications
 #
 export DAFFY_GITHUB_CLIENT_ID=
 export DAFFY_GITHUB_CLIENT_SECRET=

 # --- Sessions ---

 # generated with `pwgen -s 32 1`
 export DAFFY_SESSION_AUTH_KEY_V2=IppHcUPiJnwy8Hk5sv9m6qVFS3TP0gWq
 export DAFFY_SESSION_ENC_KEY_V2=PtwVsDVy168bBwRLSuV432s25E3ivEdy
 export DAFFY_SESSION_AUTH_KEY_V1=8svrXLgygeaG0nQ8GM6EAj3EzNMCfP6H
 export DAFFY_SESSION_ENC_KEY_V1=CXZnNBDa4LP82dA2iXkivS50EUvMfweA

## --------------------------------------------------------------------------------------------------------------------
