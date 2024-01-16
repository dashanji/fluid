## Use vineyard to accelerate kubeflow pipelines in the fluid platform

Vineyard can accelerate data sharing by utilizing shared memory compared to existing methods such as local files or S3 services. In this doc, we will show you how to use vineyard to accelerate an existing kubeflow pipeline by using the fluid platform.

### Prerequisites

- Install the argo CLI tool via the [official guide](https://github.com/argoproj/argo-workflows/releases/).


### Overview of the pipeline

The pipeline we use is a simple pipeline that trains a linear regression model on the dummy Boston Housing Dataset. It contains three steps: preprocess, train, and test.


### Prepare the environment

- Prepare a kubernetes cluster. If you don't have a kubernetes cluster on hand, you can use the following command to create a kubernetes cluster via kind(v0.20.0+):

```bash
cat <<EOF | kind create cluster -v 5 --name=kind --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: kindest/node:v1.23.0
- role: worker
  image: kindest/node:v1.23.0
- role: worker
  image: kindest/node:v1.23.0
- role: worker
  image: kindest/node:v1.23.0
EOF
```

- Prepare the docker images. You can use the following command to build the docker image:

```bash
$ GO_MODULE=on make docker-build-all
```

Then you can push these images to the registry that your kubernetes cluster can access or load these images to the cluster if your kubernetes cluster is created by kind.

```bash
$ make -C ../../ docker-push-all
```

or

```bash
$ docker images | grep fluidcloudnative | grep $(git rev-parse --short HEAD) | awk '{print $1":"$2}' | xargs kind load docker-image
```

### Install the components

- Install the argo workflow controller, which can be used as the backend of kubeflow pipeline. You can use the following command to install the argo workflow controller:

```bash
$ kubectl create ns argo
$ kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.4.8/install.yaml
```

- Install the fluid platform. You can use the following command to install the fluid platform:

```bash
$ sed -i "s/719fc87/$(git rev-parse --short HEAD)/" charts/fluid/fluid/values.yaml && helm install fluid ./charts/fluid/fluid -n fluid-system --create-namespace
```

- Install the vineyard runtime and dataset. You can use the following command to install the vineyard runtime and dataset:

```bash
$ kubectl create ns kubeflow
$ cat <<EOF | kubectl apply -f -
apiVersion: data.fluid.io/v1alpha1
kind: VineyardRuntime
metadata:
  name: vineyard
  namespace: kubeflow
spec:
  master:
    replicas: 1
  worker:
    replicas: 2
  tieredstore:
    levels:
    - mediumtype: MEM
      quota: 20Gi
---
apiVersion: data.fluid.io/v1alpha1
kind: Dataset
metadata:
  name: vineyard
  namespace: kubeflow
EOF
```

### Run the pipeline

We can use the following steps to run the pipeline:


First, we need to build the docker image for the pipeline. You can use the following command to build the image:

```bash
$ make -C tmp/kubeflow-pipeline docker-build REGISTRY={your custom registry}
```

Then you need to push the image to the registry that your kubernetes cluster can access or load the image to the cluster if your kubernetes cluster is created by kind.

```bash
$ make -C tmp/kubeflow-pipeline push-images REGISTRY={your custom registry}
```

or

```bash
$ make -C tmp/kubeflow-pipeline load-images REGISTRY={your custom registry}
```

**Notice** You need to mount a **NAS path** to the kubernetes node.
Here we mount a NAS path to the `/mnt/csi-benchmark`(shown in the `prepare-data.yaml`) path of all kubernetes nodes.
Next, we need to prepare the dataset by running the following command:

```bash
$ kubectl apply -f tmp/kubeflow-pipeline/prepare-data.yaml
```

The dataset will be stored in the host path. Also, you may need to wait for a while for the dataset to be generated and you can use the following command to check the status:

```bash
$ while ! kubectl logs -l app=prepare-data -n kubeflow | grep "preparing data time" >/dev/null; do echo "dataset unready, waiting..."; sleep 5; done && echo "dataset ready"
```

Before running the pipeline, we need to create the role and rolebinding for the pipeline as follows.

```bash
$ kubectl apply -f tmp/kubeflow-pipeline/rbac.yaml
```

After that, you can run the pipeline via the following command:

Without vineyard:

```bash
$ argo submit --watch tmp/kubeflow-pipeline/pipeline.yaml -p data_multiplier=2000 -p registry="ghcr.io/v6d-io/v6d/kubeflow-example" -n kubeflow
```

With vineyard:

```bash
$ argo submit --watch tmp/kubeflow-pipeline/pipeline-with-vineyard.yaml -p data_multiplier=2000 -p registry="ghcr.io/v6d-io/v6d/kubeflow-example" -n kubeflow
```


### Modifications to use vineyard

Compared to the original kubeflow pipeline, we could use the following command to check the differences:

```bash
$ git diff --no-index --unified=40 pipeline.py pipeline-with-vineyard.py
```

The main modifications are:
- Add the vineyard persistent volume to the pipeline. This persistent volume is used to mount the vineyard client config file to the pipeline.

Also, you can check the modifications of the source code as 
follows.

- [Save data in the preparation step](https://github.com/v6d-io/v6d/blob/main/k8s/examples/vineyard-kubeflow/preprocess/preprocess.py#L62-L72).
- [Load data in the training step](https://github.com/v6d-io/v6d/blob/main/k8s/examples/vineyard-kubeflow/train/train.py#L15-L24).
- [load data in the testing step](https://github.com/v6d-io/v6d/blob/main/k8s/examples/vineyard-kubeflow/test/test.py#L14-L20).

The main modification is to use vineyard to load and save data
rather than using local files.
