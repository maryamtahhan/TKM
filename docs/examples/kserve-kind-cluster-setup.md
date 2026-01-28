# KServe Kind Cluster Setup Guide

Complete guide for creating and configuring a Kind (Kubernetes in Docker)
cluster for KServe development and testing.

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start (5 Minutes)](#quick-start-5-minutes)
4. [Step-by-Step Installation](#step-by-step-installation)
5. [Verification](#verification)
6. [Configuration Options](#configuration-options)
7. [Troubleshooting](#troubleshooting)
8. [Next Steps](#next-steps)
9. [Cleanup](#cleanup)
10. [Additional Resources](#additional-resources)

---

## Overview

### Key Components

| Component                    | Description                        |
|------------------------------|------------------------------------|
| **Kind Cluster**             | Local Kubernetes (2 nodes)         |
| **Container Registry**       | localhost:5001 for images          |
| **KServe Controllers**       | Inference service orchestration    |
| **cert-manager**             | TLS certificate management         |
| **ClusterServingRuntimes**   | Pre-configured model servers       |

### What You'll Get

A fully functional local KServe environment with:

- **2-node Kubernetes cluster** (1 control-plane + 1 worker)
- **Local container registry** at `localhost:5001`
- **KServe v0.13+** with all controllers deployed
- **11 ClusterServingRuntimes** ready for model serving
- **Local model cache support** for faster model loading
- **Support for both Docker and Podman**

---

## Prerequisites

### Hardware & Software

- **macOS**, **Linux**, or **Windows** (with WSL2)
- **RAM**: 8GB minimum, 16GB+ recommended
- **CPU**: 4+ cores recommended
- **Disk**: 20GB+ free space
- **Architecture**: ARM64 (Apple Silicon) or AMD64 (Intel/AMD)

### Required Tools

```bash
# Install via Homebrew (macOS - recommended)
brew install go kubectl helm

# Or install individually:
brew install go        # Go 1.21+
brew install kubectl   # Kubernetes CLI
brew install helm      # Helm 3.13+
```

**For Linux users:**

- Go: <https://go.dev/dl/>
- kubectl: <https://kubernetes.io/docs/tasks/tools/>
- Helm: <https://helm.sh/docs/intro/install/>

### Container Runtime

**Option 1: Docker Desktop** (recommended for beginners)

- macOS: <https://docs.docker.com/desktop/install/mac-install/>
- Linux: <https://docs.docker.com/desktop/install/linux-install/>
- Windows: <https://docs.docker.com/desktop/install/windows-install/>

**Option 2: Podman** (lightweight alternative)

- macOS: `brew install podman` then `podman machine init && podman
  machine start`
- Linux: <https://podman.io/docs/installation>

### Verify Installations

```bash
go version       # Should show go1.21+
kubectl version --client  # Should show client version
helm version     # Should show v3.13+
docker --version # Should show Docker/Podman version
```

---

## Quick Start (5 Minutes)

For experienced users who want to get running immediately:

```bash
# 1. Clone KServe repository
git clone https://github.com/kserve/kserve.git
cd kserve

# 2. Create Kind cluster with local registry
# For Podman: export KIND_EXPERIMENTAL_PROVIDER=podman
./hack/setup/dev/manage.kind-with-registry.sh

# 3. Set environment for your architecture
# For ARM64 (Apple Silicon):
export KO_DOCKER_REPO=localhost:5001
export KO_DEFAULTPLATFORMS=linux/arm64

# For AMD64 (Intel/AMD):
export KO_DOCKER_REPO=localhost:5001
export KO_DEFAULTPLATFORMS=linux/amd64

# 4. Deploy KServe
make deploy-dev

# 5. Verify deployment
kubectl get pods -n kserve
kubectl get clusterservingruntimes

# Done! KServe is ready for InferenceServices
```

**Total time**: 5-10 minutes depending on internet speed

---

## Step-by-Step Installation

### 1. Install Required Tools

#### Option A: Homebrew (macOS - Recommended)

```bash
brew install go kubectl helm
```

#### Option B: Manual Installation

**Go:**

```bash
# Download from https://go.dev/dl/
# Or use package manager
```

**kubectl:**

```bash
# macOS
brew install kubectl

# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s
https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/
```

**Helm:**

```bash
# macOS
brew install helm

# Linux
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
| bash
```

**Verify installations:**

```bash
go version       # Should show go1.21+
kubectl version --client
helm version
```

### 2. Clone KServe Repository

```bash
# Clone the official KServe repository
git clone https://github.com/kserve/kserve.git
cd kserve

# Optional: Check out a specific version
# git checkout v0.13.0
```

### 3. Create Kind Cluster with Local Registry

KServe provides an automated setup script:

```bash
# For Docker users (default):
./hack/setup/dev/manage.kind-with-registry.sh

# For Podman users:
KIND_EXPERIMENTAL_PROVIDER=podman ./hack/setup/dev/manage.kind-with-registry.sh
```

**What this script does:**

1. **Installs Kind** (if not already installed)
   - Downloads the latest Kind binary
   - Places it in your PATH

2. **Creates Kind cluster** with custom configuration
   - Cluster name: `kind`
   - Control plane node: `kind-control-plane`
   - Worker node: `kind-worker`
   - Exposes port 31000 for ingress

3. **Sets up local registry**
   - Container name: `kind-registry`
   - Port: `5001`
   - Connected to Kind network

4. **Configures containerd**
   - Registry config at `/etc/containerd/config.toml`
   - Allows pulling from `localhost:5001`

**Expected output:**

```text
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.33.1) 🖼
 ✓ Preparing nodes 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind
```

**Verify the cluster:**

```bash
kubectl get nodes

# Expected output:
# NAME                 STATUS   ROLES           AGE   VERSION
# kind-control-plane   Ready    control-plane   2m    v1.33.1
# kind-worker          Ready    <none>          2m    v1.33.1
```

### 4. Configure Environment Variables

Set environment variables for building and deploying KServe:

**For ARM64 (Apple Silicon M1/M2/M3/M4):**

```bash
export KO_DOCKER_REPO=localhost:5001
export KO_DEFAULTPLATFORMS=linux/arm64
```

**For AMD64 (Intel/AMD):**

```bash
export KO_DOCKER_REPO=localhost:5001
export KO_DEFAULTPLATFORMS=linux/amd64
```

**For Podman users, also set:**

```bash
export KIND_EXPERIMENTAL_PROVIDER=podman
```

**Make permanent** (add to `~/.bashrc` or `~/.zshrc`):

```bash
# ARM64 example
echo 'export KO_DOCKER_REPO=localhost:5001' >> ~/.zshrc
echo 'export KO_DEFAULTPLATFORMS=linux/arm64' >> ~/.zshrc
source ~/.zshrc
```

### 5. Deploy KServe

Deploy KServe in development mode:

```bash
make deploy-dev
```

**What this deploys:**

1. **CRDs (Custom Resource Definitions) - 11 total**
   - InferenceService
   - TrainedModel
   - InferenceGraph
   - ClusterServingRuntime
   - ServingRuntime
   - ClusterStorageContainer
   - LocalModelCache
   - LocalModelNodeGroup
   - LocalModelNode
   - LLMInferenceService
   - LLMInferenceServiceConfig

2. **cert-manager** (for TLS certificates)
   - cert-manager controller
   - cert-manager webhook
   - cert-manager cainjector

3. **KServe Controllers**
   - kserve-controller-manager
   - kserve-localmodel-controller-manager
   - llmisvc-controller-manager
   - kserve-localmodelnode-agent (DaemonSet)

4. **ClusterServingRuntimes** (11 pre-configured runtimes)
   - kserve-huggingfaceserver
   - kserve-huggingfaceserver-multinode
   - kserve-lgbserver (LightGBM)
   - kserve-mlserver (SKLearn/XGBoost)
   - kserve-paddleserver
   - kserve-pmmlserver
   - kserve-sklearnserver
   - kserve-tensorflow-serving
   - kserve-torchserve
   - kserve-tritonserver
   - kserve-xgbserver

**Deployment time**: 2-5 minutes

**Monitor deployment:**

```bash
# Watch all pods in kserve namespace
kubectl get pods -n kserve -w

# Wait for all pods to be Running
# Ctrl+C when done
```

**Expected pods:**

```text
NAME                                                    READY   STATUS    AGE
cert-manager-5d7f97b46d-xxxxx                          1/1     Running   2m
cert-manager-cainjector-69d885bf55-xxxxx               1/1     Running   2m
cert-manager-webhook-8d85f86b5-xxxxx                   1/1     Running   2m
kserve-controller-manager-xxxxx-xxxxx                  2/2     Running   1m
kserve-localmodel-controller-manager-xxxxx-xxxxx       1/1     Running   1m
llmisvc-controller-manager-xxxxx-xxxxx                 1/1     Running   1m
```

**Note:** The `kserve-localmodelnode-agent` DaemonSet will show 0/0 pods
initially. Pods will be created automatically when nodes are labeled with
`kserve/localmodel=worker`.

---

## Verification

### Check All Components

```bash
# Verify all KServe pods are running
kubectl get pods -n kserve

# Verify CRDs are installed
kubectl get crd | grep serving.kserve.io

# Expected CRDs (11 total):
# clusterservingruntimes.serving.kserve.io
# clusterstoragecontainers.serving.kserve.io
# inferencegraphs.serving.kserve.io
# inferenceservices.serving.kserve.io
# llminferenceserviceconfigs.serving.kserve.io
# llminferenceservices.serving.kserve.io
# localmodelcaches.serving.kserve.io
# localmodelnodegroups.serving.kserve.io
# localmodelnodes.serving.kserve.io
# servingruntimes.serving.kserve.io
# trainedmodels.serving.kserve.io

# Verify ClusterServingRuntimes
kubectl get clusterservingruntimes

# Expected runtimes (11 total):
# kserve-huggingfaceserver
# kserve-huggingfaceserver-multinode
# kserve-lgbserver
# kserve-mlserver
# kserve-paddleserver
# kserve-pmmlserver
# kserve-sklearnserver
# kserve-tensorflow-serving
# kserve-torchserve
# kserve-tritonserver
# kserve-xgbserver
```

### Test with Simple InferenceService

Create a test InferenceService to verify everything works:

```bash
# Create test namespace
kubectl create namespace kserve-test

# Deploy simple sklearn example
cat <<EOF | kubectl apply -f -
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: sklearn-iris
  namespace: kserve-test
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: gs://kfserving-examples/models/sklearn/1.0/model
EOF

# Wait for pod to be ready
kubectl wait --for=condition=ready pod \
  -l serving.kserve.io/inferenceservice=sklearn-iris \
  -n kserve-test --timeout=300s

# Verify InferenceService is ready
kubectl get inferenceservice sklearn-iris -n kserve-test

# Cleanup test
kubectl delete namespace kserve-test
```

---

## Configuration Options

### Cluster Configuration

The Kind cluster configuration is located at:

- Script: `hack/setup/dev/manage.kind-with-registry.sh`

**Default configuration:**

```yaml
# Embedded in the script
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 31000
        hostPort: 31000
  - role: worker
```

### Custom Cluster Name

To use a different cluster name:

```bash
# Edit the script to change CLUSTER_NAME variable
# Or create cluster manually:
kind create cluster --name my-kserve-cluster --config custom-config.yaml
```

### Multi-Worker Setup

For testing distributed workloads, create a custom configuration:

**Create `custom-config.yaml`:**

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 31000
        hostPort: 31000
  - role: worker
  - role: worker
  - role: worker
```

**Create cluster:**

```bash
kind create cluster --config custom-config.yaml
```

### GPU Support (Experimental)

Kind clusters don't have direct GPU access, but you can:

1. Use GPU operator for passthrough (complex setup)
2. Use simulated GPU nodes with labels
3. Deploy to cloud Kubernetes for real GPU testing

For local GPU testing, consider using:

- **Minikube** with GPU passthrough
- **K3s** with NVIDIA container runtime
- **Cloud providers** (GKE, EKS, AKS) for production workloads

---

## Troubleshooting

### Common Issues

#### Issue: Kind not found

**Symptom:**

```text
kind: command not found
```

**Solution:**

```bash
# The script should auto-install, but you can manually install:
# macOS
brew install kind

# Linux
curl -Lo ./kind
https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

#### Issue: Cluster already exists

**Symptom:**

```text
ERROR: failed to create cluster: node(s) already exist for a cluster
with the name "kind"
```

**Solution:**

```bash
# Delete existing cluster
kind delete cluster --name kind

# Or use the management script
./hack/setup/dev/manage.kind-with-registry.sh --uninstall

# Then recreate
./hack/setup/dev/manage.kind-with-registry.sh
```

#### Issue: Registry not accessible

**Symptom:**

```text
Failed to pull image "localhost:5001/image:tag": failed to pull and
unpack image
```

**Solution:**

```bash
# Check registry is running
docker ps | grep kind-registry
# Or for Podman: podman ps | grep kind-registry

# If not running, restart it
docker start kind-registry
# Or: podman start kind-registry

# Verify registry connectivity from nodes
docker exec kind-control-plane curl localhost:5001/v2/_catalog
```

#### Issue: Pods stuck in Pending

**Symptom:**

```text
Pods remain in Pending state
```

**Solution:**

```bash
# Check node resources
kubectl describe nodes

# Check pod events
kubectl describe pod <pod-name> -n kserve

# Common causes:
# - Insufficient CPU/memory: Increase container runtime resources
#   (Docker Desktop settings or Podman machine config)
# - Image pull errors: Check registry connectivity
# - Taints: Verify worker node has no blocking taints
```

#### Issue: Wrong architecture images

**Symptom:**

```text
exec format error
```

**Solution:**

```bash
# Verify you set the correct platform
echo $KO_DEFAULTPLATFORMS

# For Apple Silicon (M1/M2/M3/M4):
export KO_DEFAULTPLATFORMS=linux/arm64

# For Intel/AMD:
export KO_DEFAULTPLATFORMS=linux/amd64

# Redeploy KServe
make undeploy-dev
make deploy-dev
```

#### Issue: kubectl context wrong

**Symptom:**

```text
The connection to the server localhost:8080 was refused
```

**Solution:**

```bash
# Set correct context
kubectl config use-context kind-kind

# Verify
kubectl config current-context
# Should show: kind-kind
```

### Debug Commands

```bash
# View cluster info
kubectl cluster-info --context kind-kind

# Check node status
kubectl get nodes -o wide

# Check all pods across namespaces
kubectl get pods --all-namespaces

# View controller logs
kubectl logs -n kserve deployment/kserve-controller-manager -c manager

# View local model controller logs
kubectl logs -n kserve deployment/kserve-localmodel-controller-manager

# Check Kind cluster list
kind get clusters

# View container runtime containers
docker ps -a | grep kind
# Or: podman ps -a | grep kind

# Check registry contents
curl http://localhost:5001/v2/_catalog
```

---

## Next Steps

After setting up your KServe Kind cluster, you can:

1. **Deploy InferenceServices**
   - Try the [sklearn
     example](https://kserve.github.io/website/latest/get_started/first_isvc/)
   - Test with your own models
   - Explore different ClusterServingRuntimes

2. **Configure Local Model Cache**
   - Cache models for faster inference startup
   - Reduce network bandwidth usage
   - Test with large language models

3. **Explore ServingRuntimes**
   - View available runtimes: `kubectl get clusterservingruntimes`
   - Create custom runtimes for your frameworks
   - Configure runtime-specific parameters

4. **Test Autoscaling**
   - Configure HPA (Horizontal Pod Autoscaler)
   - Test scale-to-zero functionality
   - Monitor resource utilization

5. **Production Deployment**
   - Move to cloud Kubernetes (GKE, EKS, AKS) <!-- codespell:ignore -->
   - Configure ingress and TLS
   - Set up monitoring and logging
   - Implement CI/CD pipelines

---

## Cleanup

### Delete KServe Only

Keep the cluster but remove KServe:

```bash
make undeploy-dev
```

### Delete Entire Cluster

Remove everything including Kind cluster and registry:

```bash
# Using the management script (recommended)
./hack/setup/dev/manage.kind-with-registry.sh --uninstall
```

**Or manually:**

```bash
# Delete Kind cluster
kind delete cluster --name kind

# Stop and remove registry
docker stop kind-registry
docker rm kind-registry
# Or for Podman:
# podman stop kind-registry
# podman rm kind-registry

# Verify cleanup
kind get clusters  # Should be empty
docker ps | grep kind  # Should show nothing
```

### Recreate Fresh Cluster

```bash
# Full cleanup and recreate
./hack/setup/dev/manage.kind-with-registry.sh --uninstall
./hack/setup/dev/manage.kind-with-registry.sh

# Redeploy KServe
export KO_DOCKER_REPO=localhost:5001
export KO_DEFAULTPLATFORMS=linux/arm64  # or linux/amd64
make deploy-dev
```

---

## Additional Resources

### Official Documentation

- [KServe Documentation](https://kserve.github.io/website/)
- [Kind Documentation](https://kind.sigs.k8s.io/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

### Example Repositories

- [KServe
  Examples](https://github.com/kserve/kserve/tree/master/docs/samples)
- [Model Serving
  Patterns](https://github.com/kserve/website/tree/main/docs/modelserving)

### Community

- [KServe Slack](https://kubeflow.slack.com/archives/CH6E58LNP)
- [KServe GitHub](https://github.com/kserve/kserve)
- [Community Meetings](https://github.com/kserve/community)

---

## Appendix: Architecture Diagram

```text
┌─────────────────────────────────────────────────────────────┐
│                      Host Machine                           │
│  ┌───────────────────────────────────────────────────────┐  │
│  │          Container Runtime (Docker/Podman)           │  │
│  │                                                       │  │
│  │  ┌────────────────────────────────────────────────┐   │  │
│  │  │            Kind Cluster (kind)                 │   │  │
│  │  │                                                │   │  │
│  │  │  ┌──────────────┐      ┌──────────────┐        │   │  │
│  │  │  │    kind-     │      │    kind-     │        │   │  │
│  │  │  │  control-    │      │   worker     │        │   │  │
│  │  │  │   plane      │      │              │        │   │  │
│  │  │  │              │      │              │        │   │  │
│  │  │  │ - API Server │      │ - Inference  │        │   │  │
│  │  │  │ - Scheduler  │      │   Pods       │        │   │  │
│  │  │  │ - Controllers│      │ - Model Cache│        │   │  │
│  │  │  └──────────────┘      └──────────────┘        │   │  │
│  │  │                                                │   │  │
│  │  │  ┌────────────────────────────────────────┐    │   │  │
│  │  │  │        kserve namespace                │    │   │  │
│  │  │  │                                        │    │   │  │
│  │  │  │  - kserve-controller-manager           │    │   │  │
│  │  │  │  - kserve-localmodel-controller        │    │   │  │
│  │  │  │  - llmisvc-controller-manager          │    │   │  │
│  │  │  │  - kserve-localmodelnode-agent         │    │   │  │
│  │  │  │  - cert-manager                        │    │   │  │
│  │  │  └────────────────────────────────────────┘    │   │  │
│  │  └────────────────────────────────────────────────┘   │  │
│  │                                                       │  │
│  │  ┌────────────────────────────────────────────────┐   │  │
│  │  │     kind-registry (localhost:5001)             │   │  │
│  │  │                                                │   │  │
│  │  │  - Container image storage                     │   │  │
│  │  │  - Accessible from both nodes                  │   │  │
│  │  └────────────────────────────────────────────────┘   │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

Network:
- Host port 5001 → kind-registry:5000
- Host port 31000 → kind-control-plane:31000 (Ingress)
- kind network connects all containers
```

---

**Last Updated**: 2026-01-28
**KServe Version**: v0.13.0+
**Kind Version**: v0.20.0+
**Kubernetes Version**: v1.33.1
