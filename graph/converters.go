package graph

import (
	"strconv"

	"github.com/sar-michal/dictionary-app/graph/model"
	"github.com/sar-michal/dictionary-app/pkg/models"
)

// Convert models Word to a GraphQL Word
func convertWord(word *models.Word) *model.Word {
	return &model.Word{
		WordID:       strconv.FormatUint(uint64(word.WordID), 10),
		PolishWord:   word.PolishWord,
		Translations: convertTranslations(word.Translations),
	}
}

// Convert a single models Translation to a GraphQL Translation
func convertTranslation(translation *models.Translation) *model.Translation {
	return &model.Translation{
		TranslationID:      strconv.FormatUint(uint64(translation.TranslationID), 10),
		EnglishTranslation: translation.EnglishTranslation,
		WordID:             strconv.FormatUint(uint64(translation.WordID), 10),
		ExampleSentences:   convertExampleSentences(translation.ExampleSentences),
	}
}

// Convert a slice of models Translation to a GraphQL Translation
func convertTranslations(translations []models.Translation) []*model.Translation {
	gqlTranslations := make([]*model.Translation, len(translations))
	for i, t := range translations {
		gqlTranslations[i] = convertTranslation(&t)
	}
	return gqlTranslations
}

// Convert a single models ExampleSentence to a GraphQL ExampleSentence
func convertExampleSentence(sentence *models.ExampleSentence) *model.ExampleSentence {
	return &model.ExampleSentence{
		SentenceID:    strconv.FormatUint(uint64(sentence.SentenceID), 10),
		SentenceText:  sentence.SentenceText,
		TranslationID: strconv.FormatUint(uint64(sentence.TranslationID), 10),
	}
}

// Convert a slice of models ExampleSentence to GraphQL ExampleSentence
func convertExampleSentences(sentences []models.ExampleSentence) []*model.ExampleSentence {
	gqlSentences := make([]*model.ExampleSentence, len(sentences))
	for i, s := range sentences {
		gqlSentences[i] = convertExampleSentence(&s)
	}
	return gqlSentences
}
