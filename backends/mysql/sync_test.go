// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"reflect"
	"testing"

	"github.com/nyl1001/sqlchemy"
)

func TestSync(t *testing.T) {
	type TableStruct1 struct {
		Id     uint64 `auto_increment:"true"`
		Name   string `width:"64" charset:"utf8"`
		Age    int    `nullable:"true" default:"12"`
		IsMale *bool  `nullable:"false" default:"true"`
	}
	type TableStruct2 struct {
		Id     uint64 `auto_increment:"true"`
		Name   string `width:"128" charset:"utf8"`
		Age    uint   `nullable:"true" default:"12"`
		Gender string `width:"8" nullable:"false" default:"male"`
	}

	sqlchemy.SetDBWithNameBackend(nil, sqlchemy.DefaultDB, sqlchemy.MySQLBackend)
	ts1 := sqlchemy.NewTableSpecFromStruct(TableStruct1{}, "table1")
	ts2 := sqlchemy.NewTableSpecFromStruct(TableStruct2{}, "table1")

	changes := sqlchemy.STableChanges{}
	changes.RemoveColumns, changes.UpdatedColumns, changes.AddColumns = sqlchemy.DiffCols(ts2.Name(), ts1.Columns(), ts2.Columns())
	backend := &SMySQLBackend{}
	sqls := backend.CommitTableChangeSQL(ts2, changes)
	want := []string{
		"ALTER TABLE `table1` MODIFY COLUMN `age` INT(10) UNSIGNED DEFAULT 12, MODIFY COLUMN `name` VARCHAR(128) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci', ADD COLUMN `gender` VARCHAR(8) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT 'male';",
	}
	if !reflect.DeepEqual(sqls, want) {
		t.Errorf("Expect: %s", want)
		t.Errorf("Got: %s", sqls)
	}
}
