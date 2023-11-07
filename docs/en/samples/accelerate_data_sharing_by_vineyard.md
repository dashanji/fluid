# DEMO - Accelerate data sharing by Vineyard

This demo introduces an example how vineyard accelerate the data sharing between python applications on the kubernetes cluster.

## Prerequisites
Before everything we are going to do, please refer to [Installation Guide](../userguide/install.md) to install Fluid on your Kubernetes Cluster, and make sure all the components used by Fluid are ready like this:

```shell
$ kubectl get pod -n fluid-system
NAME                                        READY   STATUS    RESTARTS      AGE
alluxioruntime-controller-7c54d9c76-vsrxg   1/1     Running   2 (17h ago)   18h
csi-nodeplugin-fluid-ggtjp                  2/2     Running   0             18h
csi-nodeplugin-fluid-krkbz                  2/2     Running   0             18h
dataset-controller-bdfbccd8c-8zds6          1/1     Running   0             18h
fluid-webhook-5984784577-m2xr4              1/1     Running   0             18h
fluidapp-controller-564dcd469-8dggv         1/1     Running   0             18h
```

## Configuration

**Create Vineyard Runtime and Dataset**

```yaml
$ cat <<EOF >vineyardRuntime.yaml
apiVersion: data.fluid.io/v1alpha1
kind: VineyardRuntime
metadata:
  name: vineyard
spec:
  replicas: 3
  fuse:
    global: true
    image: vineyard-mount-socket
---
apiVersion: data.fluid.io/v1alpha1
kind: Dataset
metadata:
  name: vineyard
spec:
  mounts:
  - name: vineyard
EOF
```

Create the VineyardRuntime and vineyard dataset resource object.

```shell
$ kubectl apply -f vineyardRuntime.yaml
```

Wait the vineyard dataset for ready as follows.

```shell
$ kubectl get dataset vineyard -n default
NAMESPACE   NAME       UFS TOTAL SIZE   CACHED   CACHE CAPACITY   CACHED PERCENTAGE   PHASE   AGE
default     vineyard                                                                  Bound   18m
```

** Apply the application Pod using the vineyard dataset **

```shell
$ cat << EOF > pod.yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    fuse.serverful.fluid.io/inject: "true"
    fluid.io/dataset.vineyard.sched: required
  name: application
spec:
  containers:
  - command:
    - /bin/sh
    - -c
    - |
      pip install vineyard numpy pandas;
      cat << EOF >> application.py
      import vineyard;
      import numpy as np;
      import pandas as pd;
      client = vineyard.connect("/data/vineyard.sock");
      # put a pandas dataframe to vineyard
      client.put(pd.DataFrame(np.random.randn(100, 4), columns=list('ABCD')), persist=True, name="test_dataframe");
      # put a basic data unit to vineyard
      client.put((1, 1.2345, 'xxxxabcd'), persist=True, name="test_basic_data_unit");
      client.close()
      EOF
      python application.py;
      sleep infinity;
    image: python:3.10
    imagePullPolicy: IfNotPresent
    name: application
    volumeMounts:
    - mountPath: /data
      name: vineyard-vol
  volumes:
  - name: vineyard-vol
    persistentVolumeClaim:
      claimName: vineyard
EOF
```

Create the application Pod.

```shell
$ kubectl apply -f pod.yaml
```

Now you can find the application Pod is running on the same node with vineyard instances and the vineyard socket is under the mountPath `/data`.`

