# The get-docker-image Tool
This is a tool that provides you a list of Docker Images used in a repository.

The tool needs an input of a URL that must have the repository link (.git) and chosen commit SHA hash.

`Note that the commit should not necessarily the latest commit, it can be any that was pushed at any point in time.`

The input should be provided as a URL pointing to a plaintext file. Each line will have two fields separated by a **_space_**:
- the https url of the github public repository
- the commit SHA

Example:
```bash
https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e
https://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404
```

### Used Libraries
The libraries involved in this project are:
* [go-git](https://github.com/go-git/go-git/) <br/>
  This library provides code instructions specifically for Go to do some actions on a Git repository in general (clone, commit, list tree, push etc...) <br/>
  The reason for it use here is, it provides support for in-memory as it's making the clone and accessing files inside the repository faster rather than relying on OS Filesystem directly. <br/>
* [go-git-billy](https://github.com/go-git/go-billy) <br/>
  This library provides Filesystem implementation based on memory. <br/>
  It allows interacting with files (I/O read and write) faster.

## Usage
To properly use the tool you need to provid a valid URL with the `url` flag, or an error will occur.
```bash
./get-docker-image -url https://.../file.txt
```

Example
```bash
./get-docker-image -url https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt
```

## Sample Output
After running the command, a sample output can be similar to the following
Note: it's recommended to have [jq] (https://stedolan.github.io/jq/) installed to have a better output

```javascript
{
  "Data": [
    {
      "id": 1,
      "repository": "https://github.com/app-sre/qontract-reconcile.git",
      "commit": "30af65af14a2dce962df923446afff24dd8f123e",
      "dockerfile": "dockerfiles/Dockerfile",
      "image": [
        "quay.io/app-sre/qontract-reconcile-builder:0.2.0",
        "quay.io/app-sre/qontract-reconcile-base:0.7.1"
      ]
    },
    {
      "id": 2,
      "repository": "https://github.com/app-sre/container-images.git",
      "commit": "c260deaf135fc0efaab365ea234a5b86b3ead404",
      "dockerfile": "qontract-reconcile-base/Dockerfile",
      "image": [
        "registry.access.redhat.com/ubi8/ubi:8.4",
        "registry.access.redhat.com/ubi8/ubi-minimal:8.4"
      ]
    }
  ]
}
```

## Docker
The tool is dockerized and can be run as a container.

Dockerfile is provided along the repo so you can build it and run it locally to test it out

It will simply get an alpine image with go already installed on it, it will copy the module files and download them and will place .go file.


```bash
docker build --tag get-docker-image .
```

then run it as simple as:
```bash
docker run -it get-docker-image -url [URL]
```


## Kubernetes
The Tool after being dockerized and since an image is created, we can defininetly use with Kubernetes.

As it depends on your OS, the chosen tool in my case to run Kubernetes locally is [Kind] (https://kind.sigs.k8s.io/) since it's an easy to use.

Creating a cluster is as easy as the following command
```bash
kind create cluster --name k8s-cluster
```

make sure you have [kubectl] (https://kubernetes.io/docs/tasks/tools/) installed already

The created image is hosted locally, so Kind Cluster need to be aware of that image, so we push it into the cluster
```bash
kind load docker-image get-docker-image:latest --name k8s-cluster
```


finally, testing everything out, for that purpose check the job.yaml file in the repository which contains a Kubernetes Job to run the tool.

```bash
kubectl create -f job.yaml
```

**Note:**
- the `imagePullPolicy: Never` should be set otherwise Kubernetes will look for the image outside on the Internet
- the `restartPolicy: Never` is required for a Kubernetes Job to work properly

check the job status, a successful job is indicated with a `completed` status
```bash
kubectl get jobs
```

check the pods, a successful pod is terminated with a `completed` status
```bash
kubectl get pods
```

check the logs of the listed pod
**Note:**
the pod reference in my case is docker-image-job-46sc2
```bash
kubectl logs docker-image-job-46sc2 | jq
```

if you see the JSON output then your Kubernetes Job is working as expected!