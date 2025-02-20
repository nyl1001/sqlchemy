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

package sqlite

import (
	"regexp"

	"github.com/nyl1001/pkg/errors"

	"github.com/nyl1001/sqlchemy"
)

const (
	indexPattern = `CREATE\s+INDEX\s+` + "`" + `(?P<name>\w+)` + "`" + `\s+ON\s+` + "`" + `(?P<tblname>\w+)` + "`" + `\s*\((?P<cols>` + "`" + `\w+` + "`" + `(,\s*` + "`" + `\w+` + "`" + `)*)\)`
)

var (
	indexRegexp = regexp.MustCompile(indexPattern)
)

type sSqliteTableInfo struct {
	Type string
	Name string
	Sql  string
}

func (ti *sSqliteTableInfo) parseTableIndex(ts sqlchemy.ITableSpec) (sqlchemy.STableIndex, error) {
	matches := indexRegexp.FindAllStringSubmatch(ti.Sql, -1)
	if len(matches) > 0 {
		return sqlchemy.NewTableIndex(ts, matches[0][1], sqlchemy.FetchColumns(matches[0][3]), false), nil
	}
	return sqlchemy.STableIndex{}, errors.ErrNotFound
}
