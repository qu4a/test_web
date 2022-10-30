//Типы данных вернего уровня

package models

import (
	"errors"
	"time"
)

var ErrNotRecord = errors.New("models: подходящей записи не найдено")

// новая структура, которая соответствует таблице данный MySql
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
