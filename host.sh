
#!/bin/bash
echo "($1) Building Server.."
go build -o host
echo "($1) Launching Server.."
./host -cfg "$1.toml"
echo "($1) Disposing Server.."
rm host