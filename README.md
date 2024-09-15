# hyprswitch
Hacky lil Go script to move active windows and focus in [hyprland](https://github.com/hyprwm/Hyprland) similar to [i3](https://github.com/i3/i3)'s handling grouped (tabbed) windows. 

For example, if you move focus left, this script will execute the following logic:

1. If active window is NOT in a group.
   - Execute dispatcher `movefocus l`
2. If active window IS in a group.
   - If active window IS the leftmost window in group.
     - Execute dispatcher `movefocus l`
   - If active window is NOT the leftmost window in group.
     - Execute dispatcher `changegroupactive b`


### Usage

Modify your `hyprland.conf` binds for moving windows and focus to exec the hyprswitch binary:

``` Properties
# Example
# Move focus with mainMod + arrow keys
bind = $mainMod, left, exec, /path/to/hyprswitch moveFocus l
bind = $mainMod, right, exec, /path/to/hyprswitch moveFocus r
bind = $mainMod, up, exec, /path/to/hyprswitch moveFocus u
bind = $mainMod, down, exec, /path/to/hyprswitch moveFocus d

# Move windows with mainMod + SHIFT + arrow keys
bind = $mainMod SHIFT, left, exec, /path/to/hyprswitch moveWindow l
bind = $mainMod SHIFT, right, exec, /path/to/hyprswitch moveWindow r
bind = $mainMod SHIFT, up, exec, /path/to/hyprswitch moveWindow u
bind = $mainMod SHIFT, down, exec, /path/to/hyprswitch moveWindow d
```

### Notes - Q/A

Script is a little hacky and uses `hyprctl` calls to do its work. It doesn't utilize the window address when moving windows/focus, only when determining group position. This means in rare cases if your computer is slow or you're moving *blazingly fast*, it's possible this script doesn't accurately determine if your current window is at the beginning or the end of the group.

- Why Go?
  - Cute Gopher
- Why not implement in hyprland directly?
  - Cute Gopher (I'll look into contributing directly to hyprland when I have time to not whip up something hacky)