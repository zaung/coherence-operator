#!/usr/bin/env bash

KIND=$(kind get clusters | grep operator-test)
if [[ "${KIND}" == "" ]]
then
  echo "Creating Kind cluster operator-test"
  kind create cluster --config etc/kind-config.yaml --name operator-test

  export KUBECONFIG="$(kind get kubeconfig-path --name="operator-test")"

# Install Helm
  echo "Creating Kind cluster operator-test - initialising Helm"
  helm init

# wait a few seconds so that the Tiller Pod should have appeared
  sleep 5

# Wait for Tiller to be ready...
  echo "Creating Kind cluster operator-test - waiting for Tiller"
  TILLER_POD=$(kubectl -n kube-system get pod -l name=tiller -o name)
  kubectl -n kube-system  wait --for=condition=Ready --timeout=300s ${TILLER_POD}

# The KinD cluster has RBAC so set up Helm properly
  echo "Creating Kind cluster operator-test - creating Tiller RBAC resources"
  kubectl create serviceaccount --namespace kube-system tiller
  kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
  kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

# Label the nodes so that site and rack get set in tests
  echo "Creating Kind cluster operator-test - adding node labels"
  kubectl label node operator-test-worker failure-domain.beta.kubernetes.io/zone=twighlight-zone --overwrite
  kubectl label node operator-test-worker failure-domain.beta.kubernetes.io/region=jenkins --overwrite
  kubectl label node operator-test-worker2 failure-domain.beta.kubernetes.io/zone=in-the-zone --overwrite
  kubectl label node operator-test-worker2 failure-domain.beta.kubernetes.io/region=jenkins --overwrite
  kubectl label node operator-test-worker3 failure-domain.beta.kubernetes.io/zone=zoned-out --overwrite
  kubectl label node operator-test-worker3 failure-domain.beta.kubernetes.io/region=jenkins --overwrite

# Install te rancher local storage provider so that we can use it for PVCs as KinD has no dynamic PVC out of the box
  echo "Creating Kind cluster operator-test - adding local path storage driver"
  kubectl apply -f hack/local-path-storage.yaml

# wait a few seconds so that the local storage provider Pod should have appeared
  sleep 5

# Wait for the local storage provider Pod to be ready...
  echo "Creating Kind cluster operator-test - waiting for local path storage driver pod"
  LOCAL_STORAGE_POD=$(kubectl -n local-path-storage get pod -o name)
  kubectl -n local-path-storage  wait --for=condition=Ready --timeout=300s ${LOCAL_STORAGE_POD}


  echo "Kind cluster operator-test created"
else
    echo "Kind cluster operator-test exists"
fi

