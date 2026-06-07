#!/bin/bash

LINEAR_PREFIX="ENG"

TICKET_NUMBER=$1
SESSION_NAME=$LINEAR_PREFIX$TICKET_NUMBER
WORKTREE_DIR="$HOME/dev/.worktrees"
CWD=$(basename "$PWD")
WORKTREE_PATH="$WORKTREE_DIR/${CWD%-[0-9]*}-$TICKET_NUMBER"

[ -z "$TICKET_NUMBER" ] && exit 1;

pick-wt() {
    # 1. Format text cleanly using spaces. Column 1 = Path, Column 2 = Branch
    # 2. Tell fzf to only display the second column (--with-nth=2)
    # 3. Print the raw string on selection, OR the query on zero matches
    local choice exit_code
    choice=$(
        git worktree list --porcelain \
        | awk '
            /^worktree / { path = $2 }
            /^branch / { sub("refs/heads/", "", $2); print path, $2 }
        '| fzf --height=40% \
            --reverse \
            --border \
            --with-nth=2 \
            --print-query
    )
    exit_code=$?

    if [ "$exit_code" -eq 130 ]; then
        return 1
    fi

    if [ "$exit_code" -eq 0 ] && [ -n "$choice" ]; then
        local selected_item=$(tail -n 1 <<< "$choice")
        local res="${selected_item%% *}"
        echo "$res:true"
    fi

    if [ "$exit_code" -eq 1 ]; then
        echo "$choice:false"
    fi
}

selection=$(pick-wt)
is_existing=$(awk -F':' '{print $2}' <<< "$selection")
branch_name=$(awk -F':' '{print $1}' <<< "$selection")

if [ "$is_existing" != "true" ]; then
    git worktree add -b "eng$TICKET_NUMBER/$branch_name" $WORKTREE_PATH origin/main
fi

if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
    if [ -n "$TMUX" ]; then
        tmux switch -t "$SESSION_NAME"
    else
        tmux attach -t "$SESSION_NAME"
    fi
else
    tmux new -d -s $SESSION_NAME -c $WORKTREE_PATH -n "OpenCode" "opencode; zsh"
    tmux new-window -t $SESSION_NAME -c $WORKTREE_PATH -n "Editor" "nvim .; zsh"
    tmux new-window -t $SESSION_NAME -c $WORKTREE_PATH -n "Terminal"

    if [ -n "$TMUX" ]; then
        tmux switch -t "$SESSION_NAME"
    else
        tmux attach -t "$SESSION_NAME"
    fi
fi
