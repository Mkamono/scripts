#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title ozo-control
# @raycast.mode fullOutput

# Optional parameters:
# @raycast.icon 🦆

# Documentation:
# @raycast.argument1 { "type": "dropdown", "optional":true, "placeholder": "操作", "data": [{"title":"出勤", "value":"i"}, {"title":"退勤", "value":"o"}, {"title":"公休入力", "value":"r"}, {"title":"設定ファイル削除", "value":"c"}]}
# @raycast.author Mkamono
# @raycast.authorURL https://github.com/Mkamono

ozo-control $1
