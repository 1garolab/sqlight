package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

func readBinary[T any](bytesToRead []byte) (T, error) {
	var bytesRead T
	if err := binary.Read(bytes.NewReader(bytesToRead), binary.BigEndian, &bytesRead); err != nil {
		return bytesRead, err
	}
	return bytesRead, nil
}

// TODO: make this a method over dbFile type 
// because this is all being perfomed in one single file
func readBytesTo[T any](dbFile *os.File, bytesToRead []byte, slice []byte) (T, error) {
	var bytesRead T

	_, err := dbFile.Read(bytesToRead)
	if err != nil {
		return bytesRead, err
	}

	bytesRead, err = readBinary[T](slice)
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, nil
}

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		dbHeader := make([]byte, 100) 
		dbPageSize, err := readBytesTo[uint16](databaseFile, dbHeader, dbHeader[16:18])
		if err != nil {
			log.Fatal(err)
		}

		pageHeader := make([]byte, 8) // 8 for leaf pages and 12 bytes for interior pages
		numberOfCells, err := readBytesTo[uint16](databaseFile, pageHeader, pageHeader[3:5])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("database page size: ", dbPageSize)
		fmt.Println("number of tables: ", numberOfCells)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
