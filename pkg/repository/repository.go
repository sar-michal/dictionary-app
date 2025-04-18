package repository

import (
	"github.com/sar-michal/dictionary-app/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// =============================
// GormRepository implementation
// =============================

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) GetOrCreateWord(polishWord string) (*models.Word, error) {
	word := models.Word{
		PolishWord: polishWord,
	}
	// Attempt to insert. On conflict, do nothing.
	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "polish_word"}},
		DoNothing: true,
	}).Create(&word).Error
	if err != nil {
		return nil, err
	}
	// Retrieves the word from database.
	err = r.DB.Where("polish_word = ?", polishWord).First(&word).Error
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// ListWords returns a slice of all words. Preloads translations and example sentences.
func (r *GormRepository) ListWords() ([]models.Word, error) {
	var words []models.Word
	err := r.DB.Preload("Translations.ExampleSentences").Find(&words).Error
	if err != nil {
		return nil, err
	}
	return words, nil
}

// GetWordByPolish finds a word. Preloads translations and example sentences.
func (r *GormRepository) GetWordByPolish(polishWord string) (*models.Word, error) {
	var word models.Word

	err := r.DB.
		Preload("Translations.ExampleSentences").
		Where("polish_word = ?", polishWord).
		First(&word).
		Error
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// GetWordByID finds a word. Preloads translations and example sentences.
func (r *GormRepository) GetWordByID(wordID uint) (*models.Word, error) {
	var word models.Word

	err := r.DB.
		Preload("Translations.ExampleSentences").
		First(&word, wordID).
		Error
	if err != nil {
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
	err := r.DB.Transaction(func(tx *gorm.DB) error {

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
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ListTranslations returns a slice of translations of a word.
// Preloads example sentences.
func (r *GormRepository) ListTranslations(wordID uint) ([]models.Translation, error) {
	var translations []models.Translation
	err := r.DB.
		Where("word_id = ?", wordID).
		Preload("ExampleSentences").
		Find(&translations).
		Error
	if err != nil {
		return nil, err
	}
	return translations, nil
}

// GetTranslationByID returns a translation. Preloads example sentences.
func (r *GormRepository) GetTranslationByID(translationID uint) (*models.Translation, error) {
	var translation models.Translation
	err := r.DB.
		Preload("ExampleSentences").
		First(&translation, translationID).
		Error
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *GormRepository) GetOrCreateTranslation(wordID uint, englishTranslation string) (*models.Translation, error) {
	translation := models.Translation{
		WordID:             wordID,
		EnglishTranslation: englishTranslation,
	}
	// Attempt to insert. On conflict, do nothing.
	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "word_id"}, {Name: "english_translation"}},
		DoNothing: true,
	}).Create(&translation).Error
	if err != nil {
		return nil, err
	}
	// Retrieve the translation from database.
	err = r.DB.
		Where("word_id = ? AND english_translation = ?", wordID, englishTranslation).
		First(&translation).
		Error
	if err != nil {
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
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Delete all associated example sentences
		err := tx.
			Where("translation_id = ?", translationID).
			Delete(&models.ExampleSentence{}).
			Error
		if err != nil {
			return err
		}
		// Delete the translation
		if err := tx.Delete(&models.Translation{}, translationID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
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

func (r *GormRepository) GetOrCreateExampleSentence(translationID uint, sentenceText string) (*models.ExampleSentence, error) {
	sentence := models.ExampleSentence{
		TranslationID: translationID,
		SentenceText:  sentenceText,
	}
	// Attempt to insert. On conflict, do nothing.
	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "translation_id"}, {Name: "sentence_text"}},
		DoNothing: true,
	}).Create(&sentence).Error
	if err != nil {
		return nil, err
	}
	// Retrieves the sentence from database.
	err = r.DB.
		Where("translation_id = ? AND sentence_text = ?", translationID, sentenceText).
		First(&sentence).
		Error
	if err != nil {
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

// Transaction executes the provided function within a database transaction.
func (r *GormRepository) Transaction(fn func(repo Repository) error) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := &GormRepository{DB: tx}
		return fn(txRepo)
	})
}
