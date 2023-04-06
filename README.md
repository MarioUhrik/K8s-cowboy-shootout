### K8s-cowboy-shootout
Go/GRPC/K8s simulation of a wild west shootout between 2 or more cowboys. \
Each cowboy is a distributed standalone pod, shooting at the others via GRPC.

![Architectural Diagram](k8s-cowboy-shootout.png)

### Usage

1. Clone the project `git clone git@github.com:MarioUhrik/K8s-cowboy-shootout.git`
2. Optionally, edit `helm/cowboys.json` to your satisfaction
3. Deploy into your K8s cluster via `helm template helm/ | kubectl apply -f -`
4. Check the list of deployed pods via `kubectl -n k8s-cowboy-shootout get pod` - each cowboy is a pod. Unreadiness represents death.
5. Check the logs of the cowboys to see all actions during the duel `kubectl -n k8s-cowboy-shootout logs PODNAME`
