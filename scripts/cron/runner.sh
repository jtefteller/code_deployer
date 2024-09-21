#!/bin/bash

function check_health {
	echo $(curl http://localhost:1337/health | jq -r .status)
}

STATUS=$(check_health)
echo "Server status: $STATUS"
if [ "$STATUS" != "ok" ]; then
	echo "Server is down, restarting..."
	$(pwd)/bin/code_deployer
fi
