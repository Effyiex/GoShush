
#!/bin/bash
echo "($1) Building Server.."
go build -o srv-$1
echo "($1) Launching Server.."
./srv-$1 -cfg "$1.toml"
echo "($1) Disposing Server.."
rm srv-$1
