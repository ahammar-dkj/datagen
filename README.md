
# Datagen

Datagen is service that generates fake data and sends it to Amazon SQS. This service is intended 
to be used for tests or demos.

## Building

Build static binary:

```
$ make
```

## Running

```
./datagen -queue-url <queue-url>
```

## Configuration

| Name          | Description                         | Default   | Required |
|---------------|-------------------------------------|-----------|----------|
| queue-url     | AWS SQS queue URL                   |           | Yes      |
| region        | AWS region                          | us-east-1 | No       |
| endpoint      | AWS service endpoint                |           | No       |
| interval      | Interval (ms) between messages      | 1000      | No       |

## Local SQS Queue

[Localstack](https://github.com/localstack/localstack) can be used for local development without accessing the SQS web service. 

Install Localstack and awscli-local:

```
$ pip install localstack awscli-local
```

Start SQS service with Localstack and create an SQS queue:

```
$ SERVICES=sqs TMPDIR=/private$TMPDIR localstack start
$ awslocal sqs create-queue --queue-name my_queue
```

Run datagen against local SQS queue:
```
$ ./datagen -queue-url http://localhost:4576/queue/my_queue -endpoint http://localhost:4576
```

Poll SQS queue:

```
$ ./scripts/poll-sqs.sh http://localhost:4576/queue/my_queue
```
