package confluence_test

import (
	"strings"
	"testing"

	"github.com/madeiramadeirabr/action-lifecycledoc/internal/output/confluence"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/parser/yaml"
)

func TestShouldWriteTypes(t *testing.T) {
	input := strings.NewReader(`
version: "1.0"
name: super-cool-service

confluence:
  pages:
    - spaceKey: "SPACEKEY"
      ancestorId: "123456789"
      title: Titulo

events:
  published:
    CAKE_BURNED:
      visibility: public
      module: cooker
      description: Evento disparado quando o bolo é queimado ;-;

      attributes:
        type: object # object|array|string|integer|float|bool
        nullable: true
        properties:
          cake:
            $ref: '#/types/Cake'
          guilty:
            type: array
            description: Usuários que receberam a culpa
            items:
              type: object
              properties:
                id:
                  type: string
                  description: ID do usuário
                  value: 41af6672-5b3a-4d5c-9be1-7c93dc1614e1
                name:
                  type: string
                  description: Nome do usuário
                  value: Fulano
      
      entities:
        type: object
        properties:
          cakeId:
            type: string
            value: "12354"
  
  consumed:
    CAKE_PURCHASED:
      description: Usado para inciar o processo de fazer o bolo

types:
  CakeShape:
    description: Enum dos formatos de bolo suportado
    type: string
    enum:
      - squad
      - circle
    value: circle

  CakeFlaviourEnum:
    description: Enum dos sabores possíveis do bolo
    type: string
    enum:
      - chocolate
      - banana
      - morango
      - abacaxi
    value: abacaxi
  
  CakeFlaviours:
    type: array
    items:
      $ref: '#/types/CakeFlaviourEnum'

  Cake:
    description: Representa um bolo
    type: object
    properties:
      id:
        type: string
        value: "12354"
        description: O ID do bolo
      flaviours:
        description: Nhami Nhami 😋
        $ref: '#/types/CakeFlaviours'
      shape:
        $ref: '#/types/CakeShape'
      layers:
        type: integer
        format: uint8
        description: Quantidade de camadas do bolo
        value: 5`)

	schemaResolver := schema.NewBasicResolver()
	decoder := yaml.NewDecoder()

	if err := decoder.Decode(input, schemaResolver); err != nil {
		t.Fatal(err)
	}

	t.Run("asset types", func(t *testing.T) {
		// @todo: improve this assertion
		expected := `circle|||abacaxi|||["abacaxi"]|||{"id": "12354", // string: O ID do bolo"flaviours": ["abacaxi" // CakeFlaviourEnum: Enum dos sabores possíveis do bolo], // CakeFlaviours: Nhami Nhami 😋"shape": "circle", // CakeShape: Enum dos formatos de bolo suportado"layers": 5 // integer(uint8): Quantidade de camadas do bolo}|||`

		assertTempleWriterOutput(t, confluence.TemplateRetriverFunc(newTypesTemplateMock), schemaResolver, expected)
	})

	t.Run("asset published events", func(t *testing.T) {
		// @todo: improve this assertion
		expected := `{"attributes": {"cake": {"id": "12354", // string: O ID do bolo"flaviours": ["abacaxi" // CakeFlaviourEnum: Enum dos sabores possíveis do bolo], // CakeFlaviours: Nhami Nhami 😋"shape": "circle", // CakeShape: Enum dos formatos de bolo suportado"layers": 5 // integer(uint8): Quantidade de camadas do bolo}, // Cake: Representa um bolo"guilty": [{"id": "41af6672-5b3a-4d5c-9be1-7c93dc1614e1", // string: ID do usuário"name": "Fulano" // string: Nome do usuário}] // array: Usuários que receberam a culpa},"entities": {"cakeId": "12354"}}|||`

		assertTempleWriterOutput(t, confluence.TemplateRetriverFunc(newPublishedEventsTemplateMock), schemaResolver, expected)
	})
}

func newTypesTemplateMock() string {
	return "{{range .Types}}{{.Example}}|||{{end}}"
}

func newPublishedEventsTemplateMock() string {
	return "{{range .PublishedEvents}}{{.Example}}|||{{end}}"
}

func assertStringWithNewLinesAndIdentation(t *testing.T, expected, value string) {
	t.Helper()

	expected = removeNewLineAndIdentation(expected)
	value = removeNewLineAndIdentation(value)

	if value != expected {
		t.Errorf("expected '%s', received '%s'", expected, value)
	}
}

func removeNewLineAndIdentation(s string) string {
	s = strings.ReplaceAll(s, string('\n'), "")
	return strings.ReplaceAll(s, string('\t'), "")
}

func assertTempleWriterOutput(
	t *testing.T,
	templateRetriver confluence.TemplateRetriver,
	schemaResolver schema.Resolver,
	expected string,
) {
	t.Helper()

	templateWriter := confluence.NewTemplateWriter(templateRetriver)
	writerSpy := &strings.Builder{}

	if err := templateWriter.Write(writerSpy, schemaResolver); err != nil {
		t.Fatal(err)
	}

	assertStringWithNewLinesAndIdentation(t, expected, writerSpy.String())
}
