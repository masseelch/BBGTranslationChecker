#!/bin/bash

# linux
env CGO_ENABLED=0 GOOS=linux go build -a -o build/linux/bbg-translation-checker ./cmd
zip build/linux.zip build/linux/bbg-translation-checker

# windows
env CGO_ENABLED=0 GOOS=windows go build -a -o build/win/bbg-translation-checker ./cmd
zip build/win.zip build/win/bbg-translation-checker

# mac
env CGO_ENABLED=0 GOOS=darwin go build -a -o build/mac/bbg-translation-checker ./cmd
zip build/mac.zip build/mac/bbg-translation-checker
