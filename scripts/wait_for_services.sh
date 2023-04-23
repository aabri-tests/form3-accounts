#!/bin/bash

accountapi_health_url="http://accountapi:8080/v1/health"

wait_for_service() {
    service_url=$1
    timeout_seconds=$2
    interval_seconds=$3
    max_retries=$(($timeout_seconds / $interval_seconds))

    retry_count=0
    while [ $retry_count -lt $max_retries ]
    do
        echo "Checking service health at $service_url"
        response=$(curl -s -o /dev/null -w "%{http_code}" $service_url)
        if [ "$response" == "200" ]; then
            echo "Service at $service_url is healthy"
            return 0
        fi

        echo "Service at $service_url is not healthy yet, retrying in $interval_seconds seconds..."
        retry_count=$(($retry_count + 1))
        sleep $interval_seconds
    done

    echo "Service at $service_url did not become healthy within $timeout_seconds seconds"
    return 1
}

wait_for_service $accountapi_health_url 60 5
