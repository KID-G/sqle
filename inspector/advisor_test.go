package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/stretchr/testify/assert"
	"sqle/model"
	"testing"
)

func getTestCreateTableStmt1() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_1 (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255),
PRIMARY KEY (id) USING BTREE,
KEY v1 (v1),
UNIQUE KEY v2 (v1,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`
	node, err := parseOneSql("mysql", baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt2() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_2 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255),
UNIQUE KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`
	node, err := parseOneSql("mysql", baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

type testResult struct {
	Results *InspectResults
	rules   map[string]model.Rule
}

func newTestResult() *testResult {
	return &testResult{
		Results: newInspectResults(),
		rules:   DefaultRulesMap,
	}
}

func (t *testResult) addResult(ruleName string, args ...interface{}) *testResult {
	rule, ok := t.rules[ruleName]
	if !ok {
		return t
	}
	t.Results.add(rule, args...)
	return t
}

func (t *testResult) level() string {
	return t.Results.level()
}

func (t *testResult) message() string {
	return t.Results.message()
}

func DefaultMysqlInspect() *Inspector {
	return &Inspector{
		Results: newInspectResults(),
		Rules:   DefaultRules,
		Db: model.Instance{
			DbType: "mysql",
		},
		SqlArray:      []*model.CommitSql{},
		currentSchema: "exist_db",
		allSchema:     map[string]struct{}{"exist_db": struct{}{}},
		schemaHasLoad: true,
		allTable: map[string]map[string]struct{}{
			"exist_db": map[string]struct{}{
				"exist_tb_1": struct{}{},
				"exist_tb_2": struct{}{},
			}},
		createTableStmts: map[string]*ast.CreateTableStmt{
			"exist_db.exist_tb_1": getTestCreateTableStmt1(),
			"exist_db.exist_tb_2": getTestCreateTableStmt2(),
		},
		alterTableStmts: map[string][]*ast.AlterTableStmt{},
		rollbackSqls:    []string{},
	}
}

func TestInspectResults(t *testing.T) {
	results := newInspectResults()
	results.add(DefaultRulesMap[DDL_CREATE_TABLE_NOT_EXIST])
	assert.Equal(t, "error", results.level())
	assert.Equal(t, "[error]新建表必须加入if not exists create，保证重复执行不报错", results.message())

	results.add(DefaultRulesMap[TABLE_NOT_EXIST], "not_exist_tb")
	assert.Equal(t, "error", results.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.message())
}

func runInspectCase(t *testing.T, desc string, i *Inspector, sql string, results ...*testResult) {
	stmts, err := parseSql(i.Db.DbType, sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	for n, stmt := range stmts {
		i.SqlArray = append(i.SqlArray, &model.CommitSql{
			Number: n + 1,
			Sql:    stmt.Text(),
		})
	}
	err = i.Advise()
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	if len(i.SqlArray) != len(results) {
		t.Errorf("%s test failled, error: result is unknow\n", desc)
		return
	}
	for n, sql := range i.SqlArray {
		result := results[n]
		if sql.InspectLevel != result.level() || sql.InspectResult != result.message() {
			t.Errorf("%s test failled, \n\nsql:\n %s\n\nexpect level: %s\nexpect result:\n%s\n\nactual level: %s\nactual result:\n%s\n",
				desc, sql.Sql, result.level(), result.message(), sql.InspectLevel, sql.InspectResult)
		} else {
			t.Log(fmt.Sprintf("\n\ncase:%s\nactual result:\n%s\n\n", desc, sql.InspectResult))
		}
	}
}

func TestInspector_Inspect_Message(t *testing.T) {
	runInspectCase(t, "check inspect message", DefaultMysqlInspect(),
		"use no_exist_db",
		&testResult{
			Results: &InspectResults{
				[]*InspectResult{&InspectResult{
					Level:   "error",
					Message: "schema no_exist_db 不存在",
				}},
			},
		},
	)
}

func TestInspector_Inspect_UseDatabaseStmt(t *testing.T) {
	runInspectCase(t, "use_database: ok", DefaultMysqlInspect(),
		"use exist_db",
		newTestResult(),
	)
}

func TestInspector_Advise_Select(t *testing.T) {
	runInspectCase(t, "select_from: ok", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id =1;",
		newTestResult(),
	)

	runInspectCase(t, "select_from: all columns", DefaultMysqlInspect(),
		"select * from exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(DML_DISABE_SELECT_ALL_COLUMN),
	)

	runInspectCase(t, "select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1;",
		newTestResult().addResult(DML_CHECK_INVALID_WHERE_CONDITION),
	)

	runInspectCase(t, "select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(DML_CHECK_INVALID_WHERE_CONDITION),
	)
}

func TestInspector_Advise_ObjectNotExist(t *testing.T) {
	runInspectCase(t, "use_database: database not exist", DefaultMysqlInspect(),
		"use no_exist_db",
		newTestResult().addResult(SCHEMA_NOT_EXIST, "no_exist_db"),
	)

	runInspectCase(t, "alter_table: schema not exist", DefaultMysqlInspect(),
		`
ALTER TABLE not_exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL;
`,
		newTestResult().addResult(SCHEMA_NOT_EXIST, "not_exist_db").
			addResult(TABLE_NOT_EXIST, "not_exist_db.exist_tb_1"),
	)
	runInspectCase(t, "alter_table: table not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.not_exist_tb_1 add column v5 varchar(255) NOT NULL;
`,
		newTestResult().addResult(TABLE_NOT_EXIST, "exist_db.not_exist_tb_1"),
	)

	runInspectCase(t, "create_table: schema not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists not_exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(SCHEMA_NOT_EXIST, "not_exist_db"),
	)

	runInspectCase(t, "select_from: schema not exist", DefaultMysqlInspect(),
		"select id from not_exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(SCHEMA_NOT_EXIST, "not_exist_db").
			addResult(TABLE_NOT_EXIST, "not_exist_db.exist_tb_1"),
	)
	runInspectCase(t, "select_from: table not exist", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1, exist_db.exist_tb_3 where id =1",
		newTestResult().addResult(TABLE_NOT_EXIST, "exist_db.exist_tb_3"),
	)

	runInspectCase(t, "delete: schema not exist", DefaultMysqlInspect(),
		"delete from not_exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(SCHEMA_NOT_EXIST, "not_exist_db").
			addResult(TABLE_NOT_EXIST, "not_exist_db.exist_tb_1"),
	)

	runInspectCase(t, "delete: table not exist", DefaultMysqlInspect(),
		"delete from exist_db.not_exist_tb_1 where id =1;",
		newTestResult().addResult(TABLE_NOT_EXIST, "exist_db.not_exist_tb_1"),
	)

	runInspectCase(t, "update: schema not exist", DefaultMysqlInspect(),
		"update not_exist_db.exist_tb_1 set v1='1' where id =1;",
		newTestResult().addResult(SCHEMA_NOT_EXIST, "not_exist_db").
			addResult(TABLE_NOT_EXIST, "not_exist_db.exist_tb_1"),
	)

	runInspectCase(t, "update: table not exist", DefaultMysqlInspect(),
		"update exist_db.not_exist_tb_1 set v1='1' where id =1;",
		newTestResult().addResult(TABLE_NOT_EXIST, "exist_db.not_exist_tb_1"),
	)
}

func TestInspector_Advise_ObjectExist(t *testing.T) {
	runInspectCase(t, "create_table: table exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(TABLE_EXIST, "exist_db.exist_tb_1"),
	)

	runInspectCase(t, "create_database: schema exist", DefaultMysqlInspect(),
		`
CREATE DATABASE exist_db;
`,
		newTestResult().addResult(SCHEMA_EXIST, "exist_db"),
	)
}

func TestInspector_Inspect_CreateTableStmt(t *testing.T) {
	runInspectCase(t, "create_table: ok", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: need \"if not exists\"", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CREATE_TABLE_NOT_EXIST),
	)

	runInspectCase(t, "create_table: using keyword", DefaultMysqlInspect(),
		"CREATE TABLE if not exists exist_db.`select` ("+
			"id bigint unsigned NOT NULL AUTO_INCREMENT,"+
			"v1 varchar(255) DEFAULT NULL,"+
			"v2 varchar(255) DEFAULT NULL,"+
			"PRIMARY KEY (id),"+
			"INDEX `create` (v1)"+
			")ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;",
		newTestResult().addResult(DDL_DISABLE_USING_KEYWORD, "select, create"),
	)
}

func TestInspector_InspectAlterTableStmt(t *testing.T) {
	runInspectCase(t, "alter_table: ok", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL;
`,
		newTestResult(),
	)

	runInspectCase(t, "alter_table: alter table need merge", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL;
ALTER TABLE exist_db.exist_tb_1 add primary key (id);
`,
		newTestResult(),
		newTestResult().addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE),
	)
}

func TestInspector_InspectCheck_Object_Name_Length(t *testing.T) {
	length64 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcd"
	length65 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcde"

	runInspectCase(t, "create_table: table length <= 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;`, length64),
		newTestResult(),
	)

	runInspectCase(t, "create_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
%s varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "create_table: index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
INDEX %s (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "alter_table: column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255);`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "alter_table: column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255);`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "alter_table: column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD index %s (v1);`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runInspectCase(t, "alter_table: column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME index v1_d TO %s;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)
}

func TestInspector_Inspect_Check_Primary_Key(t *testing.T) {
	runInspectCase(t, "create_table: primary key exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: primary key not exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_PRIMARY_KEY_EXIST),
	)

	runInspectCase(t, "create_table: primary key not auto increment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_PRIMARY_KEY_TYPE),
	)

	runInspectCase(t, "create_table: primary key not bigint unsigned", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_PRIMARY_KEY_TYPE),
	)
}

func TestInspector_Inspect_Check_String_Type(t *testing.T) {
	runInspectCase(t, "create_table: check char(20)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT,
	v1 char(20) DEFAULT NULL,
	v2 varchar(255) DEFAULT NULL,
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
	`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: check char(21)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT,
	v1 char(21) DEFAULT NULL,
	v2 varchar(255) DEFAULT NULL,
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
	`,
		newTestResult().addResult(DDL_CHECK_TYPE_CHAR_LENGTH),
	)
}

func TestInspector_Inspect_Check_Index(t *testing.T) {
	runInspectCase(t, "create_table: index <= 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
INDEX index_1 (id),
INDEX index_2 (id),
INDEX index_3 (id),
INDEX index_4 (id),
INDEX index_5 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: index > 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
INDEX index_1 (id),
INDEX index_2 (id),
INDEX index_3 (id),
INDEX index_4 (id),
INDEX index_5 (id),
INDEX index_6 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_INDEX_COUNT),
	)

	runInspectCase(t, "create_table: composite index columns <= 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
v3 varchar(255) DEFAULT NULL,
v4 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
INDEX index_1 (id,v1,v2,v3,v4)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: composite index columns > 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
v3 varchar(255) DEFAULT NULL,
v4 varchar(255) DEFAULT NULL,
v5 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
INDEX index_1 (id,v1,v2,v3,v4,v5)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_COMPOSITE_INDEX_MAX),
	)
}

func TestInspector_Inspect_Check_Index_Column_Type(t *testing.T) {
	runInspectCase(t, "create_table: disable index column blob (1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
b1 blob,
PRIMARY KEY (id),
INDEX index_b1 (b1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_DISABLE_INDEX_DATA_TYPE_BLOB),
	)

	runInspectCase(t, "create_table: disable index column blob (2)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
b1 blob UNIQUE KEY,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_DISABLE_INDEX_DATA_TYPE_BLOB),
	)
}

func TestInspector_Inspect_Check_Foreign_Key(t *testing.T) {
	runInspectCase(t, "create_table: has foreign key", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id),
FOREIGN KEY (id) REFERENCES exist_tb_1(id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_DISABLE_FOREIGN_KEY),
	)

}

func TestNewInspector(t *testing.T) {
	//conn, err := executor.NewConn("mysql", "root", "asd2010", "10.186.18.118", "23306", "sqle")
	//if err != nil {
	//	t.Error(err)
	//}
	//defer conn.Close()
	//sql1 := "select * from sqle.rules where name <> \"\";"
	//_ = sql1
	//sql2 := "insert into tb1 values(6,'v1');"
	//_ = sql2
	//result, err := conn.Exec(sql2)
	//if err != nil {
	//	t.Error(err)
	//}
	//ra, _ := result.RowsAffected()
	//fmt.Println("row_affects: ", ra)

	//result, err := conn.ShowMasterStatus()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//fmt.Println(result)
}