kinesis-tail
----

kinesis-tail tail tool from Amazon Kinesis Stream


# Quick Usage

```bash
$ go get -u github.com/evalphobia/kinesis-tail

$ kinesis-tail -h
Usage of kinesis-tail:
  -exclude string
        do not show the data if this word is contained
  -include string
        show the data only if this word is contained
  -interval int
        polling interval (sec) (default 1)
  -oldest
        tail from oldest record (use TRIM_HORIZON)
  -stream string
        Kinesis Stream name to tail


$ kinesis-tail -stream my-stream-name
```

## AWS user setting

set Environment variables below,

- AWS Access Key ID
    - `AWS_ACCESS_KEY_ID`
    - `AWS_ACCESS_KEY`
- AWS Secret Access Key
    - `AWS_SECRET_ACCESS_KEY`
    - `AWS_SECRET_KEY`
- Region
    - `AWS_REGION`
