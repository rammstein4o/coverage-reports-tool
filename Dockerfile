###
# Copyright (c) 2020 by PROS, Inc.  All Rights Reserved.
# This software is the confidential and proprietary information of
# PROS, Inc. ("Confidential Information").
# You may not disclose such Confidential Information, and may only
# use such Confidential Information in accordance with the terms of
# the license agreement you entered into with PROS.
###
ARG IMAGE_GO=golang:1.14-alpine
FROM ${IMAGE_GO}

WORKDIR ${GOPATH}/src/coverage-reports-tool

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o /coverage-reports-tool
