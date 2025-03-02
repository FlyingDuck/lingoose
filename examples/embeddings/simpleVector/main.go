package main

import (
	"context"
	"fmt"

	openaiembedder "github.com/henomis/lingoose/embedder/openai"
	"github.com/henomis/lingoose/index"
	indexoption "github.com/henomis/lingoose/index/option"
	simplevectorindex "github.com/henomis/lingoose/index/simpleVectorIndex"
	"github.com/henomis/lingoose/llm/openai"
	"github.com/henomis/lingoose/loader"
	"github.com/henomis/lingoose/prompt"
	"github.com/henomis/lingoose/textsplitter"
)

// download https://raw.githubusercontent.com/hwchase17/chat-your-data/master/state_of_the_union.txt

func main() {

	openaiEmbedder := openaiembedder.New(openaiembedder.AdaEmbeddingV2)

	docsVectorIndex := simplevectorindex.New("docs", ".", openaiEmbedder)
	indexIsEmpty, _ := docsVectorIndex.IsEmpty()

	if indexIsEmpty {
		err := ingestData(docsVectorIndex, openaiEmbedder)
		if err != nil {
			panic(err)
		}
	}

	query := "Describe within a paragraph what is the purpose of the NATO Alliance."
	similarities, err := docsVectorIndex.SimilaritySearch(
		context.Background(),
		query,
		indexoption.WithTopK(3),
	)
	if err != nil {
		panic(err)
	}

	for _, similarity := range similarities {
		fmt.Printf("Similarity: %f\n", similarity.Score)
		fmt.Printf("Document: %s\n", similarity.Document.Content)
		fmt.Println("Metadata: ", similarity.Document.Metadata)
		fmt.Println("----------")
	}

	documentContext := ""
	for _, similarity := range similarities {
		documentContext += similarity.Document.Content + "\n\n"
	}

	llmOpenAI := openai.NewCompletion()
	prompt1 := prompt.NewPromptTemplate(
		"Based on the following context answer to the question.\n\nContext:\n{{.context}}\n\nQuestion: {{.query}}").WithInputs(
		map[string]string{
			"query":   query,
			"context": documentContext,
		},
	)

	err = prompt1.Format(nil)
	if err != nil {
		panic(err)
	}

	output, err := llmOpenAI.Completion(context.Background(), prompt1.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}

func ingestData(docsVectorIndex *simplevectorindex.Index, openaiEmbedder index.Embedder) error {

	fmt.Printf("Ingesting data...")

	documents, err := loader.NewDirectoryLoader(".", ".txt").Load(context.Background())
	if err != nil {
		return err
	}

	textSplitter := textsplitter.NewRecursiveCharacterTextSplitter(2000, 100)

	documentChunks := textSplitter.SplitDocuments(documents)

	err = docsVectorIndex.LoadFromDocuments(context.Background(), documentChunks)
	if err != nil {
		return err
	}

	fmt.Printf("Done!\n")

	return nil
}
