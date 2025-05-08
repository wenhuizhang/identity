#!/bin/sh

reset_generated_pb_go() {
  (
    cd ..
    pb_files=$(git status --porcelain | grep identity/api | sed s/^...// | tr '\n' ' ')

    if [ -n "${pb_files}" ]; then
      echo "Resetting identity/api"
      git checkout -- "$pb_files"
    fi
  )
}

reset_generated_pb_go

cd docker &&
  docker compose -f buf-compose.yaml build --no-cache &&
  docker compose -f buf-compose.yaml run --rm -w /identity/code/api-spec buf-go run.sh
docker rmi docker-buf-go

if [ -d "../../identity/api/" ]; then
  cd ../../identity/api/ &&
    grep -rl gnostic . | xargs sed -i '' 's|github.com/google/gnostic/openapiv3|github.com/google/gnostic-models/openapiv3|g'
fi
