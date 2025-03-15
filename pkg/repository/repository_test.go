package repository_test

import (
	"log"
	"os"
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
	// Attempt a type assertion to access the underlying *gorm.DB.
	gormRepo, ok := repo.(*repository.GormRepository)
	require.True(t, ok, "Expected repository to be of type *GormRepository. Failed to begin transaction")

	tx := gormRepo.DB.Begin()
	require.NoError(t, tx.Error, "Failed to begin transaction")
	defer func() {
		require.NoError(t, tx.Rollback().Error, "Failed to rollback transaction")
	}()

	tempRepo := &repository.GormRepository{DB: tx}
	fn(tempRepo)
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
			}
			if w.PolishWord == "pies" {
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

		translation, err := txRepo.CreateTranslation(word.WordID, "elephant")
		require.NoError(t, err, "Failed to create translation for 'słoń'")

		example, err := txRepo.CreateExampleSentence(translation.TranslationID, "An elephant is a large animal")
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
