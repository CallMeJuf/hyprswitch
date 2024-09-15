package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	MOVEFOCUS  = "moveFocus"
	MOVEWINDOW = "moveWindow"
	MOVE_LEFT  = "l"
	MOVE_RIGHT = "r"
	MOVE_UP    = "u"
	MOVE_DOWN  = "d"
)

type Window struct {
	Address          string    `json:"address"`
	Mapped           bool      `json:"mapped"`
	Hidden           bool      `json:"hidden"`
	At               [2]int    `json:"at"`
	Size             [2]int    `json:"size"`
	Workspace        Workspace `json:"workspace"`
	Floating         bool      `json:"floating"`
	Pseudo           bool      `json:"pseudo"`
	Monitor          int       `json:"monitor"`
	Class            string    `json:"class"`
	Title            string    `json:"title"`
	InitialClass     string    `json:"initialClass"`
	InitialTitle     string    `json:"initialTitle"`
	PID              int       `json:"pid"`
	Xwayland         bool      `json:"xwayland"`
	Pinned           bool      `json:"pinned"`
	Fullscreen       int       `json:"fullscreen"`
	FullscreenClient int       `json:"fullscreenClient"`
	Grouped          []string  `json:"grouped"`
	Tags             []string  `json:"tags"`
	Swallowing       string    `json:"swallowing"`
	FocusHistoryID   int       `json:"focusHistoryID"`
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments provided, Example: ./switch moveFocus left")
		return
	}

	validCommands := []string{MOVEFOCUS, MOVEWINDOW}

	command := os.Args[1]

	validDirections := []string{MOVE_LEFT, MOVE_RIGHT, MOVE_UP, MOVE_DOWN}

	direction := os.Args[2]

	// Check if command is valid
	valid := false
	for _, validCommand := range validCommands {
		if command == validCommand {
			valid = true
			break
		}
	}

	if !valid {
		fmt.Println("Invalid command, valid commands are:", validCommands)
		return
	}

	// Check if direction is valid
	valid = false
	for _, validDirection := range validDirections {
		if direction == validDirection {
			valid = true
			break
		}
	}

	if !valid {
		fmt.Println("Invalid direction, valid directions are:", validDirections)
		return
	}

	cmd := exec.Command("hyprctl", "-j", "activewindow")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var workspace Window
	err = json.Unmarshal(output, &workspace)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	inGroup, err := shouldMoveWithinGroup(workspace, direction)
	if err != nil {
		fmt.Println("Error determining position within group:", err)
		return
	}

	if command == MOVEFOCUS {
		moveFocus(direction, inGroup)
	} else if command == MOVEWINDOW {
		moveWindow(direction, inGroup)
	}

}

func moveWindow(direction string, moveWithinGroup bool) {
	if moveWithinGroup {
		if direction == MOVE_LEFT {
			direction = "b"
		} else if direction == MOVE_RIGHT {
			direction = "f"
		}
	}

	hyprCommand := ""
	if moveWithinGroup {
		hyprCommand = "movegroupwindow"
	} else {
		hyprCommand = "movewindoworgroup"
	}

	move(hyprCommand, direction)
}

func moveFocus(direction string, moveWithinGroup bool) {

	if moveWithinGroup {
		if direction == MOVE_LEFT {
			direction = "b"
		} else if direction == MOVE_RIGHT {
			direction = "f"
		}
	} else {
		if direction == MOVE_LEFT {
			direction = "l"
		} else if direction == MOVE_RIGHT {
			direction = "r"
		}
	}

	hyprCommand := ""
	if moveWithinGroup {
		hyprCommand = "changegroupactive"
	} else {
		hyprCommand = "movefocus"
	}

	move(hyprCommand, direction)
}

func move(hyprCommand string, directionArg string) {
	fmt.Println("Moving", hyprCommand, "to", directionArg)
	cmd := exec.Command("hyprctl", "dispatch", hyprCommand, directionArg)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))

}

func shouldMoveWithinGroup(window Window, direction string) (bool, error) {
	if direction == MOVE_UP || direction == MOVE_DOWN || len(window.Grouped) == 0 {
		return false, nil
	}
	positionWithinGroup := -1
	for i, groupedWindow := range window.Grouped {
		if groupedWindow == window.Address {
			positionWithinGroup = i
			break
		}
	}
	if positionWithinGroup == -1 {
		return false, fmt.Errorf("window not found in group")
	}

	if direction == MOVE_LEFT && positionWithinGroup == 0 {
		return false, nil
	}

	if direction == MOVE_RIGHT && positionWithinGroup == len(window.Grouped)-1 {
		return false, nil
	}

	return true, nil
}
