#!/bin/bash/

SRV_PKG_LABEL=$1

echo_pkg() {
  echo -n "(deploy/$SRV_PKG_LABEL): "
  echo $*
}

DEPLOY_OUTPUT_ROOT="out-$SRV_PKG_LABEL"
if [ ! -d $DEPLOY_OUTPUT_ROOT ]; then
  mkdir $DEPLOY_OUTPUT_ROOT
  echo_pkg "Initialized output-root: \"$DEPLOY_OUTPUT_ROOT\"."
fi

SRV_PKG_ROOT="shush-$SRV_PKG_LABEL"
if [ ! -d "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT" ]; then
  mkdir "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT"
  echo_pkg "Initialized package-root: \"$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT\"."
fi

SRV_PKG_CONFIG="$SRV_PKG_LABEL.toml"
cp $SRV_PKG_CONFIG "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/$SRV_PKG_CONFIG"

if [ -d "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web" ]; then
  rm -r $DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web
fi
cp -r web $DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
echo_pkg -n "Compile to another infrastructure? [y/n]: "
read CHANGE_INFRASTRUCTURE
if test $CHANGE_INFRASTRUCTURE = 'y' || test $CHANGE_INFRASTRUCTURE == 'Y'; then
  printf '\n'
  echo_pkg "Enter a platform like .."
  echo_pkg "(windows/linux/darwin/netbsd/openbsd/freebsd/solaris/dragonfly/android/plan9)"
  echo_pkg -e -n "-> "
  read OS
  printf '\n'
  echo_pkg "Enter a architecture like .."
  echo_pkg "(386/amd64/arm/arm64/ppc64/ppc64le/mips/mips64/s390x/riscv64)"
  echo_pkg -e -n "-> "
  read ARCH
  printf '\n'
fi

echo_pkg "Compiling..."
env GOOS=$OS GOARCH=$ARCH go build -o "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/serv"

DEFAULT_FTP_USERNAME="effyiex"
DEFAULT_FTP_HOST="raspberrypi"
if test $2 == "/ftp" || test $2 == "-ftp"; then

  echo_ftp() {
    echo -n "(ftp-deploy/$SRV_PKG_LABEL): "
    echo $*
  }

  echo_ftp -n "Continue with default connection? [y/n]: "
  read FTP_DEFAULT_CONNECTION
  if test $FTP_DEFAULT_CONNECTION = 'y' || test $FTP_DEFAULT_CONNECTION == 'Y'; then
    FTP_USERNAME=$DEFAULT_FTP_USERNAME
    FTP_HOST=$DEFAULT_FTP_HOST  
  else
    echo_ftp -n "Enter hostname: "
    read FTP_HOST
    echo_ftp -n "Enter username: "
    read FTP_USERNAME
  fi

  echo_ftp -n "Enter Password: "
  read -s FTP_PASSWORD
  printf '\n'

  echo_ftp "Indexing local package-instance..."
  FTP_LOCAL_ROOT="$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT"
  FTP_REMOTE_ROOT="~/Documents/$SRV_PKG_ROOT"
  FTP_INDEXING=$(
    find $FTP_LOCAL_ROOT -type d -printf "mkdir $FTP_REMOTE_ROOT/%P\n"
    find $FTP_LOCAL_ROOT -type f -exec echo -n "put {} " \; -printf "$FTP_REMOTE_ROOT/%P\n"
  )
  echo "$(echo "$FTP_INDEXING")" > "$DEPLOY_OUTPUT_ROOT/latest-ftp-indexing.log"
  sleep "2s"

  echo_ftp "Deploying to \"$FTP_HOST\"..."
  ftp -n $FTP_HOST <<END_OF_PROTOCOL > "$DEPLOY_OUTPUT_ROOT/latest-ftp-deploy.log"
quote USER $FTP_USERNAME
quote PASS $FTP_PASSWORD
binary
$FTP_INDEXING
bye
END_OF_PROTOCOL
  sleep "2s"

  echo_ftp "Disposing local package-instance..."
  rm -r "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT"
  echo_ftp "Done." 

fi
