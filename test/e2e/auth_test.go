package e2e

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/strick-j/scimfe/pkg/scimfe"
)

func TestAuth_Register(t *testing.T) {
	cases := []struct {
		label   string
		req     scimfe.RegisterRequest
		wantErr string
	}{
		{
			label:   "empty payload",
			wantErr: "400 Bad Request: invalid request payload",
		},
		{
			label:   "invalid email",
			wantErr: "400 Bad Request: invalid request payload",
			req: scimfe.RegisterRequest{
				Email:    "--",
				Name:     "john",
				Password: "123456",
			},
		},
		{
			label:   "invalid name",
			wantErr: "400 Bad Request: invalid request payload",
			req: scimfe.RegisterRequest{
				Email:    "u111@example.com",
				Name:     "",
				Password: "123456",
			},
		},
		{
			label:   "invalid password",
			wantErr: "400 Bad Request: invalid request payload",
			req: scimfe.RegisterRequest{
				Email:    "u111@example.com",
				Name:     "joey",
				Password: "",
			},
		},
		{
			label: "valid creds",
			req: scimfe.RegisterRequest{
				Email:    "u111@example.com",
				Name:     "joey",
				Password: "123456",
			},
		},
		{
			label:   "duplicate registration",
			wantErr: "400 Bad Request: record already exists",
			req: scimfe.RegisterRequest{
				Email:    "u111@example.com",
				Name:     "marko",
				Password: "123456",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			sess, err := Client.Register(c.req)
			if c.wantErr != "" {
				shouldContainError(t, err, c.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, sess.User.Name, c.req.Name)
			require.Equal(t, sess.User.Email, c.req.Email)
		})
	}
}

func TestAuth_Login(t *testing.T) {
	sess, err := Client.Register(scimfe.RegisterRequest{
		Email:    "testauthlogin@mail.com",
		Name:     "testauthlogin",
		Password: "123456",
	})
	require.NoError(t, err, "failed to create a user for test case")

	cases := map[string]struct {
		wantErr  string
		creds    scimfe.Credentials
		wantUser scimfe.User
	}{
		"empty request": {
			wantErr: "400 Bad Request: invalid request payload",
		},
		"empty email": {
			creds:   scimfe.Credentials{Password: "12345"},
			wantErr: "400 Bad Request: invalid request payload",
		},
		"empty password": {
			creds:   scimfe.Credentials{Email: "testauthlogin@mail.com"},
			wantErr: "400 Bad Request: invalid request payload",
		},
		"invalid username": {
			creds:   scimfe.Credentials{Email: "badusername@mail.com", Password: "123456"},
			wantErr: "400 Bad Request: invalid username or password",
		},
		"invalid password": {
			creds:   scimfe.Credentials{Email: "testauthlogin@mail.com", Password: "badpassword"},
			wantErr: "400 Bad Request: invalid username or password",
		},
		"valid login": {
			creds:    scimfe.Credentials{Email: "testauthlogin@mail.com", Password: "123456"},
			wantUser: sess.User,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			rsp, err := Client.Login(c.creds)
			if c.wantErr != "" {
				shouldContainError(t, err, c.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, rsp)
			require.Equal(t, c.wantUser, rsp.User)
		})
	}
}

func TestAuth_Session(t *testing.T) {
	sess, err := Client.Register(scimfe.RegisterRequest{
		Email:    "testgetsession@mail.com",
		Name:     "testgetsession",
		Password: "123456",
	})
	require.NoError(t, err, "failed to create a user for test case")

	cases := map[string]struct {
		want        scimfe.SessionInfo
		wantErr     string
		token       scimfe.Token
		onBeforeRun func(t *testing.T, token *scimfe.Token)
	}{
		"empty token": {
			wantErr: "401 Unauthorized: authorization required",
		},
		"invalid token": {
			wantErr: "401 Unauthorized: authorization required",
			token:   scimfe.Token(uuid.New().String()),
		},
		"valid token": {
			token: sess.Token,
			want:  sess.Session,
		},
		"expired token": {
			token:   sess.Token,
			wantErr: "401 Unauthorized: authorization required",
			onBeforeRun: func(t *testing.T, token *scimfe.Token) {
				sess, err := Client.Login(scimfe.Credentials{
					Email:    "testgetsession@mail.com",
					Password: "123456",
				})
				require.NoError(t, err, "failed to create a user for test case")
				*token = sess.Token

				require.NoError(t, Client.Logout(sess.Token), "failed to logout from test account")
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			token := c.token
			if c.onBeforeRun != nil {
				c.onBeforeRun(t, &token)
			}

			got, err := Client.Session(token)
			if c.wantErr != "" {
				shouldContainError(t, err, c.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, c.want, *got)
		})
	}
}

func TestAuth_Logout(t *testing.T) {
	sess, err := Client.Register(scimfe.RegisterRequest{
		Email:    "testlogout@mail.com",
		Name:     "testlogout",
		Password: "123456",
	})
	require.NoError(t, err, "failed to create a user for test case")

	cases := map[string]struct {
		wantErr string
		token   scimfe.Token
	}{
		"empty token": {
			wantErr: "401 Unauthorized: authorization required",
		},
		"invalid token": {
			wantErr: "401 Unauthorized: authorization required",
			token:   scimfe.Token(uuid.New().String()),
		},
		"valid token": {
			token: sess.Token,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			err := Client.Logout(c.token)
			if c.wantErr != "" {
				shouldContainError(t, err, c.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
