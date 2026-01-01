#!/bin/bash

find . -type f ! -path '*/.git/*' -exec sh -c 'echo "$0" && cat "$0"' {} \;
