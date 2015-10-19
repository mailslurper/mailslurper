#!/bin/bash

COUNTER=0
LIMIT=10

while [ $COUNTER -lt $LIMIT ]; do
	python ./send-mail-test.py
	let COUNTER=COUNTER+1
done
