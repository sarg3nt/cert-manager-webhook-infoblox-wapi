#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

main() {
  # Set HOME to the user's home directory if it is not set
  : "${HOME:=$(eval echo "~${USER}")}"
  if [[ -z "${HOME}" ]]; then
    echo "Error: HOME environment variable is not set."
    exit 1
  fi
  
  get_latest_dev_container_version
  create_required_folders
}

#######################################
# Get the latest version of the dev container
# Arguments:
#   None
#######################################
get_latest_dev_container_version() {
  echo "************** Pull the latest version of the container ******************"
  docker pull ghcr.io/sarg3nt/go-dev-container:latest
  echo ""
}

#######################################
# Create any required missing folders if they do not exist.
# Globals:
#   HOME
# Arguments:
#   None
#######################################
create_required_folders() {
  echo "************** Create any required missing folders if they do not exist ******************"

  local directories_created=false
  if [[ ! -d "${HOME}/.docker" ]]; then
    echo "You did not have a .docker folder in your home directory, creating."
    echo "Docker may will not work properly without this folder."
    echo ""
    mkdir -p "${HOME}/.docker"
    directories_created=true
  fi

  if [[ ! -d "${HOME}/.kube" ]]; then
    echo "You did not have a .kube folder in your home directory, creating."
    echo "Kubectl and k9s will not work without this folder."
    echo ""
    mkdir -p "${HOME}/.kube"
    directories_created=true
  fi

  if [[ ! -d "${HOME}/.config/k9s" ]]; then
    echo "You did not have a .k9s folder in your home directory, creating."
    echo "K9s will use a local config."
    echo ""
    mkdir -p "${HOME}/.config/k9s"
    mkdir -p "${HOME}/.local/share/k9s"
    directories_created=true
  fi

  if [[ ! -d "${HOME}/.ssh" ]]; then
    echo "----------- WARNING: You did not have an '.ssh' folder in your home directory -----------"
    echo "  We created a '.ssh' folder for you so the mount didn't fail but you need to run 'ssh-keygen' to finish setup."
    mkdir -p "${HOME}/.ssh"
    directories_created=true
  fi

  if [[ "$directories_created" = false ]]; then
    echo "All required directories already exist."
  fi
  echo ""
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
