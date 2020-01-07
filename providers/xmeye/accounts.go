package xmeye

import (
	"context"
)

func (c *Client) FullAuthorityList(ctx context.Context) ([]string, error) {
	var result struct {
		Response
		AuthorityList []string
	}

	err := c.CallWithResult(ctx, CmdFullAuthorityListRequest, nil, &result)
	return result.AuthorityList, err
}

func (c *Client) Users(ctx context.Context) ([]User, error) {
	var result struct {
		Response
		Users []User
	}

	err := c.CallWithResult(ctx, CmdUsersRequest, nil, &result)
	return result.Users, err
}

func (c *Client) UserCreate(ctx context.Context, user User) error {
	// TODO:

	user.Password = HashPassword([]byte(user.Password))

	return nil
}

func (c *Client) UserUpdate(ctx context.Context, name string, user User) error {
	// TODO:

	user.Password = HashPassword([]byte(user.Password))

	return nil
}

func (c *Client) UserDelete(ctx context.Context, name string) error {
	// TODO:
	return nil
}

func (c *Client) UserChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	// TODO:

	oldPassword = HashPassword([]byte(oldPassword))
	newPassword = HashPassword([]byte(newPassword))

	return nil
}

func (c *Client) Groups(ctx context.Context) ([]Group, error) {
	var result struct {
		Response
		Groups []Group
	}

	err := c.CallWithResult(ctx, CmdGroupsRequest, nil, &result)
	return result.Groups, err
}

func (c *Client) GroupCreate(ctx context.Context, group Group) error {
	// TODO:
	return nil
}

func (c *Client) GroupUpdate(ctx context.Context, name string, group Group) error {
	// TODO:
	return nil
}

func (c *Client) GroupDelete(ctx context.Context, name string) error {
	// TODO:
	return nil
}
