package repository_test

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/sar-michal/dictionary-app/pkg/config"
	"github.com/sar-michal/dictionary-app/pkg/models"
	"github.com/sar-michal/dictionary-app/pkg/repository"
	"github.com/sar-michal/dictionary-app/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var repo repository.Repository

func TestMain(m *testing.M) {

	os.Setenv("GO_ENV", "test")
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := models.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
	// Inject GormRepository
	repo = &repository.GormRepository{DB: db}
	code := m.Run()
	if err := storage.CloseDB(db); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	os.Exit(code)
}

// Helper function. It creates a repository instance that uses a transaction.
// It rolls the transaction back once the function is done.
func withTransaction(t *testing.T, fn func(txRepo repository.Repository)) {
	err := repo.Transaction(func(txRepo repository.Repository) error {
		fn(txRepo)
		// Return a sentinel error to force rollback.
		return fmt.Errorf("Rollback for test")
	})
	require.Error(t, err, "Rollback should always happen")
}

// Helper function. Truncates all tables and restarts sequences associated with table columns.
func CleanupRepository(t *testing.T) {
	// Attempt a type assertion to access the underlying *gorm.DB.
	gormRepo, ok := repo.(*repository.GormRepository)
	require.True(t, ok, "Expected repository to be of type *GormRepository. Failed to cleanup database")

	err := gormRepo.DB.Exec("TRUNCATE TABLE words, translations, example_sentences RESTART IDENTITY CASCADE").Error
	require.NoError(t, err, "Failed to cleanup database")
}

func TestGetOrCreateWord(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("kot")
		require.NoError(t, err, "GetOrCreateWord should not error")
		assert.Equal(t, "kot", word.PolishWord, "PolishWord should match")
		firstID := word.WordID

		sameWord, err := txRepo.GetOrCreateWord("kot")
		require.NoError(t, err, "Second GetOrCreateWord should not error")
		assert.Equal(t, firstID, sameWord.WordID, "WordID should be consistent")
	})
}

func TestListWords(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		_, err := txRepo.GetOrCreateWord("kot")
		require.NoError(t, err, "Failed to create word kot")

		_, err = txRepo.GetOrCreateWord("pies")
		require.NoError(t, err, "Failed to create word pies")

		words, err := txRepo.ListWords()
		require.NoError(t, err, "ListWords should not error")

		// Verify that the list contains both "kot" and "pies".
		var foundKot, foundPies bool
		for _, w := range words {
			if w.PolishWord == "kot" {
				foundKot = true
			} else if w.PolishWord == "pies" {
				foundPies = true
			}
		}
		assert.True(t, foundKot, "Word 'kot' should be listed")
		assert.True(t, foundPies, "Word 'pies' should be listed")
	})
}
func TestGetWordByPolish(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		created, err := txRepo.GetOrCreateWord("lis")
		require.NoError(t, err, "Failed to create word 'lis'")

		retrieved, err := txRepo.GetWordByPolish("lis")
		require.NoError(t, err, "GetWordByPolish should not error")
		assert.Equal(t, created.WordID, retrieved.WordID, "Retrieved word should match the created word")
	})
}
func TestGetWordByID(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		created, err := txRepo.GetOrCreateWord("koń")
		require.NoError(t, err, "Failed to create word 'koń'")

		retrieved, err := txRepo.GetWordByID(created.WordID)
		require.NoError(t, err, "GetWordByID should not error")
		assert.Equal(t, created.PolishWord, retrieved.PolishWord, "Retrieved word should match the created word")
	})
}
func TestUpdateWord(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		created, err := txRepo.GetOrCreateWord("koza")
		require.NoError(t, err, "Failed to create word 'koza'")

		updated, err := txRepo.UpdateWord(created.WordID, "owca")
		require.NoError(t, err, "UpdateWord should not error")
		assert.Equal(t, "owca", updated.PolishWord, "PolishWord should be updated to 'owca'")

		retrieved, err := txRepo.GetWordByID(created.WordID)
		require.NoError(t, err, "GetWordByID should not error")
		assert.Equal(t, "owca", retrieved.PolishWord, "Retrieved word should reflect the update")
	})
}
func TestDeleteWord(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("słoń")
		require.NoError(t, err, "Failed to create word 'słoń'")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "elephant")
		require.NoError(t, err, "Failed to create translation for 'słoń'")

		example, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "An elephant is a large animal")
		require.NoError(t, err, "Failed to create example sentence")

		err = txRepo.DeleteWord(word.WordID)
		require.NoError(t, err, "DeleteWord should not error")

		_, err = txRepo.GetWordByID(word.WordID)
		assert.Error(t, err, "Expected error retrieving deleted word")

		_, err = txRepo.GetTranslationByID(translation.TranslationID)
		assert.Error(t, err, "Expected error retrieving translation after deletion")

		_, err = txRepo.GetExampleSentenceByID(example.SentenceID)
		assert.Error(t, err, "Expected error retrieving example sentence after deletion")
	})
}

func TestListTranslations(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("kot")
		require.NoError(t, err, "GetOrCreateWord should not error")

		// Create two translations for "kot": "cat" and "kitty".
		_, err = txRepo.GetOrCreateTranslation(word.WordID, "cat")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'cat'")

		_, err = txRepo.GetOrCreateTranslation(word.WordID, "kitty")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'kitty'")

		translations, err := txRepo.ListTranslations(word.WordID)
		require.NoError(t, err, "ListTranslations should not error")
		assert.Equal(t, 2, len(translations), "Expected two translations for 'kot'")
	})
}
func TestGetTranslationByID(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("pies")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "dog")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'dog'")

		retrieved, err := txRepo.GetTranslationByID(translation.TranslationID)
		require.NoError(t, err, "GetTranslationByID should not error")
		assert.Equal(t, "dog", retrieved.EnglishTranslation, "Expected English translation to be 'dog'")
	})
}
func TestGetOrCreateTranslation(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("lis")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "fox")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'fox'")
		assert.Equal(t, "fox", translation.EnglishTranslation, "Expected English translation to be 'fox'")
	})
}
func TestUpdateTranslation(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("koza")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "goat")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'goat'")

		updated, err := txRepo.UpdateTranslation(translation.TranslationID, "she-goat")
		require.NoError(t, err, "UpdateTranslation should not error")
		assert.Equal(t, "she-goat", updated.EnglishTranslation, "Expected updated English translation to be 'she-goat'")
	})
}
func TestDeleteTranslation(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("słoń")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "elephant")
		require.NoError(t, err, "GetOrCreateTranslation should not error for 'elephant'")

		example, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "An elephant is a large animal.")
		require.NoError(t, err, "GetOrCreateExampleSentence should not error")

		err = txRepo.DeleteTranslation(translation.TranslationID)
		require.NoError(t, err, "DeleteTranslation should not error")

		_, err = txRepo.GetTranslationByID(translation.TranslationID)
		assert.Error(t, err, "Expected error retrieving deleted translation")

		_, err = txRepo.GetExampleSentenceByID(example.SentenceID)
		assert.Error(t, err, "Expected error retrieving deleted example sentence")
	})
}
func TestListExampleSentences(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("kot")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "cat")
		require.NoError(t, err, "GetOrCreateTranslation should not error")

		// Create two example sentences.
		_, err = txRepo.GetOrCreateExampleSentence(translation.TranslationID, "John has a cat.")
		require.NoError(t, err, "GetOrCreateExampleSentence should not error for first sentence")

		_, err = txRepo.GetOrCreateExampleSentence(translation.TranslationID, "A cat is climbing a tree.")
		require.NoError(t, err, "GetOrCreateExampleSentence should not error for second sentence")

		sentences, err := txRepo.ListExampleSentences(translation.TranslationID)
		require.NoError(t, err, "ListExampleSentences should not error")
		assert.Equal(t, 2, len(sentences), "Expected two example sentences")
	})
}
func TestGetExampleSentenceByID(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("pies")
		require.NoError(t, err, "GetOrCreateWord should not error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "dog")
		require.NoError(t, err, "GetOrCreateTranslation should not error")

		sentence, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "A dog is running around.")
		require.NoError(t, err, "GetOrCreateExampleSentence should not error")

		retrieved, err := txRepo.GetExampleSentenceByID(sentence.SentenceID)
		require.NoError(t, err, "GetExampleSentenceByID should not error")
		assert.Equal(t, "A dog is running around.", retrieved.SentenceText, "Expected sentence text to match")
	})
}
func TestGetOrCreateExampleSentence(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("lis")
		require.NoError(t, err, "GetOrCreateWord Should Not Error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "fox")
		require.NoError(t, err, "GetOrCreateTranslation Should Not Error")

		sentence, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "The quick brown fox jumps over the lazy dog.")
		require.NoError(t, err, "GetOrCreateExampleSentence Should Not Error")
		assert.Equal(t, "The quick brown fox jumps over the lazy dog.", sentence.SentenceText,
			"Expected Sentence Text To Be 'The quick brown fox jumps over the lazy dog'")
	})
}
func TestUpdateExampleSentence(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("koza")
		require.NoError(t, err, "GetOrCreateWord Should Not Error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "goat")
		require.NoError(t, err, "GetOrCreateTranslation Should Not Error")

		sentence, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "The goat is eating grass.")
		require.NoError(t, err, "GetOrCreateExampleSentence Should Not Error")

		updated, err := txRepo.UpdateExampleSentence(sentence.SentenceID, "The goat is drinking water.")
		require.NoError(t, err, "UpdateExampleSentence Should Not Error")
		assert.Equal(t, "The goat is drinking water.", updated.SentenceText,
			"Expected Updated Sentence Text To Be 'The goat is drinking water.'")
	})
}
func TestDeleteExampleSentence(t *testing.T) {
	withTransaction(t, func(txRepo repository.Repository) {
		word, err := txRepo.GetOrCreateWord("owca")
		require.NoError(t, err, "GetOrCreateWord Should Not Error")

		translation, err := txRepo.GetOrCreateTranslation(word.WordID, "sheep")
		require.NoError(t, err, "GetOrCreateTranslation Should Not Error")

		sentence, err := txRepo.GetOrCreateExampleSentence(translation.TranslationID, "The fluffy sheep is sleeping.")
		require.NoError(t, err, "GetOrCreateExampleSentence Should Not Error")

		err = txRepo.DeleteExampleSentence(sentence.SentenceID)
		require.NoError(t, err, "DeleteExampleSentence Should Not Error")

		_, err = txRepo.GetExampleSentenceByID(sentence.SentenceID)
		assert.Error(t, err, "Expected Error Retrieving Deleted Example Sentence")
	})
}
func TestConcurrentGetOrCreateWords(t *testing.T) {
	// Clean up the database
	CleanupRepository(t)
	defer CleanupRepository(t)

	wordsToCreate := []string{
		"kot", "kot", "kot",
		"pies", "pies", "pies",
		"lis", "lis", "lis",
	}
	toCreateLen := len(wordsToCreate)
	expectedUnique := toCreateLen - 6

	var wg sync.WaitGroup
	wg.Add(len(wordsToCreate))
	errCh := make(chan error, toCreateLen)

	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := 0

	// Launch concurrent goroutines to create words.
	for _, word := range wordsToCreate {
		word := word
		go func() {
			mu.Lock()
			ready++
			if ready < toCreateLen {
				cond.Wait()
			} else {
				cond.Broadcast()
			}
			mu.Unlock()

			_, err := repo.GetOrCreateWord(word)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "GetOrCreateWord should not error in concurrent execution")
	}

	createdWords, err := repo.ListWords()
	require.NoError(t, err, "ListWords should not error")
	assert.Equal(t, expectedUnique, len(createdWords), "Expected number of words to match")
}
func TestConcurrentGetOrCreateTranslations(t *testing.T) {
	// Clean up the database
	CleanupRepository(t)
	defer CleanupRepository(t)

	word, err := repo.GetOrCreateWord("kot")
	require.NoError(t, err, "GetOrCreateWord should not error")

	translationsToCreate := []string{
		"cat", "cat", "cat",
		"kitty", "kitty", "kitty",
		"pussy", "pussy", "pussy",
	}
	toCreateLen := len(translationsToCreate)
	expectedUnique := toCreateLen - 6

	var wg sync.WaitGroup
	wg.Add(len(translationsToCreate))
	errCh := make(chan error, toCreateLen)

	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := 0

	// Launch concurrent goroutines to create translations.
	for _, trans := range translationsToCreate {
		trans := trans
		go func() {
			mu.Lock()
			ready++
			if ready < toCreateLen {
				cond.Wait()
			} else {
				cond.Broadcast()
			}
			mu.Unlock()

			_, err := repo.GetOrCreateTranslation(word.WordID, trans)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "GetOrCreateTranslation should not error in concurrent execution")
	}

	createdTranslations, err := repo.ListTranslations(word.WordID)
	require.NoError(t, err, "ListTranslations should not error")
	assert.Equal(t, expectedUnique, len(createdTranslations), "Expected number of translations to match")
}
func TestConcurrentGetOrCreateExampleSentences(t *testing.T) {
	// Clean up the database
	CleanupRepository(t)
	defer CleanupRepository(t)

	word, err := repo.GetOrCreateWord("pies")
	require.NoError(t, err, "GetOrCreateWord should not error")

	translation, err := repo.GetOrCreateTranslation(word.WordID, "dog")
	require.NoError(t, err, "GetOrCreateTranslation should not error")

	sentencesToCreate := []string{
		"The dog barks.", "The dog barks.", "The dog barks.",
		"The dog runs.", "The dog runs.", "The dog runs.",
		"The dog eats.", "The dog eats.", "The dog eats.",
	}
	toCreateLen := len(sentencesToCreate)
	expectedUnique := toCreateLen - 6
	var wg sync.WaitGroup
	wg.Add(toCreateLen)
	errCh := make(chan error, toCreateLen)
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := 0
	// Launch concurrent goroutines to create example sentences.
	for _, sentence := range sentencesToCreate {
		sentence := sentence
		go func() {
			mu.Lock()
			ready++
			if ready < toCreateLen {
				cond.Wait()
			} else {
				cond.Broadcast()
			}
			mu.Unlock()
			_, err := repo.GetOrCreateExampleSentence(translation.TranslationID, sentence)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "GetOrCreateExampleSentence should not error in concurrent execution")
	}

	examples, err := repo.ListExampleSentences(translation.TranslationID)
	require.NoError(t, err, "ListExampleSentences should not error")
	assert.Equal(t, expectedUnique, len(examples), "Expected number of example sentences to match")
}
