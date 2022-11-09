package drill_test

import (
	"testing"

	goblin "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestActivity(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
}
