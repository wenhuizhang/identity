# Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0
"""Identity SDK for Python."""

import inspect
import logging
import os
from importlib import import_module
from pkgutil import iter_modules

import agntcy.identity.node.v1alpha1
from dotenv import load_dotenv

from air import client, log


logger = logging.getLogger("identity")

if os.getenv("IDENTITY_ENABLE_LOGS", "0") == "1":
    load_dotenv()
    log.configure()


def _load_grpc_objects():
    """Load all the objects from the Python Identity SDK."""
    for _, modname, _ in iter_modules(agntcy.identity.node.v1alpha1.__path__):
        # Import the module
        module = import_module(f"agntcy.identity.node.v1alpha1.{modname}")
        # Inspect the module and set attributes on Identity SDK for each class found
        for name, obj in inspect.getmembers(module, inspect.isclass):
            setattr(IdentitySdk, name, obj)


class IdentitySdk:
    """Identity SDK for Python."""

    def __init__(self, async_mode=False):
        """Initialize the Identity SDK."""
        # Load dynamically all objects
        _load_grpc_objects()

        self.client = client.Client(async_mode)

    def get_vc_service(self):
        """Get the vc service."""
        return IdentitySdk.VcServiceStub(self.client.channel)
