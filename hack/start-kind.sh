#!/usr/bin/env bash

KIND=$(kind get clusters | grep operator-test)
if [[ "${KIND}" == "" ]]
then
  echo "Creating Kind cluster operator-test"
  kind create cluster --config etc/kind-config.yaml --name operator-test

  export KUBECONFIG="$(kind get kubeconfig-path --name="operator-test")"

  helm init

  TILLER_POD=$(kubectl -n kube-system get pod -l name=tiller -o name)
  kubectl -n kube-system  wait --for=condition=Ready ${TILLER_POD}

  kubectl create serviceaccount --namespace kube-system tiller
  kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
  kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

  kubectl label node operator-test-worker failure-domain.beta.kubernetes.io/zone=twighlight-zone --overwrite
  kubectl label node operator-test-worker failure-domain.beta.kubernetes.io/region=jenkins --overwrite
  kubectl label node operator-test-worker2 failure-domain.beta.kubernetes.io/zone=in-the-zone --overwrite
  kubectl label node operator-test-worker2 failure-domain.beta.kubernetes.io/region=jenkins --overwrite
  kubectl label node operator-test-worker3 failure-domain.beta.kubernetes.io/zone=zoned-out --overwrite
  kubectl label node operator-test-worker3 failure-domain.beta.kubernetes.io/region=jenkins --overwrite

  kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
  LOCAL_STORAGE_POD=$(kubectl -n local-path-storage get pod -o name)
  kubectl -n local-path-storage  wait --for=condition=Ready ${LOCAL_STORAGE_POD}
else
    echo "Kind cluster operator-test exists"
fi

