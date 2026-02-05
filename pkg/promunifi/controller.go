package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type controller struct {
	Info                      *prometheus.Desc
	UptimeSeconds             *prometheus.Desc
	UpdateAvailable           *prometheus.Desc
	UpdateDownloaded          *prometheus.Desc
	AutobackupEnabled         *prometheus.Desc
	WebRTCSupport             *prometheus.Desc
	IsCloudConsole            *prometheus.Desc
	DataRetentionDays         *prometheus.Desc
	DataRetention5minHours    *prometheus.Desc
	DataRetentionHourlyHours  *prometheus.Desc
	DataRetentionDailyHours   *prometheus.Desc
	DataRetentionMonthlyHours *prometheus.Desc
	UnsupportedDeviceCount    *prometheus.Desc
	InformPort                *prometheus.Desc
	HTTPSPort                 *prometheus.Desc
	PortalHTTPPort            *prometheus.Desc
}

func descController(ns string) *controller {
	labels := []string{"hostname", "site_name", "source"}
	infoLabels := []string{"version", "build", "device_type", "console_version", "hostname", "site_name", "source"}

	nd := prometheus.NewDesc

	return &controller{
		Info:                      nd(ns+"controller_info", "Controller information (always 1)", infoLabels, nil),
		UptimeSeconds:             nd(ns+"controller_uptime_seconds", "Controller uptime in seconds", labels, nil),
		UpdateAvailable:           nd(ns+"controller_update_available", "Update available (1/0)", labels, nil),
		UpdateDownloaded:          nd(ns+"controller_update_downloaded", "Update downloaded (1/0)", labels, nil),
		AutobackupEnabled:         nd(ns+"controller_autobackup_enabled", "Auto backup enabled (1/0)", labels, nil),
		WebRTCSupport:             nd(ns+"controller_webrtc_support", "WebRTC supported (1/0)", labels, nil),
		IsCloudConsole:            nd(ns+"controller_is_cloud_console", "Is cloud console (1/0)", labels, nil),
		DataRetentionDays:         nd(ns+"controller_data_retention_days", "Data retention in days", labels, nil),
		DataRetention5minHours:    nd(ns+"controller_data_retention_5min_hours", "5-minute scale retention hours", labels, nil),
		DataRetentionHourlyHours:  nd(ns+"controller_data_retention_hourly_hours", "Hourly scale retention hours", labels, nil),
		DataRetentionDailyHours:   nd(ns+"controller_data_retention_daily_hours", "Daily scale retention hours", labels, nil),
		DataRetentionMonthlyHours: nd(ns+"controller_data_retention_monthly_hours", "Monthly scale retention hours", labels, nil),
		UnsupportedDeviceCount:    nd(ns+"controller_unsupported_device_count", "Number of unsupported devices", labels, nil),
		InformPort:                nd(ns+"controller_inform_port", "Inform port number", labels, nil),
		HTTPSPort:                 nd(ns+"controller_https_port", "HTTPS port number", labels, nil),
		PortalHTTPPort:            nd(ns+"controller_portal_http_port", "Portal HTTP port number", labels, nil),
	}
}

func (u *promUnifi) exportSysinfo(r report, s *unifi.Sysinfo) {
	hostname := s.Hostname
	if hostname == "" {
		hostname = s.Name
	}

	if hostname == "" {
		hostname = s.SiteName // fallback when API omits both (e.g. remote/cloud)
	}

	labels := []string{hostname, s.SiteName, s.SourceName}
	infoLabels := []string{s.Version, s.Build, s.DeviceType, s.ConsoleVer, hostname, s.SiteName, s.SourceName}

	updateAvail := 0
	if s.UpdateAvail {
		updateAvail = 1
	}

	updateDown := 0
	if s.UpdateDown {
		updateDown = 1
	}

	autobackup := 0
	if s.Autobackup {
		autobackup = 1
	}

	webrtc := 0
	if s.HasWebRTC {
		webrtc = 1
	}

	cloud := 0
	if s.IsCloud {
		cloud = 1
	}

	r.send([]*metric{
		{u.Controller.Info, gauge, 1, infoLabels},
		{u.Controller.UptimeSeconds, gauge, s.Uptime, labels},
		{u.Controller.UpdateAvailable, gauge, updateAvail, labels},
		{u.Controller.UpdateDownloaded, gauge, updateDown, labels},
		{u.Controller.AutobackupEnabled, gauge, autobackup, labels},
		{u.Controller.WebRTCSupport, gauge, webrtc, labels},
		{u.Controller.IsCloudConsole, gauge, cloud, labels},
		{u.Controller.DataRetentionDays, gauge, s.DataRetDays, labels},
		{u.Controller.DataRetention5minHours, gauge, s.DataRet5min, labels},
		{u.Controller.DataRetentionHourlyHours, gauge, s.DataRetHour, labels},
		{u.Controller.DataRetentionDailyHours, gauge, s.DataRetDay, labels},
		{u.Controller.DataRetentionMonthlyHours, gauge, s.DataRetMonth, labels},
		{u.Controller.UnsupportedDeviceCount, gauge, s.Unsupported, labels},
		{u.Controller.InformPort, gauge, s.InformPort, labels},
		{u.Controller.HTTPSPort, gauge, s.HTTPSPort, labels},
		{u.Controller.PortalHTTPPort, gauge, s.PortalPort, labels},
	})
}
