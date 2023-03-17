package model

import (
	"github.com/atye/wikitable/bubble"
	"github.com/charmbracelet/lipgloss"
)

type table struct {
	model          bubble.Model
	data           [][]string
	originalData   [][]string
	maxColumnWidth int
}

func newTable(data [][]string, height, maxColumnWidth int) *table {
	od := make([][]string, len(data))
	copy(od, data)

	return &table{
		model:          generateModel(data, height, maxColumnWidth),
		data:           data,
		originalData:   od,
		maxColumnWidth: maxColumnWidth,
	}
}

func (t *table) remove() {
	cursor := t.model.Cursor()
	switch t.model.CursorMode() {
	case "row":
		numRows := len(t.model.Rows())
		t.removeRow(cursor + 1)
		if cursor == numRows-1 {
			t.model.SetCursor(len(t.model.Rows()) - 1)
		}
	case "column":
		numCols := len(t.model.Columns())
		t.removeColumn(cursor)
		if cursor == numCols-1 {
			t.model.SetCursor(len(t.model.Columns()) - 1)
		}
	}
}

func (t *table) removeRow(row int) {
	data := t.data

	if row == len(data)-1 {
		data = data[:len(data)-1]
	} else {
		data = append(data[:row], data[row+1:]...)
	}

	t.data = data
	t.set()
}

func (t *table) removeColumn(column int) {
	data := t.data

	if column == len(data[0])-1 {
		for i, row := range data {
			data[i] = row[:len(row)-1]
		}
	} else {
		for i, row := range data {
			data[i] = append(row[:column], row[column+1:]...)
		}
	}

	t.data = data
	t.set()
}

func (t *table) set() {
	rows := make([]bubble.Row, len(t.data[1:]))
	for i, row := range t.data[1:] {
		rows[i] = row
	}
	t.model.SetRows(rows)

	columns := make([]bubble.Column, len(t.data[0]))
	for i, col := range t.data[0] {
		colWidth := maxColumnWidth(t.data, i, t.maxColumnWidth)
		columns[i] = bubble.Column{
			Title: col,
			Width: colWidth,
		}
	}
	t.model.SetColumns(columns)
}

func (t *table) moveUp(n int) {
	t.model.MoveUp(n)
}

func (t *table) moveDown(n int) {
	t.model.MoveDown(n)
}

func (t *table) goToTop() {
	if t.model.CursorMode() == "row" {
		t.model.GotoTop()
	}
}

func (t *table) goToBottom() {
	if t.model.CursorMode() == "row" {
		t.model.GotoBottom()
	}
}

func (t *table) switchCursorMode() {
	t.model.SwitchCursorMode()
}

func (t *table) reset(height int) {
	t.model = generateModel(t.originalData, height, t.maxColumnWidth)
	t.data = t.originalData
}

func generateModel(data [][]string, height, maxColumWidth int) bubble.Model {
	var tableRows []bubble.Row
	var tableCols []bubble.Column
	for rowIndex, row := range data {
		if rowIndex == 0 {
			for colIndex, col := range row {
				colWidth := maxColumnWidth(data, colIndex, maxColumWidth)
				tableCols = append(tableCols, bubble.Column{
					Title: col,
					Width: colWidth,
				})
			}
			continue
		}
		tableRows = append(tableRows, row)
	}

	model := bubble.New(bubble.WithColumns(tableCols), bubble.WithRows(tableRows), bubble.WithHeight(height))
	s := bubble.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	model.SetStyles(s)

	return model
}

func maxColumnWidth(table [][]string, col int, maxWidth int) int {
	var width int
	for _, row := range table {
		if len(row[col]) > width {
			width = len(row[col])
		}
		if maxWidth > 0 && width >= maxWidth {
			return maxWidth
		}
	}
	return width
}
