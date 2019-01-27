#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Supply the go source file you wish to visualize as a command line argument."
    exit 1
fi

imageResults="$(docker images | grep gotrace)"

if [ ${#imageResults} -eq 0 ]; then
    printf "\nWelcome to the go runtime visualizer. A docker image will be constructed from source.\n"
    printf "Please be patient, this may take several minutes.\n\n"
    sleep 3
    docker build -t "divan/golang:gotrace" -f gotrace/runtime/Dockerfile gotrace/runtime
fi

printf "Compiling from source file '%s'..." $1 

#TODO: get it out of examples, and allow a file (and the output) to be named by the user.
docker run --rm -it \
	-e GOOS=darwin \
	-v $(pwd)/:/src divan/golang:gotrace \
		go build -o /src/gotrace/binary /src/$1

if [ "$?" -ne 0 ]; then
    printf "Failed to compile '%s'. Exiting.\n" $1
    exit 1
fi

printf "\tdone.\n"
sleep 0.5
printf "Running program and collecting trace data..."

gotrace/./binary 2> gotrace/trace.out > /dev/null

if [ "$?" -ne 0 ]; then
    printf "Runtime error during execution of '%s'. Exiting.\n" $1
    exit 1
fi

printf "\tdone.\n"
sleep 0.5
printf "\nStarting visualization.\n\n"
sleep 2

gotrace/./gotrace gotrace/trace.out