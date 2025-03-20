package repository

import (
	"github.com/sar-michal/dictionary-app/pkg/models"
)

type Repository interface {
	// GetOrCreateWord gets or creates a word in the database if it does not exist.
	GetOrCreateWord(polishWord string) (*models.Word, error)
	ListWords() ([]models.Word, error)
	GetWordByPolish(polishWord string) (*models.Word, error)
	GetWordByID(wordID uint) (*models.Word, error)
	UpdateWord(wordID uint, newPolishWord string) (*models.Word, error)
	// DeleteWord deletes a word and all its translations and example sentences.
	DeleteWord(wordID uint) error

	// GetOrCreateTranslation gets or creates a translation in the database if it does not exist.
	GetOrCreateTranslation(wordID uint, englishTranslation string) (*models.Translation, error)
	ListTranslations(wordID uint) ([]models.Translation, error)
	GetTranslationByID(translationID uint) (*models.Translation, error)
	UpdateTranslation(translationID uint, newEnglishTranslation string) (*models.Translation, error)
	// Deletes the translation and its associated example sentences
	DeleteTranslation(translationID uint) error

	// GetOrCreateExampleSentence gets or creates an example sentence in the database if it does not exist.
	GetOrCreateExampleSentence(translationID uint, sentenceText string) (*models.ExampleSentence, error)
	ListExampleSentences(translationID uint) ([]models.ExampleSentence, error)
	GetExampleSentenceByID(sentenceID uint) (*models.ExampleSentence, error)
	UpdateExampleSentence(sentenceID uint, newSentenceText string) (*models.ExampleSentence, error)
	DeleteExampleSentence(sentenceID uint) error

	// Transaction executes the provided function within a database transaction.
	Transaction(fn func(repo Repository) error) error
}
