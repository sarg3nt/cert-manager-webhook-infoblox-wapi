#!/bin/bash
#
# This script will start the dev container and open an interactive prompt into it.
#

set -euo pipefail
IFS=$'\n\t'

source "lib/sh/colors.sh"

# Check if user is using Docker Deskop for Windows or the native Docker Engine.
main() {

  # What quick two letter command do we want to use for this dev container / project.
  # If this is empty, no command will be installed.
  # Example: qq
  local docker_exec_command="cm"
  # Name of the project folder
  local project_name="cert-manager-webhook-infoblox-wapi"
  # Name of the container
  local container_name="${project_name}"
  # User being created in the container
  local container_user="vscode"

  local is_windows_docker=true
  if [[ $(docker version) != *"Server: Docker Desktop"* ]]; then
    is_windows_docker=false
  fi

  if [[ -n "$docker_exec_command" ]]; then
    # Add ${docker_exec_command} command to open the dev container to a users .bashrc file if it is not already there.
    if ! grep -F "${docker_exec_command} ()" "${HOME}/.bashrc" >/dev/null 2>&1 && [ -f "${HOME}/.bashrc" ]; then
      echo -e "${docker_exec_command} (){\n\
        docker exec -it -u ${container_user} -w /workspaces/${project_name} ${container_name} zsh\n\
      }" >>"${HOME}/.bashrc"
      echo -e "${GREEN}Created \"${docker_exec_command}\" command in your ${HOME}/.bashrc file.${NC}"
      # shellcheck source=/dev/null
      source "${HOME}/.bashrc"
    fi

    # Add ${docker_exec_command} command to open the dev container to a users .zshrc file if it is not already there.
    if ! grep -F "${docker_exec_command} ()" "${HOME}/.zshrc" >/dev/null 2>&1 && [ -f "${HOME}/.zshrc" ]; then
      echo -e "${docker_exec_command} (){\n\
        docker exec -it -u ${container_user} -w /workspaces/${project_name} ${container_name} zsh\n\
      }" >>"${HOME}/.zshrc"
      echo -e "${GREEN}Created \"${docker_exec_command}\" command in your ${HOME}/.zshrc file.${NC}"
      # shellcheck source=/dev/null
      source "${HOME}/.zshrc"
    fi
  fi

  # If the dev container is running, we assume VSCode is also running. If VSCode is not running, then open it.
  if ! docker ps | grep ${container_name}; then
    # Check if we are using "Docker for Windows". `devcontainer open` only works with "Docker for Windows".
    if [[ "${is_windows_docker}" = true ]]; then
      # Check if devcontainer cli is installed, if not let the user know and open VSCode with `code .`.
      if ! which devcontainer; then
        echo -e "${YELLOW}The 'devcontainer' command was not found."
        echo ""
        echo "Opening VS Code outside of the dev container."
        echo "Once VS Code is open, press F1 or ctrl+shift+p to open the command pallet."
        echo "Type 'devcontainer' then select 'Remote Containers: Install devcontainer CLI' from the list."
        echo -e "Windows users will need to restart their computer for this to take effect.${NC}"

        code .
      else
        # Open the dev container in another shell so it doesn't hang on the command line for 45 seconds for some unknown reason.
        (devcontainer open &) >/dev/null 2>&1
      fi
    fi

    # Looks like we are using native Docker on linux, so just do a `code .` as `devcontainer open` is not supported.
    if [[ "${is_windows_docker}" = false ]]; then
      code .
    fi
  fi

  # Here we check if the dev container has started yet.
  # Wait a max of 600 seconds (10 minutes).
  local max_wait=600
  local count=0
  local spin=("-" "\\" "|" "/")
  local rot=0
  # If not then docker_id will be empty and the while loop kicks in.
  local docker_id=""
  docker_id=$(docker container ls -f name=${container_name} -q)
  if [ -z "$docker_id" ]; then
    echo -e "${BLUE}Waiting up to 10 minutes for the dev container to start ${NO_NEW_LINE}"
  fi

  while [ -z "$docker_id" ] && ((count < max_wait)); do
    # Sleep for one second then try and get the docker_id of the dev container again.
    # Once we can get the id, the loop quits and we run the last line below, exe'ing into the container.
    sleep 1s
    docker_id=$(docker container ls -f name=${container_name} -q)
    count=$((count + 1))

    if ((count == 20)); then
      echo -ne "\b"
      echo -e "${YELLOW}The dev container is taking a while to start, VS Code could be downloading a new version or you may need to manually open it from within VS Code.${BLUE}"
    fi

    if ((count >= max_wait)); then
      exit 0
    fi
    echo -ne "\b${spin[$rot]}"
    if ((rot >= 3)); then
      rot=0
    else
      rot=$((rot + 1))
    fi
  done
  echo -ne "\b"
  echo -e "${GREEN}Dev container started, execing into it.${NC}"
  if [[ -n "$docker_exec_command" ]]; then
    echo -e "${BLUE}You can use the \"${docker_exec_command}\" command to exec into the dev container from another terminal.${NC}"
  fi
  docker exec -u "${container_user}" -w /workspaces/${project_name} -it ${container_name} zsh
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
