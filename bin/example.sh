#!/bin/bash

phrase="Sturdy Octo Goggles"
emoji=("🫎" "🥽" "🎨")

selection=${emoji[$RANDOM % ${#emoji[@]}]}
brainstorm="$phrase $selection"
cowthink -f moose "$brainstorm"
