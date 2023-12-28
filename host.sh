
#!/bin/bash
echo "($1) Building Server.."
<<<<<<< HEAD
go build -o srv-$1
echo "($1) Launching Server.."
./srv-$1 -cfg "$1.toml"
echo "($1) Disposing Server.."
rm srv-$1
=======
go build -o host
echo "($1) Launching Server.."
./host -cfg "$1.toml"
echo "($1) Disposing Server.."
rm host
>>>>>>> fe4a0115791ac06ec90bd55a1897fafbfbefddc3
