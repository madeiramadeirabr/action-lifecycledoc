# Lifecycledoc

## Descrição 

Utilitário para exportar documentação de eventos descritos em YAML como páginas do Confluence.

## Contexto de negócio

Sistemas que contém a documentação de eventos disparados e consumidos versionados juntos com o projeto.

## Squad Owner

partnertools

## Get started

Para executar o programa você precisa especificar seu Personal Access Token (PAT) do Confluence no `~/.lifecycledoc/config.yaml`. Ao executar esse utilitário pela primeira vez o mesmo criará esse arquivo em sua pasta do usuário com a seguinte estrutura:

```yaml
confluence_api_key: <TOKEN>
confluence_email: <SEU EMAIL>
confluence_host: https://acme-fake-company.atlassian.net
```

Para exportar a documentação dos eventos para o confluence basta executar o binário do programa passando o path do YAML dos eventos:
```
lifecycledoc /some/path/lifecycle.yaml
```

As demais opções do utilitário podem ser recuperados especificando a flag `-h` durante a execução do comando:
```
lifecycledoc -h
```

A especificação da sintaxe do YAML dos eventos pode ser na seguinte [página](pkg/schema/parser/yaml)

## Exit codes

* `0` - Sucesso
* `> 0` - Error

## GitHub Action
Disponibilizamos uma [action](action.yaml) para executar o `lifecycledoc` no workflow do próprio repositório para permitir automatizar a criação e atualização das documentações de eventos disparados.

### Uso
```yaml
- uses: madeiramadeirabr/action-lifecycledoc@v1
  with:
    # Especifica o caminho do arquivo que contém a definição dos eventos da aplicação
    lifecycle-file: docs/lifecycle.yaml

    # Especifica o formato do output da action.
    #
    # Formatos suportados
    # - 'github-action-json': Retorna os links das documentações como um array disponível no output `links`
    # - 'github-action-markdown': Retorna os links das documentações como um Markdown no resumo do Job do Workflow
    # - 'cli': Retorna os links gerados como logs de execução do utilitário no resumo do Job do Workflow'
    #
    # Opcional
    output-format: github-action-markdown

    # Especifica um prefixo para os títulos das documentações
    #
    # Opcional
    title-prefix: ''

    # Especifica o host do Confluence de destino
    confluence-host: https://acme-fake-company.atlassian.net

    # Especifica o token de autenticação do Confluence no formato Basic Auth
    confluence-basic-auth: ${{ secrets.CONFLUENCE_BASIC_AUTH_TOKEN }}

    # Especfica o e-mail do usuário do Confluence para autenticação com Personal Access Tokens (PAT)
    # https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
    #
    # Opcional quando o input `confluence-basic-auth` for especificado
    confluence-email: ${{ secrets.CONFLUENCE_USER_EMAIL }}

    # Especfica o Personal Access Tokens (PAT) do usuário do Confluence
    # https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
    #
    # Opcional quando o input `confluence-basic-auth` for especificado
    confluence-api-key: ${{ secrets.CONFLUENCE_PERSONAL_ACCESS_TOKEN }}
```

## Padrão de branchs

* Feature - feature/xxxxx
* Bugfix - bugfix/xxxxx

## Padrão de commmits

Esse projeto utiliza o padrão [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).