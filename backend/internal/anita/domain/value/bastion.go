package value

import "github.com/ictsc/ictsc-outlands/backend/pkg/errors"

// Bastion 踏み台サーバー
type Bastion struct {
	user     string
	password string
	host     string
	port     int
}

// NewBastion 踏み台サーバーを生成する
func NewBastion(user, password, host string, port int) (Bastion, error) {
	if len(user) < 1 || len(user) > 20 {
		return Bastion{}, errors.New(errors.ErrBadArgument, "Invalid user")
	}

	if len(password) < 1 || len(password) > 20 {
		return Bastion{}, errors.New(errors.ErrBadArgument, "Invalid password")
	}

	if len(host) < 1 || len(host) > 100 {
		return Bastion{}, errors.New(errors.ErrBadArgument, "Invalid host")
	}

	if port < 0 || port > 65535 {
		return Bastion{}, errors.New(errors.ErrBadArgument, "Invalid port")
	}

	return Bastion{
		user:     user,
		password: password,
		host:     host,
		port:     port,
	}, nil
}

// Equals 踏み台サーバーが等しいか
func (b Bastion) Equals(other Bastion) bool {
	return b.user == other.user &&
		b.password == other.password &&
		b.host == other.host &&
		b.port == other.port
}

// User ユーザー名を取得
func (b Bastion) User() string {
	return b.user
}

// Password パスワードを取得
func (b Bastion) Password() string {
	return b.password
}

// Host ホストを取得
func (b Bastion) Host() string {
	return b.host
}

// Port ポートを取得
func (b Bastion) Port() int {
	return b.port
}
