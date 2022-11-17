package protocstep

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"testing"
)

func Test_generateFile(t *testing.T) {
	err := ProtocStep(func(gen *protogen.Plugin) error {
		for _, file := range gen.Files {
			if file.Generate {
				t.Log(file)
			}
		}
		return nil
	})
	assert.NoError(t, err)
}
