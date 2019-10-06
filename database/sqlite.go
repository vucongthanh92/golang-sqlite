package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Connect to database sqlite
func Connect() *gorm.DB {
	for i := 1; i <= 3; i++ {
		db, err := gorm.Open("sqlite3", "./database/Subrails.db")
		if err == nil {
			return db
		} else {
			if i == 3 {
				panic(err.Error())
			}
		}
	}
	return nil
}
