#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE[0]}")/..
helm package ${ROOT}/helm/cluster-api-visualizer -d ${ROOT}/helm/repo
helm repo index ${ROOT}/helm/repo