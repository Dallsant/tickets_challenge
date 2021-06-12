FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y software-properties-common git  
    # psql -U postgres && \
    # create database ticketsdb;
WORKDIR /tickets_challenge






