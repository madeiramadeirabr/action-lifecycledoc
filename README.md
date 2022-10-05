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
confluence_host: https://madeiramadeira.atlassian.net
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

## Padrão de branchs

* Feature - feature/xxxxx
* Bugfix - bugfix/xxxxx

## Padrão de commmits

Esse projeto utiliza o padrão [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).