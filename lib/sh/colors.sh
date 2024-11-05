#!/bin/bash
#
# This is a library script to allow the use of colors in bash scripts.
# example usage:
# source "lib/sh/colors.sh"
# echo -e "${GREEN}This Text is green${NC}"
#

# shellcheck disable=SC2034
RED="\033[1;31m"
YELLOW="\033[1;33m"
GREEN="\033[1;32m"
BLUE="\033[1;34m"
CYAN="\033[1;36m"
NC="\033[0m"
NO_NEW_LINE='\033[0K\r'
