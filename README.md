[![Actions Status](https://github.com/Mic-U/ecsher/workflows/lint/badge.svg)](https://github.com/Mic-U/ecsher/actions)
[![Actions Status](https://github.com/Mic-U/ecsher/workflows/release/badge.svg)](https://github.com/Mic-U/ecsher/actions)

# ECSHER

## What is Ecsher?

ecsher is CLI tool describing [ECS](https://aws.amazon.com/ecs/) resources like kubectl.

```
$ ecsher get svc
Cluster: MyCluster
NAME             STATUS  LAUNCH_TYPE     SCHEDULING_STRATEGY     DESIRED RUNNING PENDING
MyService        ACTIVE                  REPLICA                 1       1       0

$ ecsher get task 
Cluster: MyCluster
NAME                                    LAUNCH_TYPE     GROUP                   CONNECTIVITY    DESIRED_STATUS LAST_STATUS      HEALTH_STATUS
32b43c46cc25464c9cc90848b9a5142d        FARGATE         service:MyService       CONNECTED       RUNNING        RUNNING          UNKNOWN 
b9cb128c24554df78ed5a019aed6fabf        FARGATE         family:nginx-fargate    CONNECTED       RUNNING        RUNNING          UNKNOWN
```

## Usage

### Prerequisities

ecsher use your AWS credentials.
Please set up AWS credentials.

- [Configuration basics \- AWS Command Line Interface](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-precedence)


### Installation

1. You can install via Homebrew.

```
$ brew tap Mic-U/ecsher    
$ brew install ecsher
```

2. You can download binary from release page.

- https://github.com/Mic-U/ecsher/releases

3. You can build binary.  
Then, move binary to the directory in PATH(for example /usr/local/bin )

```
$ go build -o ecsher ./
$ sudo mv ./ecsher /usr/local/bin/ecsher
$ ecsher get cluster
NAME            STATUS  ACTIVE_SERVICES RUNNING_TASKS   PENDING_TASKS   CONTAINER_INSTANCES
MyCluster       ACTIVE  1               1               0               0
```


### Cluster

Show List of clusters.

```
$ ecsher get cluster
NAME          STATUS  ACTIVE_SERVICES RUNNING_TASKS   PENDING_TASKS   CONTAINER_INSTANCES
MyCluster     ACTIVE  1               1               0               0
```

Show detailed info of the cluster in YAML format.

```
$ ecsher describe cluster --name MyCluster
activeservicescount: 1
attachments: []
attachmentsstatus: null
capacityproviders:
- FARGATE_SPOT
- FARGATE
...
```

You can save cluster name by set command.

```
go run main.go set cluster --name MyCluster
Cluster: MyCluster
$ cat ~/.ecsher.toml
cluster = "MyCluster"
```

### Service

Note: If you saved cluster in config file, you don't need specify `--cluster` flag.  
Show List of services in the specified cluster. 
```
$ ecsher get service --cluster MyCluster
Cluster: MyCluster
NAME             STATUS  LAUNCH_TYPE     SCHEDULING_STRATEGY     DESIRED RUNNING PENDING
MyCluster        ACTIVE                  REPLICA                 1       1       0
```

Show detailed info of the service in YAML format.

```
$ go run main.go describe service --name MyService --cluster MyCluster
capacityproviderstrategy:
- capacityprovider: FARGATE_SPOT
  base: 0
  weight: 1
...
```

`svc` is the alias of `service`.

```
$ ecsher get svc
Cluster: MyCluster
NAME             STATUS  LAUNCH_TYPE     SCHEDULING_STRATEGY     DESIRED RUNNING PENDING
MyCluster        ACTIVE                  REPLICA                 1       1       0
```

### Tasks

Note: If you saved cluster in config file, you don't need specify `--cluster` flag. 
Show List of tasks in the specified cluster. 

```
$ ecsher get task --cluster MyCluster
Cluster: MyCluster
NAME                                    LAUNCH_TYPE     GROUP                    CONNECTIVITY    DESIRED_STATUS LAST_STATUS      HEALTH_STATUS
32b43c46cc25464c9cc90848b9a5142d        FARGATE         service:MyService        CONNECTED       RUNNING        RUNNING  UNKNOWN 
```

Show detailed info of the task in YAML format.

```
$ go run main.go describe task --name 32b43c46cc25464c9cc90848b9a5142d
attachments:
- details:
  - name: subnetId
...
```

### TaskDefinition

Note: By default, ecsher shows `ACTIVE` task definitions.  
Show List of task definition families. 

```
$ ecsher get definition
FAMILY
hello
nginx-fargate
```

If you specify `--family` flag, ecsher shows list of revisions in the specified family.

```
$ ecsher get definition --family nginx-fargate
FAMILY          REVISION
nginx-fargate   1
```

Show detailed info of the task definition in YAML format.

```
$ ecsher describe definition --family nginx-fargate --revision 1
compatibilities:
- EC2
- FARGATE
containerdefinitions:
...
```

### Container Instance

Note: If you saved cluster in config file, you don't need specify `--cluster` flag. 

Show list of container instances in the specified cluster.

```
$ ecsher get instance --cluster MyCluster
NAME                              EC2_INSTANCE_ID      STATUS  DOCKER_VERSION              AGENT_VERSION  CONNECTED  REMAINING_CPU  REMAINING_MEMORY  RUNNING  PENDING
243e5612f9fe43f5af887620f9cdedfa  i-0a91a75f0affe41d4  ACTIVE  DockerVersion: 19.03.13-ce  1.50.2         true       2048           1955              0        0
40dccfd0c25947a8b34f56d00cbee1a5  i-06c277d7fa62a5b68  ACTIVE  DockerVersion: 19.03.13-ce  1.50.2         true       2048           1955              0        0
```

Show detailed info of the container instance in YAML format.

```
ecsher describe instance --name 243e5612f9fe43f5af887620f9cdedfa -c MyCluster
agentconnected: true
agentupdatestatus: ""
attachments: []
attributes:
...
```

### Logs

Note: If you saved cluster in config file, you don't need specify `--cluster` flag. 
ecsher fetch latest 100 logs.

Show service logs.
If you set `--watch` option, ecsher continues to fetch logs until you interrupt.

```
$ ecsher logs service --name MyService --watch
TIMESTAMP                          ID                                    MESSAGE
2021-02-26 18:55:32.12 +0000 UTC   b35aff85-dd30-43dd-8510-d053544cfad0  (service MyService) has reached a steady state.
2021-02-27 00:55:53.502 +0000 UTC  1ef46502-d220-43d6-9a2f-d837650ea823  (service MyService) has reached a steady state.
```

Show task logs.
If you set `--watch` option, ecsher continues to fetch logs until you interrupt.
ecsher supports only awslogs logdriver.
If the task has multiple containers, you must set `--container` option.

```
$ ecsher logs task --name 8baae5aa33db434d9888171527434fb5 --container devops
TIMESTAMP                      MESSAGE
2021-03-19 11:30:33 +0000 UTC  Current latestPublished: 2021-03-19 11:30:31.551286116 +0000 UTC
2021-03-19 11:30:33 +0000 UTC  No new articles
```