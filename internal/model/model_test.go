package model

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/atye/wikitable/bubble"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModel(t *testing.T) {
	t.Run("it sets tables from one page", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		want := data[0]
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it sets tables from multiple pages", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page,page")
		sut.input.inputs[langIndex].SetValue("en,en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		want := data[0]
		if !reflect.DeepEqual(want, modelToData(sut.tables[0].model)) {
			t.Errorf("expected %v, got %v", want, modelToData(sut.tables[0].model))
		}

		if !reflect.DeepEqual(want, modelToData(sut.tables[1].model)) {
			t.Errorf("expected %v, got %v", want, modelToData(sut.tables[0].model))
		}
	})

	t.Run("it removes first row from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column", "column2", "column3"},
			{"test2", "test2", "test2"},
			{"test3", "test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it removes middle row from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.tables[0].model.SetCursor(1)
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column", "column2", "column3"},
			{"test", "test", "test"},
			{"test3", "test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it removes last row from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.tables[0].model.SetCursor(2)
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column", "column2", "column3"},
			{"test", "test", "test"},
			{"test2", "test2", "test2"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}

		cursor := sut.tables[0].model.Cursor()
		if cursor != 1 {
			t.Errorf("expected cursor value 1, got %v", cursor)
		}
	})

	t.Run("it removes first column from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlK}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column2", "column3"},
			{"test", "test"},
			{"test2", "test2"},
			{"test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it removes middle column from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlK}))
		sut.tables[0].model.SetCursor(1)
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column", "column3"},
			{"test", "test"},
			{"test2", "test2"},
			{"test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it removes last column from a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlK}))
		sut.tables[0].model.SetCursor(2)
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))

		want := [][]string{
			{"column", "column2"},
			{"test", "test"},
			{"test2", "test2"},
			{"test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}

		cursor := sut.tables[0].model.Cursor()
		if cursor != 1 {
			t.Errorf("expected cursor value 1, got %v", cursor)
		}
	})

	t.Run("it removes first table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlT}))

		if len(sut.tables) != 2 {
			t.Errorf("expected two tables, got %d", len(sut.tables))
		}
	})

	t.Run("it removes middle table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyTab}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlT}))

		if len(sut.tables) != 2 {
			t.Errorf("expected two tables, got %d", len(sut.tables))
		}
	})

	t.Run("it removes last table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyTab}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyTab}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlT}))

		if len(sut.tables) != 2 {
			t.Errorf("expected two tables, got %d", len(sut.tables))
		}

		if sut.index != 1 {
			t.Errorf("expected table index 1, got %d", sut.index)
		}
	})

	t.Run("it fills row data", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3", "column4"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		want := [][]string{
			{"column", "column2", "column3", "column4"},
			{"test", "test", "test", ""},
			{"test2", "test2", "test2", ""},
			{"test3", "test3", "test3", ""},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it resets a table", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlD}))
		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyCtrlR}))

		want := [][]string{
			{"column", "column2", "column3"},
			{"test", "test", "test"},
			{"test2", "test2", "test2"},
			{"test3", "test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it uses maxColumnWidth", func(t *testing.T) {
		data := [][][]string{
			{
				{"column", "column2", "column3"},
				{"test", "test", "test"},
				{"test2", "test2", "test2"},
				{"test3", "test3", "test3"},
			},
		}

		fw := fakeWiki{
			GetTablesMatrixFn: func(ctx context.Context, page, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
				return data, nil
			},
		}
		sut := NewModel(fw)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.inputs[maxColumnWidthIndex].SetValue("3")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		want := [][]string{
			{"column", "column2", "column3"},
			{"test", "test", "test"},
			{"test2", "test2", "test2"},
			{"test3", "test3", "test3"},
		}
		got := modelToData(sut.tables[0].model)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("it sets error on empty page", func(t *testing.T) {
		sut := NewModel(nil)

		sut.input.inputs[pageIndex].SetValue("")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		if sut.inputErr == nil {
			t.Errorf("expected input error, got nil")
		}
	})

	t.Run("it sets error on empty lang", func(t *testing.T) {
		sut := NewModel(nil)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		if sut.inputErr == nil {
			t.Errorf("expected input error, got nil")
		}
	})

	t.Run("it sets error on invalid cleanRef", func(t *testing.T) {
		sut := NewModel(nil)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("test")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		if sut.inputErr == nil {
			t.Errorf("expected input error, got nil")
		}
	})

	t.Run("it sets error on invalid maxColumnWidth", func(t *testing.T) {
		sut := NewModel(nil)

		sut.input.inputs[pageIndex].SetValue("page")
		sut.input.inputs[langIndex].SetValue("en")
		sut.input.inputs[cleanRefIndex].SetValue("t")
		sut.input.inputs[maxColumnWidthIndex].SetValue("test")
		sut.input.focus = 4

		sut.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

		if sut.inputErr == nil {
			t.Errorf("expected input error, got nil")
		}
	})
}

func modelToData(model bubble.Model) [][]string {
	cols := model.Columns()
	rows := model.Rows()
	data := make([][]string, len(rows)+1)

	data[0] = make([]string, len(cols))
	for i, col := range cols {
		data[0][i] = col.Title
	}

	for i, row := range rows {
		data[i+1] = row
	}

	return data
}

type fakeWiki struct {
	GetTablesMatrixFn func(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error)
}

func (f fakeWiki) GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
	if f.GetTablesMatrixFn != nil {
		return f.GetTablesMatrixFn(ctx, page, lang, cleanRef, tables...)
	}
	return nil, fmt.Errorf("error")
}
