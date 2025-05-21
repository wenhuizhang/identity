#!/bin/sh
# Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


#-----------------------------------------------------------------------------#
# Globals
ROOT_DIR="/identity"
OUT_DIR="$ROOT_DIR/out"
SDK_DIR="sdk/node-go"
SPEC_DIR="$OUT_DIR/api-spec"
SWAGGER_DIR="generated"
PLATFORM_V1ALPHA1_PROTO_PATH="proto/agntcy/identity/node/v1alpha1"
PLATFORM_V1ALPHA1_GENERATED_PATH="$ROOT_DIR/generated/openapi/agntcy/identity/node/v1alpha1"
SWAGGER_SECURITY_DEFINITIONS_PATH="$ROOT_DIR/generated/openapi/agntcy/identity/node/v1alpha1/openapi.swagger.json"
#-----------------------------------------------------------------------------#

function do_cleanup_previous_version() {
  cd $OUT_DIR
  if [ -d "$SDK_DIR/client" ]; then
    rm -rf $SDK_DIR/client
  fi
  if [ -d "$SDK_DIR/models" ]; then
    rm -rf $SDK_DIR/models
  fi
  cd $ROOT_DIR
}

function do_cleanup_swagger() {
  cd $OUT_DIR
  if [ -d "$SWAGGER_DIR" ]; then
    rm -rf $SWAGGER_DIR
  fi
  cd $ROOT_DIR
}

function do_mixin() {
  cd $OUT_DIR
  for file in $(ls $1); do
    if [[ "$(printf '%s' "$file")" = *"_service"* ]]; then
      swagger mixin $1/$file $SWAGGER_SECURITY_DEFINITIONS_PATH > $1/final_$file
      mv -f $1/final_$file $1/$file
    fi
  done
  cd $ROOT_DIR
}

function do_rename() {
  cd $OUT_DIR
  for file in $(ls $1); do
    if [[ "$(printf '%s' "$file")" != *"_service"* ]]; then
      rm $1/$file
    fi
  done
  for file in $(ls $1); do
    name=${file%_service*}
    if [[ "$(printf '%s' "$file")" = *"_service"* ]]; then
      mv $1/$file $1/$name".swagger.json"
    fi
  done
  cd $ROOT_DIR
}

function do_generate() {
  cd $SPEC_DIR
  buf generate --template buf.gen.openapiv2.yaml  --output $ROOT_DIR --path $1
  cd $ROOT_DIR
}

function do_generate_go_client() {
  cd $OUT_DIR
  for file in $(ls $1); do
    swagger generate client -f $1/$file -A identity_node --template-dir $ROOT_DIR/templates --target ./$SDK_DIR
  done
  cd $ROOT_DIR
}

function do_generate_all() {
  do_generate $PLATFORM_V1ALPHA1_PROTO_PATH
  do_mixin $PLATFORM_V1ALPHA1_GENERATED_PATH
  do_rename $PLATFORM_V1ALPHA1_GENERATED_PATH
}

function do_generate_go_client_all() {
  do_generate_go_client $PLATFORM_V1ALPHA1_GENERATED_PATH
}

# Generate the code

do_cleanup_previous_version
do_generate_all
do_generate_go_client_all
do_cleanup_swagger

echo "Done"
