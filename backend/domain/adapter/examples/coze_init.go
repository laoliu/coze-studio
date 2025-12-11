/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
	"github.com/coze-dev/coze-studio/backend/infra/mysql"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
	"github.com/joho/godotenv"
)

// InitCozeConfig initializes Coze configuration for standalone programs
// This replicates the initialization done in main.go and application/base/appinfra
func InitCozeConfig(ctx context.Context) error {
	// 1. Load environment variables from .env file
	projectRoot := findProjectRoot()
	envFile := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envFile); err != nil {
		// .env file is optional, continue if not found
		fmt.Printf("Warning: .env file not found at %s, using system environment variables\n", envFile)
	}

	// 2. Initialize database connection (required for model config)
	db, err := mysql.New()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 3. Initialize storage (required for model config)
	oss, err := storage.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 4. Initialize config with DB and OSS
	if err := config.Init(ctx, db, oss); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	fmt.Println("✓ Coze configuration initialized successfully")
	return nil
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot() string {
	// Try to find go.mod in parent directories
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	// Fallback to current directory
	return "."
}
