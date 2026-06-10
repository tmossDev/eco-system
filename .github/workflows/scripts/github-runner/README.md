# GitHub Runner

This repo's deploy workflow uses:

```yaml
runs-on: [self-hosted, linux]
```

The simplest matching setup is a self-hosted runner installed as a Linux
service on the same machine that has Docker, kubectl, Helm, and k3s access.
That matches the workflow's local image flow:

```sh
docker build ...
docker save ... | sudo k3s ctr -n k8s.io images import -
```

## Install As A Host Service

Generate a fresh repository runner token in GitHub:

```text
Settings -> Actions -> Runners -> New self-hosted runner -> Linux x64
```

Then run:

```sh
export RUNNER_TOKEN='<fresh github runner token>'
./ci/github-runner/install-host-runner.sh
```

On a fresh Linux host, let the runner install its system package dependencies:

```sh
export RUNNER_TOKEN='<fresh github runner token>'
INSTALL_RUNNER_DEPENDENCIES=true ./ci/github-runner/install-host-runner.sh
```

The script downloads `actions-runner-linux-x64-2.334.0.tar.gz` and verifies it
with:

```sh
echo "048024cd2c848eb6f14d5646d56c13a4def2ae7ee3ad12122bee960c56f3d271  actions-runner-linux-x64-2.334.0.tar.gz" | shasum -a 256 -c
```

It registers the runner against `https://github.com/tmossDev/eco-system` with
the labels `self-hosted,linux`, then installs and starts the runner service.

## Required Host Access

The runner user needs:

- Docker daemon access.
- `kubectl` and `helm` on `PATH`.
- Kubeconfig access for the target k3s cluster.
- Passwordless sudo for the existing `sudo k3s ctr ... images import` workflow
  step, or the workflow must be changed to avoid `sudo`.

For passwordless k3s image import, restrict sudo as tightly as practical, for
example with a dedicated runner user and a sudoers rule for the exact k3s
binary path used on the host.

## About An `eco-ci` Namespace

It is possible to run GitHub runners inside Kubernetes in an `eco-ci`
namespace, usually with Actions Runner Controller. For this repo, that should
be a second step rather than the first runner setup because the current workflow
depends on host Docker and `sudo k3s ctr` access.

To make an in-cluster runner clean, change the image flow to push images to a
registry reachable by k3s, then deploy with normal image pulls:

```text
GitHub runner pod -> build image -> push registry -> helm upgrade --set image.tag=...
```

Without that registry change, an in-cluster runner must be privileged and mount
host Docker/k3s resources, which is workable for a homelab but much easier to
break and harder to secure.
