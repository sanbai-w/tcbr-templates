package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Gorm      *gorm.DB
	SQL       *sql.DB
	TableName string
	Connected bool // 添加连接状态标志
}

type Counter struct {
	ID        uint      `gorm:"primaryKey"`
	Count     int       `gorm:"not null;default:1"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;not null"`
}

func (Counter) TableName() string {
	// will be set dynamically via session.Table
	return "counters"
}

func Init(username, password, address, database, tableName string) (*Database, error) {
	// 检查是否有必要的数据库配置
	if username == "" || address == "" || database == "" {
		fmt.Printf("Database configuration incomplete, starting with limited functionality\n")
		return &Database{
			Gorm:      nil,
			SQL:       nil,
			TableName: tableName,
			Connected: false,
		}, nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, database)
	g, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		fmt.Printf("Please check your .env file configuration\n")
		return &Database{
			Gorm:      nil,
			SQL:       nil,
			TableName: tableName,
			Connected: false,
		}, nil
	}

	sqldb, err := g.DB()
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		return &Database{
			Gorm:      nil,
			SQL:       nil,
			TableName: tableName,
			Connected: false,
		}, nil
	}

	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	// migrate with dynamic table name
	if err := g.Table(tableName).AutoMigrate(&Counter{}); err != nil {
		fmt.Printf("Database migration failed: %v\n", err)
		return &Database{
			Gorm:      nil,
			SQL:       nil,
			TableName: tableName,
			Connected: false,
		}, nil
	}

	fmt.Printf("Database connected successfully\n")
	return &Database{
		Gorm:      g,
		SQL:       sqldb,
		TableName: tableName,
		Connected: true,
	}, nil
}

// IsConnected 返回数据库连接状态
func (d *Database) IsConnected() bool {
	return d.Connected && d.Gorm != nil
}

func (d *Database) InsertOne() error {
	if !d.IsConnected() {
		return fmt.Errorf("database not connected")
	}
	return d.Gorm.Table(d.TableName).Create(&Counter{}).Error
}

func (d *Database) Truncate() error {
	if !d.IsConnected() {
		return fmt.Errorf("database not connected")
	}
	return d.Gorm.Exec(fmt.Sprintf("TRUNCATE TABLE %s", d.TableName)).Error
}

func (d *Database) CountAll() (int64, error) {
	if !d.IsConnected() {
		return 0, fmt.Errorf("database not connected")
	}

	var count int64
	if err := d.Gorm.Table(d.TableName).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
