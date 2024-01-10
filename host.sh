#!/bin/bash

SRV_PKG_LABEL=$1
echo_pkg() {
  echo -n "($SRV_PKG_LABEL/host): "
  echo $*
}

echo_pkg "Building Server..."
go build -o srv-$SRV_PKG_LABEL

echo_pkg "Launching Server..."
./srv-$SRV_PKG_LABEL -cfg "$SRV_PKG_LABEL.toml"

echo_pkg "Disposing Server..."
rm srv-$SRV_PKG_LABEL

echo_pkg "Done."
sleep "1s"
