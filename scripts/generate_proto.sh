#!/bin/bash
# Run this script to regenerate gRPC .pb.go files

set -e

PROTO_DIR="./internal/proto"
PROTO_FILE="${PROTO_DIR}/monitor.proto"

protoc --go_out=${PROTO_DIR} \
       --go-grpc_out=${PROTO_DIR} \
       --proto_path=${PROTO_DIR} \
       ${PROTO_FILE}
