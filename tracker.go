package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
)

// styles
var (
    // Input box style
    inputBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("97")).
        Padding(1, 2).
        MarginTop(1)
)

// model structure 
type model struct {
	textInput textinput.Model
	table table.Model
	selectedCol int
	data []table.Row
	inputLabel string
	db *sql.DB
}
func fetchData(db *sql.DB) []table.Row {
	// pull rows from sql supabase db
    sqlRows, err := db.Query("SELECT class, name, duedate, priority, location, estimatedtime, actualtime, id FROM assignments ORDER BY duedate")
    if err != nil {
        log.Fatal(err)
    }
	var tableRows []table.Row

	// scan rows into table.Row array
	for sqlRows.Next() {
      var class, name, duedate, priority, location, estimatedTime, realTime, id string

      err := sqlRows.Scan(&class, &name, &duedate, &priority, &location, &estimatedTime, &realTime, &id)
      if err != nil {
          log.Fatal(err)
      }

      tableRows = append(tableRows, table.Row{class, name, duedate, priority, location, estimatedTime, realTime, id})
	}

	if err = sqlRows.Err(); err != nil {
      log.Fatal(err)
  	}
	return tableRows
	
}
func initialModel(db *sql.DB) *model {
	ti := textinput.New()
	ti.Placeholder = "type here"
	ti.CharLimit = 20
	ti.Width = 30
    columns := []table.Column{
        {Title: "Class", Width: 10},
        {Title: "Name", Width: 20},
        {Title: "Due Date", Width: 15},
        {Title: "Priority", Width: 15},
        {Title: "Location", Width: 15},
        {Title: "Estimated Time", Width: 15},
        {Title: "Real Time", Width: 15},
        {Title: "ID", Width: 10},
    }

	tableRows := fetchData(db)

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(tableRows),
        table.WithHeight(29),
    )
	t.Focus()
	// set the styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("97")).
		Bold(false)
	t.SetStyles(s)

    return &model{
		textInput: ti,
		table: t,
		selectedCol: 0,
		data: tableRows,
		inputLabel: "",
		db: db,
    }
}

func (m *model) Init() tea.Cmd {
    return textinput.Blink
}

func (m *model) updatePointer() {
	// function to add > to whatever field the cursor is currently on
	displayRows := make([]table.Row, len(m.data))
	for i, row := range(m.data){
		displayRows[i] = make([]string, len(row))
		for j, cell := range(row){
			if i == m.table.Cursor() && j == m.selectedCol{
				displayRows[i][j] = ">" + cell
			} else{
				displayRows[i][j] = cell
			}
		}
	}
	m.table.SetRows(displayRows)
}

func updateTable(id, col int, newVal string, db *sql.DB) {
	// uses this dictionary of the sql col names for the fields and uses a template to update		
	cols := []string {"class", "name", "duedate", "priority", "location", "estimatedtime", "actualtime", "id"}
	query := fmt.Sprintf("UPDATE assignments SET %s = '%s' WHERE id = %d;", cols[col], newVal, id)
 	_, err := db.Exec(query)
  	if err != nil {
      log.Printf("Update failed: %v", err)
  	}
}

func (m *model) addRow() {
	// add new row with sample data points
	_, err := m.db.Exec("INSERT INTO assignments (class, name, duedate, priority, location, estimatedtime, actualtime) VALUES ('class', 'name', 'due date', 'low', 'location', 0, 0)")
	if err != nil {
      log.Printf("Update failed: %v", err)
  	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
		case "enter":
			if m.textInput.Focused(){
				
				id, err := strconv.Atoi(m.data[m.table.Cursor()][7])
				if err != nil {
					log.Fatal(err)
				}
				updateTable(id, m.selectedCol, m.textInput.Value(), m.db)
				m.data = fetchData(m.db)
				m.updatePointer()
				m.textInput.SetValue("")
				m.textInput.Blur()
				m.table.Focus()
			}
		case "ctrl+c":
            return m, tea.Quit

		case "ctrl+w":
			if m.table.Focused() {
				m.table.Blur()
				m.textInput.Focus()
			} else{
				m.textInput.Blur()
				m.table.Focus()
			}
		// vim like motions
		case "j":
			if m.table.Focused() {
				m.table.MoveDown(1)
			}
		case "k":
			if m.table.Focused() {
				m.table.MoveUp(1)
			}

		case "l":
			if m.table.Focused() && m.selectedCol < len(m.table.Columns())-1 {
				m.selectedCol++
			}
		case "h":
			if m.table.Focused() && m.selectedCol > 0 {
				m.selectedCol--
			}
		case "A":
			if m.table.Focused(){
				m.addRow()
				m.data = fetchData(m.db)
				m.updatePointer()
			}
		}

    }
	if len(m.data) > 0 && m.table.Cursor() < len(m.data) {
		m.textInput.Placeholder = m.data[m.table.Cursor()][m.selectedCol]
	}
	m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd 
}
func (m *model) sortByDueDate() {                                                                                                                                                  
      sort.Slice(m.data, func(i, j int) bool {                                                                                                                                       
		  // sorting algo to get assignments due as soon as possible
          dateI := m.data[i][2]                                                                                                                                                      
          dateJ := m.data[j][2]                                                                                                                                                      
                                                                                                                                                                                     
          layout := "01/02/2006"                                                                                                                                                     
          timeI, errI := time.Parse(layout, dateI)                                                                                                                                   
          timeJ, errJ := time.Parse(layout, dateJ)                                                                                                                                   
                                                                                                                                                                                     
          // Handle empty or invalid dates
		  if errI != nil {                                                                                                                                                           
              return false                                                                                                                                                           
          }                                                                                                                                                                          
          if errJ != nil {                                                                                                                                                           
              return true                                                                                                                                                            
          }                                                                                                                                                                          
                                                                                                                                                                                     
          // Sort earliest dates first                                                                                                                                               
          return timeI.Before(timeJ)
      })
                                                                                                                                                                                     
      // Update the table with sorted data                                                                                                                                           
      m.updatePointer()                                                                                                                                                              
  }
func (m *model) View() string {
	var inputSection string

    textInput := m.textInput.View()
    inputContent := fmt.Sprintf("%s", textInput)

    inputSection = inputBoxStyle.Render(inputContent)
	m.updatePointer()
	tableStyled := m.table.View()

    return lipgloss.JoinVertical(
		lipgloss.Center,
		tableStyled,
		inputSection,
	)
}

func main() {

	connStr := os.Getenv("ASSIGNMENT_TRACKER_URL")
	if connStr == "" {
        log.Fatal("SUPABASE_DB_URL environment variable not set")
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot connect to Supabase:", err)
    }

    p := tea.NewProgram(initialModel(db), tea.WithAltScreen())
    p.Run()
}
