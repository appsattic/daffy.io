#!/bin/bash
## --------------------------------------------------------------------------------------------------------------------

set -e

echo "Checking ask.sh is installed ..."
if [ ! $(which ask.sh) ]; then
    echo "Please put ask.sh into ~/bin (should already be in your path from ~/.profile):"
    echo ""
    echo "    mkdir ~/bin"
    echo "    wget -O ~/bin/ask.sh https://gist.githubusercontent.com/chilts/6b547307a6717d53e14f7403d58849dd/raw/ecead4db87ad4e7674efac5ab0e7a04845be642c/ask.sh"
    echo "    chmod +x ~/bin/ask.sh"
    echo ""
    exit 2
fi
echo

# General
DAFFY_PORT=`ask.sh daffy DAFFY_PORT 'Which local port should the server listen on :'`
DAFFY_NAKED_DOMAIN=`ask.sh daffy DAFFY_NAKED_DOMAIN 'What is the naked domain (e.g. localhost:1234 or daffy.io) :'`
DAFFY_BASE_URL=`ask.sh daffy DAFFY_BASE_URL 'What is the base URL (e.g. http://localhost:1234 or https://daffy.io) :'`

# Social Providers
DAFFY_TWITTER_CONSUMER_KEY=`ask.sh daffy DAFFY_TWITTER_CONSUMER_KEY 'Enter your Twitter Consumer Key :'`
DAFFY_TWITTER_CONSUMER_SECRET=`ask.sh daffy DAFFY_TWITTER_CONSUMER_SECRET 'Enter your Twitter Consumer Secret :'`
DAFFY_GPLUS_CLIENT_ID=`ask.sh daffy DAFFY_GPLUS_CLIENT_ID 'Enter your GPlus Client Id :'`
DAFFY_GPLUS_CLIENT_SECRET=`ask.sh daffy DAFFY_GPLUS_CLIENT_SECRET 'Enter your GPlus Client Secret :'`
DAFFY_GITHUB_CLIENT_ID=`ask.sh daffy DAFFY_GITHUB_CLIENT_ID 'Enter your GitHub Client Id :'`
DAFFY_GITHUB_CLIENT_SECRET=`ask.sh daffy DAFFY_GITHUB_CLIENT_SECRET 'Enter your GitHub Client Secret :'`

 # Sessions
DAFFY_SESSION_AUTH_KEY_V2=`ask.sh daffy DAFFY_SESSION_AUTH_KEY_V2 'Enter your SESSION_AUTH_KEY_V2 :'`
DAFFY_SESSION_ENC_KEY_V2=`ask.sh daffy DAFFY_SESSION_ENC_KEY_V2 'Enter your SESSION_ENC_KEY_V2 :'`
DAFFY_SESSION_AUTH_KEY_V1=`ask.sh daffy DAFFY_SESSION_AUTH_KEY_V1 'Enter your SESSION_AUTH_KEY_V1 :'`
DAFFY_SESSION_ENC_KEY_V1=`ask.sh daffy DAFFY_SESSION_ENC_KEY_V1 'Enter your SESSION_ENC_KEY_V1 :'`

echo "Building code ..."
gb build
echo

# copy the supervisor script into place
echo "Copying supervisor config ..."
m4 \
    -D __DAFFY_PORT__=$DAFFY_PORT \
    -D __DAFFY_NAKED_DOMAIN__=$DAFFY_NAKED_DOMAIN \
    -D __DAFFY_BASE_URL__=$DAFFY_BASE_URL \
    -D __DAFFY_TWITTER_CONSUMER_KEY__=$DAFFY_TWITTER_CONSUMER_KEY \
    -D __DAFFY_TWITTER_CONSUMER_SECRET__=$DAFFY_TWITTER_CONSUMER_SECRET \
    -D __DAFFY_GPLUS_CLIENT_ID__=$DAFFY_GPLUS_CLIENT_ID \
    -D __DAFFY_GPLUS_CLIENT_SECRET__=$DAFFY_GPLUS_CLIENT_SECRET \
    -D __DAFFY_GITHUB_CLIENT_ID__=$DAFFY_GITHUB_CLIENT_ID \
    -D __DAFFY_GITHUB_CLIENT_SECRET__=$DAFFY_GITHUB_CLIENT_SECRET \
    -D __DAFFY_SESSION_AUTH_KEY_V2__=$DAFFY_SESSION_AUTH_KEY_V2 \
    -D __DAFFY_SESSION_ENC_KEY_V2__=$DAFFY_SESSION_ENC_KEY_V2 \
    -D __DAFFY_SESSION_AUTH_KEY_V1__=$DAFFY_SESSION_AUTH_KEY_V1 \
    -D __DAFFY_SESSION_ENC_KEY_V1__=$DAFFY_SESSION_ENC_KEY_V1 \
    etc/supervisor/conf.d/io-daffy.conf.m4 | sudo tee /etc/supervisor/conf.d/io-daffy.conf
echo

# restart supervisor
echo "Restarting supervisor ..."
sudo systemctl restart supervisor.service
echo

# copy the caddy conf
echo "Copying Caddy config config ..."
m4 \
    -D __DAFFY_PORT__=$DAFFY_PORT \
    -D __DAFFY_NAKED_DOMAIN__=$DAFFY_NAKED_DOMAIN \
    -D __DAFFY_BASE_URL__=$DAFFY_BASE_URL \
    -D __DAFFY_TWITTER_CONSUMER_KEY__=$DAFFY_TWITTER_CONSUMER_KEY \
    -D __DAFFY_TWITTER_CONSUMER_SECRET__=$DAFFY_TWITTER_CONSUMER_SECRET \
    -D __DAFFY_GPLUS_CLIENT_ID__=$DAFFY_GPLUS_CLIENT_ID \
    -D __DAFFY_GPLUS_CLIENT_SECRET__=$DAFFY_GPLUS_CLIENT_SECRET \
    -D __DAFFY_GITHUB_CLIENT_ID__=$DAFFY_GITHUB_CLIENT_ID \
    -D __DAFFY_GITHUB_CLIENT_SECRET__=$DAFFY_GITHUB_CLIENT_SECRET \
    -D __DAFFY_SESSION_AUTH_KEY_V2__=$DAFFY_SESSION_AUTH_KEY_V2 \
    -D __DAFFY_SESSION_ENC_KEY_V2__=$DAFFY_SESSION_ENC_KEY_V2 \
    -D __DAFFY_SESSION_AUTH_KEY_V1__=$DAFFY_SESSION_AUTH_KEY_V1 \
    -D __DAFFY_SESSION_ENC_KEY_V1__=$DAFFY_SESSION_ENC_KEY_V1 \
    etc/caddy/vhosts/io.daffy.conf.m4 | sudo tee /etc/caddy/vhosts/io.daffy.conf
echo

# restarting Caddy
echo "Restarting caddy ..."
sudo systemctl restart caddy.service
echo

## --------------------------------------------------------------------------------------------------------------------
