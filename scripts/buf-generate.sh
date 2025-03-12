#!/bin/sh

function reset_generated_pb_go() {
  (
    local pb_files
    cd ..
    pb_files=$(git status --porcelain | grep backend/internal/pkg/generated | sed s/^...// | tr '\n' ' ')

    if [ -n "${pb_files}" ]; then
      echo "Resetting backend/internal/pkg/generated"
      git checkout -- $pb_files
    fi
  )
}

reset_generated_pb_go

rm -rvf ../sdk/python/internal/generated 2>&1 || true
cd docker && \
  docker compose -f buf-compose.yaml build --no-cache && \
  docker compose -f buf-compose.yaml run --rm -w /pyramid/code/api-spec buf-go run.sh
docker rmi docker-buf-go
cd ../../backend/internal/pkg/generated/ && \
  grep -rl gnostic . | xargs sed -i '' 's/gnostic/gnostic-models/g'
