## reproduce

## docker commands for setup

initialize the docker container for ssh tunnel
```bash
$ docker run -d -p 2222:22 arvindr226/alpine-ssh
```

```bash
$ docker run \
    -v $PWD/ssh_host_ed25519_key.pub:/home/docker/.ssh/keys/ssh_host_ed25519_key.pub:ro \
    -p 3333:22 -d atmoz/sftp \
  docker::::folder1,folder2
```


expected behaviour of the go run
```bash
$ ssh -vvv 2222 root:pass@localhost -L 2000:localhost:3333 -N
```

```bash
$ sftp -v -i ./ssh_host_ed25519_key -P 2000 docker@localhost

Connected to sftp
> 
```

## expected behaviour of running go
```bash
go run main.go
```

```go
log.Printf("fi %v", fi) # show that we have folder1 and folder2
```