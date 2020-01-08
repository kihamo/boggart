package xmeye

import (
	"context"
	"errors"
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

func (c *Client) UserCreate(ctx context.Context, user User) (err error) {
	if user.Name == "" {
		return errors.New("user name is empty")
	}

	if user.Password == "" {
		return errors.New("user password is empty")
	}

	user.Password = HashPassword([]byte(user.Password))
	if user.AuthorityList == nil {
		user.AuthorityList = make([]string, 0)
	}

	_, err = c.Call(ctx, CmdUserCreateRequest, map[string]interface{}{
		"Name":      "User",
		"SessionID": c.connection.SessionIDAsString(),
		"User":      user,
	})

	return err
}

func (c *Client) UserUpdate(ctx context.Context, name string, user User) (err error) {
	if name == "" {
		return errors.New("user name is empty")
	}

	user.Password = HashPassword([]byte(user.Password))
	if user.AuthorityList == nil {
		user.AuthorityList = make([]string, 0)
	}

	_, err = c.Call(ctx, CmdUserUpdateRequest, map[string]interface{}{
		"SessionID": c.connection.SessionIDAsString(),
		"UserName":  name,
		"User":      user,
	})

	return nil
}

func (c *Client) UserDelete(ctx context.Context, name string) (err error) {
	if name == "" {
		return errors.New("user name is empty")
	}

	_, err = c.Call(ctx, CmdUserDeleteRequest, map[string]interface{}{
		"Name":      name,
		"SessionID": c.connection.SessionIDAsString(),
	})

	return err
}

func (c *Client) UserChangePassword(ctx context.Context, name, oldPassword, newPassword string) (err error) {
	if name == "" {
		return errors.New("user name is empty")
	}

	if oldPassword == "" {
		return errors.New("user old password is empty")
	}

	if newPassword == "" {
		return errors.New("user new password is empty")
	}

	oldPassword = HashPassword([]byte(oldPassword))
	newPassword = HashPassword([]byte(newPassword))

	_, err = c.Call(ctx, CmdUserChangePasswordRequest, map[string]interface{}{
		"SessionID":   c.connection.SessionIDAsString(),
		"EncryptType": "MD5",
		"NewPassWord": newPassword,
		"PassWord":    oldPassword,
		"UserName":    name,
	})

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

func (c *Client) GroupCreate(ctx context.Context, group Group) (err error) {
	if group.AuthorityList == nil {
		group.AuthorityList = make([]string, 0)
	}

	_, err = c.Call(ctx, CmdGroupCreateRequest, map[string]interface{}{
		"Name":      "Group",
		"SessionID": c.connection.SessionIDAsString(),
		"Group":     group,
	})

	return err
}

func (c *Client) GroupUpdate(ctx context.Context, name string, group Group) (err error) {
	if group.AuthorityList == nil {
		group.AuthorityList = make([]string, 0)
	}

	_, err = c.Call(ctx, CmdGroupUpdateRequest, map[string]interface{}{
		"SessionID": c.connection.SessionIDAsString(),
		"GroupName": name,
		"Group":     group,
	})

	return nil
}

func (c *Client) GroupDelete(ctx context.Context, name string) (err error) {
	_, err = c.Call(ctx, CmdGroupDeleteRequest, map[string]interface{}{
		"Name":      name,
		"SessionID": c.connection.SessionIDAsString(),
	})

	return err
}
