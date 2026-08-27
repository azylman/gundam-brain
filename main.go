package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type Options struct {
	Port      int             `json:"port"`
	AgyBin    string          `json:"agy_bin"`
	ApiKey    string          `json:"api_key"`
	Model     string          `json:"model"`
	McpConfig json.RawMessage `json:"mcp_config"`
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

func extractLatestResponse() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/root"
	}
	brainDir := filepath.Join(homeDir, ".gemini", "antigravity-cli", "brain")
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return ""
	}

	var latestDir string
	var latestTime time.Time
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err == nil && info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestDir = filepath.Join(brainDir, entry.Name())
			}
		}
	}

	if latestDir == "" {
		return ""
	}

	transcriptPath := filepath.Join(latestDir, ".system_generated", "logs", "transcript.jsonl")
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		transcriptPath = filepath.Join(latestDir, ".system_generated", "logs", "transcript_full.jsonl")
		data, err = os.ReadFile(transcriptPath)
		if err != nil {
			return ""
		}
	}

	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var step struct {
			Type    string `json:"type"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(line), &step); err == nil {
			if step.Type == "PLANNER_RESPONSE" && strings.TrimSpace(step.Content) != "" {
				return step.Content
			}
		}
	}
	return ""
}

func loadConfig() (string, string, string, string, json.RawMessage) {
	port := getEnv("PORT", "8080")
	agyBin := getEnv("AGY_BIN", "agy")
	apiKey := getEnv("GEMINI_API_KEY", getEnv("ANTIGRAVITY_API_KEY", ""))
	model := getEnv("AGY_MODEL", "Gemini 3.7 Flash (High)")
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
			if len(opts.McpConfig) > 0 {
				mcpConfig = opts.McpConfig
			}
		}
	}

	if apiKey != "" {
		ensureAgySettings(apiKey, model)
	}
	if len(mcpConfig) > 0 {
		ensureMcpConfig(mcpConfig)
	}

	return port, agyBin, apiKey, model, mcpConfig
}

func handlePrompt(agyBin, apiKey, model string, mcpConfig json.RawMessage) http.HandlerFunc {
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

		// Spawn headless Antigravity CLI in a background goroutine
		go func(prompt string) {
			log.Printf("Starting background execution for prompt: %q", prompt)
			ensureAgySettings(apiKey, model)
			if len(mcpConfig) > 0 {
				ensureMcpConfig(mcpConfig)
			}

			logFile := "/tmp/agy.log"
			_ = os.Remove(logFile)

			args := []string{"--dangerously-skip-permissions"}
			if model != "" {
				args = append(args, "--model", model)
			}
			args = append(args, "--log-file", logFile, "-p", prompt)

			cmd := exec.Command(agyBin, args...)
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
					log.Printf("Failed to spawn %s: %v", agyBin, err)
				}
			}

			outText := strings.TrimSpace(stdout.String())
			if outText == "" {
				if resp := extractLatestResponse(); resp != "" {
					outText = resp
				}
			}

			logDetails := ""
			if logData, err := os.ReadFile(logFile); err == nil && len(logData) > 0 {
				logDetails = fmt.Sprintf("\n--- LOG FILE (%s) ---\n%s", logFile, string(logData))
			}

			log.Printf("Execution finished | exit_code=%d\n--- STDOUT / RESPONSE ---\n%s\n--- STDERR ---\n%s%s",
				exitCode, outText, stderr.String(), logDetails)
		}(req.Prompt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "Prompt execution started in background",
		})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	port, agyBin, apiKey, model, mcpConfig := loadConfig()

	mux := http.NewServeMux()
	promptHandler := handlePrompt(agyBin, apiKey, model, mcpConfig)

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
	log.Printf("Gundam Brain server listening on %s (agy binary: %s, model: %s, auth: %s, mcp: %s)",
		addr, agyBin, model, authStatus, mcpStatus)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
