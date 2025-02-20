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
	"testing"
)

const tableDef = `CREATE TABLE ` + "`" + `image_properties` + "`" + ` (
  ` + "`" + `id` + "`" + ` varchar(128) CHARACTER SET ascii NOT NULL,
  ` + "`" + `image_id` + "`" + ` varchar(36) NOT NULL,
  ` + "`" + `name` + "`" + ` varchar(128) NOT NULL,
  ` + "`" + `value` + "`" + ` text,
  ` + "`" + `created_at` + "`" + ` datetime NOT NULL,
  ` + "`" + `updated_at` + "`" + ` datetime NOT NULL,
  ` + "`" + `deleted_at` + "`" + ` datetime DEFAULT NULL,
  ` + "`" + `deleted` + "`" + ` tinyint(1) NOT NULL DEFAULT '0',
  ` + "`" + `description` + "`" + ` varchar(256) DEFAULT NULL,
  ` + "`" + `external_id` + "`" + ` varchar(256) DEFAULT NULL,
  ` + "`" + `is_emulated` + "`" + ` tinyint(1) NOT NULL DEFAULT '0',
  ` + "`" + `update_version` + "`" + ` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (` + "`" + `id` + "`" + `),
  UNIQUE KEY ` + "`" + `image_id` + "`" + ` (` + "`" + `image_id` + "`" + `,` + "`" + `name` + "`" + `),
  UNIQUE KEY ` + "`" + `ix_image_properties_image_id_name` + "`" + ` (` + "`" + `image_id` + "`" + `,` + "`" + `name` + "`" + `),
  KEY ` + "`" + `ix_image_properties_image_id` + "`" + ` (` + "`" + `image_id` + "`" + `(10)),
  KEY ` + "`" + `ix_image_properties_deleted` + "`" + ` (` + "`" + `deleted` + "`" + `,` + "`" + `image_id` + "`" + `,` + "`" + `name` + "`" + `),
  KEY ` + "`" + `ix_image_properties_deleted_sorted` + "`" + ` (` + "`" + `image_id` + "`" + `,` + "`" + `deleted` + "`" + `,` + "`" + `name` + "`" + `),
  CONSTRAINT ` + "`" + `image_properties_ibfk_1` + "`" + ` FOREIGN KEY (` + "`" + `image_id` + "`" + `) REFERENCES ` + "`" + `images` + "`" + ` (` + "`" + `id` + "`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8`

type expectIndex struct {
	Name string
	Cols []string
}

func TestParseCreateTable(t *testing.T) {
	t.Logf("%s", tableDef)
	cons := parseConstraints(tableDef)
	if len(cons) != 1 {
		t.Errorf("fail to find constraints")
	}
	idxs := parseIndexes(nil, tableDef)
	if len(idxs) != 5 {
		t.Errorf("fail to find indexes")
	} else {
		expects := []struct {
			Name string
			Cols []string
		}{
			{
				Name: "image_id",
				Cols: []string{"image_id", "name"},
			},
			{
				Name: "ix_image_properties_image_id_name",
				Cols: []string{"image_id", "name"},
			},
			{
				Name: "ix_image_properties_image_id",
				Cols: []string{"image_id"},
			},
			{
				Name: "ix_image_properties_deleted",
				Cols: []string{"deleted", "image_id", "name"},
			},
			{
				Name: "ix_image_properties_deleted_sorted",
				Cols: []string{"name", "image_id", "deleted"},
			},
		}
		for i := range expects {
			if idxs[i].Name() != expects[i].Name {
				t.Errorf("expect name %s got name %s", expects[i].Name, idxs[i].Name())
			} else if !idxs[i].IsIdentical(expects[i].Cols...) {
				t.Errorf("expect column %s got name %s", expects[i].Cols, idxs[i].QuotedColumns())
			}
		}
	}
}
