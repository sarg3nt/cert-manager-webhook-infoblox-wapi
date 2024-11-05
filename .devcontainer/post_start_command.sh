#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

main() {

  # $HOME is not set in certain situations and since we will always know the home user in the dev container we hard code it.
  HOME="/home/vscode"

  git_update_diff_tool
  copy_ssh_folder
  copy_kube_config
  copy_k9s_config
  copy_docker_config
}

#######################################
# Update the git config to make vscode default merge/diff tool and add some aliases
# Arguments:
#   None
#######################################
# NOTE Required to use double quotes here because if we use single quotes the ' is removed from the aliases.
# shellcheck disable=SC2016
git_update_diff_tool() {
  echo "************** Update .gitconfig to use vscode ******************"
  grep -qxF '[merge]' ~/.gitconfig || echo "
[merge]
  tool = vscode
[mergetool \"vscode\"]
  cmd = code --wait \$MERGED
[diff]
  tool = vscode
[difftool \"vscode\"]
  cmd = code --wait --diff \$LOCAL \$REMOTE
[core]
  editor = \"code --wait\"
  pager = bat
 [alias]
   s = !git for-each-ref --format='%(refname:short)' refs/heads | fzf | xargs git switch
   c = !git for-each-ref --format='%(refname:short)' refs/heads | fzf | xargs git switch
" >>~/.gitconfig
  echo "  - Updated."
  echo ""
}

#######################################
# Copy in the user's `~/.ssh` so modifications to it in the devcontainer do not affect the host's version.
# Globals:
#   HOME
# Arguments:
#   None
#######################################
copy_ssh_folder() {
  echo "************** SSH Folder Setup ******************"
  REMOTE_CONFIG="${HOME}/.ssh-localhost"
  CONFIG="${HOME}/.ssh"
  if [[ -d "$REMOTE_CONFIG" ]]; then
    echo "  - Remote ssh folder detected, copying in."
    rm -rf "$CONFIG"
    mkdir -p "${CONFIG}" >/dev/null 2>&1
    cp -R "${REMOTE_CONFIG}/." "${CONFIG}/"
    sudo chown -R "$(id -u)" "${CONFIG}"
  else
    echo "  - No remote ssh folder detected. Could not copy in.  You probably need to set up SSH."
  fi
  echo ""
}

#######################################
# Copy in the user's `~/.kube/config` so modifications to it in the devcontainer do not affect the host's version.
# Globals:
#   HOME
# Arguments:
#   None
#######################################
copy_kube_config() {
  echo "************** k8s Config Setup ******************"
  REMOTE_CONFIG="${HOME}/.kube-localhost"
  CONFIG="${HOME}/.kube"
  if [[ -d "$REMOTE_CONFIG" ]]; then
    echo "  - Remote k8s config detected, copying in."
    rm -rf "$CONFIG"
    mkdir -p "${CONFIG}" >/dev/null 2>&1
    cp -R "${REMOTE_CONFIG}/." "${CONFIG}/"
    sudo chown -R "$(id -u)" "${CONFIG}"
  else
    echo "  - No remote k8s config detected, using defaults."
  fi
  echo ""
}

#######################################
# Copy in the user's `~/.config/k9s` if it exists, otherwise use local one.
# Globals:
#   HOME
# Arguments:
#   None
#######################################
copy_k9s_config() {
  echo "************** k9s Config Setup ******************"
  REMOTE_CONFIG="${HOME}/.config/k9s-localhost"
  CONFIG="${HOME}/.config/k9s"
  if [[ -d "$REMOTE_CONFIG" ]]; then
    echo "  - Remote k9s config detected, copying in."
    rm -rf "$CONFIG"
    mkdir -p "${CONFIG}" >/dev/null 2>&1
    cp -R "${REMOTE_CONFIG}/." "${CONFIG}/"
    sudo chown -R "$(id -u)" "${CONFIG}"
  else
    echo "  - No remote k9s config detected, using defaults."
  fi

  REMOTE_SHARE=${HOME}/.local/share/k9s-localhost
  SHARE="${HOME}/.local/share/k9s"
  if [[ -d "$REMOTE_SHARE" ]]; then
    echo "  - Remote k9s share detected, copying in."
    rm -rf "$SHARE"
    mkdir -p "$SHARE" >/dev/null 2>&1
    cp -R "${REMOTE_SHARE}/." "${SHARE}/"
    sudo chown -R "$(id -u)" "${SHARE}"
  else
    echo "  - No remote k9s share detected, using defaults."
  fi
  echo ""
}

#######################################
# Copy the user's `~/.docker/` so we can modify the config.json to remove any credStores that might be configured.
# Globals:
#   HOME
# Arguments:
#   None
#######################################
copy_docker_config() {
  echo "************** Docker Config Setup ******************"
  REMOTE_CONFIG="${HOME}/.docker-localhost"
  CONFIG="${HOME}/.docker"
  if [[ -d "$REMOTE_CONFIG" ]]; then
    echo "  - Remote Docker config detected, copying in."
    rm -rf "$CONFIG"
    mkdir -p "${CONFIG}" >/dev/null 2>&1
    cp -R "${REMOTE_CONFIG}/." "${CONFIG}/"
    sudo chown -R "$(id -u)" "${CONFIG}"

    if [[ $(jq '.credsStore' "${CONFIG}/config.json") != "null" ]]; then
      echo "  - Removing Windows credsStore from .docker/config.json."
      json_revision=$(jq 'del(.credsStore)' "${CONFIG}/config.json" 2>/dev/null)
      # -n tests if the length of json_revision is nonzero
      [ -n "${json_revision}" ] && echo -e "${json_revision}" >"${CONFIG}/config.json" 2>/dev/null
    else
      echo "  - No Windows credstore detected, skipping removal."
    fi

  else
    echo "  - No remote Docker config detected, using defaults."
  fi
  echo ""
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
