#!/bin/sh

rm -rvf ../docs 2>&1 || true
cd ../docs-src && npx docusaurus generate-proto-docs && npm run build && mv build ../docs
