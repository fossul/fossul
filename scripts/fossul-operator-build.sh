#!/bin/bash
# Needs to be run from operator directory

RELEASE="v0.7.0"
OPERATOR_IMAGE="fossul-operator"
BUNDLE_IMAGE="fossul-operator-bundle"
REPO="fossul"

make docker-build docker-push IMG=quay.io/$REPO/$OPERATOR_IMAGE:$RELEASE
if [ $? != 0 ]; then exit 1; fi

docker tag quay.io/$REPO/$OPERATOR_IMAGE:$RELEASE quay.io/$REPO/$OPERATOR_IMAGE:latest
if [ $? != 0 ]; then exit 1; fi

docker push quay.io/$REPO/$OPERATOR_IMAGE:latest
if [ $? != 0 ]; then exit 1; fi

make bundle-build BUNDLE_IMG=quay.io/$REPO/$BUNDLE_IMAGE:$RELEASE
if [ $? != 0 ]; then exit 1; fi

docker push quay.io/$REPO/$BUNDLE_IMAGE:$RELEASE
if [ $? != 0 ]; then exit 1; fi

docker tag quay.io/$REPO/$BUNDLE_IMAGE:$RELEASE quay.io/$REPO/$BUNDLE_IMAGE:latest
if [ $? != 0 ]; then exit 1; fi

docker push quay.io/$REPO/$BUNDLE_IMAGE:latest
if [ $? != 0 ]; then exit 1; fi
