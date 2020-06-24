# About this repository
I've created this project to study and learn golang and the operators framework.

## The labels-operator
The project contains a Kubernetes operator that adds some labels to all pods from namespaces configured in the custom resource.

It is developed using the [operator-framework](https://operatorframework.io/).

# Pre-requisites
- A cluster-admin access to a running OpenShift 4.x cluster
- Access to the container registry used by the ocp cluster
- git to clone this repo
- The operator-sdk just to build this project (WIP push a working image to a public registry)
- docker/podman/buildah to pull the built image to the registry (WIP push a working image to a public registry)


# Installing the labels-operator
* Clone the labels-operator repo
~~~
git clone git@github.com:git-hyagi/labels-operator.git
~~~

* Build it
~~~
cd labels-operator
operator-sdk build <registry address>/labels-operator/labels:v1
~~~

* Push the built image to the registry
~~~
docker push <registry address>/labels-operator/labels:v1
~~~

* Create the cluster objects
~~~
oc new-project labels-operator
oc create -f  deploy/crds/lab.local_labels_crd.yaml
oc create -f  deploy/service_account.yaml
oc create -f  deploy/role.yaml
oc create -f  deploy/rolebinding.yaml
oc create -f  deploy/operator.yaml
~~~

* Create a **labels-operator** `custom resource` with the projects that should have the pods with the labels
~~~
cat<<EOF> lab.local_v1_label_cr.yaml
apiVersion: lab.local/v1
kind: Label
metadata:
  name: labels-operator
  namespace: labels-operator
spec:
  projects:
  - <my project A>
  - <my project B>
EOF

oc create -f  lab.local_v1_label_cr.yaml
~~~

# Configuring the labels-operator
For now (there is an open [issue](https://github.com/git-hyagi/labels-operator/issues/8) to allow the labels customization), the only configuration available is to add or remove projects from the **labels-operator** `custom resource` which will, in turn, sync or not the labels from all pods of these projects:
~~~
oc edit labels-operator labels-operator
~~~
