package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nightlord189/ca-url-shortener/internal/entity"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"github.com/nightlord189/ca-url-shortener/internal/usecase/mock"
	"github.com/nightlord189/ca-url-shortener/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockCtrl := gomock.NewController(t)

	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	storageMock := mock.NewMockIStorage(mockCtrl)

	uc := usecase.New(storageMock, nil)

	t.Run("Internal error", func(t *testing.T) {
		t.Parallel()

		username := "user_internal_error"
		errorText := "unknown error"

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, fmt.Errorf(errorText))

		err := uc.Auth(ctx, username, "password1")
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		t.Parallel()

		username := "user_invalid_credentials"
		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&entity.User{Username: username, PasswordHash: "hash1"}, nil)

		err := uc.Auth(ctx, username, "password1")
		require.Equal(t, usecase.ErrInvalidCredentials, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		username := "user_success"
		password := "password1"
		passwordHash := pkg.GetSHA256Hash(password)

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&entity.User{Username: username, PasswordHash: passwordHash}, nil)

		err := uc.Auth(ctx, username, password)
		require.Nil(t, err)
	})

	t.Run("Create user error", func(t *testing.T) {
		t.Parallel()

		username := "user_create_error"
		password := "password1"
		passwordHash := pkg.GetSHA256Hash(password)

		errorText := "unknown error"

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
		storageMock.EXPECT().CreateUser(ctx, gomock.Cond(compareUserMatcher(entity.User{Username: username, PasswordHash: passwordHash}))).Return(fmt.Errorf(errorText))

		err := uc.Auth(ctx, username, password)
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
	})

	t.Run("Create user: success", func(t *testing.T) {
		t.Parallel()

		username := "user_create_success"
		password := "password1"
		passwordHash := pkg.GetSHA256Hash(password)

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
		storageMock.EXPECT().CreateUser(ctx, gomock.Cond(compareUserMatcher(entity.User{Username: username, PasswordHash: passwordHash}))).Return(nil)

		err := uc.Auth(ctx, username, password)
		require.Nil(t, err)
	})
}

func TestPutLink(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockCtrl := gomock.NewController(t)

	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	storageMock := mock.NewMockIStorage(mockCtrl)
	cacheMock := mock.NewMockICache(mockCtrl)

	uc := usecase.New(storageMock, cacheMock)

	t.Run("Error: no user", func(t *testing.T) {
		t.Parallel()

		username := "user_no_user"
		errorText := "unknown error"

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, fmt.Errorf(errorText))

		link, err := uc.PutLink(ctx, username, "https://example.com")
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
		assert.Empty(t, link)
	})

	t.Run("Update user error", func(t *testing.T) {
		t.Parallel()

		username := "user_update_error"
		errorText := "unknown error"
		originalURL := "https://example.com"
		shortURL := pkg.GetFNVHash(username + originalURL)

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&entity.User{Username: username}, nil)
		expectedUser := entity.User{Username: username, Links: map[string]string{shortURL: originalURL}}
		storageMock.EXPECT().UpdateUserLinks(gomock.Any(), gomock.Cond(compareUserMatcher(expectedUser))).Return(fmt.Errorf(errorText))

		link, err := uc.PutLink(ctx, username, originalURL)
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
		assert.Empty(t, link)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		username := "user_success"
		originalURL := "https://example.com"
		shortURL := pkg.GetFNVHash(username + originalURL)

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&entity.User{Username: username}, nil)
		expectedUser := entity.User{Username: username, Links: map[string]string{shortURL: originalURL}}
		storageMock.EXPECT().UpdateUserLinks(gomock.Any(), gomock.Cond(compareUserMatcher(expectedUser))).Return(nil)
		cacheMock.EXPECT().PutLink(gomock.Any(), shortURL, originalURL).Return(nil)

		link, err := uc.PutLink(ctx, username, originalURL)
		require.Nil(t, err)
		assert.Equal(t, shortURL, link)
	})
}

func TestGetOriginalLink(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockCtrl := gomock.NewController(t)

	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	storageMock := mock.NewMockIStorage(mockCtrl)
	cacheMock := mock.NewMockICache(mockCtrl)

	uc := usecase.New(storageMock, cacheMock)

	t.Run("Success from cache", func(t *testing.T) {
		t.Parallel()

		shortURL := "short1"
		originalURL := "original1"
		cacheMock.EXPECT().GetLink(gomock.Any(), shortURL).Return(originalURL, nil)

		link, err := uc.GetOriginalLink(ctx, shortURL)
		assert.Equal(t, originalURL, link)
		assert.Nil(t, err)
	})

	t.Run("Success from storage", func(t *testing.T) {
		t.Parallel()

		shortURL := "short4"
		originalURL := "original4"
		cacheMock.EXPECT().GetLink(gomock.Any(), shortURL).Return("", fmt.Errorf("not found"))
		storageMock.EXPECT().GetLink(gomock.Any(), shortURL).Return(originalURL, nil)
		cacheMock.EXPECT().PutLink(gomock.Any(), shortURL, originalURL).Return(nil)

		link, err := uc.GetOriginalLink(ctx, shortURL)
		assert.Equal(t, originalURL, link)
		require.Nil(t, err)
	})

	t.Run("Internal error", func(t *testing.T) {
		t.Parallel()

		shortURL := "short2"
		errorText := "unknown error"
		cacheMock.EXPECT().GetLink(gomock.Any(), shortURL).Return("", fmt.Errorf("not found"))
		storageMock.EXPECT().GetLink(gomock.Any(), shortURL).Return("", fmt.Errorf(errorText))

		link, err := uc.GetOriginalLink(ctx, shortURL)
		assert.Empty(t, link)
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		shortURL := "short3"
		cacheMock.EXPECT().GetLink(gomock.Any(), shortURL).Return("", fmt.Errorf("not found"))
		storageMock.EXPECT().GetLink(gomock.Any(), shortURL).Return("", usecase.ErrNotFound)

		link, err := uc.GetOriginalLink(ctx, shortURL)
		assert.Empty(t, link)
		require.NotNil(t, err)
		assert.ErrorIs(t, err, usecase.ErrNotFound)
	})
}

func compareUserMatcher(expected entity.User) func(x any) bool {
	return func(x any) bool {
		fmt.Println("MATCHER")
		userPointer := x.(*entity.User)
		user2 := *userPointer
		if user2.Username != expected.Username || user2.PasswordHash != expected.PasswordHash {
			return false
		}
		if len(user2.Links) != len(expected.Links) {
			return false
		}
		for key := range expected.Links {
			if expected.Links[key] != user2.Links[key] {
				return false
			}
		}
		return true
	}
}
