# Piper
A command-line queuing tool. It makes it easy to:
* Chain together stdout of one or many programs to stdin of one or many other programs
* Insert new items to stdin of an already running program
* Distribute a workload between multiple machines/workers

![piper-demo](https://github.com/AlfredBerg/piper/assets/18570335/c8b2e9d8-5058-4339-8309-0e40ce13dc11)


## Installation
```bash
go install github.com/AlfredBerg/piper@latest
```
A redis server is required. To install one follow https://redis.io/docs/install/install-redis/, https://hub.docker.com/_/redis or setup one with Redis Cloud.  
Specify what redis server to use with the environment variable `PIPER_REDIS_URL`, e.g. like
```
export PIPER_REDIS_URL='redis://localhost:6379/'
```
or create a configuration file (default location `$HOME/.config/.piper.yaml`) containing e.g.:
```
redis_url: "redis://localhost:6379/"
```

## Usage
```
$ piper insert -h
Insert items to one or more queues

Usage:
  piper insert [flags]

Flags:
  -h, --help            help for insert
  -i, --input string    Input file, if empty stdin
  -q, --queue strings   The queue to insert to, can be specified multiple times to insert to multiple queues

Global Flags:
      --config string   config file (default is $HOME/.piper.yaml)
```

```
$ piper stream -h
Read items from a queue

Usage:
  piper stream [flags]

Flags:
  -h, --help           help for stream
  -q, --queue string   Queue to read items from

Global Flags:
      --config string   config file (default is $HOME/.piper.yaml)

```

## Example
Read apex domains from the queue `apexes` and do subdomain discovery on them, and then find webservers and get the DNS data
```
tmux new-session -d -s 'sub-discovery' 'piper stream -q apexes | subfinder | piper insert -q subdomains'

#Fanout
tmux new-session -d -s 'subs-fanout' 'piper stream -q subdomains | anew subs | piper insert -q httpx -q dnsline' 

tmux new-session -d -s 'httpx' 'piper stream -q httpx | httpx -stream -no-color -o webservers'

tmux new-session -d -s 'dnsline' 'piper stream -q dnsline | dnsline | tee dnsdata'
```


## Other
Note that this is an `at-most-once delivery` queuing system, meaning there are no retries if a message fails to be processed.  

If a program reading from a queue crashes there can be some data loss as linux pipes are buffered. In linux the default pipe buffer is 
16 pages (4096 bytes * 16), but this program changes the pipebuffer of the directly connected pipe to the smallest possible, 1 page (4096 bytes). This buffer is lost if the input program at some point crashes deadlocks et.c.  
