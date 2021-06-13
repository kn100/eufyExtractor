// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("ScaleResults", testScaleResults)
}

func TestDelete(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsExists)
}

func TestFind(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsFind)
}

func TestBind(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsBind)
}

func TestOne(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsOne)
}

func TestAll(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsAll)
}

func TestCount(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsCount)
}

func TestHooks(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsInsert)
	t.Run("ScaleResults", testScaleResultsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("ScaleResults", testScaleResultsSliceUpdateAll)
}
