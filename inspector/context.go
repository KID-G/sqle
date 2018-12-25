package inspector

import "github.com/pingcap/tidb/ast"

type TableInfo struct {
	Size     float64
	sizeLoad bool

	// save create table parser object from db by query "show create table tb_1";
	// using in inspect and generate rollback sql
	CreateTableStmt *ast.CreateTableStmt

	// save alter table parse object from input sql;
	alterTableStmts []*ast.AlterTableStmt
}

type Context struct {
	// currentSchema will change after sql "use database"
	currentSchema string

	// originalSchemas save schemas info collected from db server
	originalSchemas map[string]struct{}
	// if schemas info has collected, set true
	schemaHasLoad bool

	// virtualSchemas save schema if create database ...
	virtualSchemas map[string] /*schema*/ struct{}

	// allTable save schema's tables info collected from db server
	// do not check every time from db server
	// if create a new schema, set schema's tables is null
	// if delete table, delete it.
	allTable   map[string] /*schema*/ map[string] /*table*/ *TableInfo
	counterDDL uint
	counterDML uint
}

func NewContext(currentSchema string) *Context {
	return &Context{
		currentSchema:   currentSchema,
		originalSchemas: map[string]struct{}{},
		allTable:        map[string]map[string]*TableInfo{},
		virtualSchemas:  map[string]struct{}{},
	}
}

func (c *Context) HasLoadSchema() bool {
	return c.schemaHasLoad
}

func (c *Context) HasSchema(name string) (has bool) {
	if c.HasLoadSchema() {
		_, has = c.originalSchemas[name]
	}
	if !has {
		_, has = c.virtualSchemas[name]
	}
	return
}

func (c *Context) LoadSchemas(schemas []string) {
	if !c.schemaHasLoad {
		for _, schema := range schemas {
			c.originalSchemas[schema] = struct{}{}
		}
		c.schemaHasLoad = true
	}
}

func (c *Context) CreateNewSchema(name string) {
	if c.HasSchema(name) {
		return
	}
	c.virtualSchemas[name] = struct{}{}
	c.allTable[name] = map[string]*TableInfo{}
}

func (c *Context) DeleteSchema(name string) {
	if !c.HasSchema(name) {
		return
	}
	delete(c.virtualSchemas, name)
	delete(c.originalSchemas, name)
	delete(c.allTable, name)
}

func (c *Context) HasLoadSchemaTables(schema string) (has bool) {
	_, has = c.allTable[schema]
	return
}

func (c *Context) LoadSchemaTables(schema string, tables []string) {
	if !c.HasSchema(schema) {
		return
	}
	if c.HasLoadSchemaTables(schema) {
		return
	}
	c.allTable[schema] = map[string]*TableInfo{}
	for _, table := range tables {
		c.allTable[schema][table] = &TableInfo{
			alterTableStmts: []*ast.AlterTableStmt{},
		}
	}
}

func (c *Context) CreateNewTable(schema, table string, info *TableInfo) {
	if !c.HasSchema(schema) {
		return
	}
	_, exist := c.GetTableInfo(schema, table)
	if exist {
		return
	}
	c.allTable[schema][table] = info
}

func (c *Context) DeleteTable(schema, table string) {
	if c.HasSchema(schema) {
		delete(c.allTable[schema], table)
	}
}

func (c *Context) GetTableInfo(schema, table string) (info *TableInfo, exist bool) {
	if c.HasSchema(schema) {
		info, exist = c.allTable[schema][table]
	} else {
		info, exist = nil, false
	}
	return
}

func (c *Context) UseSchema(schema string) {
	if c.HasSchema(schema) {
		c.currentSchema = schema
	}
}

func (c *Context) AddDDL() {
	c.counterDDL += 1
}

func (c *Context) GetDDLCounter() uint {
	return c.counterDDL
}

func (c *Context) AddDML() {
	c.counterDML += 1
}

func (c *Context) GetDMLCounter() uint {
	return c.counterDML
}