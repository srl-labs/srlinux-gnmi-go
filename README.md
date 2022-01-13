Example code demonstrating how [karimra/gnmic](https://github.com/karimra/gnmic) can be used as an API to interface with SR Linux NOS via gNMI.

Deploy an srlinux node with `containerlab dep -t srl.clab.yml`. Then run `go run .` to execute the gNMI Get request and query management ipv4 address of a running node.