FROM golang:1.23

WORKDIR /usr/src/app

services:
    web:
        build: .
        ports:
            - 3000:3000
        volumes:
            - .:/usr/src/app