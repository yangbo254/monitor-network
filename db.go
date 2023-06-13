package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type VistorRecord struct {
	gorm.Model
	Timestamp int64  `gorm:"index"`
	Src       string `gorm:"index"`
	ResultCN  int64
	ResultHK  int64
	ResultUS  int64
}

var gDb *gorm.DB

func DbInit() {
	dsn := "root:rJ0mpLY1uG79pm60nHMT@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gDb = db
	// 迁移 schema
	gDb.AutoMigrate(&VistorRecord{})
}

func DbGetVistorRecords(beginTime, endTime int64) []VistorRecord {
	if gDb == nil {
		return nil
	}
	var VistorRecords []VistorRecord
	_ = gDb.Where("Timestamp >= ? AND timestamp <= ?", beginTime, endTime).Find(&VistorRecords)
	return VistorRecords
}

func DbSetVistorRecords(timestamp int64, src string, resultCN, resultHK, resultUS int64) {
	if gDb == nil {
		return
	}
	data := VistorRecord{Timestamp: timestamp, Src: src, ResultCN: resultCN, ResultHK: resultHK, ResultUS: resultUS}
	gDb.Create(&data)
}
