package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/olekukonko/tablewriter"
)

type screen int

const (
	menuScreen screen = iota
	addTaskScreen
	listScreen
)

type ToDo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type model struct {
	choices       []string
	cursor        int
	currentScreen screen
}

func initialModel() model {
	return model{
		choices: []string{"Add task", "List Tasks", "Exit"},
	}

}

func getTodoList() []ToDo {
	todos, err := os.ReadFile("todo.json")

	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	var todoList []ToDo

	err = json.Unmarshal(todos, &todoList)

	if err != nil {
		log.Fatalf("Failed to parse json: %v\n", err)
	}

	return todoList

}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":

			switch m.cursor {

			case 0:
				m.currentScreen = addTaskScreen
			case 1:
				m.currentScreen = listScreen
			case 2:
				return m, tea.Quit

			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func menuView(m model) string {
	s := "Welcome to the Task Manager\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

func addTaskView(m model) string {
	s := "Add Task\n\n"

	return s
}

func listView(m model) string {
	s := "List view\n\n"

	s += generateTable(getTodoList())

	return s
}

func (m model) View() string {

	switch m.currentScreen {
	case addTaskScreen:
		return addTaskView(m)
	case listScreen:
		return listView(m)
	default:
		return menuView(m)
	}
}

func generateTable(todos []ToDo) string {
	var buf bytes.Buffer // Buffer para almacenar la salida

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Title", "Description"})

	for _, todo := range todos {
		table.Append([]string{todo.Title, todo.Description})
	}

	table.Render()

	return buf.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
