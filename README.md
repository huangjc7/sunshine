# Sunshine 
* Had I not seen the sun, I could have borne the shade.
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-- Siegfried Loraine Sassoon
* An application that makes kubernetes clusters shine
#### Introduction
1. After the pod restarts for the specified number of times, it will find its deployment and shut down
2. If there is a running instance corresponding to the deployment, it will give up processing the deployment
3. Reduce useless resource consumption, alarm information, and kube-scheduler scheduling pressure
#### Applicable scene
* Suitable for multi-tenant messy kubernetes cluster
#### Supported kubernetes versions
* kubernetes v1.17+
#### Use
```shell
kubectl apply -f https://raw.githubusercontent.com/huangjc7/sunshine/master/manifests/sunshine.yaml
```
#### Environment Variable

| Environment Variable                  | Description                                  | Default |
|----------------------|----------------------------------------------|---------|
| `POD_RESTART_NUMBER` |When the pod restart threshold is met, the deployment will be closed                                | `50`    |

<video src="video/1.mp4" width="800px" height="600px" controls="controls"></video>

