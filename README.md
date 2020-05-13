To run: go build cmd/main.go
To test: go test

Questions:
1.) How would you deploy your controller to a Kubernetes cluster?
    Construct a yaml file for it, with links to all needed dependencies, build a release-tag image, and deploy it either locally using minikube or kind by kubectl apply -f controller.yaml
2.) More resilient from crashing, as if the controller crashes for whatever reason, it can look at the current state of the cluster when it comes back up.
7.) It could affect the testing of the controller by not having a stable string to compare the output to, which can be averted by mocking etc.