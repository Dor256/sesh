#!/bin/bash

LINEAR_PREFIX="ENG"

TICKET_NUMBER=$1
SESSION_NAME=$LINEAR_PREFIX$TICKET_NUMBER

WORKTREE_PATH="$PWD-$TICKET_NUMBER"

OPT_BRANCH_NAME=$2


if [ -n "$OPT_BRANCH_NAME" ]; then
    git worktree add -b $OPT_BRANCH_NAME $WORKTREE_PATH
else
    git worktree add --detach $WORKTREE_PATH main
fi

tmux new -s $SESSION_NAME -c $WORKTREE_PATH -n "OpenCode" "opencode; zsh" \; new-window -c $WORKTREE_PATH -n "Editor" "nvim .; zsh" \; new-window -c $WORKTREE_PATH -n "Terminal"

