#!/bin/bash

micro --enable_stats --registry=consul --registry_address="127.0.0.1:8500" --server_name="com.pingflow.web" web --address="0.0.0.0:9000" --namespace="com.pingflow.api"
