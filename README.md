# Metrics Sender

This is the client Cli tool for the [metrics-backend](https://github.com/gorlug/metrics-backend). See that project documentation for detailed information.

For usage of the Cli you can run the command with the `--help` flag.

```bash
metrics-sender --help
```

There is also a script to run it from the local directory:

```bash
./run_local.sh --help
```

## Config

To properly run this tool a config needs to be defined. The default location is `~/.metrics-sender.yaml` and should have this format:

```yaml
url: http://metrics-backend:8080/metric
journalLogMetaFile: ~/.journalLogMetaFile.log
journalUrl: http://metrics-backend:8080/journal
```

## Usage

This tool is to be run on a fixed interval using the native Linux cron jobs. Which could for example look like this:

```
* * * * * metrics-sender ping; metrics-sender dockerPing --exclude runner-; metrics-sender journal
10 * * * * metrics-sender disk --fileSystems /dev/sda1,/mnt/storage; metrics-sender disk --fileSystems storage --zpool
```

## Build

There are ready-made build scripts for amd64 and armv7:

```bash
./build_amd64.sh
./build_armv7.sh
```

The built binary goes into the `build` directory.

Refer to the golang documentation on how to build for other architectures.

## Notes

### Go local import

https://stackoverflow.com/a/73703527

```bash
go work init
go work use .
go work use ../metrics-backend
```

### Cobra Cli Parser

https://github.com/spf13/cobra/blob/main/site/content/user_guide.md

https://github.com/dharmeshkakadia/cobra-example/

https://opensource.com/article/21/1/go-cross-compiling

### Command for finding the last upgrade on Ubuntu

```bash
zgrep -B 1 "Commandline: apt upgrade" /var/log/apt/*
```
