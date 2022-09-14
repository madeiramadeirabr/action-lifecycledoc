# Lifecycledoc

## Descrição 

Utilitário para exportar a definição do lifecycledoc versionada no projeto como páginadas do Confluence.

## Contexto de negócio

Sistemas que contém a documentação de eventos disparados e consumidos versionados juntos com o projeto.

## Squad Owner

partnertools

## Get started

Para rodar o script você precisa configurar suas credenciais do confluence em `~/.lifecycledoc/config.yaml`. Com a seguinte estrutura:

```
confluence_api_key: <TOKEN>
confluence_email: <SEU EMAIL>
confluence_host: https://madeiramadeira.atlassian.net
```

Para roda o script segue o exemplo a baixo:
```
.lifecycledoc <path do openapi>
```

## Exit codes

* `0` - Sucesso
* `> 0` - Error

## Padrão de branchs

* Feature - feature/xxxxx
* Bugfix - bugfix/xxxxx

## Padrão de commmits

Resuma de forma sucinta o que foi adicionado, removido ou refatorado