package repository

import (
	"github.com/sar-michal/dictionary-app/pkg/models"
)

type Repository interface {
	// Gets or creates a word in the database if it does not exist
	GetOrCreateWord(polishWord string) (*models.Word, error)
	ListWords() ([]models.Word, error)
	GetWordByPolish(polishWord string) (*models.Word, error)
	GetWordByID(wordID uint) (*models.Word, error)
	UpdateWord(wordID uint, newPolishWord string) (*models.Word, error)
	// Deletes a word and all its translations and example sentences
	DeleteWord(wordID uint) error

	ListTranslations(wordID uint) ([]models.Translation, error)
	GetTranslationByID(translationID uint) (*models.Translation, error)
	CreateTranslation(wordID uint, englishTranslation string) (*models.Translation, error)
	UpdateTranslation(translationID uint, newEnglishTranslation string) (*models.Translation, error)
	// Deletes the translation and its associated example sentences
	DeleteTranslation(translationID uint) error

	ListExampleSentences(translationID uint) ([]models.ExampleSentence, error)
	GetExampleSentenceByID(sentenceID uint) (*models.ExampleSentence, error)
	CreateExampleSentence(translationID uint, sentenceText string) (*models.ExampleSentence, error)
	UpdateExampleSentence(sentenceID uint, newSentenceText string) (*models.ExampleSentence, error)
	DeleteExampleSentence(sentenceID uint) error
}
