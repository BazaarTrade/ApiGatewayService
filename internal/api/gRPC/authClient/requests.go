package aClient

import (
	"context"

	"github.com/BazaarTrade/AuthProtoGen/pbA"
)

func (c *Client) Register(email, password string) error {
	req := &pbA.RegisterRequest{
		Email:    email,
		Password: password,
	}

	_, err := c.client.Register(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Login(email, password string) (int, string, error) {
	req := &pbA.LoginRequest{
		Email:    email,
		Password: password,
	}

	res, err := c.client.Login(context.Background(), req)
	if err != nil {
		return 0, "", err
	}

	return int(res.UserID), res.RefreshToken, nil
}

func (c *Client) Logout(refreshToken string) error {
	req := &pbA.LogoutRequest{
		RefreshToken: refreshToken,
	}

	_, err := c.client.Logout(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) IsRefreshTokenValid(userID int, refreshToken string) (bool, string, error) {
	req := &pbA.IsRefreshTokenValidRequest{
		UserID:       int64(userID),
		RefreshToken: refreshToken,
	}

	res, err := c.client.IsRefreshTokenValid(context.Background(), req)
	if err != nil {
		return false, "", err
	}

	return res.IsValid, res.NewRefreshToken, nil
}

func (c *Client) ChangePassword(email, oldPassword, newPassword string, refreshToken string) error {
	req := &pbA.ChangePasswordRequest{
		Email:        email,
		OldPassword:  oldPassword,
		NewPassword:  newPassword,
		RefreshToken: refreshToken,
	}

	_, err := c.client.ChangePassword(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}
