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
	// Find all translations of the word
	var translations []models.Translation
	err := tx.Where("word_id = ?", wordID).Find(&translations).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// Delete all associated example sentences
	for _, t := range translations {
		err = tx.
			Where("translation_id = ?", t.TranslationID).
			Delete(&models.ExampleSentence{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Delete all translations of the word
	err = tx.Where("word_id = ?", wordID).Delete(&models.Translation{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// Delete the word
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
	translation := models.Translation{
		WordID:             wordID,
		EnglishTranslation: englishTranslation,
	}

	if err := r.DB.Create(&translation).Error; err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *GormRepository) UpdateTranslation(translationID uint, newEnglishTranslation string) (*models.Translation, error) {
	translation, err := r.GetTranslationByID(translationID)
	if err != nil {
		return nil, err
	}

	translation.EnglishTranslation = newEnglishTranslation

	if err := r.DB.Save(translation).Error; err != nil {
		return nil, err
	}
	return translation, nil
}

func (r *GormRepository) DeleteTranslation(translationID uint) error {
	if err := r.DB.Delete(&models.Translation{}, translationID).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) ListExampleSentences(translationID uint) ([]models.ExampleSentence, error) {
	var sentences []models.ExampleSentence
	err := r.DB.
		Where("translation_id = ?", translationID).
		Find(&sentences).
		Error
	if err != nil {
		return nil, err
	}
	return sentences, nil
}

func (r *GormRepository) GetExampleSentenceByID(sentenceID uint) (*models.ExampleSentence, error) {
	var sentence models.ExampleSentence

	if err := r.DB.First(&sentence, sentenceID).Error; err != nil {
		return nil, err
	}
	return &sentence, nil
}

func (r *GormRepository) CreateExampleSentence(translationID uint, sentenceText string) (*models.ExampleSentence, error) {
	sentence := models.ExampleSentence{
		TranslationID: translationID,
		SentenceText:  sentenceText,
	}
	if err := r.DB.Create(&sentence).Error; err != nil {
		return nil, err
	}
	return &sentence, nil
}

func (r *GormRepository) UpdateExampleSentence(sentenceID uint, newSentenceText string) (*models.ExampleSentence, error) {
	sentence, err := r.GetExampleSentenceByID(sentenceID)
	if err != nil {
		return nil, err
	}

	sentence.SentenceText = newSentenceText

	if err := r.DB.Save(sentence).Error; err != nil {
		return nil, err
	}
	return sentence, nil
}

func (r *GormRepository) DeleteExampleSentence(sentenceID uint) error {
	if err := r.DB.Delete(&models.ExampleSentence{}, sentenceID).Error; err != nil {
		return err
	}
	return nil
}
