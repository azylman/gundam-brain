package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PromptRequest struct {
	Prompt         string `json:"prompt"`
	ConversationID string `json:"conversation_id"`
}

type Options struct {
	Port           int             `json:"port"`
	AgyBin         string          `json:"agy_bin"`
	ApiKey         string          `json:"api_key"`
	Model          string          `json:"model"`
	SystemPrompt   string          `json:"system_prompt"`
	TimeoutMinutes int             `json:"timeout_minutes"`
	McpConfig      json.RawMessage `json:"mcp_config"`
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func ensureAgySettings(apiKey, model string) {
	if apiKey == "" {
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/root"
	}
	configDir := filepath.Join(homeDir, ".gemini", "antigravity-cli")
	_ = os.MkdirAll(configDir, 0755)

	settingsPath := filepath.Join(configDir, "settings.json")
	settings := map[string]interface{}{
		"modelProvider": "gemini",
		"model":         model,
	}
	if data, err := os.ReadFile(settingsPath); err == nil {
		_ = json.Unmarshal(data, &settings)
		settings["modelProvider"] = "gemini"
		if model != "" {
			settings["model"] = model
		}
	}
	if out, err := json.MarshalIndent(settings, "", "  "); err == nil {
		_ = os.WriteFile(settingsPath, out, 0644)
	}
}

func ensureSystemRules(customPrompt string) {
	if strings.TrimSpace(customPrompt) == "" {
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/root"
	}
	rulesDir := filepath.Join(homeDir, ".gemini", "rules")
	_ = os.MkdirAll(rulesDir, 0755)

	overrideContent := "# User Custom Instructions\n" + strings.TrimSpace(customPrompt) + "\n"
	ruleFile := filepath.Join(rulesDir, "user_override.md")
	_ = os.WriteFile(ruleFile, []byte(overrideContent), 0644)
	log.Printf("Configured custom user rules in %s", ruleFile)
}

func ensureMcpConfig(rawConfig json.RawMessage) {
	if len(rawConfig) == 0 {
		return
	}
	trimmed := strings.TrimSpace(string(rawConfig))
	if trimmed == "" || trimmed == `""` || trimmed == "null" {
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/root"
	}
	configDir := filepath.Join(homeDir, ".gemini", "config")
	_ = os.MkdirAll(configDir, 0755)
	targetPath := filepath.Join(configDir, "mcp_config.json")

	var configContent []byte
	var strVal string
	if err := json.Unmarshal(rawConfig, &strVal); err == nil && strVal != "" {
		configContent = []byte(strVal)
	} else {
		configContent = rawConfig
	}

	var js map[string]interface{}
	if err := json.Unmarshal(configContent, &js); err == nil {
		if formatted, err := json.MarshalIndent(js, "", "  "); err == nil {
			configContent = formatted
		}
	}

	if err := os.WriteFile(targetPath, configContent, 0644); err != nil {
		log.Printf("Failed to write %s: %v", targetPath, err)
	} else {
		log.Printf("Configured MCP servers in %s", targetPath)
	}
}

func extractResponseAndError(convID string) (string, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/root"
	}

	targetDirs := []string{
		filepath.Join(homeDir, ".gemini", "antigravity-cli", "brain", convID),
		filepath.Join(homeDir, ".gemini", "antigravity", "brain", convID),
	}

	var lastResponse string
	var lastError string

	for _, dir := range targetDirs {
		for _, name := range []string{"transcript_full.jsonl", "transcript.jsonl"} {
			tPath := filepath.Join(dir, ".system_generated", "logs", name)
			data, err := os.ReadFile(tPath)
			if err != nil {
				continue
			}

			lines := strings.Split(string(data), "\n")
			for i := len(lines) - 1; i >= 0; i-- {
				line := strings.TrimSpace(lines[i])
				if line == "" {
					continue
				}
				var step struct {
					Type      string          `json:"type"`
					Status    string          `json:"status"`
					Error     json.RawMessage `json:"error"`
					Content   string          `json:"content"`
					Thinking  string          `json:"thinking"`
					ToolCalls json.RawMessage `json:"tool_calls"`
				}
				if err := json.Unmarshal([]byte(line), &step); err == nil {
					if (step.Status == "ERROR" || len(step.Error) > 0) && lastError == "" {
						if errStr := strings.TrimSpace(string(step.Error)); errStr != "" && errStr != "null" {
							lastError = errStr
						}
					}
					if step.Type == "PLANNER_RESPONSE" && lastResponse == "" {
						if strings.TrimSpace(step.Content) != "" {
							lastResponse = step.Content
						} else if len(step.ToolCalls) > 2 && string(step.ToolCalls) != "[]" && string(step.ToolCalls) != "null" {
							lastResponse = fmt.Sprintf("[Tool Call Requested]: %s", string(step.ToolCalls))
						}
					}
				}
			}
			if lastResponse != "" || lastError != "" {
				return lastResponse, lastError
			}
		}
	}
	return lastResponse, lastError
}

func loadConfig() (string, string, string, string, string, int, json.RawMessage) {
	port := getEnv("PORT", "8080")
	agyBin := getEnv("AGY_BIN", "agy")
	apiKey := getEnv("GEMINI_API_KEY", getEnv("ANTIGRAVITY_API_KEY", ""))
	model := getEnv("AGY_MODEL", "Gemini 3.7 Flash (High)")
	systemPrompt := getEnv("SYSTEM_PROMPT", "")
	timeoutMinutes := 15
	if tm := os.Getenv("TIMEOUT_MINUTES"); tm != "" {
		if val, err := strconv.Atoi(tm); err == nil && val > 0 {
			timeoutMinutes = val
		}
	}
	var mcpConfig json.RawMessage

	// Read Home Assistant add-on options if available
	if data, err := os.ReadFile("/data/options.json"); err == nil {
		var opts Options
		if err := json.Unmarshal(data, &opts); err == nil {
			if opts.Port != 0 {
				port = fmt.Sprintf("%d", opts.Port)
			}
			if strings.TrimSpace(opts.AgyBin) != "" {
				agyBin = opts.AgyBin
			}
			if strings.TrimSpace(opts.ApiKey) != "" {
				apiKey = opts.ApiKey
			}
			if strings.TrimSpace(opts.Model) != "" {
				model = opts.Model
			}
			if strings.TrimSpace(opts.SystemPrompt) != "" {
				systemPrompt = opts.SystemPrompt
			}
			if opts.TimeoutMinutes > 0 {
				timeoutMinutes = opts.TimeoutMinutes
			}
			if len(opts.McpConfig) > 0 {
				mcpConfig = opts.McpConfig
			}
		}
	}

	if apiKey != "" {
		ensureAgySettings(apiKey, model)
	}
	ensureSystemRules(systemPrompt)
	if len(mcpConfig) > 0 {
		ensureMcpConfig(mcpConfig)
	}

	return port, agyBin, apiKey, model, systemPrompt, timeoutMinutes, mcpConfig
}

func handlePrompt(agyBin, apiKey, model, systemPrompt string, timeoutMinutes int, mcpConfig json.RawMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req PromptRequest
		if err := json.Unmarshal(body, &req); err != nil || strings.TrimSpace(req.Prompt) == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid payload: 'prompt' field is required and cannot be empty",
			})
			return
		}

		convID := strings.TrimSpace(req.ConversationID)
		if convID == "" {
			convID = uuid.New().String()
		}

		// Spawn headless Antigravity CLI in a background goroutine with bounded execution timeout
		go func(prompt, convID string) {
			log.Printf("Starting background execution for prompt: %q (conversation: %s, timeout: %d minutes)", prompt, convID, timeoutMinutes)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
			defer cancel()

			args := []string{"--dangerously-skip-permissions"}
			if model != "" {
				args = append(args, "--model", model)
			}
			args = append(args, "--conversation", convID, "-p", prompt)

			cmd := exec.CommandContext(ctx, agyBin, args...)
			cmd.Dir = "/app"
			cmd.Stdin = strings.NewReader("")
			if apiKey != "" {
				cmd.Env = append(os.Environ(),
					"GEMINI_API_KEY="+apiKey,
				)
			}

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					exitCode = -1
					log.Printf("Process error: %v", err)
				}
			}

			outText, errDetail := extractResponseAndError(convID)
			if outText == "" {
				outText = strings.TrimSpace(stdout.String())
			}

			logDetails := ""
			if errDetail != "" {
				logDetails = fmt.Sprintf("\n--- ERROR DETAILS ---\n%s", errDetail)
			}

			log.Printf("Execution finished | conversation=%s exit_code=%d\n--- STDOUT / RESPONSE ---\n%s\n--- STDERR ---\n%s%s",
				convID, exitCode, outText, stderr.String(), logDetails)
		}(req.Prompt, convID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":          "accepted",
			"conversation_id": convID,
			"message":         "Prompt execution started in background",
		})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	port, agyBin, apiKey, model, systemPrompt, timeoutMinutes, mcpConfig := loadConfig()

	mux := http.NewServeMux()
	promptHandler := handlePrompt(agyBin, apiKey, model, systemPrompt, timeoutMinutes, mcpConfig)

	mux.HandleFunc("POST /", promptHandler)
	mux.HandleFunc("POST /prompt", promptHandler)
	mux.HandleFunc("POST /api/prompt", promptHandler)
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/health", handleHealth)

	addr := ":" + port
	authStatus := "no API key configured"
	if apiKey != "" {
		authStatus = "Gemini API key configured"
	}
	mcpStatus := "no MCP servers configured"
	if len(mcpConfig) > 0 {
		mcpStatus = "custom MCP config loaded"
	}
	log.Printf("Gundam Brain server listening on %s (agy binary: %s, model: %s, timeout: %dm, auth: %s, mcp: %s)",
		addr, agyBin, model, timeoutMinutes, authStatus, mcpStatus)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
