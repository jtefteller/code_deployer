#!/usr/bin/env bash

function check_health {
	echo $(curl http://localhost:1337/health | jq -r .status)
}

STATUS=$(check_health)
echo "Server status: $STATUS"
if [ "$STATUS" != "ok" ]; then
	echo "Server is down, restarting..."
	/root/code_deployer/bin/code_deployer
fi
