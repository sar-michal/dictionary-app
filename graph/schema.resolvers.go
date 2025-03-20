package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sar-michal/dictionary-app/graph/model"
	"github.com/sar-michal/dictionary-app/pkg/models"
	"github.com/sar-michal/dictionary-app/pkg/repository"
)

// CreateWord is the resolver for the createWord field.
func (r *mutationResolver) CreateWord(ctx context.Context, polishWord string) (*model.Word, error) {
	validWord, err := validateInput(polishWord)
	if err != nil {
		return nil, fmt.Errorf("failed to validate polish word: %w", err)
	}

	word, err := r.Repo.GetOrCreateWord(validWord)
	if err != nil {
		return nil, fmt.Errorf("failed to create word: %w", err)
	}
	return convertWord(word), nil
}

// UpdateWord is the resolver for the updateWord field.
func (r *mutationResolver) UpdateWord(ctx context.Context, wordID string, newPolishWord string) (*model.Word, error) {
	id, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid wordID: %w", err)
	}

	validWord, err := validateInput(newPolishWord)
	if err != nil {
		return nil, fmt.Errorf("failed to validate new polish word: %w", err)
	}

	word, err := r.Repo.UpdateWord(uint(id), validWord)
	if err != nil {
		return nil, fmt.Errorf("failed to update word: %w", err)
	}
	return convertWord(word), nil
}

// DeleteWord is the resolver for the deleteWord field.
func (r *mutationResolver) DeleteWord(ctx context.Context, wordID string) (bool, error) {
	id, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invaild wordID: %w", err)
	}

	if err := r.Repo.DeleteWord(uint(id)); err != nil {
		return false, fmt.Errorf("failed to delete word: %w", err)
	}
	return true, nil
}

// CreateTranslationWithWord is the resolver for the CreateTranslationWithWord field.
func (r *mutationResolver) CreateTranslationWithWord(ctx context.Context, polishWord string, englishTranslation string, exampleSentences []string) (*model.Translation, error) {
	validWord, err := validateInput(polishWord)
	if err != nil {
		return nil, fmt.Errorf("failed to validate polish word: %w", err)
	}

	validTranslation, err := validateInput(englishTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to validate english translation: %w", err)
	}

	var validSentences []string
	for _, sentence := range exampleSentences {
		validSentence, err := validateInput(sentence)
		if err != nil {
			return nil, fmt.Errorf("failed to validate example sentence: %w", err)
		}
		validSentences = append(validSentences, validSentence)
	}

	var resultTranslation *models.Translation
	err = r.Repo.Transaction(func(txRepo repository.Repository) error {
		word, err := txRepo.GetOrCreateWord(validWord)
		if err != nil {
			return fmt.Errorf("failed to get or create word: %w", err)
		}

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, validTranslation)
		if err != nil {
			return fmt.Errorf("failed to create translation: %w", err)
		}

		for _, validSentence := range validSentences {
			_, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, validSentence)
			if err != nil {
				return fmt.Errorf("failed to create example sentence: %w", err)
			}
		}

		translation, err = txRepo.GetTranslationByID(translation.TranslationID)
		if err != nil {
			return fmt.Errorf("failed to retrieve translation: %w", err)
		}
		resultTranslation = translation
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return convertTranslation(resultTranslation), nil
}

// CreateTranslation is the resolver for the CreateTranslation field.
func (r *mutationResolver) CreateTranslation(ctx context.Context, wordID string, englishTranslation string, exampleSentences []string) (*model.Translation, error) {
	id, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid wordID: %w", err)
	}

	validTranslation, err := validateInput(englishTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to validate english translation: %w", err)
	}

	var validSentences []string
	for _, sentence := range exampleSentences {
		validSentence, err := validateInput(sentence)
		if err != nil {
			return nil, fmt.Errorf("failed to validate example sentence: %w", err)
		}
		validSentences = append(validSentences, validSentence)
	}

	var resultTranslation *models.Translation
	err = r.Repo.Transaction(func(txRepo repository.Repository) error {
		translation, err := txRepo.GetOrCreateTranslation(uint(id), validTranslation)
		if err != nil {
			return fmt.Errorf("failed to create translation: %w", err)
		}

		for _, validSentence := range validSentences {
			_, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, validSentence)
			if err != nil {
				return fmt.Errorf("failed to create example sentence: %w", err)
			}
		}

		translation, err = txRepo.GetTranslationByID(translation.TranslationID)
		if err != nil {
			return fmt.Errorf("failed to retrieve translation: %w", err)
		}
		resultTranslation = translation
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return convertTranslation(resultTranslation), nil
}

// UpdateTranslation is the resolver for the updateTranslation field.
func (r *mutationResolver) UpdateTranslation(ctx context.Context, translationID string, newEnglishTranslation string) (*model.Translation, error) {
	id, err := strconv.ParseUint(translationID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid translationID: %w", err)
	}

	validTranslation, err := validateInput(newEnglishTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to validate english translation: %w", err)
	}

	translation, err := r.Repo.UpdateTranslation(uint(id), validTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to update translation: %w", err)
	}
	return convertTranslation(translation), nil
}

// DeleteTranslation is the resolver for the deleteTranslation field.
func (r *mutationResolver) DeleteTranslation(ctx context.Context, translationID string) (bool, error) {
	id, err := strconv.ParseUint(translationID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid translationID: %w", err)
	}

	if err := r.Repo.DeleteTranslation(uint(id)); err != nil {
		return false, fmt.Errorf("failed to delete translation: %w", err)
	}
	return true, nil
}

// CreateExampleSentence is the resolver for the CreateExampleSentence field.
func (r *mutationResolver) CreateExampleSentence(ctx context.Context, translationID string, sentenceText string) (*model.ExampleSentence, error) {
	id, err := strconv.ParseUint(translationID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid translationID: %w", err)
	}

	validSentence, err := validateInput(sentenceText)
	if err != nil {
		return nil, fmt.Errorf("failed to validate example sentence: %w", err)
	}

	sentence, err := r.Repo.GetOrCreateExampleSentence(uint(id), validSentence)
	if err != nil {
		return nil, fmt.Errorf("failed to create example sentence: %w", err)
	}
	return convertExampleSentence(sentence), nil
}

// UpdateExampleSentence is the resolver for the updateExampleSentence field.
func (r *mutationResolver) UpdateExampleSentence(ctx context.Context, sentenceID string, newSentenceText string) (*model.ExampleSentence, error) {
	id, err := strconv.ParseUint(sentenceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sentenceID: %w", err)
	}

	validSentence, err := validateInput(newSentenceText)
	if err != nil {
		return nil, fmt.Errorf("failed to validate example sentence: %w", err)
	}

	updatedSentence, err := r.Repo.UpdateExampleSentence(uint(id), validSentence)
	if err != nil {
		return nil, fmt.Errorf("failed to update example sentence: %w", err)
	}

	return convertExampleSentence(updatedSentence), nil
}

// DeleteExampleSentence is the resolver for the deleteExampleSentence field.
func (r *mutationResolver) DeleteExampleSentence(ctx context.Context, sentenceID string) (bool, error) {
	id, err := strconv.ParseUint(sentenceID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid sentenceID: %w", err)
	}

	if err := r.Repo.DeleteExampleSentence(uint(id)); err != nil {
		return false, fmt.Errorf("failed to delete example sentence: %w", err)
	}

	return true, nil
}

// Words is the resolver for the words field.
func (r *queryResolver) Words(ctx context.Context) ([]*model.Word, error) {
	words, err := r.Repo.ListWords()
	if err != nil {
		return nil, fmt.Errorf("failed to list words: %w", err)
	}

	return convertWords(words), nil
}

// WordByPolish is the resolver for the wordByPolish field.
func (r *queryResolver) WordByPolish(ctx context.Context, polishWord string) (*model.Word, error) {
	validWord, err := validateInput(polishWord)
	if err != nil {
		return nil, fmt.Errorf("failed to validate polish word: %w", err)
	}

	word, err := r.Repo.GetWordByPolish(validWord)
	if err != nil {
		return nil, fmt.Errorf("failed to get word by polish: %w", err)
	}
	return convertWord(word), nil
}

// WordByID is the resolver for the wordByID field.
func (r *queryResolver) WordByID(ctx context.Context, wordID string) (*model.Word, error) {
	id, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid wordID: %w", err)
	}

	word, err := r.Repo.GetWordByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get word by id: %w", err)
	}
	return convertWord(word), nil
}

// Translations is the resolver for the translations field.
func (r *queryResolver) Translations(ctx context.Context, wordID string) ([]*model.Translation, error) {
	id, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid wordID: %w", err)
	}

	translations, err := r.Repo.ListTranslations(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to list translations: %w", err)
	}

	return convertTranslations(translations), nil
}

// TranslationByID is the resolver for the translationByID field.
func (r *queryResolver) TranslationByID(ctx context.Context, translationID string) (*model.Translation, error) {
	id, err := strconv.ParseUint(translationID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid translationID: %w", err)
	}

	translation, err := r.Repo.GetTranslationByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get translation by ID: %w", err)
	}

	return convertTranslation(translation), nil
}

// ExampleSentences is the resolver for the exampleSentences field.
func (r *queryResolver) ExampleSentences(ctx context.Context, translationID string) ([]*model.ExampleSentence, error) {
	id, err := strconv.ParseUint(translationID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid translationID: %w", err)
	}

	sentences, err := r.Repo.ListExampleSentences(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to list example sentences: %w", err)
	}

	return convertExampleSentences(sentences), nil
}

// ExampleSentenceByID is the resolver for the exampleSentenceByID field.
func (r *queryResolver) ExampleSentenceByID(ctx context.Context, sentenceID string) (*model.ExampleSentence, error) {
	id, err := strconv.ParseUint(sentenceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sentenceID: %w", err)
	}

	sentence, err := r.Repo.GetExampleSentenceByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get example sentence by ID: %w", err)
	}

	return convertExampleSentence(sentence), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
