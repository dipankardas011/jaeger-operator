# Plan for the Jaeger v2 controller

Author: Dipankar Das

## Core problem
the existing repos manage jaeger-v1 deployment, which consists 
of multiple different binaries AND uses CLI flags / env vars for configuration. 
In Jaeger-v2 we are dealing with a single binary that is configurable to 
fit different deployment modes using config files. So the existing 
operators/chars won't work for v2, we need to modify them.


## Current Plan
Date: 2024-02-24

- [x] first need to think of what all the ports needed for it to be open!
- [x] what all arguments needs to be used and a demo deployment
- [ ] why are there controller's reconcile func() there are tracing enabled
- [ ] things to look for where is the manager.go and is there a single manager aka controller?
- [ ] understand the overall flow and working of jaeger
- [ ] explore the controllers already present for v1 and try to come up with the all-in-one controller
- [ ] talk with mentor on this
- [ ] once all looks good go for the implementation

## Key explanations
* build-in config will always run with in-memory storage, if you need a different storage you need to pass explicit config
* default configuration is in Jaeger repo `cmd/jaeger/internal/all-in-one.yaml`
* the webhooks are there inside the jaeger-operator to identify if any deploemynent has a annotation if yes it will inject the sidecar of jaeger otherwise it will not
* The Controller will requeue the Request to be processed again if the returned error is non-nil or Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
* interesting thing
```go

// ReconcileNamespace reconciles a Namespace object
type ReconcileNamespace struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client

	// this avoid the cache, which we need to bypass because the default client will attempt to place
	// a watch on Namespace at cluster scope, which isn't desirable to us...
	rClient client.Reader

	scheme *runtime.Scheme
}
```

* another interesting utils function
```go
func util.Truncate(format string, max int, values ...interface{}) string

Truncate will shorten the length of the instance name so that it contains at most max chars when combined with the fixed part If the fixed part is already bigger than the max, this function is noop.

// [`util.Truncate` on pkg.go.dev](https://pkg.go.dev/github.com/jaegertracing/jaeger-operator/pkg/util#Truncate)
```

* Docs on how to use the existing jaeger v1 operator [Refer](https://www.jaegertracing.io/docs/1.54/operator/)
* TODO: find the yaml template for the jaeger v2 configuration

## Tasks

### first how to run the basic `all-in-one` image

Refer: https://hub.docker.com/layers/jaegertracing/jaeger

```bash
docker run --rm -p 16686:16686 jaegertracing/jaeger:latest
```

### Ports needed to be open

Refer: https://github.com/jaegertracing/jaeger/blob/ef4791ec7b0761f7a6f9ac3cace54a2039483d8e/cmd/jaeger/Dockerfile#L9-L36

### Explore sub-commands
**components**: Outputs available components in this collector distribution
```bash
docker run --rm jaegertracing/jaeger:latest components
```

**docs**: Generates documentation
```bash
$ docker volume create demo                                                      
demo

$ docker run -v demo:/tmp:rw jaegertracing/jaeger:latest docs --dir="/tmp"
2024/02/24 08:11:00 application version: git-commit=ef4791ec7b0761f7a6f9ac3cace54a2039483d8e, git-version=v1.54.0, build-date=2024-02-22T14:58:48Z
2024/02/24 08:11:00 Generating documentation in /tmp

$ docker inspect demo                                                     
[
    {
        "CreatedAt": "2024-02-24T13:40:37+05:30",
        "Driver": "local",
        "Labels": null,
        "Mountpoint": "/var/lib/docker/volumes/demo/_data",
        "Name": "demo",
        "Options": null,
        "Scope": "local"
    }
]

$ ls -la /var/lib/docker/volumes/demo/_data
lsd: /var/lib/docker/volumes/demo/_data: Permission denied (os error 13).


$ sudo ls -la /var/lib/docker/volumes/demo/_data
total 40
drwxrwxrwt 1 root  root  406 Feb 24 13:41 .
drwx-----x 1 root  root   10 Feb 24 13:40 ..
-rw-r--r-- 1 10001 root  973 Feb 24 13:41 jaeger_completion_bash.md
-rw-r--r-- 1 10001 root  753 Feb 24 13:41 jaeger_completion_fish.md
-rw-r--r-- 1 10001 root  831 Feb 24 13:41 jaeger_completion.md
-rw-r--r-- 1 10001 root  720 Feb 24 13:41 jaeger_completion_powershell.md
-rw-r--r-- 1 10001 root 1013 Feb 24 13:41 jaeger_completion_zsh.md
-rw-r--r-- 1 10001 root  459 Feb 24 13:41 jaeger_components.md
-rw-r--r-- 1 10001 root  457 Feb 24 13:41 jaeger_docs.md
-rw-r--r-- 1 10001 root 1719 Feb 24 13:41 jaeger.md
-rw-r--r-- 1 10001 root 1356 Feb 24 13:41 jaeger_validate.md
-rw-r--r-- 1 10001 root  291 Feb 24 13:41 jaeger_version.md

```

**Debugging**: you can override the entrypoint by
```bash
docker run --rm -it --entrypoint /bin/sh jaegertracing/jaeger:latest
```
