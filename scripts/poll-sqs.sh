#!/bin/bash -e

queue_url=$1
if [ -z "queue_url" ]; then
  echo "usage: poll-sqs.sh <queue-url>"
  exit 1
fi

echo "polling queue: $queue_url"
while : ; do
  messages=$(awslocal sqs receive-message --queue-url $queue_url --visibility-timeout 1 --wait-time-seconds 20 --max-number-of-messages 10 | jq -c .Messages[])
  echo "received messages: $(echo $messages | jq -cr | wc -l | awk '{print $1}')"
  echo $messages | jq .

  for h in $(echo $messages | jq -r .ReceiptHandle); do
    awslocal sqs delete-message --queue-url $queue_url --receipt-handle $h
    echo "deleted message: $h"
  done

  sleep 5
done
