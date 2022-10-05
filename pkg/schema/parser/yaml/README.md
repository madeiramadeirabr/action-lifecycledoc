# Definição de eventos no formato YAML
A definição dos eventos utiliza um subconjunto de keywords do JSON Schema com algumas keywords específicas para tentar reduzir a curva de aprendizagem.

## Estrutura
Toda definição de eventos deve contem a seguinte estrutura base:
```yaml
version: "1.0"
name: project-name

confluence:
  pages:

events:
  published:

  consumed:

types:
```

### Campos obrigatórios

#### version (string)
Especifica a versão do schema utilizada na documentação dos eventos.

Versões suportadas
 - 1.0

#### name (string)
Especifica o nome do Projeto/Sistema que dispara esses eventos.

#### confluence (map)
Especifica as configurações para o Confluence, por meio das seguintes propriedades:

| Propriedade | Tipo | Descrição |
| ----------- | ---- | --------- |
| `pages` | `ConfluencePage` array  | Especifica as páginas que devem ser geradas/atualizadas |

#### events (map)
Especifica os eventos publicados e consumidos pelo Projeto/Sistema, por meio das seguintes propriedades:

| Propriedade | Tipo | Descrição |
| ----------- | ---- | --------- |
| `published` | `PublishedEvent` map | Especifca os eventos publicados. A chave de cada item do mapa deve ser o nome do evento disparado. |
| `consumed` | `ConsumedEvent` map | Especifica os eventos consumidos. A chave de cada item do mapa deve ser o nome do evento consumido. |

#### types (map)
Especifica os tipos dados usados nos eventos publicados por de `TypeObject`s. Os tipos definidos existem para facilitar a definição dos eventos publicados, por meio de referencia ao mesmos (`$ref` keyword), assim evitando repetir a mesma definição entre eventos.

A chave de cada item do mapa deve ser o identificador (nome) de cada tipo.

Cada tipo declarado tem uma caminho de referência no schema, com o seguinte padrão:
`#/types/TypeName`, onde o `#/types/` é uma constante e o `TypeName` é o identificador do tipo. O caminho deve ser resolvido automaticamente `schema.Resolver`.

## Schema

### ConfluencePage
Possui as seguintes propriedades:
| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `spaceKey` | string | Sim | Especifica o código do seu espaço no Confluence. |
| `ancestorId` | string | Sim | Especifica qual é a pagina "pai". |
| `title` | string | Não | Sobrescreve o título padrão da página: "Life Cycle Events: {project name}" |

### TypeObject
Define um tipo, sua declaração depende do seu tipo. Há 2 grupos de propriedades na definição dessa tipo: **propriedades comuns** e **propriedades específicas**

Esse tipo empresta alguns definições do JSON Schema.

#### Properidades comuns
| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `type`  | string | Sim | Indica o tipo da definição. O valor do mesmo define as propriedades específicas |
| `description` | string | Não | Adicionar uma descrição para a definição. |
| `nullable` | boolean | Não | Indicia se a definição permite valors nulos |
| `$ref` | string | Não | Referencia outro tipo definido. Ao usar essar keyword o tipo será ignorado, uma vez que o tipo dessa definição é o tipo referenciado. |

##### Tipos suportados
| Keyword | Descrição |
| ------- | --------- |
| `integer` | Números interos com sinal |
| `number` | Números decimais com sinal |
| `string` | Texto |
| `boolean` | Verdairo ou falso |
| `array` | Array de `TypeObject`s |
| `object` | Um "mapa" de outros `TypeObject`s |

Os tipos `integer`, `number`, `string` e `boolean` são chamados de `Scalar` types.

#### Properidades específicas

##### Tipo Scalar
| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `value` | Scalar | Sim | Especifica um valor de exemplo para a definição |
| `enum` | `Scalar` array | Não | Especifica valores possíveis para a definição |
| `format` | string | Não | Texto livre que especifica/delimite os valores possíveis |

##### array
| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `items` | `TypeObject` | Sim | Especifica o tipo dos items do array |

##### object
A chave de cada item do mapa deve ser o identificador da propriedade e o valor deve ser do tipo `TypeObject`.

### PublishedEvent
Define um evento publicado, possui as seguintes propriedades:

| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `visibility` | string[public,protected,private] | Sim | Especifica a visibilidade do evento. |
| `module` | string | Não | Especifica o módulo do sistema que gerou do evento. |
| `description` | string | Não | Descrição do evento. |
| `attributes` | `TypeObject` | Sim | Especifica as propriedades do evento. |
| `entities` | `TypeObject` | Sim | Especifica as entidades do presentes no evento. |

### ConsumedEvent
Define um evento consumdo, possui as seguintes propriedades:

| Keyword | Tipo | Obrigatório | Descrição |
| ------- | ---- | ----------- | --------- |
| `description` | string | Sim | Descre o uso do determinado evento pela aplicação. |

# Exemplo
```yaml
version: "1.0"
name: super-cool-service

confluence:
  pages:
    - spaceKey: "SPACE_KEY"
      ancestorId: "SOME_ID"

events:
  published:
    CAKE_BURNED:
      visibility: public
      module: cooker
      description: Evento disparado quando o bolo é queimado ;-;

      attributes:
        type: object
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
        value: 5
```