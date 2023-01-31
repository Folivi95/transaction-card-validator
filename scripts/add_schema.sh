#!/bin/sh

apk --no-cache add curl;
curl -d '{"schema": "[{\"type\":\"record\",\"namespace\":\"namespace\",\"name\":\"name\",\"doc\":\"Somedata\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}]", "schemaType": "AVRO"}' -X POST -H "Content-Type: application/json" localhost:8081/subjects/test_transaction-v1/versions;