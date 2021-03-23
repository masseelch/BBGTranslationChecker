#!/bin/bash

# linux
printf "linux\n  build ... "
env CGO_ENABLED=0 GOOS=linux go build -a -o build/linux/bbg-translation-checker ./cmd
printf "\xE2\x9C\x94\n  compress ... "
zip -j build/linux.zip build/linux/bbg-translation-checker > /dev/null
printf "\xE2\x9C\x94\n\n"

# windows
printf "windows\n  build ... "
env CGO_ENABLED=0 GOOS=windows go build -a -o build/win/bbg-translation-checker.exe ./cmd
printf "\xE2\x9C\x94\n  compress ... "
zip -j build/win.zip build/win/bbg-translation-checker.exe > /dev/null
printf "\xE2\x9C\x94\n\n"

# mac
printf "mac\n  build ... "
env CGO_ENABLED=0 GOOS=darwin go build -a -o build/mac/bbg-translation-checker ./cmd
printf "\xE2\x9C\x94\n  compress ... "
zip -j build/mac.zip build/mac/bbg-translation-checker > /dev/null
printf "\xE2\x9C\x94\n"
