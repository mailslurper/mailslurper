#!/usr/bin/env bash

#
# Usage:
# ./package.sh -v 1.9.1 -t osx or ./package.sh --version 1.9.1 --target osx
#
# Targets:
# 		osx, linux, windows
#
# Note: This is setup to deploy one at a time. I can't seem
# to get the cross-compile working, so I'm compiling on each platform
# individually. The windows build is assuming Cygwin or some other
# shell environment with Unixy Bash-style command support.
#

if [ "$#" -lt 4 ]; then
	echo "Please provide a version number and target for this deployment package. i.e. ./package.sh -v 1.0.0 -t osx"
	exit 1
fi

#
# Get the command line args
#
while [[ $# > 1 ]]
do
	key="$1"

	case $key in
		-v|--version)
			VERSION="$2"
			shift
			;;

		-t|--target)
			TARGET="$2"
			shift
			;;
	esac

	shift
done

echo "Packaging MailSlurper v$VERSION for $TARGET..."
#exit 1

ZIPFILENAME="mailslurper-$VERSION-$TARGET.zip"

#
# Generate compiled assets
#
cd ../cmd/mailslurper
go generate
cd ../../bin

#
# Create deploy directory
#
if [ -d "./deploy" ]; then
	#
	# Remove previous builds
	#
	rm ./deploy/*
fi

if [ ! -d "./deploy" ]; then
	mkdir deploy
fi

#
# Copy non-executable assets to the deploy folder
#
cp ../LICENSE ./deploy
cp ../cmd/mailslurper/config.json ./deploy
cp ../logo/logo.png ./deploy
cp ../README.md ./deploy
cp ./create-mssql.sql ./deploy
cp ./create-mysql.sql ./deploy

#
# Compile for the various targets. Copy to the deploy folder
#

# OSX
if [ $TARGET = "osx" ]; then
	cd ../cmd/mailslurper
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper ../../bin/deploy

	cd ../createcredentials
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w"
	mv ./createcredentials ../../bin/deploy

	cd ../../bin/deploy
	zip -r -X $ZIPFILENAME *
	cd ..
fi

# Linux
if [ $TARGET = "linux" ]; then
	cd ../cmd/mailslurper
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper ../../bin/deploy

	cd ../createcredentials
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
	mv ./createcredentials ../../bin/deploy

	cd ../../bin/deploy
	zip -r -X $ZIPFILENAME *
	cd ..
fi

# Windows
if [ $TARGET = "windows" ]; then
	cd ../cmd/mailslurper
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper.exe ../../bin/deploy

	cd ../createcredentials
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"
	mv ./createcredentials.exe ../../bin/deploy

	cd ../../bin/deploy
	zip -r -X $ZIPFILENAME *
	cd ..
fi

#
# Clean up
#
rm ./deploy/LICENSE
rm ./deploy/config.json
rm ./deploy/logo.png
rm ./deploy/README.md
rm ./deploy/create-mssql.sql
rm ./deploy/create-mysql.sql

if [ -f "./deploy/mailslurper.exe" ]; then
	rm ./deploy/mailslurper.exe
fi

if [ -f "./deploy/createcredentials.exe" ]; then
	rm ./deploy/createcredentials.exe
fi

if [ -f "./deploy/mailsurper" ]; then
	rm ./deploy/mailslurper
fi

if [ -f "./deploy/createcredentials" ]; then
	rm ./deploy/createcredentials
fi

echo "Package complete."
