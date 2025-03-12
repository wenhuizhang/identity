#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

PROTO_PACKAGE_NAME="agntcy.pyramid.v1alpha1"
PROTO_FILE_PATH="agntcy.pyramid/v1alpha1/"

function get_module_name_from_package() {
  echo $(dirname "$1" | xargs basename)
}

echo ""
echo " _____         _____      ______          _        "
echo "|  __ \       |_   _|     | ___ \        | |       "
echo "| |  \/ ___     | | ___   | |_/ / __ ___ | |_ ___  "
echo "| | __ / _ \    | |/ _ \  |  __/ '__/ _ \| __/ _ \ "
echo "| |_\ \ (_) |   | | (_) | | |  | | | (_) | || (_) |"
echo " \____/\___/    \_/\___/  \_|  |_|  \___/ \__\___/ "
echo ""

source "${PyramID_ROOT}/protoc.sh"
cd ${PyramID_ROOT}

protoc_install

cd "${PyramID_ROOT}/local"

type_files=$(find . -path "*/app/*/types/types.go")
packages=""

for file in $type_files; do
  dir="${file#./}"
  package=$(dirname "$dir")
  packages="$packages $package"
done

packages=$(echo "$packages" | sed 's/\s$//' | sed 's/^\s//')

cd "${PyramID_ROOT}/local/github.com/agntcy/pyramid"

go get github.com/gogo/protobuf/proto
go mod tidy

packages_comma_separated=$(echo "$packages" | tr ' ' ',')

# Detect GO enums
go-enum-to-proto \
  --packages="${packages_comma_separated}" \
  --output-dir="${PyramID_ROOT}/local"

# Tag the GO enums by changing them to structs so that go-to-protobuf
# can reference them by name and not by the underlying type (ex: int)
go-enum-patch --patch="${PyramID_ROOT}/local/enums.json" --type=go

go-to-protobuf \
  --apimachinery-packages="" \
  --proto-import="${PyramID_ROOT}/third_party/protos" \
  --output-dir="${PyramID_ROOT}/local" \
  --packages="${packages_comma_separated}" \
  --keep-gogoproto=false \
  -v=8

# Change the enums detected earlier from proto messages to actual proto enums
go-enum-patch --patch="${PyramID_ROOT}/local/enums.json" --type=proto

cd "${PyramID_ROOT}"

for package in $packages; do
  mkdir -p "local/output"
  module_name=$(get_module_name_from_package "${package}")
  cp "local/${package}/generated.proto" "local/output/${module_name}.proto"
done

protos=$(find local/output -iname "*.proto")

for m in $protos; do
  sed -i 's/syntax = "proto2";/syntax = "proto3";/g' "${m}"
  sed -i 's|go_package = [^ ]\+|go_package = "github.com/agntcy/pyramid/internal/pkg/generated/agntcy.pyramid/v1alpha1;pyramid_sdk_go";|g' "${m}"
done

for package in $packages; do
  proto_file=$(get_module_name_from_package "${package}")
  package=$(echo "$package" | sed 's|\.|\\.|g')
  import=$(echo "$package" | sed 's|/|\\.|g')
  for m in $protos; do
    sed -i "s|${import}|${PROTO_PACKAGE_NAME}|g" "${m}"
    sed -i "s|${package}/generated.proto|${PROTO_FILE_PATH}${proto_file}.proto|g" "${m}"
  done
done

cp -r "${PyramID_ROOT}/local/output/." "${PyramID_ROOT}/code/api-spec/proto/agntcy.pyramid/v1alpha1"

echo ""
echo "______ _   _______   _____                           _       "
echo "| ___ \ | | |  ___| |  __ \                         | |      "
echo "| |_/ / | | | |_    | |  \/ ___ _ __   ___ _ __ __ _| |_ ___ "
echo "| ___ \ | | |  _|   | | __ / _ \ '_ \ / _ \ '__/ _  | __/ _ \ "
echo "| |_/ / |_| | |     | |_\ \  __/ | | |  __/ | | (_| | ||  __/"
echo "\____/ \___/\_|      \____/\___|_| |_|\___|_|  \__,_|\__\___|"
echo ""

rm -rvf ${PyramID_ROOT}/code/pyramid/internal/pkg/generated 2>&1 || true

cd "${PyramID_ROOT}/code/api-spec"

# Go
/usr/local/bin/buf generate --debug -v

# Python
/usr/local/bin/buf generate --include-imports --template buf.gen.python.yaml --debug -v

# Openapi
/usr/local/bin/buf generate --template buf.gen.openapi.yaml --output ../docs-src/static/api/openapi/v1alpha1 --path proto/agntcy.pyramid/v1alpha1

# Proto
/usr/local/bin/buf generate --template buf.gen.doc.yaml --output ../docs-src/static/api/proto/v1alpha1 --path proto/agntcy.pyramid/v1alpha1
