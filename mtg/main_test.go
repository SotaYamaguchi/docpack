package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configJSON  string
		wantErr     bool
		wantProject map[string]string
	}{
		{
			name: "正常なconfig",
			configJSON: `{
				"projects": {
					"project-a": "PREFIX_A",
					"project-b": "PREFIX_B"
				}
			}`,
			wantErr: false,
			wantProject: map[string]string{
				"project-a": "PREFIX_A",
				"project-b": "PREFIX_B",
			},
		},
		{
			name:       "不正なJSON",
			configJSON: `{"projects": {`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 一時ファイルを作成
			tmpfile, err := os.CreateTemp("", "config*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Remove(tmpfile.Name()) }()

			if _, err := tmpfile.Write([]byte(tt.configJSON)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// テスト実行
			config, err := loadConfig(tmpfile.Name())

			if tt.wantErr {
				if err == nil {
					t.Errorf("loadConfig() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(config.Projects) != len(tt.wantProject) {
				t.Errorf("projects count = %v, want %v", len(config.Projects), len(tt.wantProject))
				return
			}

			for key, want := range tt.wantProject {
				if got, ok := config.Projects[key]; !ok || got != want {
					t.Errorf("projects[%s] = %v, want %v", key, got, want)
				}
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("loadConfig() error = nil, want error for non-existent file")
	}
}

func TestResolvePrefix(t *testing.T) {
	// テスト用の設定ファイルを作成
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	configData := map[string]any{
		"projects": map[string]string{
			"test-project": "TEST_PREFIX",
		},
	}
	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		project    string
		prefix     string
		configPath string
		want       string
		wantErr    bool
	}{
		{
			name:       "プロジェクト名で解決",
			project:    "test-project",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "TEST_PREFIX",
			wantErr:    false,
		},
		{
			name:       "プレフィックスを直接指定",
			project:    "",
			prefix:     "DIRECT_PREFIX",
			configPath: tmpfile.Name(),
			want:       "DIRECT_PREFIX",
			wantErr:    false,
		},
		{
			name:       "存在しないプロジェクト",
			project:    "nonexistent",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "",
			wantErr:    true,
		},
		{
			name:       "project と prefix 両方未指定",
			project:    "",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePrefix(tt.project, tt.prefix, tt.configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("resolvePrefix() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("resolvePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("resolvePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenameFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// テスト用のファイルを作成
	testFiles := []string{
		"PREFIX_main.txt",
		"PREFIX_main_document.pdf",
		"PREFIX_other.txt",
	}

	for _, fname := range testFiles {
		fpath := filepath.Join(tmpDir, fname)
		if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// テスト実行
	currentDate := "20260320"
	err = renameFiles("PREFIX", tmpDir, currentDate, "")
	if err != nil {
		t.Errorf("renameFiles() error = %v", err)
		return
	}

	// 結果確認
	expectedFiles := map[string]bool{
		"PREFIX_20260320.txt":          true,
		"PREFIX_20260320_document.pdf": true,
		"PREFIX_other.txt":             true,
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	foundFiles := make(map[string]bool)
	for _, entry := range entries {
		foundFiles[entry.Name()] = true
	}

	for fname := range expectedFiles {
		if !foundFiles[fname] {
			t.Errorf("Expected file %s not found", fname)
		}
	}

	for fname := range foundFiles {
		if !expectedFiles[fname] {
			t.Errorf("Unexpected file %s found", fname)
		}
	}
}

func TestRenameFilesWithSuffix(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// テスト用のファイルを作成
	testFile := "PREFIX_main.txt"
	fpath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// テスト実行
	currentDate := "20260320"
	suffix := "_MTG後"
	err = renameFiles("PREFIX", tmpDir, currentDate, suffix)
	if err != nil {
		t.Errorf("renameFiles() error = %v", err)
		return
	}

	// 結果確認
	expectedFile := "PREFIX_20260320_MTG後.txt"
	if _, err := os.Stat(filepath.Join(tmpDir, expectedFile)); os.IsNotExist(err) {
		t.Errorf("Expected file %s not found", expectedFile)
	}
}

func TestCollectFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// テスト用のファイルを作成
	testFiles := []string{
		"PREFIX_20260320.txt",
		"PREFIX_20260320_document.pdf",
		"OTHER_file.txt",
	}

	for _, fname := range testFiles {
		fpath := filepath.Join(tmpDir, fname)
		if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// テスト実行
	destFolder := filepath.Join(tmpDir, "PREFIX_資料送付_20260320")
	err = collectFiles("PREFIX", tmpDir, destFolder)
	if err != nil {
		t.Errorf("collectFiles() error = %v", err)
		return
	}

	// 結果確認: 移動先ディレクトリが作成されたか
	if _, err := os.Stat(destFolder); os.IsNotExist(err) {
		t.Errorf("Destination folder %s not created", destFolder)
		return
	}

	// 結果確認: ファイルが移動されたか
	expectedInDest := []string{
		"PREFIX_20260320.txt",
		"PREFIX_20260320_document.pdf",
	}

	for _, fname := range expectedInDest {
		destPath := filepath.Join(destFolder, fname)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s not found in destination", fname)
		}

		// 元の場所には存在しないことを確認
		srcPath := filepath.Join(tmpDir, fname)
		if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
			t.Errorf("File %s should have been moved from source", fname)
		}
	}

	// OTHER_file.txt は移動されていないことを確認
	otherPath := filepath.Join(tmpDir, "OTHER_file.txt")
	if _, err := os.Stat(otherPath); os.IsNotExist(err) {
		t.Errorf("File OTHER_file.txt should not have been moved")
	}
}

func TestLoadConfigWithMailTemplates(t *testing.T) {
	configJSON := `{
		"projects": {
			"test-project": "TEST_PREFIX"
		},
		"mail_templates": {
			"test-project": {
				"prep": {
					"to": ["customer@example.com"],
					"cc": ["team@example.com"],
					"bcc": [],
					"subject": "テスト件名",
					"body": "テスト本文"
				},
				"memo": {
					"to": ["customer@example.com"],
					"cc": [],
					"bcc": [],
					"subject": "",
					"body": "議事録本文"
				}
			}
		}
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := loadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	if config.MailTemplates == nil {
		t.Error("MailTemplates should not be nil")
		return
	}

	projectTemplate, ok := config.MailTemplates["test-project"]
	if !ok {
		t.Error("test-project template not found")
		return
	}

	// prep テンプレートのチェック
	if len(projectTemplate.Prep.To) != 1 || projectTemplate.Prep.To[0] != "customer@example.com" {
		t.Errorf("prep.To = %v, want [customer@example.com]", projectTemplate.Prep.To)
	}
	if projectTemplate.Prep.Subject != "テスト件名" {
		t.Errorf("prep.Subject = %v, want テスト件名", projectTemplate.Prep.Subject)
	}
	if projectTemplate.Prep.Body != "テスト本文" {
		t.Errorf("prep.Body = %v, want テスト本文", projectTemplate.Prep.Body)
	}

	// memo テンプレートのチェック
	if projectTemplate.Memo.Subject != "" {
		t.Errorf("memo.Subject = %v, want empty string", projectTemplate.Memo.Subject)
	}
	if projectTemplate.Memo.Body != "議事録本文" {
		t.Errorf("memo.Body = %v, want 議事録本文", projectTemplate.Memo.Body)
	}
}

func TestGetMailTemplate(t *testing.T) {
	configJSON := `{
		"projects": {
			"test-project": "TEST_PREFIX"
		},
		"mail_templates": {
			"test-project": {
				"prep": {
					"to": ["customer@example.com"],
					"cc": [],
					"bcc": [],
					"subject": "テスト件名",
					"body": "テスト本文"
				}
			}
		}
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		project     string
		mailType    string
		wantErr     bool
		wantSubject string
	}{
		{
			name:        "正常なテンプレート取得",
			project:     "test-project",
			mailType:    "prep",
			wantErr:     false,
			wantSubject: "テスト件名",
		},
		{
			name:     "存在しないプロジェクト",
			project:  "nonexistent",
			mailType: "prep",
			wantErr:  true,
		},
		{
			name:     "不正なメールタイプ",
			project:  "test-project",
			mailType: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := getMailTemplate(tmpfile.Name(), tt.project, tt.mailType)

			if tt.wantErr {
				if err == nil {
					t.Error("getMailTemplate() error = nil, wantErr true")
				}
				return
			}

			if err != nil {
				t.Errorf("getMailTemplate() error = %v, wantErr false", err)
				return
			}

			if template.Subject != tt.wantSubject {
				t.Errorf("Subject = %v, want %v", template.Subject, tt.wantSubject)
			}
		})
	}
}

func TestFormatMailOutput(t *testing.T) {
	template := &MailTemplate{
		To:      []string{"customer@example.com", "another@example.com"},
		Cc:      []string{"team@example.com"},
		Bcc:     []string{"bcc@example.com"},
		Subject: "テスト件名",
		Body:    "テスト本文\n複数行あります",
	}

	output := formatMailOutput(template)

	expectedLines := []string{
		"To: customer@example.com, another@example.com",
		"Cc: team@example.com",
		"Bcc: bcc@example.com",
		"件名: テスト件名",
		"",
		"テスト本文",
		"複数行あります",
	}

	for _, expected := range expectedLines {
		if !contains(output, expected) {
			t.Errorf("Output should contain: %q", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
