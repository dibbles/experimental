#!/bin/bash

##### Version specs
# These defaults are known compatible versions
export TEKTON_VERSION="0.7.0"
export TEKTON_TRIGGERS_VERSION="0.1.0"

# To prevent Git Hub rate limiting when pulling images\
export GITHUB_TOKEN=

##### Dashboard specs
export DASHBOARD_INSTALL_NS="tekton-pipelines"

# Note that to receive webhooks, your github must be able to http POST to your Tekton installation. 
# Our initial testing has used Docker Desktop and GitHub Enterprise. 

# Set this to your github - used to create webhooks
export GITHUB_URL="https://github.ibm.com"

# This is the repo you want to set up a webhook for. See github.com/mnuttall/simple for a public copy of this repo. 
export GITHUB_REPO="https://github.ibm.com/MNUTTALL/simple" 