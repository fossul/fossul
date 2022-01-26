#!/bin/bash
# Needs to be run from operator directory

OPERATOR_IMAGE="fossul-operator:latest"
BUNDLE_IMAGE="fossul-operator-bundle:latest"
REPO="fossul"

sudo make docker-build docker-push IMG=quay.io/$REPO/$OPERATOR_IMAGE
if [ $? != 0 ]; then exit 1; fi

sudo make bundle-build BUNDLE_IMG=quay.io/$REPO/$BUNDLE_IMAGE
if [ $? != 0 ]; then exit 1; fi

sudo docker push quay.io/$REPO/$BUNDLE_IMAGE
if [ $? != 0 ]; then exit 1; fi
