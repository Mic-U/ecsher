[![Actions Status](https://github.com/Mic-U/ecsher/workflows/lint/badge.svg)](https://github.com/Mic-U/ecsher/actions)

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