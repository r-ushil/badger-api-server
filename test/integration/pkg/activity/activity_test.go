package drill_test

import (
	"testing"

	goblin "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestActivity(t *testing.T) {

	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("ActivityServer", func() {
		g.It("Should be able to retrieve all activities")

		g.It("Should be able to retrieve an activity from a valid activity id")

		g.It("Should have an error when retrieving an activity from an invalid activity id")

		g.It("Should be able to retrieve all activies within a certain time frame")
	})

}
