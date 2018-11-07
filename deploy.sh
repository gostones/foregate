#!/usr/bin/env bash

#args:

# DRY_RUN=echo

name=$1
[[ $# -lt 1 ]] && name="foregate"

webport=5080

##
function deploy() {
    local domain="run.aws-usw02-pr.ice.predix.io"

    echo "Pushing service $domain $name ..."

    #$DRY_RUN cf delete-route $domain --hostname $name -f

    $DRY_RUN cf push $name -f manifest.yml -d $domain --hostname $name --no-start; if [ $? -ne 0 ]; then
        return 1
    fi

    $DRY_RUN cf set-env $name FG_WEB http://localhost:${webport}
    $DRY_RUN cf set-env $name PG_VSCHEME https
    $DRY_RUN cf set-env $name PG_VHOST ${name}.${domain}

    $DRY_RUN cf start $name
}

#
echo "### Deploying ..."

deploy; if [ $? -ne 0 ]; then
    echo "#### Deploy failed"
    exit 1
fi

exit 0
##
