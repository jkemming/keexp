[private]
default:
  @just --list

build:
  go build jkemming.com/keexp

update:
  go get -u
  go mod tidy
