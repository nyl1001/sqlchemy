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

package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/golang-plus/uuid"

	"github.com/nyl1001/pkg/gotypes"
	"github.com/nyl1001/pkg/jsonutils"
	"github.com/nyl1001/pkg/util/timeutils"
	"yunion.io/x/log"

	"github.com/nyl1001/sqlchemy"
	_ "github.com/nyl1001/sqlchemy/backends"
)

func uuid4() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}

func now() string {
	return timeutils.MysqlTime(time.Time{})
}

// SCompondStruct is an example struct
type SCompondStruct struct {
	Id  string
	Age int
}

func (cs *SCompondStruct) IsZero() bool {
	return len(cs.Id) == 0 && cs.Age == 0
}

func (cs *SCompondStruct) Equals(obj gotypes.ISerializable) bool {
	comp, ok := obj.(*SCompondStruct)
	if !ok {
		return false
	}
	return cs.Age == comp.Age && cs.Id == comp.Id
}

func (cs *SCompondStruct) String() string {
	return jsonutils.Marshal(cs).String()
}

func init() {
	comp := SCompondStruct{}
	compType := reflect.TypeOf(&comp)
	gotypes.RegisterSerializable(compType, func() gotypes.ISerializable {
		return &SCompondStruct{}
	})
}

type TestTable struct {
	Id        string               `primary:"true" width:"128" charset:"ascii" nullable:"false"`
	Name      string               `width:"64" charset:"utf8" primary:"true"`
	Gender    string               `width:"10" charset:"ascii"`
	Age       uint8                `default:"18"`
	Info      jsonutils.JSONObject `nullable:"false"`
	Compond   *SCompondStruct      `width:"1024"`
	CreatedAt time.Time            `nullable:"false" created_at:"true" clickhouse_partition_by:"toYYYYMM(created_at)"`
	UpdatedAt time.Time            `nullable:"false" updated_at:"true"`
	Version   int                  `default:"0" nullable:"false" auto_version:"true"`
	DeletedAt time.Time            ``
	Deleted   bool                 `nullable:"false" default:"false"`
	Notes     string               `width:"32" default:"default notes"`
}

type Ticket struct {
	Id     int       `auto_increment:"true"`
	UserId string    `width:"128" charset:"ascii" nullable:"false"`
	Date   time.Time `nullable:"false"`
}

func (t *TestTable) BeforeInsert() {
	t.Id = uuid4()
	dict := jsonutils.NewDict()
	dict.Add(jsonutils.NewString("Test"), "name")
	t.Info = dict
	t.Compond = &SCompondStruct{Id: "123456", Age: 24}
}

func (t *Ticket) BeforeInsert() {
	t.Date = time.Now().UTC()
}

type AgentTable struct {
	UserId string `primary:"true" width:"128" charset:"ascii"`
	Age    int
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run main.go <backend> ...")
		fmt.Println("backend: sqlite|mysql|clickhouse")
		fmt.Println("       go run main.go sqlite <filename>")
		fmt.Println("       go run main.go mysql <user> <passwd> <host> <database> [port]")
		fmt.Println("       go run main.go clickhouse <>")
		os.Exit(1)
	}
	// db, err := sql.Open("mysql", "testgo:openstack@tcp(127.0.0.1:3306)/testgo?charset=utf8&parseTime")
	dbName := sqlchemy.DBName("mydb")

	var db *sql.DB
	var err error
	var backend sqlchemy.DBBackendName
	switch os.Args[1] {
	case "sqlite", "sqlite3":
		dbfile := os.Args[2]
		db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbfile))
		backend = sqlchemy.SQLiteBackend
	case "mysql":
		user := os.Args[2]
		pass := os.Args[3]
		host := os.Args[4]
		dbname := os.Args[5]
		port := "3306"
		if len(os.Args) > 6 {
			port = os.Args[6]
		}
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime", user, pass, host, port, dbname))
		backend = sqlchemy.MySQLBackend
	case "clickhouse":
		host := os.Args[2]
		dbname := os.Args[3]
		port := "9000"
		if len(os.Args) > 4 {
			port = os.Args[4]
		}
		db, err = sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?database=%s&read_timeout=10&write_timeout=20", host, port, dbname))
		backend = sqlchemy.ClickhouseBackend
	default:
		panic(fmt.Sprintf("unsupported backend %s", os.Args[1]))
	}
	if err != nil {
		panic(fmt.Sprintf("Open DB failed: %s", err))
	}

	sqlchemy.SetDBWithNameBackend(db, dbName, backend)
	defer sqlchemy.CloseDB()

	tablespec := sqlchemy.NewTableSpecFromStructWithDBName(TestTable{}, "testtable", dbName)
	{
		err := tablespec.Sync()
		if err != nil {
			panic(fmt.Sprintf("table.sync %s", err.Error()))
		}
	}
	{
		err := tablespec.CheckSync()
		if err != nil {
			panic(fmt.Sprintf("table.checkSync %s", err.Error()))
		}
	}

	agespec := sqlchemy.NewTableSpecFromStructWithDBName(AgentTable{}, "age_tbl", dbName)
	agespec.Sync()
	agespec.CheckSync()

	ticketSpec := sqlchemy.NewTableSpecFromStructWithDBName(Ticket{}, "ticket_tbl", dbName)
	ticketSpec.Sync()
	ticketSpec.CheckSync()

	t1 := tablespec.Instance()
	t2 := agespec.Instance()
	// t3 := ticketSpec.Instance()
	q := t1.Query(t1.Field("id"), t2.Field("age")).
		Join(t2, sqlchemy.AND(sqlchemy.Equals(t1.Field("id"), t2.Field("user_id")), sqlchemy.GE(t2.Field("age"), 20))).
		Filter(sqlchemy.Like(t1.Field("Id"), "123%")).
		Limit(10).
		Asc(t2.Field("age"))
	fmt.Println(q.String())
	fmt.Println(q.Variables())

	subq := q.SubQuery()

	q2 := subq.Query().Desc("id")
	fmt.Println(q2.String())
	fmt.Println(q2.Variables())

	dt1 := TestTable{}
	dt1.Name = "Test"
	dt1.Gender = "male"
	dt1.Deleted = true
	// dt1.Notes = "not null notes"
	err = tablespec.Insert(&dt1)
	if err != nil {
		log.Errorf("Insert error: %s", err)
	}

	fmt.Println("dt1 after insert: ", dt1)

	dt2 := Ticket{}
	dt2.UserId = dt1.Id
	err = ticketSpec.Insert(&dt2)
	if err != nil {
		log.Errorf("Insert ticket fail %s", err)
	}
	fmt.Println(dt2)

	count := q.Count()
	fmt.Println("Count: ", count)
	count = q2.Count()
	fmt.Println("Count: ", count)

	q = t1.Query().Desc(t1.Field("created_at")).IsNotEmpty("info").IsNotNull("compond").Limit(10)
	mapData, err := q.AllStringMap()
	if err != nil {
		log.Errorf("query first %s", err)
	} else {
		for _, v := range mapData {
			fmt.Println(v)
		}
	}

	dt3 := TestTable{}
	err = q.First(&dt3)
	if err != nil {
		log.Errorf("First error %s", err)
	} else {
		fmt.Println("Before update", dt3)
	}

	/*session, err := tablespec.PrepareUpdate(dt3)
	if err != nil {
		log.Errorf("Fail to prepare update %s", err)
	}else {
		dt3.Name = "New name 4"
		dt3.Compond = &SCompondStruct{Id:"998822333", Age: 80}
		// dt3.Compond.Age = 80
		// dt3.Compond.Id = "998822333"
		diff, err := session.SaveUpdate(dt3)
		if err != nil {
			log.Errorf("SaveUpdate fail %s", err)
		}else {
			log.Infof("Update difference: %s", sqlchemy.UpdateDiffString(diff))
		}
	}*/

	_, err = tablespec.Update(&dt3, func() error {
		dt3.Name = "New name 4"
		dt3.Age = 10
		dt3.Compond = &SCompondStruct{Id: "998822333", Age: 80}
		return nil
	})
	if err != nil {
		log.Errorf("update fail %s", err)
	}

	dt3.Age = 1
	target := TestTable{Id: dt3.Id, Name: dt3.Name}
	err = tablespec.Increment(dt3, &target)
	if err != nil {
		log.Errorf("incremental faild %s", err)
	} else {
		log.Infof("Increment: %d %d", target.Age, dt3.Age)
	}

	q = t1.Query().Equals("id", dt3.Id)
	err = q.First(&dt3)
	if err != nil {
		log.Errorf("First error %s", err)
	} else {
		fmt.Println("After update: ", jsonutils.Marshal(dt3))
	}

	log.Infof("Start SQuery ALL")

	dt4 := make([]TestTable, 0)
	err = q.GT("version", 0).All(&dt4)
	if err != nil {
		log.Errorf("query all fail %s", err)
	} else {
		log.Infof("SQuery all no error %d", len(dt4))
		for _, v := range dt4 {
			fmt.Println("dt4 ", jsonutils.Marshal(v))
		}
	}

	qId1 := t1.Query(t1.Field("id"))
	t3 := ticketSpec.Instance()
	qId2 := t3.Query(sqlchemy.CAST(t3.Field("id"), "TEXT", "id"))

	{
		union, err := sqlchemy.UnionAllWithError(qId1, qId2)
		if err != nil {
			log.Errorf("Union with error: %s", err)
		} else {
			q := union.Query()
			fmt.Println(q.String())

			type sID struct {
				Id string
			}
			idList := make([]sID, 0)
			err := q.All(&idList)
			if err != nil {
				log.Errorf("fail to query idList %s", err)
			} else {
				log.Infof("Test: %s", jsonutils.Marshal(idList))
			}
		}
	}

}
