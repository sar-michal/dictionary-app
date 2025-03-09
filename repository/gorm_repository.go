package repository

import (
	"github.com/sar-michal/dictionary-app/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) GetOrCreateWord(polishWord string) (*models.Word, error) {
	var word models.Word
	result := r.DB.
		Where("polish_word = ?", polishWord).
		FirstOrCreate(&word, models.Word{PolishWord: polishWord})

	if result.Error != nil {
		return nil, result.Error
	}

	return &word, nil
}

func (r *GormRepository) ListWords() ([]models.Word, error) {
	var words []models.Word
	result := r.DB.Find(&words)

	if result.Error != nil {
		return nil, result.Error
	}

	return words, nil
}

func (r *GormRepository) GetWordByPolish(polishWord string) (*models.Word, error) {
}

func (r *GormRepository) GetWordByID(wordID uint) (*models.Word, error) {
}

func (r *GormRepository) UpdateWord(wordID uint, newPolishWord string) (*models.Word, error) {
}

func (r *GormRepository) DeleteWord(wordID uint) error {
}

func (r *GormRepository) ListTranslations(wordID uint) ([]models.Translation, error) {
}

func (r *GormRepository) GetTranslationByID(translationID uint) (*models.Translation, error) {
}

func (r *GormRepository) CreateTranslation(wordID uint, englishTranslation string) (*models.Translation, error) {
}

func (r *GormRepository) UpdateTranslation(translationID uint, newEnglishTranslation string) (*models.Translation, error) {
}

func (r *GormRepository) DeleteTranslation(translationID uint) error {
}

func (r *GormRepository) ListExampleSentences(translationID uint) ([]models.ExampleSentence, error) {
}

func (r *GormRepository) GetExampleSentenceByID(sentenceID uint) (*models.ExampleSentence, error) {
}

func (r *GormRepository) CreateExampleSentence(translationID uint, sentenceText string) (*models.ExampleSentence, error) {
}

func (r *GormRepository) UpdateExampleSentence(sentenceID uint, newSentenceText string) (*models.ExampleSentence, error) {
}

func (r *GormRepository) DeleteExampleSentence(sentenceID uint) error {
}
