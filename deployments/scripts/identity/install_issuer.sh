#!/bin/sh
# Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

set -e

# Global vars
######################################################################
BINARY_ARCH=$(uname -m)
BINARY_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
BINARY_NAME="identity"
BINARY_VERSION="0.0.1-beta.3"
LATEST_BINARY_URI="https://github.com/agntcy/identity/releases/latest/download/${BINARY_NAME}_${BINARY_VERSION}_${BINARY_OS}_${BINARY_ARCH}.tar.gz"

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
check_sudo() {
  if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root or with sudo."

    exit 1
  fi
}

# Install the CLI
######################################################################
do_install() {
  echo $LATEST_BINARY_URI
  # Download the latest binary
  if $CURL_IS_AVAILABLE > /dev/null; then
    $SUDO_CURL_COMMAND
  elif $WGET_IS_AVAILABLE > /dev/null; then
    $SUDO_WGET_COMMAND
  else
    echo "Curl or wget must be present for this installer to work."

    exit 1
  fi

  # Set permissions & ownership
  $SUDO_CHMOD_COMMAND
}
######################################################################

# Check if the user can sudo
check_sudo

# Install
do_install
