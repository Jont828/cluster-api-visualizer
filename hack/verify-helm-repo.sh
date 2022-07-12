set -o nounset
set -o errexit
set -o pipefail
set -x

ROOT=$(dirname "${BASH_SOURCE}")/..

for chart in "cluster-api-visualizer"; do
  LATEST_CHART=$(ls ${ROOT}/helm/repo/${chart}*.tgz | sort -rV | head -n 1)
  LATEST_VERSION=$(echo $LATEST_CHART | grep -Eoq [0-9]+\.[0-9]+\.[0-9]+)
  MATCH_STRING="version: $LATEST_VERSION"
  # verify that the current version in Chart.yaml is the most recent, packaged chart in the repo
  if ! grep -q "${MATCH_STRING}" ${ROOT}/helm/${chart}/Chart.yaml; then
    echo The version of the $chart helm chart checked into the git repository does not match the latest packaged version $LATEST_VERSION in the repo
    exit 1
  fi
  rm -Rf ${ROOT}/chart_verify
  mkdir ${ROOT}/chart_verify
  tar -xf $LATEST_CHART -C chart_verify
  diff -r ${ROOT}/chart_verify/cluster-api-visualizer/ ${ROOT}/helm/cluster-api-visualizer/ --exclude README.md || exit 1
  rm -Rf ${ROOT}/chart_verify
done

curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
helm template ./helm/cluster-api-visualizer > /dev/null

exit 0