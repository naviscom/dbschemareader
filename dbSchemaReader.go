package dbschemareader
//package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

)
type SessionTableSpecs struct {
	TableName				string
	TokenColumnName 		string
	SessionIdColumnName		string
	CreateTimingColumnName	string
	ExpiryTimingColumnName	string
	StatusColumnName		string
	AgentColumnName			string
	DeviceColumnName		string
	LocationColumnName		string
	SecurityColumnName		string
	FkUserColumn			string
	FkUserColumn2			string
	RefUserTableName		string
	RefUserTableColumn		string
}

type UserTableSpecs struct {
	TableName				string
	AuthColumnName			string
	DateColumnName			string
	StatusColumnName		string
	LoginNameColumnName		string
	SecurityColumnName		string
}

type Table_Struct struct {
	Table_name          			string
	Table_Columns       			[]Table_columns
	IndexDetails        			[]index_name_details
	CompositeForeignKeys			[]CompositeForeignKeysAndReferences
	CompositeUniqueConstraints		[]CompositeUniqueConstraint
	CheckConstraints				[]CheckConstraint
	CompositePrimaryConstraint		string
	ForeignKeys						[]foreign_key_details
	OutputFileName      			string
	FunctionSignature   			string
	FunctionSignature2  			string
	FunctionSignature3  			string
	IsSessionsTable					bool
	IsUserTable						bool
	SessionTableSpecs				SessionTableSpecs
	UserTableSpecs					UserTableSpecs
}

type CompositeForeignKeysAndReferences struct {
	ConstraintName       			string
	CompositeForeignKeys			string
	CompositeForeignKeysReferences 	string
	CompositeForeignKeysOnClause	string
	FK_Related_TableName					string
	FK_Related_SingularTableName			string
	FK_Related_TableName_Singular_Object	string
	FK_Related_TableName_Plural_Object		string
	FK_Related_Table_Column					string
}

type CompositeUniqueConstraint struct {
	ConstraintName       	string
	ConstraintTableName		string
	ConstraintColumns     	string
}

type CheckConstraint struct {
	CheckConstraintName			string
	CheckConstraintValue		string
	CheckConstraintTableName	string
	CheckConstraintColumnName	string
}

type Table_columns struct {
	Column_name     string
	ColumnType      string
	PrimaryFlag     bool
	UniqueFlag      bool
	ForeignFlag		bool
	Not_Null		bool
	DefaultValue	string
	RelatedTable	foreign_key_details
	InLineCheck		string
	ColumnNameParams string
}

type index_name_details struct {
	IndexName       string
	IndexColumn     []string
}

type foreign_key_details struct {
	FK_Column       						string
	FK_Related_TableName					string
	FK_Related_SingularTableName			string
	FK_Related_TableName_Singular_Object	string
	FK_Related_TableName_Plural_Object		string
	FK_Related_Table_Column					string
	FK_OnClause								string
}

///////////////////////////////////////////////////

type FK_Hierarchy struct {
	TableName				string
	RelatedTablesLevels		[]RelatedTables
}

type RelatedTables struct {
	Hierarchy_TableName	string
	RelatedTableList	[]RelatedTable
}

type RelatedTable struct {
	FK_Related_TableName					string
	FK_Related_SingularTableName			string
	FK_Related_TableName_Singular_Object	string
	FK_Related_TableName_Plural_Object		string
	FK_Related_Table_Column					string
}

// Alternative simpler approach if you prefer not to use regex
func extractColumnNameSimple(constraintExpr string) string {
	expr := strings.TrimSpace(constraintExpr)
	
	// Find the first operator (ordered by length to avoid substring issues)
	operators := []string{">=", "<=", "<>", "!=", ">", "<", "="}
	
	minIndex := len(expr)
	foundOperator := ""
	
	// Find the earliest occurring operator
	for _, op := range operators {
		if idx := strings.Index(expr, op); idx != -1 && idx < minIndex {
			minIndex = idx
			foundOperator = op
		}
	}
	
	if foundOperator != "" {
		columnName := strings.TrimSpace(expr[:minIndex])
		columnName = strings.TrimSpace(columnName)
		return columnName
	}
	
	return strings.TrimSpace(expr)
}

func extractColumnNamesWithParenthesis(line []string, columnIndex int) string {
	var element, KeyColumn string
	var nestLevel,y int
	var nestComplete bool
	for {
		//fmt.Println("line from extractColumnNamesWithParenthesis",line, columnIndex)	
		if columnIndex < len(line){
			element = element + line[columnIndex]
			//fmt.Println("line from extractColumnNamesWithParenthesis", "element: ", element ,"columnIndex: ", columnIndex)
		}else {
			break
		}
		nestLevel = 0
		for y=0; y < len(element); y++{			
			if element[y:y+1] == "(" {
				nestLevel++
			} else if element[y:y+1] == ")" {
				nestLevel--
				if nestLevel == 0 {
					nestComplete = true
					break
				}else{
					nestComplete = false
				}
			}
		}
		if nestComplete {
			KeyColumn = element
			break
		}else{
			columnIndex++
			if columnIndex < len(line) {
				line[columnIndex] = " "+line[columnIndex]
			}
		}
	}
	if nestComplete {
		if KeyColumn[len(KeyColumn)-1:] == "," {
			KeyColumn = KeyColumn[0:len(KeyColumn)-1]
		}	
	}
	// fmt.Println("KeyColumn: ", KeyColumn)
	return KeyColumn
}

func columnNameCleanUp(columnName string) string {

	if columnName[len(columnName)-1:] == ";" {
		columnName = columnName[0:len(columnName)-1]
	}
	if columnName[len(columnName)-1:] == "," {
		columnName = columnName[0:len(columnName)-1]
	}
	if columnName[0:1] == "'" {
		columnName = columnName[1:len(columnName)-1]
	}	
	if columnName[0:1] == `"` {
		columnName = columnName[1:len(columnName)-1]
	}
	if columnName[0:1] == "(" {
		columnName = columnName[1:len(columnName)-1]
		if columnName[0:1] == "`" {
			columnName = columnName[1:len(columnName)-1]
		}
		if columnName[0:1] == "[" {
			columnName = columnName[1:len(columnName)-1]
		}	
		if columnName[0:1] == "'" {
			columnName = columnName[1:len(columnName)-1]
		}
		if columnName[0:1] == `"` {
			columnName = columnName[1:len(columnName)-1]
		}
		if strings.Contains(columnName, "[") || strings.Contains(columnName, "]") {
			// Remove all square brackets from the string
			columnName = strings.ReplaceAll(columnName, "[", "")
			columnName = strings.ReplaceAll(columnName, "]", "")
		}
		if strings.Contains(columnName, `"`) {
			// Remove all " from the string
			columnName = strings.ReplaceAll(columnName, `"`, "")
		}
	}
	return columnName
}

func ReadSchema(filePath string, tableX []Table_Struct)  ([]Table_Struct, []FK_Hierarchy) {
//func ReadSchema(filePath string, tableX []sqlcq_test.Table_Struct)  ([]sqlcq_test.Table_Struct, []FK_Hierarchy) {
	fmt.Println("Time: ",time.Now())
	// fmt.Println("filePath: ",filePath, "tableX: ",tableX)
	// var tableX []Table_Struct
	var table Table_Struct
	// fmt.Println(filePath)
	readFile, err := os.Open(filePath)
	if err != nil {
	  fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		res1 := strings.Split(fileScanner.Text(), " ")		
		// fmt.Println("res1: ",res1, "len(res1): ",len(res1))
		if len(res1) > 1 {
			if res1[0] == "CREATE" && res1[1] == "TABLE" {
				// fmt.Println(`Inside "CREATE" && res1[1] == "TABLE"`)
				// fmt.Println(res1[2], strings.TrimSpace(res1[2][1:len(res1[2])-1]))
				//fmt.Println("table.Table_name: ", table.Table_name)
				table.Table_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
				table.FunctionSignature2 = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
				// fmt.Println("table.Table_name: ", table.Table_name, "table.FunctionSignature2: ", table.FunctionSignature2)
				if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
					table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-3])+"y"
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-3]+"y")
					// fmt.Println("1-table.OutputFileName: ", table.OutputFileName, "table.FunctionSignature: ", table.FunctionSignature)
				} else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
					table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-1])
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-1])
					// fmt.Println("2-table.OutputFileName: ", table.OutputFileName, "table.FunctionSignature: ", table.FunctionSignature)
				} else {
					table.OutputFileName = table.Table_name
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
					// fmt.Println("3-table.OutputFileName: ", table.OutputFileName, "table.FunctionSignature: ", table.FunctionSignature)
				}
				table.Table_Columns = nil
			}
			if res1[0] == "" && res1[1] == "" && strings.TrimSpace(res1[2][0:1]) == `"` {
				var tabColumns Table_columns

				// fmt.Println(`Inside """ && res1[1] == "" && strings.TrimSpace(res1[2][0:1]"`)
				// fmt.Println("res1: ", res1)
				tabColumns.Column_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
				tabColumns.Column_name = strings.ReplaceAll(tabColumns.Column_name, "__", "_")
				tabColumns.ColumnType = strings.TrimSpace(res1[3][0:])
				tabColumns.ColumnType = strings.TrimSpace(strings.ReplaceAll(tabColumns.ColumnType, ",", ""))
				tabColumns.ForeignFlag = false
				// fmt.Println("tabColumns.Column_name: ", tabColumns.Column_name, "tabColumns.ColumnType: ", tabColumns.ColumnType, "tabColumns.ForeignFlag: ", tabColumns.ForeignFlag)
				///////////////////
				column_name_slice := strings.Split(tabColumns.Column_name,"_")
				// fmt.Println(column_name_slice)
				for k := 0; k < len(column_name_slice); k++ {
					if column_name_slice[k] == "id" {
						column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k]))
					}else{
						column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k][0:1]))+strings.TrimSpace(column_name_slice[k][1:])
					}
				}
				tabColumns.ColumnNameParams = strings.Join(column_name_slice,"")
				// fmt.Println("tabColumns.ColumnNameParams: ", tabColumns.ColumnNameParams)
				var primaryWordIndex, uniqueWordIndex, notWordIndex, nullWordIndex, defaultWordIndex int
				// fmt.Println("res1: ", res1)
				for i, element := range res1 {
					if element == "PRIMARY" {
						primaryWordIndex = i
					}
					if element == "UNIQUE" || element == "UNIQUE," {
						uniqueWordIndex = i
					}
					if element == "NOT" {
						notWordIndex = i
					}
					if element == "NULL" || element == "NULL," {
						nullWordIndex = i
					}
					if element == "DEFAULT" {
						defaultWordIndex = i
					}
				}
				if primaryWordIndex > 0 {
					tabColumns.PrimaryFlag = true
					// fmt.Println("PrimaryFlag")
				}else{ tabColumns.PrimaryFlag = false}

				if uniqueWordIndex > 0 {
					tabColumns.UniqueFlag = true
					// fmt.Println("UniqueFlag")
				}else{ tabColumns.UniqueFlag = false}

				if defaultWordIndex > 0 {
					// fmt.Println("DefaultValue", "-",strings.TrimSpace(res1[defaultWordIndex+1]), "-")
					tabColumns.DefaultValue = strings.ReplaceAll(res1[defaultWordIndex+1], ",", "")
					    if strings.HasPrefix(tabColumns.DefaultValue, "(") && strings.HasSuffix(tabColumns.DefaultValue, ")") {
        					tabColumns.DefaultValue = strings.TrimPrefix(tabColumns.DefaultValue, "(")
       						tabColumns.DefaultValue = strings.TrimSuffix(tabColumns.DefaultValue, ")")
						}
					// fmt.Println("DefaultValue", "-",tabColumns.DefaultValue, "-")
				}

				if notWordIndex > 0 && nullWordIndex > 0 {
					tabColumns.Not_Null = true
				}else{
					tabColumns.Not_Null = false
				}
				//fmt.Println("tabColumns: ", tabColumns)
				table.Table_Columns = append(table.Table_Columns, tabColumns)
			}
			if res1[0] == "CREATE" && res1[1] == "INDEX" {
				var onIndex int //usingIndex, indexIndex
				for x, element := range res1{
					if element == "ON"{
						onIndex = x
					}
				}
					
				// fmt.Println(`Inside "CREATE" && res1[1] == "INDEX"`)
				// fmt.Println(res1)				
				for i:=0; i<len(tableX); i++{
					if tableX[i].Table_name == strings.TrimSpace(res1[3][1:len(res1[3])-1]) { 
						var index index_name_details
						index.IndexName = strings.TrimSpace(res1[3][1:len(res1[3])-1]) + strconv.Itoa(rand.Intn(90000))
			
						indexColumnName := extractColumnNamesWithParenthesis(res1, onIndex+2)
						indexColumnName = columnNameCleanUp(indexColumnName)
						indexColumnName = strings.ReplaceAll(indexColumnName, `"`, "")
						// indexColumnName = strings.ReplaceAll(indexColumnName, " ", "")
						// fmt.Println("indexColumnName: ", indexColumnName)
						index.IndexColumn = append(index.IndexColumn,   indexColumnName)
						tableX[i].IndexDetails = append(tableX[i].IndexDetails, index)
						//index.IndexColumn = nil
						//index.IndexName = ""
					}
				}
			}
			if res1[0] == "ALTER" && res1[1] == "TABLE" {
				// fmt.Println(`Inside "ALTER" && res1[2] == "TABLE"`)
				// fmt.Println(res1)
				var tableWordIndex, referenceWordIndex, keyWordIndex, onWordIndex, checkWordIndex, uniqueWordIndex, constraintWordIndex, foreignWordIndex int //, tableWordIndex, foreignWordIndex,  primaryWordIndex, onlyWordIndex,  int
				for i, element := range res1 {
					if element == "ALTER" {
						//alterWordIndex = i
					}
					if element == "TABLE" {
						tableWordIndex = i					
					}
					// if element == "ADD" {
					// 	//addWordIndex = i
					// }
					if element == "FOREIGN" {
						foreignWordIndex = i
					}
					if element == "KEY" {
						keyWordIndex = i
					}
					if element == "REFERENCES" {
						referenceWordIndex = i
					}
					if element == "ON" {
						onWordIndex = i
					}
					if element == "CHECK" {
						checkWordIndex = i
					}
					if element == "UNIQUE" {
						uniqueWordIndex = i
					}
					if element == "CONSTRAINT" {
						constraintWordIndex = i
					}
				}
				// fmt.Println(tableWordIndex, referenceWordIndex, keyWordIndex, onWordIndex, checkWordIndex, uniqueWordIndex, constraintWordIndex, foreignWordIndex)
				tableName := strings.ReplaceAll(res1[tableWordIndex+1],`"`, "")
				var indexOfTableX int
				for i:=0; i<len(tableX); i++{
					if tableX[i].Table_name == tableName {
						indexOfTableX = i
					}
				}
				if constraintWordIndex > 0 && uniqueWordIndex > 0 {
					var compositeUniqueConstraint CompositeUniqueConstraint
					compositeUniqueConstraint.ConstraintName = res1[constraintWordIndex+1]
					compositeUniqueConstraint.ConstraintTableName = tableX[indexOfTableX].Table_name
					compositeUniqueConstraint.ConstraintColumns = extractColumnNamesWithParenthesis(res1, uniqueWordIndex +1)
					compositeUniqueConstraint.ConstraintColumns = columnNameCleanUp(compositeUniqueConstraint.ConstraintColumns)
					tableX[indexOfTableX].CompositeUniqueConstraints = append(tableX[indexOfTableX].CompositeUniqueConstraints, compositeUniqueConstraint)
					// fmt.Println("last appended CompositeUniqueConstraints: ", tableX[indexOfTableX].CompositeUniqueConstraints[len(tableX[indexOfTableX].CompositeUniqueConstraints)-1])
				}
				if constraintWordIndex > 0 && checkWordIndex > 0 {
					var checkConstraint CheckConstraint
					checkConstraint.CheckConstraintName = res1[constraintWordIndex+1]
					checkConstraint.CheckConstraintTableName = tableX[indexOfTableX].Table_name
					// fmt.Println("res1: ", res1,  "res1[checkWordIndex +1]: ", res1[checkWordIndex +1])
					checkConstraint.CheckConstraintColumnName = strings.ReplaceAll(res1[checkWordIndex +1], "(", "")
					checkConstraint.CheckConstraintColumnName = extractColumnNameSimple(checkConstraint.CheckConstraintColumnName)
					// fmt.Println("checkConstraint.CheckConstraintColumnName: ", checkConstraint.CheckConstraintColumnName)
					checkConstraint.CheckConstraintValue = extractColumnNamesWithParenthesis(res1, checkWordIndex +1)
					checkConstraint.CheckConstraintValue = columnNameCleanUp(checkConstraint.CheckConstraintValue)
					CheckConstraintValueArray := strings.Split(strings.TrimSpace(checkConstraint.CheckConstraintValue)," ")
					checkConstraint.CheckConstraintColumnName = strings.TrimSpace(CheckConstraintValueArray[0])
					// if checkConstraint.CheckConstraintTableName == "users" {
					// 	fmt.Println("res1: ", res1)
					// 	fmt.Println("res1[checkWordIndex +1]: ", res1[checkWordIndex +1])
					// 	fmt.Println("checkConstraint.CheckConstraintName: ", checkConstraint.CheckConstraintName)
					// 	fmt.Println("checkConstraint.CheckConstraintTableName: ", checkConstraint.CheckConstraintTableName)
					// 	fmt.Println("checkConstraint.CheckConstraintColumnName: ", checkConstraint.CheckConstraintColumnName)
					// 	fmt.Println("checkConstraint.CheckConstraintValue: ", checkConstraint.CheckConstraintValue)
	
					// }

					tableX[indexOfTableX].CheckConstraints = append(tableX[indexOfTableX].CheckConstraints, checkConstraint)
					//fmt.Println("last appended CheckConstraints: ", tableX[indexOfTableX].CheckConstraints[len(tableX[indexOfTableX].CheckConstraints)-1])
				}
				if foreignWordIndex > 0 && keyWordIndex > 0 && referenceWordIndex > 0 {
					var fkDetails foreign_key_details
					var compositeForeignKeysAndReferences CompositeForeignKeysAndReferences
					var onClause string
					foreignKeyColumn := extractColumnNamesWithParenthesis(res1, keyWordIndex+1)
					foreignKeyColumn = columnNameCleanUp(foreignKeyColumn)
					foreignKeyColumn = strings.ReplaceAll(foreignKeyColumn, " ", "")
					// fmt.Println("foreignKeyColumn: ", foreignKeyColumn)

					relatedTable := res1[referenceWordIndex+1]
					relatedTable = columnNameCleanUp(relatedTable)
					// fmt.Println("relatedTable: ", relatedTable)

					relatedTableColumn := extractColumnNamesWithParenthesis(res1, referenceWordIndex+2)
					relatedTableColumn = columnNameCleanUp(relatedTableColumn)
					relatedTableColumn = strings.ReplaceAll(relatedTableColumn, " ", "")
					// fmt.Println("relatedTableColumn: ", relatedTableColumn)

					if onWordIndex > 0 {
						onClauseSlice := res1[onWordIndex+1:]
						onClause = strings.Join(onClauseSlice, " ")
						onClause = strings.ReplaceAll(onClause, ";", "")
						// fmt.Println("onClause: ", onClause)
					}

					if strings.Contains(foreignKeyColumn, ",") {
						compositeForeignKeysAndReferences.CompositeForeignKeys = foreignKeyColumn
						compositeForeignKeysAndReferences.CompositeForeignKeysReferences = relatedTable + " " + relatedTableColumn
						compositeForeignKeysAndReferences.CompositeForeignKeysOnClause = onClause

						compositeForeignKeysAndReferences.FK_Related_TableName = relatedTable
						if strings.TrimSpace(relatedTable[len(relatedTable)-3:]) == `ies` {
							compositeForeignKeysAndReferences.FK_Related_SingularTableName = strings.TrimSpace(relatedTable[:len(relatedTable)-3])+"y"
						}else if strings.TrimSpace(relatedTable[len(relatedTable)-1:]) == `s` {
							compositeForeignKeysAndReferences.FK_Related_SingularTableName = strings.TrimSpace(relatedTable[:len(relatedTable)-1])
						}else {
							compositeForeignKeysAndReferences.FK_Related_SingularTableName = relatedTable
						}	
						if strings.TrimSpace(relatedTable[len(relatedTable)-3:]) == `ies` {
							compositeForeignKeysAndReferences.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(relatedTable[0:1]))+strings.TrimSpace(relatedTable[1:len(relatedTable)-3]+"y")
						} else if strings.TrimSpace(relatedTable[len(relatedTable)-1:]) == `s` {
							compositeForeignKeysAndReferences.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(relatedTable[0:1]))+strings.TrimSpace(relatedTable[1:len(relatedTable)-1])
						} else {
							compositeForeignKeysAndReferences.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(relatedTable[0:1]))+strings.TrimSpace(relatedTable[1:])
						}						
						compositeForeignKeysAndReferences.FK_Related_TableName_Plural_Object = strings.ToUpper(strings.TrimSpace(relatedTable[0:1]))+strings.TrimSpace(relatedTable[1:])
						compositeForeignKeysAndReferences.FK_Related_Table_Column = relatedTableColumn


						tableX[indexOfTableX].CompositeForeignKeys = append(tableX[indexOfTableX].CompositeForeignKeys, compositeForeignKeysAndReferences)
						// fmt.Println("last appended CompositeForeignKeys: ", tableX[indexOfTableX].CompositeForeignKeys[len(tableX[indexOfTableX].CompositeForeignKeys)-1])
					} else {
						fkDetails.FK_Column = foreignKeyColumn
						fkDetails.FK_Related_TableName = relatedTable
						fkDetails.FK_Related_Table_Column = relatedTableColumn
						fkDetails.FK_OnClause = onClause

						if strings.TrimSpace(fkDetails.FK_Related_TableName[len(fkDetails.FK_Related_TableName)-3:]) == `ies` {
							fkDetails.FK_Related_SingularTableName = strings.TrimSpace(fkDetails.FK_Related_TableName[:len(fkDetails.FK_Related_TableName)-3])+"y"
						}else if strings.TrimSpace(fkDetails.FK_Related_TableName[len(fkDetails.FK_Related_TableName)-1:]) == `s` {
							fkDetails.FK_Related_SingularTableName = strings.TrimSpace(fkDetails.FK_Related_TableName[:len(fkDetails.FK_Related_TableName)-1])
						}else {
							fkDetails.FK_Related_SingularTableName = fkDetails.FK_Related_TableName
						}	
						if strings.TrimSpace(fkDetails.FK_Related_TableName[len(fkDetails.FK_Related_TableName)-3:]) == `ies` {
							fkDetails.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(fkDetails.FK_Related_TableName[0:1]))+strings.TrimSpace(fkDetails.FK_Related_TableName[1:len(fkDetails.FK_Related_TableName)-3]+"y")
						} else if strings.TrimSpace(fkDetails.FK_Related_TableName[len(fkDetails.FK_Related_TableName)-1:]) == `s` {
							fkDetails.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(fkDetails.FK_Related_TableName[0:1]))+strings.TrimSpace(fkDetails.FK_Related_TableName[1:len(fkDetails.FK_Related_TableName)-1])
						} else {
							fkDetails.FK_Related_TableName_Singular_Object = strings.ToUpper(strings.TrimSpace(fkDetails.FK_Related_TableName[0:1]))+strings.TrimSpace(fkDetails.FK_Related_TableName[1:])
						}						
						fkDetails.FK_Related_TableName_Plural_Object = strings.ToUpper(strings.TrimSpace(fkDetails.FK_Related_TableName[0:1]))+strings.TrimSpace(fkDetails.FK_Related_TableName[1:])

						tableX[indexOfTableX].ForeignKeys = append(tableX[indexOfTableX].ForeignKeys, fkDetails)
						for j:=0; j<len(tableX[indexOfTableX].Table_Columns); j++ {
							if tableX[indexOfTableX].Table_Columns[j].Column_name == fkDetails.FK_Column {
								tableX[indexOfTableX].Table_Columns[j].ForeignFlag = true
							}
						}
					}
				}
				// time.Sleep(2 * time.Second)
			}
			if res1[0] == "ALTER" && res1[3] == "ADD" && res1[4] == "COLUMN" {
				// fmt.Println(`Inside "ALTER" && res1[3] == "ADD" && res1[4] == "COLUMN"`)
				// fmt.Println(res1)
				for i:=0; i<len(tableX); i++{
					var tabColumns Table_columns
					if tableX[i].Table_name == strings.TrimSpace(res1[2][1:len(res1[2])-1]) { 
						tabColumns.Column_name = strings.TrimSpace(res1[5][1:len(res1[5])-1])
						tabColumns.ColumnType = strings.TrimSpace(res1[6][0:])
						tabColumns.ForeignFlag = false
						///////////////////
						column_name_slice := strings.Split(tabColumns.Column_name,"_")
						for k := 0; k < len(column_name_slice); k++ {
							if column_name_slice[k] == "id" {
								column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k]))
							}else{
								column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k][0:1]))+strings.TrimSpace(column_name_slice[k][1:])
						
							}
						}
						tabColumns.ColumnNameParams = strings.Join(column_name_slice,"")
						//////////////////
						if len(res1) > 7 {
							if res1[7] == `PRIMARY` {
								tabColumns.PrimaryFlag = true
							} else{
								tabColumns.PrimaryFlag = false
							}
							if res1[7] == `UNIQUE` {
								tabColumns.UniqueFlag = true
							} else{
								tabColumns.UniqueFlag = false            
							}
						} else {
							tabColumns.PrimaryFlag = false
							tabColumns.UniqueFlag = false
						}
						tableX[i].Table_Columns = append(tableX[i].Table_Columns, tabColumns)
					}
				}
			}
		}
		if len(res1) == 1 {
			if res1[0] == ");" {
				// fmt.Println(`Inside ");"`)
				// fmt.Println(res1)
				tableX = append(tableX, table)
			}
		}
	}
	// fmt.Println("I am here")



	for i:=0; i<len(tableX); i++{
		// fmt.Println("table Name: ", tableX[i].Table_name, "OutputFileName: ", tableX[i].OutputFileName,  "FunctionSignature: ", tableX[i].FunctionSignature, "FunctionSignature2: ", tableX[i].FunctionSignature2)
		for j:=0; j<len(tableX[i].Table_Columns); j++{
			// fmt.Println("    column name: ", tableX[i].Table_Columns[j].Column_name, "Type: " , tableX[i].Table_Columns[j].ColumnType, "PrimaryKeyFlag: ", tableX[i].Table_Columns[j].PrimaryFlag, "UniqueKeyFlag: ", tableX[i].Table_Columns[j].UniqueFlag, "ForeignKeyFlag: ", tableX[i].Table_Columns[j].ForeignFlag, "ColumnNameForParams: ", tableX[i].Table_Columns[j].ColumnNameParams)
		}
		for j:=0; j<len(tableX[i].IndexDetails); j++{
			// fmt.Println("    index name: ", tableX[i].IndexDetails[j].IndexName)
			for k:=0; k<len(tableX[i].IndexDetails[j].IndexColumn); k++{
				// fmt.Println("    index column name: ", tableX[i].IndexDetails[j].IndexColumn[k])
			}
		}
		for j:=0; j<len(tableX[i].ForeignKeys); j++{
			// fmt.Println("tableX[i].Table_name: ",tableX[i].Table_name)
			// fmt.Println("    FK_Column: ", tableX[i].ForeignKeys[j].FK_Column, "FK_Related_TableName: ", tableX[i].ForeignKeys[j].FK_Related_TableName, "FK_Related_SingularTableName: ", tableX[i].ForeignKeys[j].FK_Related_SingularTableName, "FK_Related_Table_Column: ", tableX[i].ForeignKeys[j].FK_Related_Table_Column)
		}
		for j:=0; j<len(tableX[i].CompositeForeignKeys); j++{
			// fmt.Println("tableX[i].Table_name: ",tableX[i].Table_name)
			// fmt.Println("    FK_Related_TableName: ", tableX[i].CompositeForeignKeys[j].FK_Related_TableName, "FK_Related_SingularTableName: ", tableX[i].CompositeForeignKeys[j].FK_Related_SingularTableName, "FK_Related_Table_Column: ", tableX[i].CompositeForeignKeys[j].FK_Related_Table_Column)
		}
	}




	////////////////////Building Relationship Hierarchy//////////////////////////
	var relatedTables RelatedTables
	var fk_Hierarchy FK_Hierarchy
	var fk_HierarchyX []FK_Hierarchy
	// fmt.Println("Now I will start populating the relationship tree")
	for i:=0; i<len(tableX); i++{
		// fmt.Println("Table No.: ",i+1,"tableX.Table_name", tableX[i].Table_name, "len(tableX[i].ForeignKeys): ", len(tableX[i].ForeignKeys))
		relatedTables.Hierarchy_TableName = ""
		relatedTables.RelatedTableList = nil
		fk_Hierarchy.RelatedTablesLevels = nil
		if len(tableX[i].ForeignKeys) > 0 {
			// fmt.Println("inside if  len(tableX[i].ForeignKeys) > 0", "tableX[i].ForeignKeys: ", tableX[i].ForeignKeys)
			// fmt.Println("I am going inside the loop..")
			// time.Sleep(2 * time.Second)
			for j :=0; j < len(tableX[i].ForeignKeys); j++ {
				var relatedTable RelatedTable
				// fmt.Println("at the start of loop---->relatedTable: ",relatedTable)
				// time.Sleep(2 * time.Second)
				relatedTable.FK_Related_TableName = tableX[i].ForeignKeys[j].FK_Related_TableName
				relatedTable.FK_Related_SingularTableName = tableX[i].ForeignKeys[j].FK_Related_SingularTableName
				relatedTable.FK_Related_Table_Column = tableX[i].ForeignKeys[j].FK_Related_Table_Column
				relatedTable.FK_Related_TableName_Singular_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Singular_Object
				relatedTable.FK_Related_TableName_Plural_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Plural_Object
				relatedTables.Hierarchy_TableName = tableX[i].Table_name
				relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
				// fmt.Println("towards the end of loop---->relatedTable: ",relatedTable)
				// fmt.Println("towards the end of loop---->relatedTables: ",relatedTables)
				// time.Sleep(2 * time.Second)
			}
			fk_Hierarchy.TableName = tableX[i].Table_name
			fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
			// fmt.Println("fk_Hierarchy.RelatedTablesLevels: ", fk_Hierarchy.RelatedTablesLevels)
			fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
			// fmt.Println("fk_HierarchyX[len(fk_HierarchyX)-1]: ", fk_HierarchyX[len(fk_HierarchyX)-1])
			// time.Sleep(2 * time.Second)
		// 	////following is the new addition/////
		// 	if len(tableX[i].CompositeForeignKeys) > 0 {
		// 		for j :=0; j < len(tableX[i].CompositeForeignKeys); j++ {
		// 			var relatedTable RelatedTable
		// 			// fmt.Println("at the start of loop---->relatedTable: ",relatedTable)
		// 			// time.Sleep(2 * time.Second)
		// 			// compositeForeignKeysReferences := strings.Split(tableX[i].CompositeForeignKeys[j].CompositeForeignKeysReferences, " ")

		// 			relatedTable.FK_Related_TableName = tableX[i].CompositeForeignKeys[j].FK_Related_TableName
		// 			relatedTable.FK_Related_SingularTableName = tableX[i].CompositeForeignKeys[j].FK_Related_SingularTableName
		// 			relatedTable.FK_Related_TableName_Singular_Object = tableX[i].CompositeForeignKeys[j].FK_Related_TableName_Singular_Object
		// 			relatedTable.FK_Related_TableName_Plural_Object = tableX[i].CompositeForeignKeys[j].FK_Related_TableName_Plural_Object
		// 			relatedTable.FK_Related_Table_Column = tableX[i].CompositeForeignKeys[j].FK_Related_Table_Column
		// 			relatedTables.Hierarchy_TableName = tableX[i].Table_name
		// 			relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
		// 			// fmt.Println("towards the end of loop---->relatedTable: ",relatedTable)
		// 			// fmt.Println("towards the end of loop---->relatedTables: ",relatedTables)
		// 			// time.Sleep(2 * time.Second)
		// 		}
		// 		fk_Hierarchy.TableName = tableX[i].Table_name
		// 		fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
		// 		// fmt.Println("fk_Hierarchy.RelatedTablesLevels: ", fk_Hierarchy.RelatedTablesLevels)
		// 		fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
		// 		// fmt.Println("fk_HierarchyX[len(fk_HierarchyX)-1]: ", fk_HierarchyX[len(fk_HierarchyX)-1])
		// 		// time.Sleep(2 * time.Second)
		// 	}
		// } else if len(tableX[i].ForeignKeys) == 0 && len(tableX[i].CompositeForeignKeys) > 0 {
		// 	if len(tableX[i].CompositeForeignKeys) > 0 {
		// 		for j :=0; j < len(tableX[i].CompositeForeignKeys); j++ {
		// 			var relatedTable RelatedTable
		// 			// fmt.Println("at the start of loop---->relatedTable: ",relatedTable)
		// 			// time.Sleep(2 * time.Second)
		// 			// compositeForeignKeysReferences := strings.Split(tableX[i].CompositeForeignKeys[j].CompositeForeignKeysReferences, " ")
		// 			// relatedTable.FK_Related_TableName = compositeForeignKeysReferences[0]
		// 			// relatedTable.FK_Related_SingularTableName = compositeForeignKeysReferences[0]
		// 			// relatedTable.FK_Related_TableName_Singular_Object = compositeForeignKeysReferences[0]
		// 			// relatedTable.FK_Related_TableName_Plural_Object = compositeForeignKeysReferences[0]
		// 			// relatedTable.FK_Related_Table_Column = compositeForeignKeysReferences[1]
	
		// 			relatedTable.FK_Related_TableName = tableX[i].CompositeForeignKeys[j].FK_Related_TableName
		// 			relatedTable.FK_Related_SingularTableName = tableX[i].CompositeForeignKeys[j].FK_Related_SingularTableName
		// 			relatedTable.FK_Related_TableName_Singular_Object = tableX[i].CompositeForeignKeys[j].FK_Related_TableName_Singular_Object
		// 			relatedTable.FK_Related_TableName_Plural_Object = tableX[i].CompositeForeignKeys[j].FK_Related_TableName_Plural_Object
		// 			relatedTable.FK_Related_Table_Column = tableX[i].CompositeForeignKeys[j].FK_Related_Table_Column
	
		// 			relatedTables.Hierarchy_TableName = tableX[i].Table_name
		// 			relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
		// 			// fmt.Println("towards the end of loop---->relatedTable: ",relatedTable)
		// 			// fmt.Println("towards the end of loop---->relatedTables: ",relatedTables)
		// 			// time.Sleep(2 * time.Second)
		// 		}
		// 		fk_Hierarchy.TableName = tableX[i].Table_name
		// 		fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
		// 		// fmt.Println("fk_Hierarchy.RelatedTablesLevels: ", fk_Hierarchy.RelatedTablesLevels)
		// 		fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
		// 		// fmt.Println("fk_HierarchyX[len(fk_HierarchyX)-1]: ", fk_HierarchyX[len(fk_HierarchyX)-1])
		// 		// time.Sleep(2 * time.Second)
		// 	}

		} else {
			// fmt.Println("inside else  len(tableX[i].ForeignKeys) > 0")
			fk_Hierarchy.TableName = tableX[i].Table_name
			fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
			// fmt.Println("fk_Hierarchy.RelatedTablesLevels: ", fk_Hierarchy.RelatedTablesLevels)
			fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
			// fmt.Println("fk_HierarchyX[len(fk_HierarchyX)-1]: ", fk_HierarchyX[len(fk_HierarchyX)-1])
			// time.Sleep(2 * time.Second)
		}
		// fmt.Println("Now i am outside..........")
		// time.Sleep(3 * time.Second)
		var c int
		var d int
		var e int
		c = 0
		// fmt.Println("i: ",i,"tableX[i].Table_name: ", tableX[i].Table_name)
		for k :=0; k < len(fk_HierarchyX); k++{
			// fmt.Println("*********************************************")
			// fmt.Println("inside for k :=0; k < len(fk_HierarchyX); k++")
			// fmt.Println("*********************************************")
			// fmt.Println("fk_HierarchyX[k].TableName: ", fk_HierarchyX[k].TableName, "k: ", k)
			d = len(fk_HierarchyX[k].RelatedTablesLevels) //2
			e = d - c //3
			// fmt.Println("c: ", c, "d: ", d,  "e: ", e)
			if fk_HierarchyX[k].TableName == tableX[i].Table_name {
				// fmt.Println("inside if fk_HierarchyX[k].TableName == tableX[i].Table_name")
				// fmt.Println("k: ", k,"i: ", i,"	fk_HierarchyX[k].TableName: ",fk_HierarchyX[k].TableName, "tableX[i].Table_name: ", tableX[i].Table_name)
				for l :=len(fk_HierarchyX[k].RelatedTablesLevels)-e; l < len(fk_HierarchyX[k].RelatedTablesLevels); l++{					
					// fmt.Println("	**************************************************************************************************************")
					// fmt.Println("	inside for l :=len(fk_HierarchyX[k].RelatedTablesLevels)-e; l < len(fk_HierarchyX[k].RelatedTablesLevels); l++")
					// fmt.Println("	**************************************************************************************************************")
					// fmt.Println("	fk_HierarchyX[k].RelatedTablesLevels[l]): ",fk_HierarchyX[k].RelatedTablesLevels[l])
					for m:=0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++{ //carrier, serviceareatype, site
						// fmt.Println("		***************************************************************************************")
						// fmt.Println("		inside for m:=0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++")
						// fmt.Println("		***************************************************************************************")
						// fmt.Println("		m: ", m, "len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList: ", len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList))
						// fmt.Println("		fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m]: ", fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m])
						if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							// fmt.Println("		I am breaking this because this is a self referencing case - Hierarchy_TableName == FK_Related_TableName")
							break
						}
						// // Check if the related table is the root table (would create a cycle)
						// if fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName == fk_HierarchyX[k].TableName {
						// 	// fmt.Println("		Avoiding cycle back to root table: ", fk_HierarchyX[k].TableName)
						// 	continue
						// }

						for z :=0; z < len(tableX); z++ {
							// fmt.Println("			**************************************")
							// fmt.Println("			inside for z :=0; z < len(tableX); z++")
							// fmt.Println("			**************************************")
							// fmt.Println("			tableX[z].Table_name: ", tableX[z].Table_name)
							// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName: ", fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName)
							if tableX[z].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {	
								// fmt.Println("			inside if tableX[z].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName")

								// // Simple cycle detection: check if this table already exists at any level
								// cycleDetected := false
								
								// // Check if this is the root table
								// if tableX[z].Table_name == fk_HierarchyX[k].TableName {
								// 	fmt.Println("			Cycle back to root table detected: ", tableX[z].Table_name)
								// 	cycleDetected = true
								// }
								
								// // Check if this table has already been processed in the current hierarchy
								// if !cycleDetected {
								// 	for checkLevel := 0; checkLevel < len(fk_HierarchyX[k].RelatedTablesLevels); checkLevel++ {
								// 		if fk_HierarchyX[k].RelatedTablesLevels[checkLevel].Hierarchy_TableName == tableX[z].Table_name {
								// 			fmt.Println("			Table ", tableX[z].Table_name, " already exists at level ", checkLevel, " - avoiding cycle")
								// 			cycleDetected = true
								// 			break
								// 		}
								// 	}
								// }
								
								// if cycleDetected {
								// 	fmt.Println("			Skipping to prevent cycle")
								// 	continue
								// }

								if len(tableX[z].ForeignKeys) > 0 {
									// fmt.Println("			inside if len(tableX[z].ForeignKeys) > 0")
									relatedTables.RelatedTableList = nil
									relatedTables.Hierarchy_TableName = ""
									for y :=0; y < len(tableX[z].ForeignKeys); y++{
										var relatedTable RelatedTable
										// fmt.Println("				*****************************************************")
										// fmt.Println("				inside for y :=0; y < len(tableX[z].ForeignKeys); y++")
										// fmt.Println("				*****************************************************")
										// fmt.Println("				z: ",z,"y: ",y,"	tableX[z].ForeignKeys[y]: ", tableX[z].ForeignKeys[y])
										relatedTable.FK_Related_TableName = tableX[z].ForeignKeys[y].FK_Related_TableName
										relatedTable.FK_Related_SingularTableName = tableX[z].ForeignKeys[y].FK_Related_SingularTableName
										relatedTable.FK_Related_Table_Column = tableX[z].ForeignKeys[y].FK_Related_Table_Column
										relatedTable.FK_Related_TableName_Singular_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Singular_Object
										relatedTable.FK_Related_TableName_Plural_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Plural_Object
										relatedTables.Hierarchy_TableName = tableX[z].Table_name	
										relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
										// fmt.Println("				relatedTables.Hierarchy_TableName: ", relatedTables.Hierarchy_TableName)
										// fmt.Println("				relatedTables.Hierarchy_TableName: ", relatedTables.RelatedTableList[len(relatedTables.RelatedTableList)-1])
										// time.Sleep(2 * time.Second)
									}
									fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
									// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels: ", fk_HierarchyX[k].RelatedTablesLevels)

								// 	if len(tableX[z].CompositeForeignKeys) > 0 {
								// 		for y :=0; y < len(tableX[z].CompositeForeignKeys); y++{
								// 			var relatedTable RelatedTable
								// 			// fmt.Println("				*****************************************************")
								// 			// fmt.Println("				inside for y :=0; y < len(tableX[z].ForeignKeys); y++")
								// 			// fmt.Println("				*****************************************************")
								// 			// fmt.Println("				z: ",z,"y: ",y,"	tableX[z].ForeignKeys[y]: ", tableX[z].ForeignKeys[y])
	
								// 			relatedTable.FK_Related_TableName = tableX[z].CompositeForeignKeys[y].FK_Related_TableName
								// 			relatedTable.FK_Related_SingularTableName = tableX[z].CompositeForeignKeys[y].FK_Related_SingularTableName
								// 			relatedTable.FK_Related_TableName_Singular_Object = tableX[z].CompositeForeignKeys[y].FK_Related_TableName_Singular_Object
								// 			relatedTable.FK_Related_TableName_Plural_Object = tableX[z].CompositeForeignKeys[y].FK_Related_TableName_Plural_Object
								// 			relatedTable.FK_Related_Table_Column = tableX[z].CompositeForeignKeys[y].FK_Related_Table_Column
							
								// 			relatedTables.Hierarchy_TableName = tableX[z].Table_name	
								// 			relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
								// 			// fmt.Println("				relatedTables: ", relatedTables)
								// 			// time.Sleep(2 * time.Second)
								// 		}
								// 		fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
								// 		// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels: ", fk_HierarchyX[k].RelatedTablesLevels)	
								// 	}
								// } else if len(tableX[z].ForeignKeys) == 0 && len(tableX[z].CompositeForeignKeys) > 0 {
								// 	if len(tableX[z].CompositeForeignKeys) > 0 {
								// 		for y :=0; y < len(tableX[z].CompositeForeignKeys); y++{
								// 			var relatedTable RelatedTable
								// 			// fmt.Println("				*****************************************************")
								// 			// fmt.Println("				inside for y :=0; y < len(tableX[z].ForeignKeys); y++")
								// 			// fmt.Println("				*****************************************************")
								// 			// fmt.Println("				z: ",z,"y: ",y,"	tableX[z].ForeignKeys[y]: ", tableX[z].ForeignKeys[y])
		
								// 			relatedTable.FK_Related_TableName = tableX[z].CompositeForeignKeys[y].FK_Related_TableName
								// 			relatedTable.FK_Related_SingularTableName = tableX[z].CompositeForeignKeys[y].FK_Related_SingularTableName
								// 			relatedTable.FK_Related_TableName_Singular_Object = tableX[z].CompositeForeignKeys[y].FK_Related_TableName_Singular_Object
								// 			relatedTable.FK_Related_TableName_Plural_Object = tableX[z].CompositeForeignKeys[y].FK_Related_TableName_Plural_Object
								// 			relatedTable.FK_Related_Table_Column = tableX[z].CompositeForeignKeys[y].FK_Related_Table_Column
	
								// 			relatedTables.Hierarchy_TableName = tableX[z].Table_name	
								// 			relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
								// 			// fmt.Println("				relatedTables: ", relatedTables)
								// 			// time.Sleep(2 * time.Second)
								// 		}
								// 		fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
								// 		// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels: ", fk_HierarchyX[k].RelatedTablesLevels)										
								// 	}
								}
								break
							}
							// time.Sleep(1 * time.Second)

						}
						// time.Sleep(1 * time.Second)

					}
					// time.Sleep(1 * time.Second)

				}
				// fmt.Println("d: ", d)
				c = d
				// time.Sleep(1 * time.Second)

			}
			// time.Sleep(1 * time.Second)

		}
		// time.Sleep(2 * time.Second)
	}




	//////Print DBSchemaReader////
	for i := 0; i < len(tableX); i++ {
		// fmt.Println("Table_name: ",tableX[i].Table_name)
		for i := 0; i < len(tableX[i].Table_Columns); i++ {
			// fmt.Println("tableX[i].Table_Columns[j]: ",tableX[i].Table_Columns[j])
		}
	}

	//////Identifying user management table////
	// isUserTable := false					
	for i := 0; i < len(tableX); i++ {
		hasUserColumns := false
		hasAuthColumns := false
		hasDateColumns := false
		// hasStatusColumns := false
		// hasLoginNameColumns := false
		// hasSecurityColumns := false
		// hasOtherColumns := false
			hasUserName := strings.Contains(strings.ToLower(tableX[i].Table_name), "user") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "users") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "account") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "accounts") ||
						// Authentication related
						strings.Contains(strings.ToLower(tableX[i].Table_name), "auth") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "authentication") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "login") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "signin") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "member") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "members") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "membership") ||
						// Profile related
						strings.Contains(strings.ToLower(tableX[i].Table_name), "profile") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "profiles") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "person") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "persons") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "people") ||
						// Employee/Customer related
						strings.Contains(strings.ToLower(tableX[i].Table_name), "employee") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "employees") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "customer") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "customers") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "client") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "clients") ||
						// Identity related
						strings.Contains(strings.ToLower(tableX[i].Table_name), "identity") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "identities") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "credential") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "credentials")	
			if hasUserName {
				tableX[i].UserTableSpecs.TableName = tableX[i].Table_name
				for x, column := range tableX[i].Table_Columns {
					// fmt.Println("tableX[i].Table_Columns[x]: ",tableX[i].Table_Columns[x])
					// fmt.Println("1column: ", column)
					colName := strings.ToLower(column.Column_name)
					// User Table Authentication columns
					if colName == "password" || colName == "password_hash" || colName == "passwordhash" ||
					colName == "pass_hash" || colName == "passhash" || colName == "pwd" || colName == "pwd_hash" ||
					colName == "hashed_password" || colName == "hashedpassword" || colName == "encrypted_password" ||
					colName == "encryptedpassword" {
						hasAuthColumns = true
						tableX[i].UserTableSpecs.AuthColumnName = tableX[i].Table_Columns[x].Column_name
					}
					// Date columns
					if colName == "registration_date" || colName == "registrationdate" || colName == "signup_date" ||
					colName == "created_date" || colName == "join_date" || colName == "joindate" ||
					colName == "last_login" || colName == "lastlogin" || colName == "last_access" ||
					colName == "password_created_at" || colName == "password_changed_at" || colName == "password_updated_at" ||
					colName == "passwordcreatedat" || colName == "passwordchangedat" || colName == "passwordupdatedat" ||
					colName == "createdat" || colName == "changedat" || colName == "updatedat" ||
					colName == "created_at" || colName == "changed_at" || colName == "updated_at" {
						hasDateColumns = true
						tableX[i].UserTableSpecs.DateColumnName = tableX[i].Table_Columns[x].Column_name
					}
					// Status/Role columns
					if colName == "status" || colName == "user_status" || colName == "account_status" ||
					colName == "role" || colName == "user_role" || colName == "role_id" ||
					colName == "is_active" || colName == "isactive" || colName == "active" ||
					colName == "is_verified" || colName == "isverified" || colName == "verified" ||
					colName == "is_deleted" || colName == "isdeleted" || colName == "deleted" ||
					colName == "is_blocked" || colName == "isblocked" || colName == "blocked" {
						// hasStatusColumns = true
						// tableX[i].UserTableSpecs.StatusColumnName = tableX[i].Table_Columns[x].Column_name
					}
					// Email/Username columns
					if colName == "email" || colName == "email_address" || colName == "emailaddress" ||
					colName == "username" || colName == "user_name" || colName == "login_name" || colName == "loginname" ||
					colName == "login" {
						// hasLoginNameColumns = true
						// tableX[i].UserTableSpecs.LoginNameColumnName = tableX[i].Table_Columns[x].Column_name
					}
					// Token/Security columns
					if colName == "reset_token" || colName == "verification_code" || colName == "verificationcode" ||
					colName == "verification_token" || colName == "verificationtoken" || colName == "api_token" ||
					colName == "refresh_token" || colName == "access_token" || colName == "is_email_verified" ||
					colName == "isemailverified" || colName == "email_verification_code" || colName == "emailverificationcode" ||
					colName == "verification_code_validity" || colName == "verificationcode_alidity" ||
					colName == "code_validity" || colName == "codevalidity" {
						// hasSecurityColumns = true
						// tableX[i].UserTableSpecs.SecurityColumnName = tableX[i].Table_Columns[x].Column_name
					}
					// Other columns
					if colName == "first_name" || colName == "firstname" || colName == "fname" ||
					colName == "last_name" || colName == "lastname" || colName == "lname" ||
					colName == "full_name" || colName == "fullname" || colName == "display_name" ||
					colName == "displayname" || colName == "phone" || colName == "phone_number" ||
					colName == "phonenumber" || colName == "mobile" || colName == "mobile_number" ||
					colName == "mobilenumber" || colName == "contact_number" || colName == "profile_picture" ||
					colName == "avatar" || colName == "photo" || colName == "bio" || colName == "biography" ||
					colName == "about" || colName == "address" || colName == "street_address" || colName == "city" ||
					colName == "state" || colName == "country" || colName == "postal_code" || colName == "zip_code" ||
					colName == "employee_id" || colName == "employeeid" || colName == "emp_id" ||
					colName == "department" || colName == "department_id" || colName == "position" {
						// hasOtherColumns = true
					}
				}
				// if hasAuthColumns && hasDateColumns && hasStatusColumns && hasLoginNameColumns && hasSecurityColumns && hasOtherColumns {
				if hasAuthColumns && hasDateColumns {
					hasUserColumns = true
				}				
			}
			if hasUserName && hasUserColumns {
				tableX[i].IsUserTable = true
				break
			}
	}

	//////Identifying session management table////
	for i := 0; i < len(tableX); i++ {
		hasSessionColumns := false
		hasTokenColumns := false
		hasCreateTimingColumns := false
		hasExpiryTimingColumns := false
		hasSessionStatusColumns := false
		// hasAgentColumns := false
		// hasClientIPColumns :=false
		// hasDeviceColumns := false
		// hasLocationColumns := false
		// hasSecurityColumns := false
		// isSessionsTable := false					
		// 1. Check table name (existing check)
		hasSessionName := strings.Contains(strings.ToLower(tableX[i].Table_name), "session") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "sessions") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "user_session") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "usersession") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "users_sessions") || 
						strings.Contains(strings.ToLower(tableX[i].Table_name), "usersession") ||
						// Authentication session variants
						strings.Contains(strings.ToLower(tableX[i].Table_name), "auth_session") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "auth_sessions") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "authsession") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "authsessions") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "login_session") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "login_sessions") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "loginsession") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "loginsessions") ||
						// Token-based sessions
						strings.Contains(strings.ToLower(tableX[i].Table_name), "session_token") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "sessiontoken") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "access_token") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "accesstoken") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "refresh_token") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "refreshtoken") ||
						// Activity/tracking sessions
						strings.Contains(strings.ToLower(tableX[i].Table_name), "activity_session") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "user_activity") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "login_log") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "loginlog") ||
						// Security/audit sessions
						strings.Contains(strings.ToLower(tableX[i].Table_name), "security_session") ||
						strings.Contains(strings.ToLower(tableX[i].Table_name), "audit_session")
		if hasSessionName {
			for x, column := range tableX[i].Table_Columns {
				// fmt.Println("tableX[i].Table_Columns[x]: ",tableX[i].Table_Columns[x])
				// fmt.Println("1column: ", column)
				colName := strings.ToLower(column.Column_name)
				// Token columns
				if colName == "session_token" || colName == "sessiontoken" || colName == "token" ||
					colName == "access_token" || colName == "accesstoken" ||
					colName == "refresh_token" || colName == "refreshtoken" {
					hasTokenColumns = true
					tableX[i].SessionTableSpecs.TokenColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// SessionId columns
				if colName == "session_id" || colName == "sessionid" || colName == "id" {
					//hasSessionIdColumns = true
					tableX[i].SessionTableSpecs.SessionIdColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Token Create Timing columns
				if  ((colName == "created_at" || colName == "createdat" || strings.Contains(colName, "create")) && column.ColumnType == "timestamptz") || 
				((colName == "created_at" || colName == "createdat" || strings.Contains(colName, "create")) && column.ColumnType == "date") {
					hasCreateTimingColumns = true
					tableX[i].SessionTableSpecs.CreateTimingColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Token Expiry Timing columns
				if ((colName == "expires_at" || colName == "expiresat" || colName == "expiry" ||colName == "expires" || colName == "expiration" ||
				colName == "expiration_time" || strings.Contains(colName, "expir")) && column.ColumnType == "timestamptz") || 
				((colName == "expires_at" || colName == "expiresat" || colName == "expiry" ||colName == "expires" || colName == "expiration" ||
				colName == "expiration_time" || strings.Contains(colName, "expir")) && column.ColumnType == "date") {
					hasExpiryTimingColumns = true
					tableX[i].SessionTableSpecs.ExpiryTimingColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Status/Role columns
				if colName == "is_active" || colName == "isactive" || colName == "active" ||
				colName == "is_blocked" || colName == "isblocked" || colName == "blocked" ||
				colName == "is_valid" || colName == "isvalid" || colName == "valid" ||
				colName == "is_revoked" || colName == "isrevoked" || colName == "revoked" {
					hasSessionStatusColumns = true
					tableX[i].SessionTableSpecs.StatusColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// User Client IP columns
				if colName == "ip_address" || colName == "client" || colName == "client_ip" ||
				colName == "clientip" || colName == "ipaddress" || colName == "ip" {
					// hasClientIPColumns = true
					// tableX[i].SessionTableSpecs.AgentColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// User Agent columns
				if colName == "user_agent" || colName == "useragent" || colName == "agent" {
					// hasAgentColumns = true
					// tableX[i].SessionTableSpecs.AgentColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Device columns
				if colName == "device" || colName == "device_type" || colName == "devicetype" ||
				colName == "browser" || colName == "browser_name" || colName == "browsername" ||
				colName == "platform" || colName == "os" || colName == "operating_system" {
					// hasDeviceColumns = true
					// tableX[i].SessionTableSpecs.DeviceColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Location columns
				if colName == "location" || colName == "country" || colName == "city" ||
				colName == "timezone" || colName == "time_zone" {
					// hasLocationColumns = true
					// tableX[i].SessionTableSpecs.LocationColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}
				// Security columns
				if colName == "fingerprint" || colName == "device_fingerprint" ||
				colName == "csrf_token" || colName == "csrftoken" {
					// hasSecurityColumns = true
					// tableX[i].SessionTableSpecs.SecurityColumnName = tableX[i].Table_Columns[x].ColumnNameParams
				}					
			}
			if hasTokenColumns && hasCreateTimingColumns && hasExpiryTimingColumns && hasSessionStatusColumns {
			// hasAgentColumns && hasClientIPColumns && hasDeviceColumns && hasLocationColumns && hasSecurityColumns {
				hasSessionColumns = true
			}				
		}
		var refTableHasUserName bool
		var hasUserReference bool
				
		for  _, column := range tableX[i].ForeignKeys  {
			// Check if this foreign key references a user table
			refTableName := strings.ToLower(column.FK_Related_TableName)
			refTableHasUserName = strings.Contains(refTableName, "user") || 
						strings.Contains(refTableName, "users") || 
						strings.Contains(refTableName, "account") || 
						strings.Contains(refTableName, "accounts") ||
						// Authentication related
						strings.Contains(refTableName, "auth") ||
						strings.Contains(refTableName, "authentication") ||
						strings.Contains(refTableName, "login") ||
						strings.Contains(refTableName, "signin") ||
						strings.Contains(refTableName, "member") ||
						strings.Contains(refTableName, "members") ||
						strings.Contains(refTableName, "membership") ||
						// Profile related
						strings.Contains(refTableName, "profile") ||
						strings.Contains(refTableName, "profiles") ||
						strings.Contains(refTableName, "person") ||
						strings.Contains(refTableName, "persons") ||
						strings.Contains(refTableName, "people") ||
						// Employee/Customer related
						strings.Contains(refTableName, "employee") ||
						strings.Contains(refTableName, "employees") ||
						strings.Contains(refTableName, "customer") ||
						strings.Contains(refTableName, "customers") ||
						strings.Contains(refTableName, "client") ||
						strings.Contains(refTableName, "clients") ||
						// Identity related
						strings.Contains(refTableName, "identity") ||
						strings.Contains(refTableName, "identities") ||
						strings.Contains(refTableName, "credential") ||
						strings.Contains(refTableName, "credentials")
			if refTableHasUserName {
				for _, element := range tableX[i].Table_Columns {
					if element.Column_name == column.FK_Column {
						tableX[i].SessionTableSpecs.RefUserTableName = refTableName
						tableX[i].SessionTableSpecs.RefUserTableColumn = strings.ToLower(column.FK_Related_Table_Column)
						tableX[i].SessionTableSpecs.FkUserColumn = element.ColumnNameParams
						tableX[i].SessionTableSpecs.FkUserColumn2 = column.FK_Column 
					}
				}
				hasUserReference = true
				break
			}	
		}
		if hasSessionName && hasSessionColumns && hasUserReference {
			// isSessionsTable = true
			tableX[i].IsSessionsTable = true
			break
		}
		if tableX[i].IsSessionsTable {
			fmt.Println("tableX[i].SessionTableSpecs", tableX[i].SessionTableSpecs)		
		}

	}

    // Create a map to easily look up tables by name
    tableMap := make(map[string]FK_Hierarchy)
    for _, hierarchy := range fk_HierarchyX {
        tableMap[hierarchy.TableName] = hierarchy
    }
    
    // For each table, print its full hierarchy
    for i, hierarchy := range fk_HierarchyX {
        fmt.Println("i:", i, "Table_name:", hierarchy.TableName)
        
        // Print the direct relationships
        for j := 0; j < len(hierarchy.RelatedTablesLevels); j++ {
            level := hierarchy.RelatedTablesLevels[j]
            fmt.Println("    j:", j, "    Hierarchy_TableName:", level.Hierarchy_TableName)
            
            // For each related table, print its relationships
            for k := 0; k < len(level.RelatedTableList); k++ {
                relatedTable := level.RelatedTableList[k]
                fmt.Println("        k:", k, "    Direct dependency:", relatedTable.FK_Related_TableName)
                
                // Now traverse up the chain to print the full hierarchy
				// time.Sleep(3 * time.Second)
                // printTableChain(relatedTable.FK_Related_TableName, tableMap, 3)
            }
        }
        fmt.Println() // Add an empty line for better readability
    }

	
	return tableX, fk_HierarchyX
}

