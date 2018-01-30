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
rm ./www/www.go
go generate

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
cp ./LICENSE ./deploy
cp ./config.json ./deploy
cp ./MailSlurperLogo.ico ./deploy
cp ./MailSlurperLogo.png ./deploy
cp ./README.md ./deploy
cp -R ./scripts ./deploy

#
# Compile for the various targets. Copy to the deploy folder
#

# OSX
if [ $TARGET = "osx" ]; then
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper ./deploy

	cd deploy
	zip -r -X $ZIPFILENAME *

	rm ./mailslurper
	cd ..
fi

# Linux
if [ $TARGET = "linux" ]; then
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper ./deploy

	cd deploy
	zip -r -X $ZIPFILENAME *

	rm ./mailslurper
	cd ..
fi

# Windows
if [ $TARGET = "windows" ]; then
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"
	mv ./mailslurper.exe ./deploy

	cd deploy
	zip -r -X $ZIPFILENAME *

	rm ./mailslurper.exe
	cd ..
fi

#
# Clean up
#
rm ./deploy/LICENSE
rm ./deploy/config.json
rm ./deploy/MailSlurperLogo.ico
rm ./deploy/MailSlurperLogo.png
rm ./deploy/README.md
rm -R ./deploy/scripts

echo "Package complete."
