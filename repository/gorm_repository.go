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
	err := r.DB.
		Where("polish_word = ?", polishWord).
		FirstOrCreate(&word, models.Word{PolishWord: polishWord}).
		Error

	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *GormRepository) ListWords() ([]models.Word, error) {
	var words []models.Word

	if err := r.DB.Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}

func (r *GormRepository) GetWordByPolish(polishWord string) (*models.Word, error) {
	var word models.Word

	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *GormRepository) GetWordByID(wordID uint) (*models.Word, error) {
	var word models.Word

	if err := r.DB.First(&word, wordID).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *GormRepository) UpdateWord(wordID uint, newPolishWord string) (*models.Word, error) {
	word, err := r.GetWordByID(wordID)
	if err != nil {
		return nil, err
	}

	word.PolishWord = newPolishWord

	if err := r.DB.Save(word).Error; err != nil {
		return nil, err
	}
	return word, nil
}

func (r *GormRepository) DeleteWord(wordID uint) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// find all translations of the word
	var Translations []models.Translation
	err := tx.Where("word_id = ?", wordID).Find(&Translations).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// delete all associated example sentences
	for _, t := range Translations {
		err = tx.
			Where("translation_id = ?", t.TranslationID).
			Delete(&models.ExampleSentence{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// delete all translations of the word
	err = tx.Where("word_id = ?", wordID).Delete(&models.Translation{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// delete the word
	err = tx.Delete(&models.Word{}, wordID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *GormRepository) ListTranslations(wordID uint) ([]models.Translation, error) {
	var translations []models.Translation

	if err := r.DB.Where("word_id = ?", wordID).Find(&translations).Error; err != nil {
		return nil, err
	}
	return translations, nil
}

func (r *GormRepository) GetTranslationByID(translationID uint) (*models.Translation, error) {
	var translation models.Translation

	if err := r.DB.First(&translation, translationID).Error; err != nil {
		return nil, err
	}
	return &translation, nil
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
