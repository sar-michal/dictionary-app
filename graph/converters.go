package graph

import (
	"strconv"

	"github.com/sar-michal/dictionary-app/graph/model"
	"github.com/sar-michal/dictionary-app/pkg/models"
)

// Convert models Word to a GraphQL Word
func convertWord(w *models.Word) *model.Word {
	return &model.Word{
		WordID:       strconv.FormatUint(uint64(w.WordID), 10),
		PolishWord:   w.PolishWord,
		Translations: convertTranslations(w.Translations),
	}
}

// Convert a slice of models Translation to a GraphQL Translation
func convertTranslations(translations []models.Translation) []*model.Translation {
	gqlTranslations := make([]*model.Translation, len(translations))
	for i, t := range translations {
		gqlTranslations[i] = &model.Translation{
			TranslationID:      strconv.FormatUint(uint64(t.TranslationID), 10),
			EnglishTranslation: t.EnglishTranslation,
			WordID:             strconv.FormatUint(uint64(t.WordID), 10),
			ExampleSentences:   convertExampleSentences(t.ExampleSentences),
		}
	}
	return gqlTranslations
}

// Convert a slice of models ExampleSentence to GraphQL ExampleSentence
func convertExampleSentences(sentences []models.ExampleSentence) []*model.ExampleSentence {
	gqlSentences := make([]*model.ExampleSentence, len(sentences))
	for i, s := range sentences {
		gqlSentences[i] = &model.ExampleSentence{
			SentenceID:    strconv.FormatUint(uint64(s.SentenceID), 10),
			SentenceText:  s.SentenceText,
			TranslationID: strconv.FormatUint(uint64(s.TranslationID), 10),
		}
	}
	return gqlSentences
}
