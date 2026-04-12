package internal

import (
	"imap-sync/config"
)

var (
	lang string
	Data map[string]map[string]string
)

func InitLocalizer() {
	lang = config.Conf.Language

	switch lang {
	case "en":
		Data = map[string]map[string]string{
			"index": {
				"source_details":      "Source Details",
				"destination_details": "Destination Details",
				"server":              "Server",
				"account":             "Account",
				"account_name":        "Username",
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
				"queue":          "Queue",
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
			"notify": {
				"success":     "Successful",
				"success_msg": " successfully synchronized.",
				"fail":        "Failed",
				"fail_msg":    " failed synchronizztion.",
			},
			"validation": {
				"source_connecting":   "Connecting to source server...",
				"source_success":      "Source connection successful",
				"dest_connecting":     "Connecting to destination server...",
				"dest_success":        "Destination connection successful",
				"dest_waiting":        "Waiting...",
				"dest_skipped":        "Skipped",
				"connection_failed":   "Could not reach server",
				"auth_failed":         "Wrong password or username",
				"missing_field":       "Required field is empty",
				"source_label":        "Source",
				"dest_label":          "Destination",
				"stats_loading":       "Fetching mailbox statistics...",
				"stats_success":       "Comparison ready",
				"stats_error":         "Could not fetch statistics",
				"folder":              "Folder",
				"source_count":        "Source Messages",
				"source_size":         "Source Size",
				"dest_count":          "Dest Messages",
				"dest_size":           "Dest Size",
				"total":               "Total",
				"reset":               "Reset",
				"sync_added":          "Task added to queue",
				"activity_log":        "Activity Log",
				"activity_validating": "Validating credentials...",
				"activity_syncing":    "Sync started, added to queue",
				"activity_completed":  "Completed successfully",
				"activity_failed":     "Failed",
				"activity_pending":    "Waiting in queue...",
				"activity_no_events":  "No activity yet. Start a migration to see live updates.",
				"activity_clear":      "Clear",
			},
		}
	case "tr":
		Data = map[string]map[string]string{
			"index": {
				"source_details":      "Kaynak Bilgileri",
				"destination_details": "Hedef Bilgileri",
				"server":              "Sunucu",
				"account":             "Hesap",
				"account_name":        "Kullanıcı Adı",
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
				"queue":          "İşlem Kuyruğu",
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
			"notify": {
				"success":     "Başarılı",
				"success_msg": " mail adresleri arasındaki senkronizasyon başarıyla tamamlandı.",
				"fail":        "Başarısız",
				"fail_msg":    " mail adresleri arasındaki senkronizasyon başarısız oldu.",
			},
			"validation": {
				"source_connecting":   "Kaynak sunucuya bağlanılıyor...",
				"source_success":      "Kaynak bağlantısı başarılı",
				"dest_connecting":     "Hedef sunucuya bağlanılıyor...",
				"dest_success":        "Hedef bağlantısı başarılı",
				"dest_waiting":        "Bekliyor...",
				"dest_skipped":        "Atlandı",
				"connection_failed":   "Sunucuya erişilemedi",
				"auth_failed":         "Şifre veya kullanıcı adı hatalı",
				"missing_field":       "Gerekli alan boş",
				"source_label":        "Kaynak",
				"dest_label":          "Hedef",
				"stats_loading":       "Posta kutusu bilgileri alınıyor...",
				"stats_success":       "Karşılaştırma hazır",
				"stats_error":         "İstatistikler alınamadı",
				"folder":              "Klasör",
				"source_count":        "Kaynak Adet",
				"source_size":         "Kaynak Boyut",
				"dest_count":          "Hedef Adet",
				"dest_size":           "Hedef Boyut",
				"total":               "Toplam",
				"reset":               "Sıfırla",
				"sync_added":          "İşlem kuyruğa eklendi",
				"activity_log":        "Aktivite",
				"activity_validating": "Kimlik bilgileri doğrulanıyor...",
				"activity_syncing":    "Senkronizasyon başlatıldı, kuyruğa eklendi",
				"activity_completed":  "Başarıyla tamamlandı",
				"activity_failed":     "Başarısız",
				"activity_pending":    "Kuyrukta bekliyor...",
				"activity_no_events":  "Henüz aktivite yok. Canlı güncellemeleri görmek için bir taşıma başlatın.",
				"activity_clear":      "Temizle",
			},
		}
	}
}
