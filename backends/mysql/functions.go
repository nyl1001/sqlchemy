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
	"fmt"

	"github.com/nyl1001/sqlchemy"
)

// GROUP_CONCAT2 represents the SQL function GROUP_CONCAT
func (mysql *SMySQLBackend) GROUP_CONCAT2(name string, sep string, field sqlchemy.IQueryField) sqlchemy.IQueryField {
	return sqlchemy.NewFunctionField(name, fmt.Sprintf("GROUP_CONCAT(%%s SEPARATOR '%s')", sep), field)
}
