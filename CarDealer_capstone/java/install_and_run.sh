#!/bin/sh
usage()
{
  echo "Usage: $0 OPTION" >&2
  echo "Where  OPTION := { install | run }"
  exit 1
}

if [ "$#" -gt 1 -o "$#" -eq 0 ]; then
  usage
fi
#1 install; 2 run sample
OPTION=2

check_opt()
{
  if [ "$1" = "install" ]; then
    OPTION=1
  elif [ "$1" = "run" ]; then
    OPTION=2
  else
    usage
  fi
}

install()
{
  echo "Maven install the project"
  mvn install
  mvn install:install-file -Dfile=./libs/grpc-netty-1.11.0.jar -DgroupId=io.grpc -DartifactId=grpc-netty -Dversion=1.11.0 -Dpackaging=jar
}

run()
{
  echo "Run the CarDemo sample"
  mvn exec:java -Dexec.mainClass="main.CarDemo"  
}

if [ "$#" -eq 1 ]; then
  check_opt $1
fi

if [ "$OPTION" -eq 1 ]; then
  install
else
  run
fi

exit 0


