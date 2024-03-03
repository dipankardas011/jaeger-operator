# Plan for the Jaeger v2 controller

Author: Dipankar Das

## Core problem
the existing repos manage jaeger-v1 deployment, which consists 
of multiple different binaries AND uses CLI flags / env vars for configuration. 
In Jaeger-v2 we are dealing with a single binary that is configurable to 
fit different deployment modes using config files. So the existing 
operators/chars won't work for v2, we need to modify them.


## Current Plan

- [x] first need to think of what all the ports needed for it to be open!
- [x] what all arguments needs to be used and a demo deployment
- [x] why are there controller's reconcile func() there are tracing enabled
    - because then we can enable tracing of the jeager operator itself for any traces and spans
- [x] understand the overall flow and working of jaeger
    - somewhat good progress so far
    - came to know about what all things are needed for the jaeger operator
- [ ] explore the controllers already present for v1 and try to come up with the all-in-one controller
    - `[Blocker]` need to figure out the configuration of v2 struct
- [ ] talk with mentor on this
- [ ] once all looks good go for the implementation / demo
    - created a small demo of the controller with the version v2alpha1
- [ ] When selected use this guide

## Key explanations
* build-in config will always run with in-memory storage, if you need a different storage you need to pass explicit config

* default configuration is in Jaeger repo `cmd/jaeger/internal/all-in-one.yaml`

* the webhooks are there inside the jaeger-operator to identify if any deploemynent has a annotation if yes it will inject the sidecar of jaeger otherwise it will not

* The Controller will requeue the Request to be processed again if the returned error is non-nil or Result.Requeue is true, otherwise upon completion it will remove the work from the queue.

* Interesting thing
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

* a sidecar is still a valid deployment strategy, but I would expect people to run OTEL collector as the sidecar, not configure Jaeger-v2 binary into that role (even though it is technically possible to do).

* v1 does not use config files, it's configured via CLI flags (although viper might also support a config file, we just never used or recommended it).
v2 does not use CLI flags at all, only yaml config.

* I think it will certainly be helpful to reuse jaeger configs as is, i.e. one could just run the binary manually with the same config as passed to the operator.  reusing Jaeger configurations as they are, enabling the binary to run manually with the same configuration as passed to the operator.

* I assume that's what otel-operator is doing too. But Jaeger operator does more things, if I am not mistaken, like preparing the storage. Which, frankly, I am not sure it should be doing. Orchestrating our own builds like es-rollover / es-cleaner - that I understand, but orchestracting Cassandra / Elasticsearch should not be in scope, people should use more official tools for that.

* Yuri assume that's what otel-operator is doing too. But Jaeger operator does more things, if I am not mistaken, like preparing the storage. Which, frankly, I am not sure it should be doing. Orchestrating our own builds like es-rollover / es-cleaner - that I understand, but orchestracting Cassandra / Elasticsearch should not be in scope, people should use more official tools for that.

* jaeger v2 binary has support for storage backend other than memory [`2024-03-01`]
    ```go
    type Config struct {
            Memory        map[string]memoryCfg.Configuration   `mapstructure:"memory"`
            Badger        map[string]badgerCfg.NamespaceConfig `mapstructure:"badger"`
            GRPC          map[string]grpcCfg.Configuration     `mapstructure:"grpc"`
            Elasticsearch map[string]esCfg.Configuration       `mapstructure:"elasticsearch"`
            // TODO add other storage types here
            // TODO how will this work with 3rd party storage implementations?
            //      Option: instead of looking for specific name, check interface.
    }
    ```
    > **Note**: it is refering the source code
    > Badger is a built-in single-host database. GRPC is an extensibility solution where the actual storage backend can be implemented as a remote GRPC service.

* help chart refers to upgrade of existing jaeger-helm-chart. brew we don't have today, it's for installing a binary on Macs (mostly for running as all-in-one, but configuration is left to the user)

* there is something which we need to do before any modification aka make sure all the data are stable check the jaeger operator src to better understand

* A good example for the sidecar thing can be found in this [Refer](https://github.com/jpkrohling/opentelemetry-collector-deployment-patterns/tree/main/pattern-3-kubernetes)

* Created a working demo on the jaeger by default configuration [Check there](./operator)

* webhooks in the v1 are used to detected any annotations so that using the mutating webhoook we can deploy the jaeger sidecar by refering to the closes deployment we check the inject thing for name or namespace.
    > **NOTE**: the name has higher priority than namespace

* Frontend and UI configurations `V1` [Refer](https://www.jaegertracing.io/docs/1.54/frontend-ui/#configuration)

* CLI flags `V1` [Refer](https://www.jaegertracing.io/docs/1.54/cli/#jaeger-all-in-one-prometheus)

* Compare the newly created default jaeger via operator and the manifest you cretaed and compare the changes

  from the jaeger manifest I created
  ```yaml

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name:  jaeger-all-in-one
    namespace: jaeger
    labels:
      app:  jaeger
  spec:
    selector:
      matchLabels:
        app: jaeger
    replicas: 1
    template:
      metadata:
        labels:
          app:  jaeger
      spec:
        containers:
        - name:  jaeger
          image:  jaegertracing/jaeger:latest
          args: ['--config', '/configs/memory.yaml']
          
          volumeMounts:
          - name: config
            mountPath: /configs

          ports:
          - containerPort: 5775
          - containerPort: 4317
          - containerPort: 4318
          - containerPort: 6831
          - containerPort: 6832
          - containerPort: 5778
          - containerPort: 16686
          - containerPort: 14268
          - containerPort: 9411
        volumes:
          - name: config
            configMap:
              name: jaeger-configuration
              items:
              - key: "memory.yaml"
                path: "memory.yaml"
  ```

  from the jaeger v1 operator deployed the simple deployment
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    annotations:
      deployment.kubernetes.io/revision: "1"
      linkerd.io/inject: disabled
      prometheus.io/port: "14269"
      prometheus.io/scrape: "true"
    generation: 1
    labels:
      app: jaeger
      app.kubernetes.io/component: all-in-one
      app.kubernetes.io/instance: simplest
      app.kubernetes.io/managed-by: jaeger-operator
      app.kubernetes.io/name: simplest
      app.kubernetes.io/part-of: jaeger
    name: simplest
    namespace: default
  spec:
    progressDeadlineSeconds: 600
    replicas: 1
    revisionHistoryLimit: 10
    selector:
      matchLabels:
        app: jaeger
        app.kubernetes.io/component: all-in-one
        app.kubernetes.io/instance: simplest
        app.kubernetes.io/managed-by: jaeger-operator
        app.kubernetes.io/name: simplest
        app.kubernetes.io/part-of: jaeger
    strategy:
      type: Recreate
    template:
      metadata:
        annotations:
          linkerd.io/inject: disabled
          prometheus.io/port: "14269"
          prometheus.io/scrape: "true"
          sidecar.istio.io/inject: "false"
        creationTimestamp: null
        labels:
          app: jaeger
          app.kubernetes.io/component: all-in-one
          app.kubernetes.io/instance: simplest
          app.kubernetes.io/managed-by: jaeger-operator
          app.kubernetes.io/name: simplest
          app.kubernetes.io/part-of: jaeger
      spec:
        containers:
        - args:
          - --sampling.strategies-file=/etc/jaeger/sampling/sampling.json
          env:
          - name: SPAN_STORAGE_TYPE
            value: memory
          - name: METRICS_STORAGE_TYPE
          - name: COLLECTOR_ZIPKIN_HOST_PORT
            value: :9411
          - name: JAEGER_DISABLED
            value: "false"
          - name: COLLECTOR_OTLP_ENABLED
            value: "true"
          image: jaegertracing/all-in-one:1.54.0
          imagePullPolicy: IfNotPresent
          livenessProbe:
            failureThreshold: 5
            httpGet:
              path: /
              port: 14269
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 15
            successThreshold: 1
            timeoutSeconds: 1
          name: jaeger
          ports:
          - containerPort: 5775
            name: zk-compact-trft
            protocol: UDP
          - containerPort: 5778
            name: config-rest
            protocol: TCP
          - containerPort: 6831
            name: jg-compact-trft
            protocol: UDP
          - containerPort: 6832
            name: jg-binary-trft
            protocol: UDP
          - containerPort: 9411
            name: zipkin
            protocol: TCP
          - containerPort: 14267
            name: c-tchan-trft
            protocol: TCP
          - containerPort: 14268
            name: c-binary-trft
            protocol: TCP
          - containerPort: 16685
            name: grpc-query
            protocol: TCP
          - containerPort: 16686
            name: query
            protocol: TCP
          - containerPort: 14269
            name: admin-http
            protocol: TCP
          - containerPort: 14250
            name: grpc
            protocol: TCP
          - containerPort: 4317
            name: grpc-otlp
            protocol: TCP
          - containerPort: 4318
            name: http-otlp
            protocol: TCP
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /
              port: 14269
              scheme: HTTP
            initialDelaySeconds: 1
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          volumeMounts:
          - mountPath: /etc/jaeger/sampling
            name: simplest-sampling-configuration-volume
            readOnly: true
        volumes:
        - configMap:
            defaultMode: 420
            items:
            - key: sampling
              path: sampling.json
            name: simplest-sampling-configuration
          name: simplest-sampling-configuration-volume
  ```

* `TODO: category`; find the yaml template for the jaeger v2 configuration so that we can think of the operator extracting it rather than it being
    ```go
    map[string]any
    ```

* `TODO: In-Progress`; try out the new configurations options for v2
    came to use the existing v2 binary for the controller

* `TODO: category`; how gonna we deploy for the sidecar thing?
    need to plan about the webhook stuff

* more detailed understanding for jaeger as in general and Opentelemtry

    1. Jaeger-query: is the how we get the traces aka UI is the best exampl
    2. Jaeger-collector: is the service which recieves the traces aka spans from the application. we specify the url and port where the collector is running it can be otel or jaeger specific
    3. Jaeger-exporter: is the service which stores all the spans and traces for persistance like memory, elasticsearch and cassandra
    4. jaeger-ingester: some high I/O based application can use messages queues like kafka

    and that the flow is like this
    ```
    application
        |
        V
    jaeger-collector
        |
        V
    jaeger-ingester (like streaming service aka kafka)
        |
        V
    jaeger-processor
        |
        V
    jaeger-exporter / storage-backend
    ```
    and the jaeger querier directly fetches the data from store-backend or jaeger exporter


* the elastic search demo
    ```bash
    cd jaeger-operator
    make es

    # then deploy the jaeger component
    # check the elasticisearch-v1.yaml file
    ```
    > **Note**: check the file named `elasticsearch-v1.yaml` which contains the resources created by the controller


* the casandra demo
    ```bash
    cd jaeger-operator
    make casendra

    # as as before
    # check the examples folder it contains the deployment
    ```
    > **Note**: check the file named `cassandra-v1.yaml` which contains the resources created by the controller


* `TODO: conform is required` why are these required can we drop them off?
    refering to the jaeger v1 operator
    ```
            env:
            - name: SPAN_STORAGE_TYPE
              value: memory
            - name: METRICS_STORAGE_TYPE
            - name: COLLECTOR_ZIPKIN_HOST_PORT
              value: :9411
            - name: JAEGER_DISABLED
              value: "false"
            - name: COLLECTOR_OTLP_ENABLED
              value: "true"
    ```

## Progress documentation

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


### Kubernetes manifest for the simple pod

Its inside [Manifest](./manifests/all-in-one.yaml)

### New Jaeger operator

check the `operator` folder

#### v2alpha1
* added the a simple demo on how to get the memory storage which is default config to working using the k8s operator
* it uses pods and services for deployment and it deploys the jaeger service per namespace
* make sure the name and the namespace are different b/w jaeger crd

