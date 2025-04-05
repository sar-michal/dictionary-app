# Dictionary App

## Table of Contents
- [Description](#description)
- [Entity Relationship Diagram](#entity-relationship-diagram)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Running Tests](#running-tests)
- [GraphQL API](#graphql-api)
  - [Word operations](#word-operations)
  - [Translation operations](#translation-operations)
  - [Example sentence operations](#example-sentence-operations)

## Description

The Dictionary App is an application designed to collect translations of Polish words into English in a relational, normalized database. It provides a public GraphQL API that allows end users to manage translations - create new ones, retrieve, modify, and remove existing ones. Users can send a Polish word and its English translation along with example sentences demonstrating the word's usage. The system supports multiple translations for a single word (e.g., "pisać" can be translated as "write" and "type"). The application is designed to handle concurrent processing of translations.

## Entity Relationship Diagram
![Dictionary-App-ERD](https://github.com/user-attachments/assets/956a5a5d-ecd3-4d06-b0ed-e8dd0e306e2f)

## Dependencies

- Go 1.24.0
- PostgreSQL
- Docker
- gqlgen (GraphQL implementation for go)
- godotenv (Environment variable loader)
- testify (Testing framework)
- GORM (ORM library for Go)

## Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/sar-michal/dictionary-app.git
    cd dictionary-app
    ```
2. **Install dependencies:**
    ```sh
    go mod tidy
    ```
3. **Set up environment variables:**  
    Create a .env file in the project root. For this purpose, the .env.example file was provided. An example .env file configuration:
    ```properties
    DB_HOST=localhost
    DB_USER=user
    DB_PASSWORD=pass
    DB_NAME=db
    DB_PORT=5432
    DB_SSLMODE=disable
    ```
4. **Run PostgreSQL container using Docker:**
    ```sh
    docker-compose up -d
    ```
5. **Run the application**
   ```sh
   go run cmd/main.go
   ```
## Running Tests

1. **Set up test environment variables:**  
    Create a .env file in the project root. For this purpose, the .env.test.example file was provided. An example .env.test file configuration:
    ```properties
    DB_HOST=localhost
    DB_USER=testuser
    DB_PASSWORD=testpass
    DB_NAME=testdb
    DB_PORT=5431
    DB_SSLMODE=disable
    ```
2. **Run PostgreSQL test container using Docker:**
    ```sh
    docker-compose --file compose.test.yml --env-file .env.test up -d
    ```
3. **Adjust the hardcoded Config in /pkg/repository/repository_test.go (if .env.test differs from the example one):**
    ```go
    // Hardcoded config to prevent accidents
	config := &config.Config{
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		Port:     "5431",
		SSLMode:  "disable",
	}
    ```
4. **Run the tests:**
    ```sh
    go test ./...
    ```
5. **Tear down the test environment:**
    ```sh
    docker-compose --file compose.test.yml down
    ```
## GraphQL API
Example Queries and Mutations:  

### Word operations

#### CreateNewWord
```graphql
mutation CreateNewWord {
    createWord(polishWord: "kot") {
        wordID
        polishWord
        translations {
            translationID
            englishTranslation
        }
    }
}
```

#### UpdateWord
```graphql
mutation UpdateWord {
    updateWord(wordID: "1", newPolishWord: "pies") {
        wordID
        polishWord
    }
}
```

#### DeleteWord
```graphql
mutation DeleteWord {
    deleteWord(wordID: "1")
}
```

#### GetAllWordsWithDetails
```graphql
query GetAllWordsWithDetails {
    words {
        wordID
        polishWord
        translations {
            translationID
            englishTranslation
            exampleSentences {
                sentenceID
                sentenceText
            }
        }
    }
}
```

#### GetAllWords
```graphql
query GetAllWords {
    words {
        polishWord
        wordID
    }
}
```

#### GetWordByPolish
```graphql
query GetWordByPolish {
    wordByPolish(polishWord: "kot") {
        wordID
        polishWord
        translations {
            englishTranslation
            exampleSentences {
                sentenceText
            }
        }
    }
}
```

#### GetWordByID
```graphql
query GetWordByID {
    wordByID(wordID: "1") {
        wordID
        polishWord
    }
}
```

### Translation operations

#### CreateTranslationWithWord
```graphql
mutation CreateTranslationWithWord {
    createTranslationWithWord(
        polishWord: "wilk", 
        englishTranslation: "wolf", 
        exampleSentences: ["The wolf is scary.", "The wolf howls."]
    ) {
        translationID
        englishTranslation
        wordID
        exampleSentences {
            sentenceID
            sentenceText
        }
    }
}
```

#### AddTranslationToWord
```graphql
mutation AddTranslationToWord {
    createTranslation(
        wordID: "1", 
        englishTranslation: "fox", 
        exampleSentences: ["The fox is swift.", "The fox is cunning."]
    ) {
        translationID
        englishTranslation
        wordID
        exampleSentences {
            sentenceID
            sentenceText
        }
    }
}
```

#### UpdateTranslation
```graphql
mutation UpdateTranslation {
    updateTranslation(translationID: "1", newEnglishTranslation: "feline") {
        translationID
        englishTranslation
        wordID
    }
}
```

#### DeleteTranslation
```graphql
mutation DeleteTranslation {
    deleteTranslation(translationID: "1")
}
```

#### GetTranslationsForWord
```graphql
query GetTranslationsForWord {
    translations(wordID: "1") {
        translationID
        englishTranslation
        wordID
        exampleSentences {
            sentenceID
            sentenceText
        }
    }
}
```

#### GetTranslationByID
```graphql
query GetTranslationByID {
    translationByID(translationID: "1") {
        translationID
        englishTranslation
        wordID
    }
}
```

### Example sentence operations

#### AddExampleSentence
```graphql
mutation AddExampleSentence {
    createExampleSentence(
        translationID: "1", 
        sentenceText: "The cat lounges in the sun."
    ) {
        sentenceID
        sentenceText
        translationID
    }
}
```

#### UpdateExampleSentenceText
```graphql
mutation UpdateExampleSentenceText {
    updateExampleSentence(
        sentenceID: "1", 
        newSentenceText: "The cat naps in the afternoon sun."
    ) {
        sentenceID
        sentenceText
        translationID
    }
}
```

#### DeleteExampleSentence
```graphql
mutation DeleteExampleSentence {
    deleteExampleSentence(sentenceID: "1")
}
```

#### GetExampleSentences
```graphql
query GetExampleSentences {
    exampleSentences(translationID: "1") {
        sentenceID
        sentenceText
        translationID
    }
}
```

#### GetExampleSentenceByID
```graphql
query GetExampleSentenceByID {
    exampleSentenceByID(sentenceID: "1") {
        sentenceID
        sentenceText
        translationID
    }
}
```
