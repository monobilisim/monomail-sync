package internal

import (
	"imap-sync/config"
	"testing"
)

func setLang(lang string) {
	config.Conf.Language = lang
}

func TestInitLocalizer_English(t *testing.T) {
	setLang("en")
	InitLocalizer()

	if Data == nil {
		t.Fatal("Data should not be nil after InitLocalizer")
	}

	requiredKeys := []string{"index", "login", "admin", "table", "validation"}
	for _, key := range requiredKeys {
		if _, ok := Data[key]; !ok {
			t.Errorf("Data missing key %q", key)
		}
	}

	if Data["index"]["sync"] != "Start Synchronization" {
		t.Errorf("index sync = %q, want 'Start Synchronization'", Data["index"]["sync"])
	}
	if Data["index"]["validate"] != "Validate Credentials" {
		t.Errorf("index validate = %q, want 'Validate Credentials'", Data["index"]["validate"])
	}

	v := Data["validation"]
	validationKeys := []string{
		"source_connecting", "source_success",
		"dest_connecting", "dest_success",
		"connection_failed", "auth_failed", "missing_field",
		"folder", "total", "reset",
		"activity_log", "activity_clear",
	}
	for _, key := range validationKeys {
		if v[key] == "" {
			t.Errorf("validation[%q] should not be empty", key)
		}
	}
}

func TestInitLocalizer_Turkish(t *testing.T) {
	setLang("tr")
	InitLocalizer()

	if Data == nil {
		t.Fatal("Data should not be nil after InitLocalizer")
	}

	requiredKeys := []string{"index", "login", "admin", "table", "validation"}
	for _, key := range requiredKeys {
		if _, ok := Data[key]; !ok {
			t.Errorf("Data missing key %q for Turkish", key)
		}
	}

	if Data["index"]["sync"] != "Senkronizasyonu Başlat" {
		t.Errorf("Turkish sync = %q, want 'Senkronizasyonu Başlat'", Data["index"]["sync"])
	}
	if Data["validation"]["source_connecting"] != "Kaynak sunucuya bağlanılıyor..." {
		t.Errorf("Turkish source_connecting = %q", Data["validation"]["source_connecting"])
	}
	if Data["validation"]["folder"] != "Klasör" {
		t.Errorf("Turkish folder = %q, want 'Klasör'", Data["validation"]["folder"])
	}
}

func TestInitLocalizer_KeyParity(t *testing.T) {
	setLang("en")
	InitLocalizer()
	enData := make(map[string]map[string]string)
	for k, v := range Data {
		enData[k] = make(map[string]string)
		for kk, vv := range v {
			enData[k][kk] = vv
		}
	}

	setLang("tr")
	InitLocalizer()
	trData := Data

	for section, enMap := range enData {
		trMap, ok := trData[section]
		if !ok {
			t.Errorf("Turkish missing section %q", section)
			continue
		}

		for key := range enMap {
			if _, ok := trMap[key]; !ok {
				t.Errorf("Turkish section %q missing key %q", section, key)
			}
		}

		for key := range trMap {
			if _, ok := enMap[key]; !ok {
				t.Errorf("English section %q missing key %q (present in Turkish)", section, key)
			}
		}
	}
}

func TestInitLocalizer_NoEmptyValues(t *testing.T) {
	for _, l := range []string{"en", "tr"} {
		setLang(l)
		InitLocalizer()

		for section, m := range Data {
			for key, val := range m {
				if val == "" {
					t.Errorf("[%s] Data[%q][%q] is empty", l, section, key)
				}
			}
		}
	}
}
