package clickhouse

import (
	"testing"

	"github.com/nyl1001/sqlchemy"
	"github.com/nyl1001/sqlchemy/backends/tests"
)

func TestQuery(t *testing.T) {
	t.Run("query all fields", func(t *testing.T) {
		tests.BackendTestReset(sqlchemy.ClickhouseBackend)
		q := tests.GetTestTable().Query()
		want := "SELECT `t1`.`col0`, `t1`.`col1` FROM `test` AS `t1`"
		tests.AssertGotWant(t, q.String(), want)
	})

	t.Run("query regexp field", func(t *testing.T) {
		tests.BackendTestReset(sqlchemy.ClickhouseBackend)
		testTable := tests.GetTestTable()
		q := testTable.Query(testTable.Field("col0")).Regexp("col1", "^ab$")
		want := "SELECT `t1`.`col0` FROM `test` AS `t1` WHERE match(`t1`.`col1`,  ? )"
		tests.AssertGotWant(t, q.String(), want)
	})
}
