# SPADE: **S**elective and **PA**rtial **DE**cryption using Functional Encryption
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/hosseinabdinf/SPADE)](https://goreportcard.com/report/github.com/hosseinabdinf/SPADE)

SPADE is an FE-based scheme that enables running selective and partial decryption queries over
vectors of ciphertexts.

To see two real-world applications of SPADE, please have a look at the usescases directory.

## Notes

1. Curators are required to specify both the number of users and the maximum number of entries
   allowed for each plaintext vector. Users must encrypt their data as a vector of integers,
   ensuring that none of the values are **zero**.
2. Be careful when defining the users' data vector using `make()`;
   This method assigns the **zero** value to the elements during initialization.

## Changing the protobuf structure

If you want to modify the protobuf structure, please first change the following file:

    ./spadeproto/spade.proto

and then run the proto compiler command as follows to generate the new protobuf files:

    protoc --go_out=. --go-grpc_out=. spade.proto 

## Testing Instruction

    go test     

## Benchmarking Instruction

    go test -benchtime=10x -bench=BenchmarkSpade -benchmem -run=^$

