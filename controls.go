package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ButtonAction represents the action a button performs
type ButtonAction string

const (
	AddSphereAction   ButtonAction = "add_sphere"
	AddSpriteAction   ButtonAction = "add_sprite"
	ClearAllAction    ButtonAction = "clear_all"
	PauseResumeAction ButtonAction = "pause_resume"
	ResetAction       ButtonAction = "reset"
	GravityAction     ButtonAction = "gravity"
	BounceAction      ButtonAction = "bounce"
	SizeAction        ButtonAction = "size"
	ColorAction       ButtonAction = "color"
)

// Button represents an interactive button
type Button struct {
	Label   string
	Action  ButtonAction
	Active  bool
	Focused bool
	Width   int
	Height  int
}

// ButtonMsg is sent when a button is activated
type ButtonMsg struct {
	Action ButtonAction
}

// ControlPanel manages the interactive control panel with responsive layouts
type ControlPanel struct {
	buttons      []Button
	focused      int
	width        int
	height       int
	buttonStyles ButtonStyles

	// Parameter display values
	gravityText string
	sizeText    string
	colorText   string

	// Responsive layout mode
	compactMode      bool
	ultraCompactMode bool
}

// ButtonStyles defines the visual styles for buttons with enhanced polish
type ButtonStyles struct {
	Normal  lipgloss.Style
	Focused lipgloss.Style
	Active  lipgloss.Style
	Hover   lipgloss.Style
}

// NewControlPanel creates a new interactive control panel
func NewControlPanel(width, height int) *ControlPanel {
	// Simplified button styles for horizontal layout
	buttonStyles := ButtonStyles{
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E6ED")).
			Background(lipgloss.Color("#2C3E50")).
			Padding(0, 1).
			MarginRight(1),
		Focused: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3498DB")).
			Padding(0, 1).
			MarginRight(1).
			Bold(true),
		Active: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#E74C3C")).
			Padding(0, 1).
			MarginRight(1).
			Bold(true),
		Hover: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F9FA")).
			Background(lipgloss.Color("#5DADE2")).
			Padding(0, 1).
			MarginRight(1).
			Bold(true),
	}

	// Create buttons - core 5 buttons as per original design
	buttons := []Button{
		{Label: "Add Sphere", Action: AddSphereAction, Width: 12},
		{Label: "Add Sprite", Action: AddSpriteAction, Width: 12},
		{Label: "Clear All", Action: ClearAllAction, Width: 11},
		{Label: "Pause", Action: PauseResumeAction, Width: 7},
		{Label: "Reset", Action: ResetAction, Width: 7},
	}

	return &ControlPanel{
		buttons:      buttons,
		focused:      0,
		width:        width,
		height:       height,
		buttonStyles: buttonStyles,
	}
}

// Init implements tea.Model interface
func (cp *ControlPanel) Init() tea.Cmd {
	return nil
}

// Update handles control panel updates
func (cp *ControlPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "right":
			if cp.ultraCompactMode {
				// In ultra compact mode, only navigate between essential buttons
				essentialButtons := []int{0, 1, 2, 3} // Add Sphere, Add Sprite, Clear, Pause
				currentPos := -1
				for i, idx := range essentialButtons {
					if idx == cp.focused {
						currentPos = i
						break
					}
				}
				if currentPos == -1 {
					cp.focused = essentialButtons[0]
				} else {
					cp.focused = essentialButtons[(currentPos+1)%len(essentialButtons)]
				}
			} else {
				cp.focused = (cp.focused + 1) % len(cp.buttons)
			}
		case "shift+tab", "left":
			if cp.ultraCompactMode {
				// In ultra compact mode, only navigate between essential buttons
				essentialButtons := []int{0, 1, 2, 3} // Add Sphere, Add Sprite, Clear, Pause
				currentPos := -1
				for i, idx := range essentialButtons {
					if idx == cp.focused {
						currentPos = i
						break
					}
				}
				if currentPos == -1 {
					cp.focused = essentialButtons[len(essentialButtons)-1]
				} else {
					cp.focused = essentialButtons[(currentPos-1+len(essentialButtons))%len(essentialButtons)]
				}
			} else {
				cp.focused = (cp.focused - 1 + len(cp.buttons)) % len(cp.buttons)
			}
		case "enter", " ":
			// Activate focused button
			return cp, cp.activateButton(cp.focused)
		}
	case tea.MouseMsg:
		// Handle mouse clicks on buttons
		if msg.Type == tea.MouseLeft {
			buttonIndex := cp.getButtonAtPosition(msg.X, msg.Y)
			if buttonIndex >= 0 {
				cp.focused = buttonIndex
				return cp, cp.activateButton(buttonIndex)
			}
		}
	}
	return cp, nil
}

// activateButton creates a command to activate the button
func (cp *ControlPanel) activateButton(index int) tea.Cmd {
	if index >= 0 && index < len(cp.buttons) {
		return func() tea.Msg {
			return ButtonMsg{Action: cp.buttons[index].Action}
		}
	}
	return nil
}

// getButtonAtPosition returns the button index at the given position, or -1 if none
func (cp *ControlPanel) getButtonAtPosition(x, y int) int {
	// This is a simplified implementation - in a real app you'd calculate exact positions
	// For now, we'll use a simple heuristic based on button order
	if y == 3 { // Main button row
		buttonWidth := 15 // approximate button width including margins
		buttonIndex := x / buttonWidth
		if buttonIndex >= 0 && buttonIndex < len(cp.buttons) {
			return buttonIndex
		}
	}
	return -1
}

// View renders the control panel with minimal clutter
func (cp *ControlPanel) View() string {
	var lines []string

	// Single row layout - combine title, buttons, and hints in minimum space
	if cp.ultraCompactMode {
		// Ultra compact: Everything in 2 lines max

		// Line 1: Essential buttons only
		var buttonParts []string
		essentialButtons := []int{0, 1, 2, 3} // Add Sphere, Add Sprite, Clear, Pause

		for _, idx := range essentialButtons {
			var buttonText string
			switch cp.buttons[idx].Action {
			case AddSphereAction:
				buttonText = "●"
			case AddSpriteAction:
				buttonText = "◆"
			case ClearAllAction:
				buttonText = "Clear"
			case PauseResumeAction:
				if cp.buttons[idx].Label == "Resume" {
					buttonText = "▶"
				} else {
					buttonText = "⏸"
				}
			}

			if idx == cp.focused {
				buttonText = "→" + buttonText + "←"
			}

			var style lipgloss.Style
			if idx == cp.focused {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498DB")).Bold(true)
			} else {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E6ED"))
			}

			buttonParts = append(buttonParts, style.Render(buttonText))
		}

		// Combine controls and params in one line
		controlsLine := strings.Join(buttonParts, " ") + " | " +
			fmt.Sprintf("⚙️%s 📏%s 🎨%s", cp.gravityText, cp.sizeText, cp.colorText)
		lines = append(lines, controlsLine)

		// Line 2: Essential keys only
		keyHints := "Keys: A●  S◆  C=Clear  P=Pause  F=Perf  TAB=Navigate"
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
		lines = append(lines, keyStyle.Render(keyHints))

	} else if cp.compactMode {
		// Compact: 3 lines max

		// Line 1: Title
		title := "🎮 CONTROLS"
		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true).Align(lipgloss.Center)
		lines = append(lines, titleStyle.Width(cp.width).Render(title))

		// Line 2: All buttons in one row
		var buttonParts []string
		for i, button := range cp.buttons {
			var buttonText string
			if cp.compactMode {
				buttonText = cp.getCompactLabel(button)
			} else {
				buttonText = button.Label
			}

			if i == cp.focused {
				buttonText = "→" + buttonText + "←"
			}

			var style lipgloss.Style
			if i == cp.focused {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498DB")).Bold(true)
			} else {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E6ED"))
			}

			buttonParts = append(buttonParts, style.Render(buttonText))
		}

		buttonRow := strings.Join(buttonParts, " ")
		lines = append(lines, buttonRow)

		// Line 3: Parameters and key hints combined
		paramStatus := fmt.Sprintf("⚙️%s 📏%s 🎨%s", cp.gravityText, cp.sizeText, cp.colorText)
		keyHints := " | Keys: A●  S◆  C=Clear  P=Pause  G=Gravity  B=Bounce  Z=Size  X=Color  F=Perf"
		combinedLine := paramStatus + keyHints

		paramStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F39C12"))
		lines = append(lines, paramStyle.Render(combinedLine))

	} else {
		// Normal mode: Still compact but more readable

		// Line 1: Title
		title := "🎮 PHYSICS CONTROLS"
		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true).Align(lipgloss.Center)
		lines = append(lines, titleStyle.Width(cp.width).Render(title))

		// Line 2: All buttons
		var buttonParts []string
		for i, button := range cp.buttons {
			buttonText := button.Label

			if i == cp.focused {
				buttonText = "→" + buttonText + "←"
			}

			var style lipgloss.Style
			if i == cp.focused {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498DB")).Bold(true)
			} else {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E6ED"))
			}

			buttonParts = append(buttonParts, style.Render(buttonText))
		}

		buttonRow := strings.Join(buttonParts, " ")
		lines = append(lines, buttonRow)

		// Line 3: Parameters
		paramStatus := fmt.Sprintf("⚙️%s 📏%s 🎨%s", cp.gravityText, cp.sizeText, cp.colorText)
		paramStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F39C12"))
		lines = append(lines, paramStyle.Render(paramStatus))

		// Line 4: Key hints
		keyHints := "Keys: A=Add●  S=Add◆  C=Clear  P=Pause  R=Reset  G=Gravity  B=Bounce  Z=Size  X=Color  F=Perf  T=Test  L=Limit  TAB=Navigate"
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
		lines = append(lines, keyStyle.Render(keyHints))
	}

	return strings.Join(lines, "\n")
}

// getCompactLabel returns shortened button labels for compact modes
func (cp *ControlPanel) getCompactLabel(button Button) string {
	switch button.Action {
	case AddSphereAction:
		return "●"
	case AddSpriteAction:
		return "◆"
	case ClearAllAction:
		return "Clear"
	case PauseResumeAction:
		if button.Label == "Resume" {
			return "▶"
		}
		return "⏸"
	case ResetAction:
		return "↻"
	case GravityAction:
		return "⬇"
	case BounceAction:
		return "🏀"
	case SizeAction:
		return "📏"
	case ColorAction:
		return "🎨"
	default:
		return button.Label
	}
}

// UpdatePauseButton updates the pause button label based on current state
func (cp *ControlPanel) UpdatePauseButton(paused bool) {
	for i := range cp.buttons {
		if cp.buttons[i].Action == PauseResumeAction {
			if paused {
				cp.buttons[i].Label = "Resume"
			} else {
				cp.buttons[i].Label = "Pause"
			}
			break
		}
	}
}

// SetButtonActive sets a button's active state
func (cp *ControlPanel) SetButtonActive(action ButtonAction, active bool) {
	for i := range cp.buttons {
		if cp.buttons[i].Action == action {
			cp.buttons[i].Active = active
			break
		}
	}
}

// UpdateParameterDisplay updates the parameter display text
func (cp *ControlPanel) UpdateParameterDisplay(gravityText, sizeText, colorText string) {
	cp.gravityText = gravityText
	cp.sizeText = sizeText
	cp.colorText = colorText
}

// UpdateResponsiveMode sets the appropriate layout mode based on available space
func (cp *ControlPanel) UpdateResponsiveMode(width, height int) {
	cp.width = width
	cp.height = height

	// More aggressive responsive mode based on available space
	if height <= 3 || width < 40 {
		// Ultra compact: Only essential controls, minimal space
		cp.ultraCompactMode = true
		cp.compactMode = false
	} else if height <= 5 || width < 70 {
		// Compact: Reduced spacing and shorter labels
		cp.ultraCompactMode = false
		cp.compactMode = true
	} else {
		// Normal: Full layout with all visual flourishes
		cp.ultraCompactMode = false
		cp.compactMode = false
	}
}
