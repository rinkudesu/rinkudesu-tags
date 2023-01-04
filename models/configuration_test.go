package models

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfiguration_AllDefaults(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("TAGS_DB", "")

	config := NewConfiguration()

	assert.Equal(t, "/api", config.BasePath)
	assert.Equal(t, logrus.InfoLevel, config.LogLevel)
	assert.Nil(t, config.TrustedProxies)
	assert.Equal(t, "localhost:5000", config.ListenAddress)
	assert.Equal(t, "rinkudesu", config.SsoClientId)
}

func TestNewConfiguration_CustomValues(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("TAGS_BASE-PATH", "/test")
	_ = os.Setenv("TAGS_LOG-LEVEL", "warn")
	_ = os.Setenv("TAGS_DB", "postgres://postgres:postgres@localhost:5432/postgres")
	_ = os.Setenv("TAGS_PROXY", "192.168.0.1,10.0.0.1,10.0.0.2")
	_ = os.Setenv("TAGS_ADDRESS", "192.168.0.1:80")
	_ = os.Setenv("TAGS_AUTHORITY", "http://localhost")
	_ = os.Setenv("TAGS_CLIENTID", "not-rinkudesu")

	config := NewConfiguration()

	assert.Equal(t, "/test", config.BasePath)
	assert.Equal(t, logrus.WarnLevel, config.LogLevel)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres", config.DbConnection)
	assert.Equal(t, []string{"192.168.0.1", "10.0.0.1", "10.0.0.2"}, config.TrustedProxies)
	assert.Equal(t, "192.168.0.1:80", config.ListenAddress)
	assert.Equal(t, "http://localhost", config.SsoAuthority)
	assert.Equal(t, "not-rinkudesu", config.SsoClientId)
}
