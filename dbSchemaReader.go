package dbschemareader
// package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	// "time"
)

type Table_Struct struct {
	Table_name          			string
	Table_Columns       			[]table_columns
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
}

type CompositeForeignKeysAndReferences struct {
	ConstraintName       			string
	CompositeForeignKeys			string
	CompositeForeignKeysReferences 	string
	CompositeForeignKeysOnClause	string
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
}

type table_columns struct {
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
	FK_Related_Table_Column					string
	FK_Related_SingularTableName			string
	FK_Related_TableName_Plural_Object		string
	FK_Related_TableName_Singular_Object	string
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
	fmt.Println("filePath: ",filePath, "tableX: ",tableX)
	// var tableX []Table_Struct
	var table Table_Struct
	var tabColumns table_columns
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
				// if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
				// 	table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-3]+"y")
				// } else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
				// 	table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-1])
				// } else {
				// 	table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
				// }
				table.Table_Columns = nil
			}
			if res1[0] == "" && res1[1] == "" && strings.TrimSpace(res1[2][0:1]) == `"` {
				// fmt.Println(`Inside """ && res1[1] == "" && strings.TrimSpace(res1[2][0:1]"`)
				// fmt.Println(res1)
				tabColumns.Column_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
				tabColumns.Column_name = strings.ReplaceAll(tabColumns.Column_name, "__", "_")
				tabColumns.ColumnType = strings.TrimSpace(res1[3][0:])
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
				if len(res1) > 4 {
					if res1[4] == `PRIMARY` {
					tabColumns.PrimaryFlag = true
					} else{
					tabColumns.PrimaryFlag = false
					}
					if res1[4] == `UNIQUE` {
					tabColumns.UniqueFlag = true
					} else{
					tabColumns.UniqueFlag = false            
					}
				} else {
					tabColumns.PrimaryFlag = false
					tabColumns.UniqueFlag = false
				}
				table.Table_Columns = append(table.Table_Columns, tabColumns)
			}
			if res1[0] == "CREATE" && res1[1] == "INDEX" {
				var onIndex int //usingIndex, indexIndex
				for x, element := range res1{
					if element == "ON"{
						onIndex = x
					}
					// if element == "INDEX"{
					// 	indexIndex = x
					// }
					// if element == "USING"{
					// 	usingIndex = x
					// }
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

						// for m := 4; m < len(res1); m++ {            
						// 	indexColumnName := res1[m]
						// 	if strings.TrimSpace(indexColumnName[0:1]) == `(` {
						// 		indexColumnName = strings.TrimSpace(indexColumnName[2:len(indexColumnName)-1])
						// 	} else if strings.TrimSpace(indexColumnName[0:1]) == `"` {
						// 		indexColumnName = strings.TrimSpace(indexColumnName[1:len(indexColumnName)-1])
						// 	}
						// 	if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `)` {
						// 		indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-2])
						// 	} else if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `"` {
						// 		indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-1])
						// 	}
						// 	// fmt.Println(indexColumnName)
						// 	index.IndexColumn = append(index.IndexColumn,   indexColumnName)
						// }
						tableX[i].IndexDetails = append(tableX[i].IndexDetails, index)
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
					// if element == "DEFERRABLE" {
					// 	//deferrableWordIndex = i
					// }
					// if element == "PRIMARY" {
					// 	primaryWordIndex = i
					// }
					// if element == "ONLY" {
					// 	onlyWordIndex = i
					// }
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
					checkConstraint.CheckConstraintValue = extractColumnNamesWithParenthesis(res1, checkWordIndex +1)
					checkConstraint.CheckConstraintValue = columnNameCleanUp(checkConstraint.CheckConstraintValue)
					tableX[indexOfTableX].CheckConstraints = append(tableX[indexOfTableX].CheckConstraints, checkConstraint)
					// fmt.Println("last appended CheckConstraints: ", tableX[indexOfTableX].CheckConstraints[len(tableX[indexOfTableX].CheckConstraints)-1])
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
						onClause := strings.Join(onClauseSlice, " ")
						onClause = strings.ReplaceAll(onClause, ";", "")
						// fmt.Println("onClause: ", onClause)
					}

					if strings.Contains(foreignKeyColumn, ",") {
						compositeForeignKeysAndReferences.CompositeForeignKeys = foreignKeyColumn
						compositeForeignKeysAndReferences.CompositeForeignKeysReferences = relatedTableColumn
						compositeForeignKeysAndReferences.CompositeForeignKeysOnClause = onClause
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
					var tabColumns table_columns
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
			// fmt.Println("    FK_Column: ", tableX[i].ForeignKeys[j].FK_Column, "FK_Related_TableName: ", tableX[i].ForeignKeys[j].FK_Related_TableName, "FK_Related_SingularTableName: ", tableX[i].ForeignKeys[j].FK_Related_SingularTableName, "FK_Related_Table_Column: ", tableX[i].ForeignKeys[j].FK_Related_Table_Column)
		}
	}
	//////////////////////////////////////////////
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
			for j :=0; j < len(tableX[i].ForeignKeys); j++{
				var relatedTable RelatedTable
				// fmt.Println("at the start of loop---->tableX[i].Table_name: ", tableX[i].Table_name)
				// fmt.Println("at the start of loop---->tableX[i].ForeignKeys: ", tableX[i].ForeignKeys[j])
				// fmt.Println("at the start of loop---->relatedTables: ",relatedTables)
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
		}else{
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
						// fmt.Println("		fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList: ", fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList)
						if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							// fmt.Println("		I am breaking because this is a self referencing case - Hierarchy_TableName == FK_Related_TableName")
							break
						}
						for z :=0; z < len(tableX); z++{
							// fmt.Println("			**************************************")
							// fmt.Println("			inside for z :=0; z < len(tableX); z++")
							// fmt.Println("			**************************************")
							// fmt.Println("			tableX[z].Table_name: ", tableX[z].Table_name)
							// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName: ", fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName)
							if tableX[z].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {	
								// fmt.Println("			inside if tableX[z].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName")
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
										// fmt.Println("				relatedTables: ", relatedTables)
										// time.Sleep(2 * time.Second)
									}
									fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
									// fmt.Println("			fk_HierarchyX[k].RelatedTablesLevels: ", fk_HierarchyX[k].RelatedTablesLevels)
								}
								break
							}
						}
					}
				}
				// fmt.Println("d: ", d)
				c = d
			}
		}
	}
	//////Print DBSchemaReader////
	for i := 0; i < len(tableX); i++ {
		// fmt.Println("Table_name: ",tableX[i].Table_name)
		for j := 0; j < len(tableX[i].Table_Columns); j++ {
			// fmt.Println("	Column_name: ",tableX[i].Table_Columns[j].Column_name)
		}
		for j := 0; j < len(tableX[i].ForeignKeys); j++ {
			// fmt.Println("	ForeignKeys: ",tableX[i].ForeignKeys[j].FK_Column, "FK_Related_TableName: ",tableX[i].ForeignKeys[j].FK_Related_TableName,  "FK_Related_Table_Column: ",tableX[i].ForeignKeys[j].FK_Related_Table_Column)
		}
	} 

    // Create a map to easily look up tables by name
    tableMap := make(map[string]FK_Hierarchy)
    for _, hierarchy := range fk_HierarchyX {
        tableMap[hierarchy.TableName] = hierarchy
    }
    
    // // For each table, print its full hierarchy
    // for i, hierarchy := range fk_HierarchyX {
    //     fmt.Println("i:", i, "Table_name:", hierarchy.TableName)
        
    //     // Print the direct relationships
    //     for j := 0; j < len(hierarchy.RelatedTablesLevels); j++ {
    //         level := hierarchy.RelatedTablesLevels[j]
    //         fmt.Println("    j:", j, "    Hierarchy_TableName:", level.Hierarchy_TableName)
            
    //         // For each related table, print its relationships
    //         for k := 0; k < len(level.RelatedTableList); k++ {
    //             relatedTable := level.RelatedTableList[k]
    //             fmt.Println("        k:", k, "    Direct dependency:", relatedTable.FK_Related_TableName)
                
    //             // Now traverse up the chain to print the full hierarchy
	// 			// time.Sleep(3 * time.Second)
    //             // printTableChain(relatedTable.FK_Related_TableName, tableMap, 3)
    //         }
    //     }
    //     fmt.Println() // Add an empty line for better readability
    // }


	return tableX, fk_HierarchyX
}

// func main() {
// 	// go run . ~/Documents/workspaces/dbschemas/old_catalyst_schema.sql
// 	filePath := os.Args[1]
// 	var tableX []Table_Struct
// 	var fk_HierarchyX []FK_Hierarchy
// 	// _, _ = ReadSchema(filePath, tableX)
// 	tableX, fk_HierarchyX = ReadSchema(filePath, tableX)
// 	fmt.Println("From dbschemareader---->len(tableX): ", len(tableX), len(fk_HierarchyX))
// 	fmt.Println("From dbschemareader---->len(fk_HierarchyX): ", len(fk_HierarchyX))
// }





//////Following code shall not be deleted. Here exists some parts which can be used later on

// func main() {
// 	// go run . ~/Documents/workspaces/dbschemas/old_catalyst_schema.sql
// 	filePath := os.Args[1]
// 	var tableX []Table_Struct
// 	var fk_HierarchyX []FK_Hierarchy
// 	// _, _ = ReadSchema(filePath, tableX)
// 	tableX, fk_HierarchyX = ReadSchema(filePath, tableX)
// 	fmt.Println("len(tableX): ", len(tableX), len(fk_HierarchyX))
// 	fmt.Println("len(fk_HierarchyX): ", len(fk_HierarchyX))
	// //////Print DBSchemaReader////
	// for i := 0; i < len(tableX); i++ {
	// 	fmt.Println("Table_name: ",tableX[i].Table_name)
	// 	for j := 0; j < len(tableX[i].Table_Columns); j++ {
	// 		fmt.Println("	Column_name: ",tableX[i].Table_Columns[j].Column_name)
	// 	}
	// 	for j := 0; j < len(tableX[i].ForeignKeys); j++ {
	// 		fmt.Println("	ForeignKeys: ",tableX[i].ForeignKeys[j].FK_Column, "FK_Related_TableName: ",tableX[i].ForeignKeys[j].FK_Related_TableName,  "FK_Related_Table_Column: ",tableX[i].ForeignKeys[j].FK_Related_Table_Column)
	// 	}
	// } 
	// for i := 0; i < len(fk_HierarchyX); i++ {
	// 	fmt.Println("i: ",i," Table_name: ",fk_HierarchyX[i].TableName)
	// 	for j := 0; j < len(fk_HierarchyX[i].RelatedTablesLevels); j++ {
	// 		fmt.Println("	j: ",j,"	fk_HierarchyX[i].RelatedTablesLevels[j].Hierarchy_TableName: ", fk_HierarchyX[i].RelatedTablesLevels[j].Hierarchy_TableName)
	// 		for k := 0; k < len(fk_HierarchyX[i].RelatedTablesLevels[j].RelatedTableList); k++ {
	// 			fmt.Println("		k: ",k,"	fk_HierarchyX[i].RelatedTablesLevels[j].RelatedTableList[k].FK_Related_TableName: ",fk_HierarchyX[i].RelatedTablesLevels[j].RelatedTableList[k].FK_Related_TableName)
	// 		}
	// 	}
	// }
	///////////////////////////////////////////////
    // Create a map to easily look up tables by name
    // tableMap := make(map[string]FK_Hierarchy)
    // for _, hierarchy := range fk_HierarchyX {
    //     tableMap[hierarchy.TableName] = hierarchy
    // }
    
    // // For each table, print its full hierarchy
    // for i, hierarchy := range fk_HierarchyX {
    //     fmt.Println("i:", i, "Table_name:", hierarchy.TableName)
        
    //     // Print the direct relationships
    //     for j := 0; j < len(hierarchy.RelatedTablesLevels); j++ {
    //         level := hierarchy.RelatedTablesLevels[j]
    //         fmt.Println("    j:", j, "    Hierarchy_TableName:", level.Hierarchy_TableName)
            
    //         // For each related table, print its relationships
    //         for k := 0; k < len(level.RelatedTableList); k++ {
    //             relatedTable := level.RelatedTableList[k]
    //             fmt.Println("        k:", k, "    Direct dependency:", relatedTable.FK_Related_TableName)
                
    //             // Now traverse up the chain to print the full hierarchy
	// 			// time.Sleep(3 * time.Second)
    //             // printTableChain(relatedTable.FK_Related_TableName, tableMap, 3)
    //         }
    //     }
    //     fmt.Println() // Add an empty line for better readability
    // }
	//////////////////////////////////////////////
// }

// func printTableChain(tableName string, tableMap map[string]FK_Hierarchy, indentLevel int) {
//     // Check if the table exists in our map
//     hierarchy, exists := tableMap[tableName]
//     if !exists {
//         return
//     }
    
//     // Get the indent string based on the level
//     indent := strings.Repeat("    ", indentLevel)
    
//     // For each level in this table's hierarchy
//     for j := 0; j < len(hierarchy.RelatedTablesLevels); j++ {
//         level := hierarchy.RelatedTablesLevels[j]
        
//         // For each related table at this level
//         for k := 0; k < len(level.RelatedTableList); k++ {
//             relatedTable := level.RelatedTableList[k]
//             fmt.Println(indent + "Indirect dependency:", relatedTable.FK_Related_TableName)
            
//             // Recursively print this table's dependencies
//             printTableChain(relatedTable.FK_Related_TableName, tableMap, indentLevel+1)
//         }
//     }
// }
