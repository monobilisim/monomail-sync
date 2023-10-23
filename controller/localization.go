package controller

import "flag"

var (
	lang = flag.String("lang", "en", "Language")
	Data map[string]map[string]string
)

func InitLocalizer() {

	switch *lang {
	case "en":
		Data = map[string]map[string]string{
			"index": {
				"source_details":      "Source Details",
				"destination_details": "Destination Details",
				"server":              "Server",
				"account":             "Account",
				"password":            "Password",
				"validate":            "Validate Credentials",
				"sync":                "Start Synchronization",
				"user_queue":          "User Table",
			},
			"login": {
				"sign_in":     "Sign In",
				"description": "Sign in to access admin panel",
				"username":    "Username",
				"password":    "Password",
			},
			"admin": {
				"index":          "Index",
				"source_server":  "Source Server",
				"source_account": "Source Account",
				"dest_server":    "Destination Server",
				"dest_account":   "Destination Account",
				"status":         "Status",
				"actions":        "Actions",
			},
			"table": {
				"index":          "Index",
				"source_server":  "Source Server",
				"source_account": "Source Account",
				"dest_server":    "Destination Server",
				"dest_account":   "Destination Account",
				"status":         "Status",
				"actions":        "Actions",
			},
		}
	case "tr":
		Data = map[string]map[string]string{
			"index": {
				"source_details":      "Kaynak Bilgileri",
				"destination_details": "Hedef Bilgileri",
				"server":              "Sunucu",
				"account":             "Hesap",
				"password":            "Parola",
				"validate":            "Bilgileri Doğrula",
				"sync":                "Senkronizasyonu Başlat",
				"user_queue":          "Kullanıcı İşlem Kuyruğu",
			},
			"login": {
				"sign_in":     "Giriş Yap",
				"description": "Admin paneline erişebilmek için giriş yapın.",
				"username":    "Kullanıcı Adı",
				"password":    "Parola",
			},
			"admin": {
				"index":          "Sıra",
				"source_server":  "Kaynak Sunucu",
				"source_account": "Kaynak Hesap",
				"dest_server":    "Hedef Sunucu",
				"dest_account":   "Hedef Hesap",
				"status":         "Durum",
				"actions":        "Eylemler",
			},
			"table": {
				"index":          "Sıra",
				"source_server":  "Kaynak Sunucu",
				"source_account": "Kaynak Hesap",
				"dest_server":    "Hedef Sunucu",
				"dest_account":   "Hedef Hesap",
				"status":         "Durum",
				"actions":        "Eylemler",
			},
		}
	}
}
