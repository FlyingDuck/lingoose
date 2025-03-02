package huggingface

import (
	"context"
	"fmt"
	"os"
)

const APIBaseURL = "https://api-inference.huggingface.co/models/"

const (
	ErrHuggingFaceCompletion = "huggingface completion error"
)

type HuggingFaceMode int

const (
	HuggingFaceModeCoversational HuggingFaceMode = iota
	HuggingFaceModeTextGeneration
)

type HuggingFace struct {
	mode        HuggingFaceMode
	token       string
	model       string
	temperature float32
	maxLength   *int
	minLength   *int
	topK        *int
	topP        *float32
	verbose     bool
}

func New(model string, temperature float32, verbose bool) *HuggingFace {
	return &HuggingFace{
		mode:        HuggingFaceModeCoversational,
		token:       os.Getenv("HUGGING_FACE_HUB_TOKEN"),
		model:       model,
		temperature: temperature,
		verbose:     verbose,
	}
}

// WithModel sets the model to use for the LLM
func (h *HuggingFace) WithModel(model string) *HuggingFace {
	h.model = model
	return h
}

// WithTemperature sets the temperature to use for the LLM
func (h *HuggingFace) WithTemperature(temperature float32) *HuggingFace {
	h.temperature = temperature
	return h
}

// WithMaxLength sets the maxLength to use for the LLM
func (h *HuggingFace) WithMaxLength(maxLength int) *HuggingFace {
	h.maxLength = &maxLength
	return h
}

// WithMinLength sets the minLength to use for the LLM
func (h *HuggingFace) WithMinLength(minLength int) *HuggingFace {
	h.minLength = &minLength
	return h
}

// WithToken sets the API key to use for the LLM
func (h *HuggingFace) WithToken(token string) *HuggingFace {
	h.token = token
	return h
}

// WithVerbose sets the verbose flag to use for the LLM
func (h *HuggingFace) WithVerbose(verbose bool) *HuggingFace {
	h.verbose = verbose
	return h
}

// WithTopK sets the topK to use for the LLM
func (h *HuggingFace) WithTopK(topK int) *HuggingFace {
	h.topK = &topK
	return h
}

// WithTopP sets the topP to use for the LLM
func (h *HuggingFace) WithTopP(topP float32) *HuggingFace {
	h.topP = &topP
	return h
}

// WithMode sets the mode to use for the LLM
func (h *HuggingFace) WithMode(mode HuggingFaceMode) *HuggingFace {
	h.mode = mode
	return h
}

// Completion returns the completion for the given prompt
func (h *HuggingFace) Completion(ctx context.Context, prompt string) (string, error) {

	var output string
	var outputs []string
	var err error
	switch h.mode {
	case HuggingFaceModeTextGeneration:
		outputs, err = h.textgenerationCompletion(ctx, []string{prompt})
		if err == nil {
			output = outputs[0]
		}
	case HuggingFaceModeCoversational:
		fallthrough
	default:
		output, err = h.conversationalCompletion(ctx, prompt)
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrHuggingFaceCompletion, err)
	}

	return output, nil
}

// BatchCompletion returns the completion for the given prompts
func (h *HuggingFace) BatchCompletion(ctx context.Context, prompts []string) ([]string, error) {

	var outputs []string
	var err error
	switch h.mode {
	case HuggingFaceModeTextGeneration:
		outputs, err = h.textgenerationCompletion(ctx, prompts)
	case HuggingFaceModeCoversational:
		fallthrough
	default:
		return nil, fmt.Errorf("batch completion not supported for conversational mode")
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrHuggingFaceCompletion, err)
	}

	return outputs, nil
}
