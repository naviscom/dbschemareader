package dbschemareader

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"


)

type Table_Struct struct {
	Table_name          string
	Table_Columns       []table_columns
	IndexDetails        []index_name_details
	ForeignKeys			[]foreign_key_details
	OutputFileName      string
	FunctionSignature   string
	FunctionSignature2  string
	FunctionSignature3  string
}

type table_columns struct {
	Column_name     string
	PrimaryFlag     bool
	UniqueFlag      bool
	ForeignFlag		bool
	ColumnType      string
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
}

type FK_Hierarchy struct{
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

func ReadSchema(filePath string, tableX []Table_Struct)  ([]Table_Struct, []FK_Hierarchy) {
	// fmt.Println(filePath, tableX)
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
		// fmt.Println(len(res1), res1)
		if len(res1) > 1 {
			if res1[0] == "CREATE" && res1[1] == "TABLE" {
				// fmt.Println(`Inside "CREATE" && res1[1] == "TABLE"`)
				// fmt.Println(res1)
				table.Table_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
				table.FunctionSignature2 = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
				if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
					table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-3])+"y"
				}else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
					table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-1])
				}else {
					table.OutputFileName = table.Table_name
				}
				if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-3]+"y")
					} else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-1])
					} else {
					table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
				}
				table.Table_Columns = nil
			}
			if res1[0] == "" && res1[1] == "" && strings.TrimSpace(res1[2][0:1]) == `"` {
				// fmt.Println(`Inside """ && res1[1] == "" && strings.TrimSpace(res1[2][0:1]"`)
				// fmt.Println(res1)
				tabColumns.Column_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
				tabColumns.ColumnType = strings.TrimSpace(res1[3][0:])
				tabColumns.ForeignFlag = false
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
				// fmt.Println(`Inside "CREATE" && res1[1] == "INDEX"`)
				// fmt.Println(res1)				
				for i:=0; i<len(tableX); i++{
					if tableX[i].Table_name == strings.TrimSpace(res1[3][1:len(res1[3])-1]) { 
					var index index_name_details
					index.IndexName = strings.TrimSpace(res1[3][1:len(res1[3])-1]) + strconv.Itoa(rand.Intn(90000))
					for m :=4; m<len(res1); m++ {            
						indexColumnName := res1[m]
						if strings.TrimSpace(indexColumnName[0:1]) == `(` {
						indexColumnName = strings.TrimSpace(indexColumnName[2:len(indexColumnName)-1])
						} else if strings.TrimSpace(indexColumnName[0:1]) == `"` {
						indexColumnName = strings.TrimSpace(indexColumnName[1:len(indexColumnName)-1])
						}
						if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `)` {
						indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-2])
						} else if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `"` {
						indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-1])
						}
						// fmt.Println(indexColumnName)
						index.IndexColumn = append(index.IndexColumn,   indexColumnName)
					}
					tableX[i].IndexDetails = append(tableX[i].IndexDetails, index)
					}
				}
			}
			if res1[0] == "ALTER" && res1[4] == "FOREIGN" {
				// fmt.Println(`Inside "ALTER" && res1[4] == "FOREIGN"`)
				// fmt.Println(res1)
				for i:=0; i<len(tableX); i++{
					var fkDetails foreign_key_details
					if tableX[i].Table_name == strings.TrimSpace(res1[2][1:len(res1[2])-1]) { 
						fkDetails.FK_Column = res1[6]
						fkDetails.FK_Column = strings.TrimSpace(fkDetails.FK_Column[2:len(fkDetails.FK_Column)-2])
						fkDetails.FK_Related_TableName = res1[8]
						fkDetails.FK_Related_TableName = strings.TrimSpace(fkDetails.FK_Related_TableName[1:len(fkDetails.FK_Related_TableName)-1])


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


						fkDetails.FK_Related_Table_Column = res1[9]
						fkDetails.FK_Related_Table_Column = strings.TrimSpace(fkDetails.FK_Related_Table_Column[2:len(fkDetails.FK_Related_Table_Column)-3])
						tableX[i].ForeignKeys = append(tableX[i].ForeignKeys, fkDetails)
						for j:=0; j<len(tableX[i].Table_Columns); j++ {
							if tableX[i].Table_Columns[j].Column_name == fkDetails.FK_Column { tableX[i].Table_Columns[j].ForeignFlag = true}
						}
					}
				}						
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
						// fmt.Println(column_name_slice)
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
	var relatedTable RelatedTable
	var relatedTables RelatedTables
	var fk_Hierarchy FK_Hierarchy
	var fk_HierarchyX []FK_Hierarchy
	for i:=0; i<len(tableX); i++{
		relatedTables.RelatedTableList = nil
		fk_Hierarchy.RelatedTablesLevels = nil
		if len(tableX[i].ForeignKeys) > 0 {
			for j :=0; j < len(tableX[i].ForeignKeys); j++{
				relatedTable.FK_Related_TableName = tableX[i].ForeignKeys[j].FK_Related_TableName
				relatedTable.FK_Related_SingularTableName = tableX[i].ForeignKeys[j].FK_Related_SingularTableName
				relatedTable.FK_Related_Table_Column = tableX[i].ForeignKeys[j].FK_Related_Table_Column
				relatedTable.FK_Related_TableName_Singular_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Singular_Object
				relatedTable.FK_Related_TableName_Plural_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Plural_Object
				// fmt.Println("tableX[i].Table_name: ",tableX[i].Table_name,"relatedTable: ",relatedTable)
				
				relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
				//  fmt.Println("tableX[i].Table_name: ",tableX[i].Table_name,"relatedTables: ",relatedTables)
			}
			fk_Hierarchy.TableName = tableX[i].Table_name
			fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
			fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
		}else{
			fk_Hierarchy.TableName = tableX[i].Table_name
			// fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
			fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
		}
		// fmt.Println("tableX[i].Table_name: ",tableX[i].Table_name,"fk_Hierarchy.RelatedTablesLevels: ",fk_Hierarchy.RelatedTablesLevels, "fk_Hierarchy: ",fk_Hierarchy)
		var c int
		var d int
		var e int
		c = 0
		// fmt.Println("i: ",i,"tableX[i].Table_name: ", tableX[i].Table_name)
		for k :=0; k < len(fk_HierarchyX); k++{
			d = len(fk_HierarchyX[k].RelatedTablesLevels) //3
			e = d - c //3
			if fk_HierarchyX[k].TableName == tableX[i].Table_name {
				// fmt.Println("	k: ",k,"	fk_HierarchyX[k].TableName: ",fk_HierarchyX[k].TableName)
				for l :=len(fk_HierarchyX[k].RelatedTablesLevels)-e; l < len(fk_HierarchyX[k].RelatedTablesLevels); l++{
					if l == 0 {
						fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName = fk_HierarchyX[k].TableName
					}
					// fmt.Println("		l: ",l,"	fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName: ",fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName)
					for m:=0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++{ //carrier, serviceareatype, site
						// fmt.Println("			m: ",m,"	fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName: ",fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName)
						for z :=0; z < len(tableX); z++{
							if tableX[z].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {	
								if len(tableX[z].ForeignKeys) > 0 {
									relatedTables.RelatedTableList = nil
									relatedTables.Hierarchy_TableName = ""
									for y :=0; y < len(tableX[z].ForeignKeys); y++{
										// fmt.Println("				z: ",z,"y: ",y,"	tableX[z].ForeignKeys[y].FK_Related_TableName: ", tableX[z].ForeignKeys[y].FK_Related_TableName)
										relatedTable.FK_Related_TableName = tableX[z].ForeignKeys[y].FK_Related_TableName
										relatedTable.FK_Related_SingularTableName = tableX[z].ForeignKeys[y].FK_Related_SingularTableName
										relatedTable.FK_Related_Table_Column = tableX[z].ForeignKeys[y].FK_Related_Table_Column
										relatedTable.FK_Related_TableName_Singular_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Singular_Object
										relatedTable.FK_Related_TableName_Plural_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Plural_Object	
										
										relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)	
									}
									fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
									fk_HierarchyX[k].RelatedTablesLevels[len(fk_HierarchyX[k].RelatedTablesLevels)-1].Hierarchy_TableName = fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName
								}
							}
						}
					}			
				}
				c = d
			}
		}				
	}
	return tableX, fk_HierarchyX
}