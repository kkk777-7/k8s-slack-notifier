# k8s-slack-notifier
Nofity k8s event to slack channel  
ex) pod Create, Delete

![MiConv com__notifier-demo](https://user-images.githubusercontent.com/86363983/198680141-937d6cef-2bb2-48cf-9f92-26c192acd5fb.gif)

My Practice Custom Controller!

# Prerequisites
The following must be installed
- kubectl
- Kubebuilder
- kind
- Go
- Docker
- tilt
- ctlptl

# Quick start
```
$ make start
$
$ vi examples/slack.yaml // setup your environment
$
$ kubectl create ns k8s-slack-notifier-system
$ kubectl create secret generic slack-secret --from-file=examples/slack.yaml -n k8s-slack-notifier-system
$ tilt up
$
$ kubectl apply -f examples/statefulset.yaml
```

# Teardown
```
$ tilt down
$ make stop 
```