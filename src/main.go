package main

import (
	//"encoding/csv"
	//"flag"
	//"io/ioutil"
	//"github.com/BurntSushi/toml"
	//"unicode/utf8"
    //"database/sql"
    //_ "github.com/lib/pq"
	//"bytes"
	"encoding/binary"
	"flag"
	"io/ioutil"
		
	//"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"strings"
	//"github.com/LindsayBradford/go-dbf/godbf"
)

/*
type Config struct {
    LinhaInicioTabela       int    `toml:"linha_inicio_tabela"`
    HostPostgreSQL          string `toml:"host_postgresql"`
    PortPostgreSQL          string `toml:"port_postgresql"`
    DatabasePostgreSQL      string `toml:"database_postgresql"`
	UserPostgreSQL          string `toml:"user_postgresql"`
	PasswordUserPostgreSQL  string `toml:"password_user_postgresql"`
    ArquivoExcel            string `toml:"arquivo_excel"`
    PlanilhaExcel           string `toml:"planilha_excel"`
}

var conf Config


func init(){

/*	
	sql = strings.TrimSuffix(sql, ", ")
    fmt.Println(sql);

    link := getLink()

    _, err = link.Query(sql)


}

/*
func getLink() (*sql.DB) {
    conStr := "host="      + conf.HostPostgreSQL
    conStr += " port="     + conf.PortPostgreSQL
    conStr += " user="     + conf.UserPostgreSQL
    conStr += " password=" + conf.PasswordUserPostgreSQL
    conStr += " dbname="   + conf.DatabasePostgreSQL
    conStr += " sslmode=disable"

    //fmt.Println(conStr)
    // db, err := sql.Open("postgres", "host=localhost port=5432 user=aprenda password=golang dbname=blog sslmode=disable")
    
    db, err := sql.Open("postgres", conStr)

    if err != nil {
        panic(err)
    }

    //err = db.Ping()

    return db
}
/**/


func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
	   if e != nil { return e }
	   if strings.EqualFold(filepath.Ext(d.Name()),  ext) {
		  a = append(a, s)
	   }
	   return nil
	})
	return a
}

type databaseHeader struct {
	recordLength      uint16
	headerBlockSize   uint16
	fileType          uint8
	dataBlockSizeCode byte // 1 K, 2 K, 3K or 4K//
	recordCount       uint32
	blocksUsedCount   uint16
	blocksTotalCount  uint16
	lastBlockInUse    uint16
	fieldCount        uint8
	keyFieldsCount    uint8
}

// blockHeader contains the block record information
type blockHeader struct {
	nextBlockNumber  uint16
	prevBlockNumber  uint16
	offsetLastRecord uint16
}

type fieldDescription struct {
	ordinal   int
	fieldType uint8
	length    uint8
	name      string
}

var fields map[byte]fieldDescription

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readByteLittleEnd(fileHandle *os.File) (uint8, error) {
	var result byte
	input8 := make([]byte, 1)

	_, err := fileHandle.Read(input8)
	check(err)

	result = input8[0]
	return result, err

}

func readLongLittleEnd(fileHandle *os.File) (uint32, error) {
	var result uint32
	input32 := make([]byte, 4)

	_, err := fileHandle.Read(input32)
	check(err)

	result = binary.LittleEndian.Uint32(input32)
	return result, err
}

func readShortLittleEnd(fileHandle *os.File) (uint16, error) {
	var result uint16
	input16 := make([]byte, 2)

	_, err := fileHandle.Read(input16)
	check(err)

	result = binary.LittleEndian.Uint16(input16)
	return result, err
}

func setupDatabaseHeader(inFile *os.File) (databaseHeader, error) {
	var err error
	var header databaseHeader
	header.recordLength, err = readShortLittleEnd(inFile)
	check(err)
	header.headerBlockSize, err = readShortLittleEnd(inFile)
	check(err)
	header.fileType, err = readByteLittleEnd(inFile)
	check(err)
	header.dataBlockSizeCode, err = readByteLittleEnd(inFile)
	check(err)
	header.recordCount, err = readLongLittleEnd(inFile)
	check(err)
	header.blocksUsedCount, err = readShortLittleEnd(inFile)
	check(err)
	header.blocksTotalCount, err = readShortLittleEnd(inFile)
	check(err)
	_, err = readShortLittleEnd(inFile) // Throw away the first block code
	check(err)
	header.lastBlockInUse, err = readShortLittleEnd(inFile)
	check(err)
	// Go to the field count
	_, err = inFile.Seek(0x0021, 0)
	check(err)
	header.fieldCount, err = readByteLittleEnd(inFile)
	check(err)
	header.keyFieldsCount, err = readByteLittleEnd(inFile)
	check(err)
	return header, err
}

func sendFieldNamesToFile(fields map[byte]fieldDescription, outputFileName string) error {
	var outString string
	numFields := len(fields)

	for i := 0; i < numFields; i++ {
		field := fields[byte(i)]
		outString = outString + field.name + ","
	}

	outString = outString[0 : len(outString)-1]
	outString = outString + "\n"
	err := ioutil.WriteFile(outputFileName, []byte(outString), 0644)
	check(err)

	return err
}

func pullFieldDescs(inFile *os.File, header databaseHeader) error {
	// Go to 0x78 to start file lengths

	_, err := inFile.Seek(120, 0)
	check(err)

	//fieldDescs := make([]fieldDescription, header.fieldCount)
	//fields := make(map[byte]fieldDescription)
	var fieldCounter byte
	fieldCounter = 0
	maxCount := header.fieldCount

	// Fetch the type and length

	var currentField fieldDescription

	for fieldCounter < maxCount {
		currentField = fields[fieldCounter]
		currentField.fieldType, err = readByteLittleEnd(inFile)
		//fieldDescs[fieldCounter].fieldType, err = readByteLittleEnd(inFile)
		check(err)

		currentField.length, err = readByteLittleEnd(inFile)
		check(err)
		fields[fieldCounter] = currentField
		fieldCounter++
	}

	// fetch the names
	var offset int64
	offset = int64(203) + int64(header.fieldCount*6)
	_, err = inFile.Seek(offset, 0)
	check(err)

	fieldCounter = 0
	var valueRead byte
	var fieldNameBytes []byte
	for fieldCounter < maxCount {
		currentField = fields[fieldCounter]
		for {
			valueRead, err = readByteLittleEnd(inFile)
			check(err)

			if valueRead == 0x00 {
				break
			} else {
				fieldNameBytes = append(fieldNameBytes, valueRead)
			}
		}
		currentField.name = string(fieldNameBytes)
		fieldNameBytes = fieldNameBytes[:0]
		fields[fieldCounter] = currentField
		fieldCounter++
	}

	return err
}

func main(){
	dbFiles := find("paradox", ".db")


	dbFile:=dbFiles[10]
	println (dbFile)

	inputFileNamePtf := flag.String("i", "paradox/Cotas.DB", "The input file name for processing")
	outputFileNamePtf := flag.String("o", "./output3.csv", "The output file name to push the data to")

	//var currentOffset int64
	flag.Parse()

	inFile, err := os.Open(*inputFileNamePtf)
	check(err)
	defer inFile.Close()
	
	dbDatabaseHead, err := setupDatabaseHeader(inFile) // Go get the database header
	check(err)

	fields = make(map[byte]fieldDescription) // Pull the Field Descriptions
	err = pullFieldDescs(inFile, dbDatabaseHead)
	check(err)

	err = sendFieldNamesToFile(fields, *outputFileNamePtf)
	check(err)

/*
	numFields := len(fields)

	for i := 0; i < numFields; i++ {
		field := fields[byte(i)]
		outString := "campo " + field.name + " tipo " + string(field.fieldType)
		println(outString)
	}
/**/
}