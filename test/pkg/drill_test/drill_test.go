package drill_test

import (
	"badger-api/pkg/drill"
	"testing"

	goblin "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestDrill(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	g.Describe("DrillsStructs", func() {
		var drill = drill.Drill{Id: "3", Name: "cones", Description: "sprinting drill"}

		// Passing Test
		g.It("Should retrieve id", func() {
			Expect(drill.GetId()).To(Equal("3"))
		})
		// Failing Test
		g.It("Should retreive name", func() {
			Expect(drill.GetName()).To(Equal("cones"))
		})
		g.It("Should retrieve description", func() {
			Expect(drill.GetDescription()).To(Equal("sprinting drill"))
		})
	})
}
