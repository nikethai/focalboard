// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
package app

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/nikethai/focalboard/server/auth"
	"github.com/nikethai/focalboard/server/services/config"
	"github.com/nikethai/focalboard/server/services/metrics"
	"github.com/nikethai/focalboard/server/services/permissions/mmpermissions"
	mmpermissionsMocks "github.com/nikethai/focalboard/server/services/permissions/mmpermissions/mocks"
	permissionsMocks "github.com/nikethai/focalboard/server/services/permissions/mocks"
	"github.com/nikethai/focalboard/server/services/store/mockstore"
	"github.com/nikethai/focalboard/server/services/webhook"
	"github.com/nikethai/focalboard/server/ws"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore/mocks"
)

type TestHelper struct {
	App          *App
	Store        *mockstore.MockStore
	FilesBackend *mocks.FileBackend
	logger       mlog.LoggerIFace
	API          *mmpermissionsMocks.MockAPI
}

func SetupTestHelper(t *testing.T) (*TestHelper, func()) {
	ctrl := gomock.NewController(t)
	cfg := config.Configuration{}
	store := mockstore.NewMockStore(ctrl)
	filesBackend := &mocks.FileBackend{}
	auth := auth.New(&cfg, store, nil)
	logger, _ := mlog.NewLogger()
	sessionToken := "TESTTOKEN"
	wsserver := ws.NewServer(auth, sessionToken, false, logger, store)
	webhook := webhook.NewClient(&cfg, logger)
	metricsService := metrics.NewMetrics(metrics.InstanceInfo{})

	mockStore := permissionsMocks.NewMockStore(ctrl)
	mockAPI := mmpermissionsMocks.NewMockAPI(ctrl)
	permissions := mmpermissions.New(mockStore, mockAPI, mlog.CreateConsoleTestLogger(t))

	appServices := Services{
		Auth:             auth,
		Store:            store,
		FilesBackend:     filesBackend,
		Webhook:          webhook,
		Metrics:          metricsService,
		Logger:           logger,
		SkipTemplateInit: true,
		Permissions:      permissions,
	}
	app2 := New(&cfg, wsserver, appServices)

	tearDown := func() {
		app2.Shutdown()
		if logger != nil {
			_ = logger.Shutdown()
		}
	}

	return &TestHelper{
		App:          app2,
		Store:        store,
		FilesBackend: filesBackend,
		logger:       logger,
		API:          mockAPI,
	}, tearDown
}
