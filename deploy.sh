#!/bin/bash/

PKG_FLAG_FTP_DEPLOY=false
PKG_FLAG_TLS_DEPLOY=false
for arg in "$@"
do
  if test $arg = "-ftp" || test $arg = "/ftp"; then
    PKG_FLAG_FTP_DEPLOY=true
  elif test $arg = "-tls" || test $arg = "/tls"; then
    PKG_FLAG_TLS_DEPLOY=true
  fi
done

SRV_PKG_LABEL=$1

ECHO_PKG_INSTITUTION="deploy"
echo_pkg() {
  echo -n "($SRV_PKG_LABEL/$ECHO_PKG_INSTITUTION): "
  echo $*
  if [ "$(echo -n "$*" | tail -c 1)" = " " ]; then 
    printf ' '
  fi
}

DEPLOY_OUTPUT_ROOT="pkg-$SRV_PKG_LABEL"
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
cp $SRV_PKG_CONFIG "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/config.toml"

if [ -d "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web" ]; then
  rm -r $DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web
fi
cp -r web $DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/web

PKG_OS=$(go env GOOS)
PKG_ARCH=$(go env GOARCH)
while [ true ]; do
  echo_pkg -n "Compile to another infrastructure? [y/n]: "
  read CHANGE_INFRASTRUCTURE
  if test $CHANGE_INFRASTRUCTURE = 'y' || test $CHANGE_INFRASTRUCTURE = 'Y'; then
    printf '\n'
    echo_pkg "Enter a platform like .."
    echo_pkg "(windows/linux/darwin/netbsd/openbsd/freebsd/solaris/dragonfly/android/plan9)"
    echo_pkg -e -n "-> "
    read PKG_OS
    printf '\n'
    echo_pkg "Enter a architecture like .."
    echo_pkg "(386/amd64/arm/arm64/ppc64/ppc64le/mips/mips64/s390x/riscv64)"
    echo_pkg -e -n "-> "
    read PKG_ARCH
    printf '\n'
    break
  elif test $CHANGE_INFRASTRUCTURE = 'n' || test $CHANGE_INFRASTRUCTURE = 'N'; then
    break
  fi
done

echo_pkg "Compiling server..."
env GOOS=$PKG_OS GOARCH=$PKG_ARCH go build -o "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/srv-$SRV_PKG_LABEL"
sleep "1s"
echo_pkg "Done."

if [ $PKG_FLAG_TLS_DEPLOY = true ]; then

  ECHO_PKG_INSTITUTION="tls-deploy"

  if [ ! -d "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/tls" ]; then 
    echo_pkg "Initializing TLS-Root..."
    mkdir "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/tls"
  fi

  echo_pkg -e -n "Certificate (PEM) > "
  read TLS_IN_CERTIFICATE
  echo_pkg -e -n "Key (PEM) > "
  read TLS_IN_KEY

  echo_pkg "Writing Certificate..."
  echo "$TLS_IN_CERTIFICATE" > "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/tls/certificate.pem"
  echo_pkg "Writing Key..."
  echo "$TLS_IN_KEY" > "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT/tls/key.pem"

  sleep "1s"
  echo_pkg "Done."

fi

DEFAULT_FTP_ENDPOINT_ROOT="Effyiex-Software"
DEFAULT_FTP_USERNAME="effyiex"
DEFAULT_FTP_PI_HOST="Effyiex-PI"
DEFAULT_FTP_MIN_HOST="Effyiex-Minimal"
if [ PKG_FLAG_FTP_DEPLOY = true ]; then

  ECHO_PKG_INSTITUTION="ftp-deploy"

  while [ true ]; do
    echo_pkg "Continue with which connection? (pi/min/custom/cancel)"
    echo_pkg -e -n "-> "
    read FTP_IN_CONNECTION
    FTP_USERNAME=$DEFAULT_FTP_USERNAME
    if test $FTP_IN_CONNECTION = "pi"; then
      FTP_HOST=$DEFAULT_FTP_PI_HOST
    elif test $FTP_IN_CONNECTION = "min"; then
      FTP_HOST=$DEFAULT_FTP_MIN_HOST  
    elif test $FTP_IN_CONNECTION = "custom"; then
      echo_pkg -n "Enter hostname: "
      read FTP_HOST
      echo_pkg -n "Enter username: "
      read FTP_USERNAME
    elif test $FTP_IN_CONNECTION = "cancel"; then
      echo_pkg "Canceled."
      exit 0
    fi
    ping $FTP_HOST &>/dev/null; if [ $? = 0 ]; then
      break
    else
      echo_pkg "Connection failed. Try again."
    fi
  done

  echo_pkg -n "Enter Password: "
  printf ' '
  read -s FTP_PASSWORD
  printf '\n'

  echo_pkg "Indexing local package-instance..."
  FTP_LOCAL_ROOT="$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT"
  FTP_REMOTE_ROOT="$DEFAULT_FTP_ENDPOINT_ROOT/$SRV_PKG_ROOT"
  FTP_INDEXING=$(
    printf "mkdir $DEFAULT_FTP_ENDPOINT_ROOT\n"
    find $FTP_LOCAL_ROOT -type d -printf "mkdir $FTP_REMOTE_ROOT/%P\n"
    find $FTP_LOCAL_ROOT -type f -exec echo -n "put {} " \; -printf "$FTP_REMOTE_ROOT/%P\n"
  )
  echo "$(echo "$FTP_INDEXING")" > "$DEPLOY_OUTPUT_ROOT/latest-ftp-indexing.log"
  sleep "1s"

  echo_pkg "Initiating Python-FTP-Client..."
  echo "
from ftplib import FTP, all_errors as FTP_ERRORS
print(\"Creating FTP-Client...\")
ftp = FTP(\"$FTP_HOST\")
print(\"Logging in...\")
ftp.login(\"$FTP_USERNAME\", \"$FTP_PASSWORD\")
print(\"Starting transfer...\")
for indexed in \"\"\"$FTP_INDEXING\"\"\".split(\"\n\"):
  args = indexed.split(\" \")
  if len(args) < 2:
    continue
  elif args[0] == \"mkdir\":
    print(\" > MKD \" + args[1])
    try:
      print(\" < \" + ftp.mkd(args[1]))
    except FTP_ERRORS as e:
      print(\" < Error: \" + str(e).split(None, 1)[0])
  elif len(args) < 3:
    continue
  elif args[0] == \"put\":
    print(\" > STOR \" + args[2])
    try:
      print(\" < \" + ftp.storbinary(\"STOR \" + args[2], open(args[1], \"rb\")))  
    except FTP_ERRORS as e:
      print(\" < Error: \" + str(e).split(None, 1)[0])
ftp.close()
  " > "$DEPLOY_OUTPUT_ROOT/ftpclient.py"
  sleep "1s"

  echo_pkg "Deploying to \"$FTP_HOST\"..."
  py "$DEPLOY_OUTPUT_ROOT/ftpclient.py" > "$DEPLOY_OUTPUT_ROOT/latest-ftp-deploy.log"
  sleep "1s"

  if test $3 = "/keep" || test $3 = "-keep"; then
    echo_pkg "Keeping Python-FTP-Client."
    echo_pkg "Keeping local package-instance."
  else
    echo_pkg "Disposing Python-FTP-Client..."
    rm "$DEPLOY_OUTPUT_ROOT/ftpclient.py"
    echo_pkg "Disposing local package-instance..."
    rm -r "$DEPLOY_OUTPUT_ROOT/$SRV_PKG_ROOT"
  fi

  sleep "1s" 
  echo_pkg "Done." 

fi
