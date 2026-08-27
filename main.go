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
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type Options struct {
	Port   int    `json:"port"`
	AgyBin string `json:"agy_bin"`
	ApiKey string `json:"api_key"`
	Model  string `json:"model"`
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

func loadConfig() (string, string, string, string) {
	port := getEnv("PORT", "8080")
	agyBin := getEnv("AGY_BIN", "agy")
	apiKey := getEnv("GEMINI_API_KEY", getEnv("ANTIGRAVITY_API_KEY", ""))
	model := getEnv("AGY_MODEL", "gemini-2.5-flash")

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
		}
	}

	if apiKey != "" {
		ensureAgySettings(apiKey, model)
	}

	return port, agyBin, apiKey, model
}

func handlePrompt(agyBin, apiKey, model string) http.HandlerFunc {
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

			logDetails := ""
			if logData, err := os.ReadFile(logFile); err == nil && len(logData) > 0 {
				logDetails = fmt.Sprintf("\n--- LOG FILE (%s) ---\n%s", logFile, string(logData))
			}

			log.Printf("Execution finished | exit_code=%d\n--- STDOUT ---\n%s\n--- STDERR ---\n%s%s",
				exitCode, stdout.String(), stderr.String(), logDetails)
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
	port, agyBin, apiKey, model := loadConfig()

	mux := http.NewServeMux()
	promptHandler := handlePrompt(agyBin, apiKey, model)

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
	log.Printf("Gundam Brain server listening on %s (agy binary: %s, model: %s, auth: %s)", addr, agyBin, model, authStatus)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
