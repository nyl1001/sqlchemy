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
	"database/sql"
	"testing"

	"yunion.io/x/jsonutils"

	"github.com/nyl1001/pkg/tristate"

	"github.com/nyl1001/sqlchemy"
)

var (
	triCol         = NewTristateColumn("", "field", nil, false)
	notNullTriCol  = NewTristateColumn("", "field", map[string]string{sqlchemy.TAG_NULLABLE: "false"}, false)
	boolCol        = NewBooleanColumn("field", nil, false)
	notNullBoolCol = NewBooleanColumn("field", map[string]string{sqlchemy.TAG_NULLABLE: "false"}, false)
	integerCol     = NewIntegerColumn("field", nil, false)
	floatCol       = NewFloatColumn("field", nil, false)
	textCol        = NewTextColumn("field", nil, false)
	notNullTextCol = NewTextColumn("field", map[string]string{sqlchemy.TAG_NULLABLE: "false"}, false)
	defTextCol     = NewTextColumn("field", map[string]string{sqlchemy.TAG_DEFAULT: "new!"}, false)
	dateCol        = NewDateTimeColumn("field", nil, false)
	notNullDateCol = NewDateTimeColumn("field", map[string]string{sqlchemy.TAG_NULLABLE: "false"}, false)
	compCol        = NewCompoundColumn("field", nil, false)
)

func TestColumns(t *testing.T) {
	cases := []struct {
		in   sqlchemy.IColumnSpec
		want string
	}{
		{
			in:   &triCol,
			want: "`field` INTEGER",
		},
		{
			in:   &notNullTriCol,
			want: "`field` INTEGER",
		},
		{
			in:   &boolCol,
			want: "`field` INTEGER",
		},
		{
			in:   &notNullBoolCol,
			want: "`field` INTEGER NOT NULL",
		},
		{
			in:   &integerCol,
			want: "`field` INTEGER",
		},
		{
			in:   &floatCol,
			want: "`field` REAL",
		},
		{
			in:   &textCol,
			want: "`field` TEXT COLLATE NOCASE",
		},
		{
			in:   &notNullTextCol,
			want: "`field` TEXT NOT NULL COLLATE NOCASE",
		},
		{
			in:   &defTextCol,
			want: "`field` TEXT DEFAULT 'new!' COLLATE NOCASE",
		},
		{
			in:   &dateCol,
			want: "`field` TEXT COLLATE NOCASE",
		},
		{
			in:   &notNullDateCol,
			want: "`field` TEXT NOT NULL COLLATE NOCASE",
		},
		{
			in:   &compCol,
			want: "`field` TEXT COLLATE NOCASE",
		},
	}
	for _, c := range cases {
		got := c.in.DefinitionString()
		if got != c.want {
			t.Errorf("got %s want %s", got, c.want)
		}
	}
}

func TestConvertValue(t *testing.T) {
	cases := []struct {
		in   interface{}
		want interface{}
		col  sqlchemy.IColumnSpec
	}{
		{
			in:   true,
			want: 1,
			col:  &boolCol,
		},
		{
			in:   false,
			want: 0,
			col:  &boolCol,
		},
		{
			in:   tristate.True,
			want: 1,
			col:  &triCol,
		},
		{
			in:   tristate.False,
			want: 0,
			col:  &triCol,
		},
		{
			in:   tristate.None,
			want: sql.NullInt32{},
			col:  &triCol,
		},
		{
			in:   23,
			want: 23,
			col:  &integerCol,
		},
		{
			in:   jsonutils.NewDict(),
			want: `{}`,
			col:  &compCol,
		},
	}
	for _, c := range cases {
		got := c.col.ConvertFromValue(c.in)
		if got != c.want {
			t.Errorf("%s [%#v] want: %#v got: %#v", c.col.DefinitionString(), c.in, c.want, got)
		}
	}
}
func TestConvertString(t *testing.T) {
	cases := []struct {
		in   string
		want interface{}
		col  sqlchemy.IColumnSpec
	}{
		{
			in:   `true`,
			want: 1,
			col:  &boolCol,
		},
		{
			in:   "false",
			want: 0,
			col:  &boolCol,
		},
		{
			in:   "true",
			want: 1,
			col:  &triCol,
		},
		{
			in:   "false",
			want: 0,
			col:  &triCol,
		},
		{
			in:   "none",
			want: sql.NullInt32{},
			col:  &triCol,
		},
		{
			in:   "23",
			want: int64(23),
			col:  &integerCol,
		},
		{
			in:   "0.01",
			want: 0.01,
			col:  &floatCol,
		},
	}
	for _, c := range cases {
		got := c.col.ConvertFromString(c.in)
		if got != c.want {
			t.Errorf("%s [%s] want: %#v got: %#v", c.col.DefinitionString(), c.in, c.want, got)
		}
	}
}
