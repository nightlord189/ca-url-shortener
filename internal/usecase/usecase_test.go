package usecase_test

import (
	"context"
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"github.com/nightlord189/ca-url-shortener/internal/usecase/mock"
	"github.com/nightlord189/ca-url-shortener/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
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
		storageMock.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
			require.NotNil(t, user)
			assert.Equal(t, entity.User{Username: username, PasswordHash: passwordHash}, *user)
			return fmt.Errorf(errorText)
		})

		err := uc.Auth(ctx, username, password)
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), errorText)
	})

	t.Run("Create user: success", func(t *testing.T) {
		t.Parallel()

		username := "user_create_success"
		password := "password1"
		passwordHash := pkg.GetSHA256Hash(password)

		// to prevent intersections with previous test
		storageMock := mock.NewMockIStorage(mockCtrl)
		uc := usecase.New(storageMock, nil)

		storageMock.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
		storageMock.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
			require.NotNil(t, user)
			assert.Equal(t, entity.User{Username: username, PasswordHash: passwordHash}, *user)
			return nil
		})

		err := uc.Auth(ctx, username, password)
		require.Nil(t, err)
	})
}
