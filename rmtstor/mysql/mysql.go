package mysql

import (
	"fmt"

	"gopkg.in/oleiade/reflections.v1"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlConfig struct {
	Database string
	Host     string
	Password string
	Username string
}

type IMysql interface {
	Connect(config MysqlConfig)
	Close()
	CreateBill(bill *Bill)
	CreateFile(file File)
	CreateUser(user *User)
	FindOrCreateUser(user *User) *User
	FetchBills(interface{}, interface{}) []*Bill
	FetchUser(interface{}, interface{}) *User
	Migrate()
	UpdateBill(where interface{}, args interface{}, bill *Bill)
}

type Mysql struct {
	db *gorm.DB
}

func (db *Mysql) Close() {
	db.db.Close()
}

func (db *Mysql) Connect(config MysqlConfig) {
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Database)
	_db, err := gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect database")
	}

	db.db = _db
}

func (db *Mysql) CreateBill(bill *Bill) {
	if bill.RequestContent != "" {
		bill.MultiID = bill.CreateMultiID()
	}
	db.db.Create(bill)
}

func (db *Mysql) CreateFile(file File) {
	db.db.Create(&file)
}

func (db *Mysql) CreateUser(user *User) {
	userRow := &User{}
	db.db.AutoMigrate(userRow)
	db.db.Where("username = ?", user.Username).First(userRow)
	if userRow == nil {
		db.db.Create(user)
	}
}

func (db *Mysql) FindOrCreateUser(user *User) *User {
	exisitingUser := &User{}
	db.db.Where(user).First(exisitingUser)

	if exisitingUser.ID == "" {
		db.db.Create(user)
		return user
	} else {
		return exisitingUser
	}
}

func (db *Mysql) FetchBills(query interface{}, args interface{}) []*Bill {
	var records []*Bill
	if query != nil {
		db.db.Where(query, args).Find(&records)
	} else {
		db.db.Find(&records)
	}
	return records
}

func (db *Mysql) FetchUser(query interface{}, args interface{}) *User {
	user := &User{}
	if query != nil {
		db.db.Where(query, args).Find(user)
	} else {
		db.db.Find(user)
	}
	return user
}

func (db *Mysql) Migrate() {
	db.db.AutoMigrate(&Bill{}, &File{}, &User{})
}

func (db *Mysql) UpdateBill(where interface{}, args interface{}, bill *Bill) {
	billInDb := Bill{}
	if where != nil && args != nil {
		db.db.Where(where, args)
	}
	db.db.First(&billInDb)

	fields, _ := reflections.Fields(bill)
	for _, f := range fields {
		currentValue, _ := reflections.GetField(billInDb, f)
		updatedValue, _ := reflections.GetField(bill, f)
		if updatedValue != nil && updatedValue != "" && currentValue != updatedValue {
			reflections.SetField(&billInDb, f, updatedValue)
		}
	}
	db.db.Save(&billInDb)
}

func New(config MysqlConfig) IMysql {
	mysql := &Mysql{}
	mysql.Connect(config)
	mysql.Migrate()
	return mysql
}
