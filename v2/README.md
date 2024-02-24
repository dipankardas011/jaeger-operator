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

- [ ] first need to think of what all the ports needed for it to be open!
- [ ] what all arguments needs to be open?
- [ ] explore the controllers already present for v1 and try to come up with the all-in-one controller
- [ ] talk with mentor on this
- [ ] once all looks good go for the implementation

## Key explanations
* build-in config will always run with in-memory storage, if you need a different storage you need to pass explicit config
* default configuration is in Jaeger repo `cmd/jaeger/internal/all-in-one.yaml`

## Tasks

### first how to run the basic `all-in-one` image

Refer: https://hub.docker.com/layers/jaegertracing/jaeger

```bash
docker run --rm -p 16686:16686 jaegertracing/jaeger:latest
```

### Ports needed to be open

Refer: https://github.com/jaegertracing/jaeger/blob/ef4791ec7b0761f7a6f9ac3cace54a2039483d8e/cmd/jaeger/Dockerfile#L9-L36

