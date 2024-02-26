#!/bin/bash

check_aicommits_installed() {
    which aicommits >/dev/null 2>&1
}

regular_git_commit() {
    read -p "Enter your commit message: " commit_message
    git commit -m "$commit_message"
}

if check_aicommits_installed; then
    /usr/bin/expect ./aicommits.exp
else
    echo ""
    echo "aicommits not installed."
    echo "follow instruction from https://github.com/Nutlope/aicommits to install aicommits"
    echo "Falling back to regular git commit."
    echo ""
    regular_git_commit
fi
