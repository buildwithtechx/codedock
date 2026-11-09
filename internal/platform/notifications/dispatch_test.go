package notifications

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubNotifSettingsRepo struct {
	NotificationSettingsRepository
	settings      *NotificationSettings
	subscriptions map[string][]NotificationSubscription
}

func (s *stubNotifSettingsRepo) GetNotificationSettings(context.Context) (*NotificationSettings, error) {
	return s.settings, nil
}

func (s *stubNotifSettingsRepo) ListSubscriptions(_ context.Context, userID, _ string) ([]NotificationSubscription, error) {
	if s.subscriptions != nil {
		return s.subscriptions[userID], nil
	}
	return nil, nil
}

type stubDeliveriesRepo struct {
	NotificationDeliveryRepository
	recorded []*NotificationDelivery
}

func (s *stubDeliveriesRepo) RecordDelivery(_ context.Context, delivery *NotificationDelivery) error {
	s.recorded = append(s.recorded, delivery)
	return nil
}

type stubDeviceTokenRepo struct {
	DeviceTokenRepository
	tokens []*DeviceToken
}

func (s *stubDeviceTokenRepo) ListAllTokens(context.Context) ([]*DeviceToken, error) {
	return s.tokens, nil
}

func TestDispatcherRecordsSuccessfulDelivery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	deliveries := &stubDeliveriesRepo{}
	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: &NotificationSettings{
		DiscordEnabled:    true,
		DiscordWebhookURL: server.URL,
	}}, nil, deliveries, nil, nil)
	err := dispatcher.Send(&NotificationEvent{EventType: "test_global_discord", Title: "hello"})
	if err != nil {
		t.Fatalf("send test notification: %v", err)
	}
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 recorded delivery, got %d", len(deliveries.recorded))
	}
	recorded := deliveries.recorded[0]
	if recorded.Channel != "discord" || !recorded.Success || recorded.Error != "" {
		t.Fatalf("unexpected delivery record: %+v", recorded)
	}
}

func TestDispatcherRecordsFailedDelivery(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: &NotificationSettings{}}, nil, deliveries, nil, nil)
	err := dispatcher.Send(&NotificationEvent{EventType: "test_global_slack", Title: "hello"})
	if err == nil {
		t.Fatal("expected send to fail for unconfigured provider")
	}
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 recorded delivery, got %d", len(deliveries.recorded))
	}
	recorded := deliveries.recorded[0]
	if recorded.Channel != "slack" || recorded.Success {
		t.Fatalf("unexpected delivery record: %+v", recorded)
	}
	if recorded.Error == "" {
		t.Fatal("expected error text on failed delivery")
	}
}

func TestDispatcherFanoutRealEvent(t *testing.T) {
	discordServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer discordServer.Close()

	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer slackServer.Close()

	deliveries := &stubDeliveriesRepo{}
	settings := &NotificationSettings{
		NotificationAlerts: true,
		DiscordEnabled:     true,
		DiscordWebhookURL:  discordServer.URL,
		SlackEnabled:       true,
		SlackWebhookURL:    slackServer.URL,
	}

	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: settings}, nil, deliveries, nil, nil)
	err := dispatcher.Send(&NotificationEvent{
		EventType: EventDeploySucceeded,
		Title:     "Deploy Succeeded",
		Message:   "Service web deployed",
		Level:     "success",
		URL:       "https://dashboard.example.com",
	})
	if err != nil {
		t.Fatalf("expected successful fanout, got: %v", err)
	}

	if len(deliveries.recorded) != 2 {
		t.Fatalf("expected 2 recorded deliveries (discord + slack), got %d", len(deliveries.recorded))
	}
}

func TestDispatcherSilencedWhenAlertsDisabled(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	settings := &NotificationSettings{
		NotificationAlerts: false,
		DiscordEnabled:     true,
		DiscordWebhookURL:  "https://discord.com/api/webhooks/test",
	}

	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: settings}, nil, deliveries, nil, nil)
	err := dispatcher.Send(&NotificationEvent{
		EventType: EventDeployFailed,
		Title:     "Deploy Failed",
	})
	if err != nil {
		t.Fatalf("expected nil on silenced alerts, got: %v", err)
	}

	if len(deliveries.recorded) != 0 {
		t.Fatalf("expected 0 deliveries when alerts silenced, got %d", len(deliveries.recorded))
	}
}

func TestDispatcherFanoutExpoPush(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	settings := &NotificationSettings{
		NotificationAlerts: true,
	}
	deviceTokens := &stubDeviceTokenRepo{
		tokens: []*DeviceToken{
			{Token: "ExponentPushToken[device-mobile-1]", Platform: "expo"},
			{Token: "unsupported-raw-token", Platform: "ios"}, // Non-expo token skipped
		},
	}

	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: settings}, nil, deliveries, nil, deviceTokens)
	_ = dispatcher.Send(&NotificationEvent{
		EventType: EventServerUnreachable,
		Title:     "Server unreachable: srv-1",
		Message:   "Docker connection lost",
		Level:     "error",
	})

	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 delivery record for expo token, got %d", len(deliveries.recorded))
	}
	if deliveries.recorded[0].Channel != "expo" {
		t.Fatalf("expected channel expo, got %s", deliveries.recorded[0].Channel)
	}
}

func TestDispatcherEnforcesUserCategorySubscription(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	settings := &NotificationSettings{
		NotificationAlerts: true,
	}

	notifRepo := &stubNotifSettingsRepo{
		settings: settings,
		subscriptions: map[string][]NotificationSubscription{
			"user-opted-out": {
				{Category: "deployment", Enabled: false},
			},
			"user-opted-in": {
				{Category: "deployment", Enabled: true},
			},
		},
	}

	deviceTokens := &stubDeviceTokenRepo{
		tokens: []*DeviceToken{
			{UserID: "user-opted-out", Token: "ExponentPushToken[user-1-phone]", Platform: "expo"},
			{UserID: "user-opted-in", Token: "ExponentPushToken[user-2-phone]", Platform: "expo"},
		},
	}

	dispatcher := NewDispatcherService(nil, notifRepo, nil, deliveries, nil, deviceTokens)
	_ = dispatcher.Send(&NotificationEvent{
		EventType: EventDeploySucceeded,
		Title:     "Deploy Succeeded",
		Message:   "Production release completed",
	})

	// Only user-opted-in's token should have received the notification
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected exactly 1 delivery record for opted-in user, got %d", len(deliveries.recorded))
	}
}

func TestDispatcherSpecificEventOverridesCategory(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	settings := &NotificationSettings{
		NotificationAlerts: true,
	}

	notifRepo := &stubNotifSettingsRepo{
		settings: settings,
		subscriptions: map[string][]NotificationSubscription{
			"user-granular": {
				{Category: "deployment", Enabled: true},
				{Category: "deploy.succeeded", Enabled: false},
			},
		},
	}

	deviceTokens := &stubDeviceTokenRepo{
		tokens: []*DeviceToken{
			{UserID: "user-granular", Token: "ExponentPushToken[user-granular-phone]", Platform: "expo"},
		},
	}

	dispatcher := NewDispatcherService(nil, notifRepo, nil, deliveries, nil, deviceTokens)

	// deploy.succeeded is specifically disabled -> 0 deliveries
	_ = dispatcher.Send(&NotificationEvent{
		EventType: EventDeploySucceeded,
		Title:     "Deploy Succeeded",
	})
	if len(deliveries.recorded) != 0 {
		t.Fatalf("expected 0 deliveries for silenced specific event, got %d", len(deliveries.recorded))
	}

	// deploy.failed falls back to category 'deployment' which is enabled -> 1 delivery
	_ = dispatcher.Send(&NotificationEvent{
		EventType: EventDeployFailed,
		Title:     "Deploy Failed",
	})
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 delivery for enabled category, got %d", len(deliveries.recorded))
	}
}

func TestDispatcherGlobalExpoProviderConfigured(t *testing.T) {
	dispatcher := &DispatcherService{}

	s1 := &NotificationSettings{ExpoEnabled: true, ExpoPushToken: "ExponentPushToken[test-channel-token]"}
	if !dispatcher.isProviderConfigured(s1, "expo") {
		t.Fatalf("expected expo provider to be configured")
	}

	s2 := &NotificationSettings{ExpoEnabled: false, ExpoPushToken: "ExponentPushToken[test-channel-token]"}
	if dispatcher.isProviderConfigured(s2, "expo") {
		t.Fatalf("expected disabled expo provider not to be configured")
	}

	s3 := &NotificationSettings{ExpoEnabled: true, ExpoPushToken: ""}
	if dispatcher.isProviderConfigured(s3, "expo") {
		t.Fatalf("expected empty token expo provider not to be configured")
	}
}
