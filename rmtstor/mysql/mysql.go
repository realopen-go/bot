package mysql

import (
	"gopkg.in/oleiade/reflections.v1"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type IMysql interface {
	Connect()
	Close()
	CreateBill(bill *Bill)
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

func (db *Mysql) Connect() {
	_db, err := gorm.Open("mysql", "root:1234@/realopen?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.db = _db
}

func (db *Mysql) CreateBill(bill *Bill) {
	db.db.Create(bill)
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
	db.db.AutoMigrate(&Bill{}, &User{})
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

func New() IMysql {
	mysql := &Mysql{}
	mysql.Connect()
	mysql.Migrate()
	return mysql
}
