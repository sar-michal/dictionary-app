package models

import (
	"gorm.io/gorm"
)

type Word struct {
	WordID       uint          `gorm:"primaryKey"`
	PolishWord   string        `gorm:"uniqueIndex;not null"`
	Translations []Translation `gorm:"foreignKey:WordID"`
}

type Translation struct {
	TranslationID      uint              `gorm:"primaryKey"`
	WordID             uint              `gorm:"not null"`
	EnglishTranslation string            `gorm:"not null"`
	ExampleSentences   []ExampleSentence `gorm:"foreignKey:TranslationID"`
}
type ExampleSentence struct {
	SentenceID    uint   `gorm:"primaryKey"`
	TranslationID uint   `gorm:"not null"`
	SentenceText  string `gorm:"not null"`
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Word{}, &Translation{}, &ExampleSentence{})
	if err != nil {
		return err
	}
	return nil
}
