#!/bin/sh
# Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

set -e

# Global vars
######################################################################
BINARY_ARCH=$(uname -m)
if [ "$BINARY_ARCH" = "x86_64" ]; then
  BINARY_ARCH="amd64"
elif [ "$BINARY_ARCH" = "aarch64" ]; then
  BINARY_ARCH="arm64"
fi
BINARY_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
BINARY_NAME="identity"
LATEST_BINARY_URI="https://github.com/agntcy/identity/releases/latest/download/${BINARY_NAME}_${BINARY_OS}_${BINARY_ARCH}"

INSTALL_LOCATION="/usr/local/bin"
SUDO="sudo"
SUDO_CHMOD_COMMAND="${SUDO} chmod 755 ${INSTALL_LOCATION}/${BINARY_NAME}"
SUDO_WGET_COMMAND="${SUDO} wget -q ${LATEST_BINARY_URI} -O ${INSTALL_LOCATION}/${BINARY_NAME}"
SUDO_CURL_COMMAND="${SUDO} curl -sL ${LATEST_BINARY_URI} -o ${INSTALL_LOCATION}/${BINARY_NAME}"
CURL_IS_AVAILABLE="command -v curl"
WGET_IS_AVAILABLE="command -v wget"
######################################################################

# Check if the user can sudo
########################################################################
check_can_sudo() {
  echo "Sudo permissions are required to install the CLI to ${INSTALL_LOCATION}."
  echo "Checking if the user can sudo..."

  if ! sudo -v > /dev/null; then
    echo "The user cannot sudo and will not be able to install the CLI."

    exit 1
  fi
}

# Install the CLI
######################################################################
do_install() {
  echo "Installing the CLI to ${INSTALL_LOCATION}..."

  # Download the latest binary
  if $CURL_IS_AVAILABLE > /dev/null; then
    $SUDO_CURL_COMMAND
  elif $WGET_IS_AVAILABLE > /dev/null; then
    $SUDO_WGET_COMMAND
  else
    echo "Curl or wget must be present for this installer to work."

    exit 1
  fi

  echo "Downloaded the latest binary to ${INSTALL_LOCATION}/${BINARY_NAME}."
  echo "Setting permissions and ownership for ${INSTALL_LOCATION}/${BINARY_NAME}..."

  # Set permissions & ownership
  $SUDO_CHMOD_COMMAND

  echo "Permissions and ownership set for ${INSTALL_LOCATION}/${BINARY_NAME}."
  echo "Installation complete. You can now use the CLI by running '${BINARY_NAME}' from your terminal."
}
######################################################################

# Check if the user can sudo
check_can_sudo

# Install
do_install
