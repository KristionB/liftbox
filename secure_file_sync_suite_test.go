package secure_file_sync_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSecureFileSync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SecureFileSync Suite")
}
