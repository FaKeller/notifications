package notificationtypes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NotificationTypes Suite")
}
