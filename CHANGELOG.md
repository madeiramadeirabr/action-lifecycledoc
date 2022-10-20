# Changelog
Todas as mudanças notáveis deste este projeto serão documentadas neste arquivo.

O formato baseia-se em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
e este projeto segue [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

---

## [1.0.0] - 20-10-2022
### Acrescentado
- Adicionado GitHub Docker Container Action para permitir configurar workflows para criar/atualizar documentações de eventos
- Adicionado geração de arquivo de somas `sha256` dos binários criados no build do projeto

---

## [0.3.0] - 20-10-2022
### Acrescentado
- Adicionado opção `github-action-markdown` para formatar output do utilitário como markdown

## Alterado
- Alterado opção `github-action` para `github-action-json` da flag `outputFormat` para especificar melhor o seu resultado

---

## [0.2.0] - 19-10-2022
### Acrescentado
- Adicionado configuração `confluence_basic_auth` para especificar header de authenticação do Confluence diretamente

---

## [0.1.0] - 19-10-2022
### Acrescentado
- Tipo Scalar para o schema
- Tipo Array para o schema
- Tipo Object para o schema
- Tipo Reference para o schema
- Tipo ConfluencePage para o schema
- Tipo PublishedEvent para o schema
- Tipo ConsumedEvent para o schema
- Schema Resolver básico (apenas referências de tipos)
- YAML Schema Decoder
- JSONC (JSON com Comentários) Encoder
- API REST Client do Confluence básico (com suporte ao pacote `context`)
- Output da documentação renderizada para o Confluence

[Unreleased]: https://github.com/madeiramadeirabr/action-lifecycledoc/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/madeiramadeirabr/action-lifecycledoc/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/madeiramadeirabr/action-lifecycledoc/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/madeiramadeirabr/action-lifecycledoc/releases/tag/v0.1.0