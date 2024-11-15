package api

import (
	"encoding/json"
	"net/http"

	"github.com/vrutik2809/compile-x/core"
	"github.com/vrutik2809/compile-x/core/runner"
)

// Define a struct to represent the JSON request payload
type RequestPayload struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// Define a struct to represent the JSON response
type ResponsePayload struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleCodeExecution(w http.ResponseWriter, r *http.Request) {
	// Set response content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Only accept POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponsePayload{
			Message: "Only POST method is allowed",
			Data:    nil,
		})
		return
	}

	// Decode the incoming JSON payload
	var reqPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Message: "Invalid JSON payload",
			Data:    nil,
		})
		return
	}

	var sourceCode core.SourceCode
	if reqPayload.Language == "java" {
		sourceCode = core.NewSourceCode(core.JAVA_22, reqPayload.Code)
	} else if reqPayload.Language == "cpp" {
		sourceCode = core.NewSourceCode(core.CPP_17_20, reqPayload.Code)
	} else if reqPayload.Language == "python" {
		sourceCode = core.NewSourceCode(core.PYTHON_3_12, reqPayload.Code)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Message: "Invalid language",
			Data:    nil,
		})
		return
	}
	cli, err := core.GetDockerClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Message: "Internal server error",
			Data:    nil,
		})
		return
	}
	runner := runner.NewSandBoxedRunner(cli)
	output := runner.Run(sourceCode)

	// Respond with JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{
		Message: "success",
		Data: map[string]interface{}{
			"output":   output.GetOutput(),
			"exitCode": output.GetExitCode(),
		},
	})
}