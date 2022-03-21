#!/usr/bin/env bash

# This script sends a message to Slack via the webhook URL

set -Eeou pipefail

function main() {
    self_check "$@"
    post "$@"
}

function self_check() {
    set +e
    curl --version &> /dev/null
    if [[ $? -ne 0 ]]; then
        echo "curl is not installed"
        exit 1
    fi
    set -e

    set +u
    if [[ -z "${SLACK_WEBHOOK_URL}" ]]; then
        echo "SLACK_WEBHOOK_URL is not set"
        usage
        exit 1
    fi

    if [[ $# -eq 0 ]]; then
        echo "slack.sh needs at least one argument"
        usage
        exit 1
    fi
    set -u
}

function post() {
    set +eu
    read -d '' body << EOF
{
    "url":                    "${SLACK_WEBHOOK_URL}",
    "arg1":                   "${1}",
    "arg2":                   "${2}",
    "arg3":                   "${3}",
    "release_repo":           "${RELEASE_REPO}",
    "release_build_number":   "${RELEASE_BUILD_NUMBER}",
    "bitbucket_repo":         "${BITBUCKET_REPO_FULL_NAME}",
    "bitbucket_build_number": "${BITBUCKET_BUILD_NUMBER}"
}
EOF
    set -eu
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Secret-Header: ${SLACK_PASSWORD}" \
        -d "$body" \
        https://vcui5pmeoh.execute-api.eu-central-1.amazonaws.com/SlackNotifier
    echo
}

function usage() {
    echo "
    Slack Notifier

    Usage:

    Set Slack webhook
    $ export SLACK_WEBHOOK_URL=your_slack_webhook_url

    Simple info
    $ ./slack.sh 'Simple info message'

    Message with title
    $ ./slack.sh 'Title' 'Info message'

    Message with title and level (lvel = 0 -csuccess, levle != 0 - fail)
    $ ./slack.sh 'Deployment report' 'Deployment suceeded' 0

    Multiline message with level and title (separate lines with ',')
    $ ./slack.sh 'Deployment failed!' 'env: staging,commit: add configs,file: 31112020-123-v1.2.3,pipe: http://pipe.com/123' 1
    "
}

main "$@"
