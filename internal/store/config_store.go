package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zalando/go-keyring"
)

const (
	appName = "QTe"
	keyName = "encryption-key"
)

type connectionInfo struct {
	Driver string `json:"driver"`
	Url    string `json:"url"`
}

type connectionStore map[string]connectionInfo

// ConnectionConfig represents the JSON file structure
type ConnectionConfig struct {
	key         []byte
	Connections connectionStore `json:"connections"` // Name -> Encrypted DSN
}

type ConnectionView struct {
	Driver string
	Name   string
	Url    string
}

type ConnectionsLoadedMsg struct {
	Conns []ConnectionView
	Err   error
}

var connectionCfg *ConnectionConfig
var once sync.Once

//var mu *sync.Mutex

func ConnectionsCfg() *ConnectionConfig {
	once.Do(func() {
		cfg, err := loadConfig()
		if err != nil {
			panic("could not load/create the connections.json file: " + err.Error())
		}

		key, err := getOrGenerateKey()
		if err != nil {
			panic("could not get/generate the required key from os.keyring: " + err.Error())
		}

		cfg.key = key
		connectionCfg = cfg
	})

	return connectionCfg
}

// SaveConnection encrypts the DSN and saves it to the JSON file
func (cfg *ConnectionConfig) SaveConnection(name, driver, dsn string) error {
	encryptedDSN, err := encrypt(dsn, cfg.key)
	if err != nil {
		return err
	}

	cfg.Connections[name] = connectionInfo{
		Driver: driver,
		Url:    encryptedDSN,
	}
	return writeConfig(cfg)
}

// LoadConnection reads the JSON file and decrypts the requested DSN
func (cfg *ConnectionConfig) LoadConnection(name string) (string, string, error) {
	info, exists := cfg.Connections[name]
	if !exists {
		return "", "", errors.New("connection not found")
	}

	url, err := decrypt(info.Url, cfg.key)
	return info.Driver, url, err
}

func (cfg *ConnectionConfig) LoadAll() tea.Cmd {
	return func() tea.Msg {
		i := 0
		result := make([]ConnectionView, len(cfg.Connections))

		for name, info := range cfg.Connections {
			decryptedUrl, err := decrypt(info.Url, cfg.key)
			if err != nil {
				return ConnectionsLoadedMsg{Err: err}
			}

			result[i] = ConnectionView{Name: name, Url: decryptedUrl, Driver: info.Driver}
			i = i + 1
		}

		return ConnectionsLoadedMsg{Conns: result, Err: nil}
	}
}

func (cfg *ConnectionConfig) Exists(name, driver, url string) (bool, error) {
	info, exists := cfg.Connections[name]
	if !exists {
		return false, nil
	}

	decryptedUrl, err := decrypt(info.Url, cfg.key)
	if err != nil {
		return false, err
	}

	return decryptedUrl == url && info.Driver == driver, nil
}

// --- CONFIGURATION & FILE PERMISSIONS ---

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// Windows: %AppData%\Roaming\qte\connections.json
	// macOS: ~/Library/Application Support/qte/connections.json
	// Linux: ~/.config/qte/connections.json
	appDir := filepath.Join(configDir, appName)

	// Create directory with 0700 permissions (read/write/execute for owner only)
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", err
	}

	configFilePath := filepath.Join(appDir, "connections.json")
	if _, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_EXCL, 0600); err != nil && !errors.Is(err, fs.ErrExist) {
		return "", err
	}

	return configFilePath, nil
}

func loadConfig() (*ConnectionConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &ConnectionConfig{Connections: make(connectionStore)}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return empty config if file doesn't exist
		}
		return nil, err
	}

	if len(data) == 0 {
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func writeConfig(cfg *ConnectionConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write file with 0600 permissions (read/write for owner only)
	return os.WriteFile(path, data, 0600)
}

// --- OS KEYRING MANAGEMENT ---

func getOrGenerateKey() ([]byte, error) {
	// Try to get existing key from OS keyring
	keyBase64, err := keyring.Get(appName, keyName)
	if err == nil {
		return base64.StdEncoding.DecodeString(keyBase64)
	}

	// If not found, generate a new 32-byte (256-bit) AES key
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	// Save the new key to the OS keyring
	keyBase64 = base64.StdEncoding.EncodeToString(key)
	if err := keyring.Set(appName, keyName, keyBase64); err != nil {
		return nil, err
	}

	return key, nil
}

// --- ENCRYPTION (AES-GCM) ---

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(cryptoText string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
