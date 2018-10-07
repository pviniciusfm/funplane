#!/usr/bin/env bash -xe

kubectl delete namespace fanplane
kubectl delete role fanplane-role
kubectl delete clusterrolebinding fanplane-clusterrole
kubectl delete sa/buildkite-fanplane-account
kubectl delete crd gateways.fanplane.io
kubectl delete crd envoybootstraps.fanplane.io
