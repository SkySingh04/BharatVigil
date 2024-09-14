#!/bin/bash

installBharatVigil() {
    echo "Installing BharatVigil..."

    # Clone the repository containing docker-compose.yml
    git clone https://github.com/SkySingh04/BharatVigil.git /tmp/BharatVigilApp || {
        echo "Failed to clone the Bharat Vigil repository."
        exit 1
    }

    cd /tmp/BharatVigilApp || exit

    echo "Running Docker Compose..."

    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        echo "Docker Compose not found. Please install Docker Compose and try again."
        exit 1
    fi

    # Start the app using Docker Compose
    docker-compose up -d

    # Create an alias for starting the app with Docker Compose
    alias_cmd="alias bharatvigil-start='cd /tmp/BharatVigilApp && docker-compose up -d'"
    
    # Add alias to shell configuration file based on the user's shell
    current_shell="$(basename "$SHELL")"
    if [[ "$current_shell" == "zsh" ]]; then
        if ! grep -q "alias bharatvigil-start=" "$HOME/.zshrc"; then
            echo "$alias_cmd" >> "$HOME/.zshrc"
        fi
        source "$HOME/.zshrc"
    elif [[ "$current_shell" == "bash" ]]; then
        if ! grep -q "alias bharatvigil-start=" "$HOME/.bashrc"; then
            echo "$alias_cmd" >> "$HOME/.bashrc"
        fi
        source "$HOME/.bashrc"
    else
        if ! grep -q "alias bharatvigil-start=" "$HOME/.profile"; then
            echo "$alias_cmd" >> "$HOME/.profile"
        fi
        source "$HOME/.profile"
    fi

    echo "Alias 'bharatvigil-start' added! Use 'bharatvigil-start' to run the app."
}

installBharatVigil
