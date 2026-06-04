#!/bin/bash

LINEAR_PREFIX="ENG"

TICKET_NUMBER=$1
SESSION_NAME=$LINEAR_PREFIX$TICKET_NUMBER

WORKTREE_PATH="${PWD%-[0-9]*}-$TICKET_NUMBER"

OPT_BRANCH_NAME=$2


if [ -n "$OPT_BRANCH_NAME" ]; then
    git worktree add -b "eng$TICKET_NUMBER/$OPT_BRANCH_NAME" $WORKTREE_PATH
else
    git worktree add --detach $WORKTREE_PATH main
fi


tmux new -d -s $SESSION_NAME -c $WORKTREE_PATH -n "OpenCode" "opencode; zsh"
tmux new-window -t $SESSION_NAME -c $WORKTREE_PATH -n "Editor" "nvim .; zsh"
tmux new-window -t $SESSION_NAME -c $WORKTREE_PATH -n "Terminal"

if [ -n "$TMUX" ]; then
    tmux switch -t "$SESSION_NAME"
else
    tmux attach -t "$SESSION_NAME"
fi
