# <Nome do projeto>

### Descrição

### Contexto de negócio

### Squad

### Requisitos (Se aplicavel)

### Configurações para uso (Se aplicavel)

### Exemplos de uso (Se aplicavel)

### Exit codes (Status Code, Result)

### Monitoria e Logs (Caso seja aplicação - obrigatório)
<!-- Descrição simples da regras de negócio que estão sendo logadas -->
<!-- Descrição das monitorias + link para elas -->
### Alertas (Caso seja aplicação - obrigatório)
<!-- Descrição de quando deve ser alertado -->
<!-- Listagem de todos os alertas configurados -->
<!-- Canais que são notificados os alertas -->

### Checklist de Engenharia

*Arquitetura*

- [ ] A aplicação nasceu com base em uma arquitetura de referência?
- [ ] A aplicação possui CI/CD seguindo as práticas descritas no seguinte documento [CI/CD ](https://madeiramadeira.atlassian.net/wiki/spaces/SRE/pages/2560131648)
- [ ] O banco de dados que a aplicação utiliza, é consumido apenas por essa aplicação?
- [ ] A chamadas síncronas que faço em outros serviços, utiliza o padrão Retry? [Padrão Retry](https://docs.microsoft.com/pt-br/azure/architecture/patterns/retry)
- [ ] As APIs são versionadas e toda mudança de contrato gera mudança na versão?
- [ ] Possui HealthChecks para validar todas as dependências externas? <!-- add documentação de explicação -->
- [ ] Seguem as práticas dos 12 fatores? <!-- add link para explicação -->
- [ ] Os eventos seguem as práticas descritas no documento ? [Eventos na MadeiraMadeira](https://madeiramadeira.atlassian.net/wiki/spaces/CAR/pages/2141291341)
- [ ] As comunicações assíncronas entre serviços, ocorre via evento?
- [ ] É publicado seguindo o fluxo de Gestão de Mudanças (GMUD)?

*Observabilidade*
- [ ] Os logs seguem o padrão RFC da MadeiraMadeira ? [RFC Logs MadeiraMadeira](https://madeiramadeira.atlassian.net/wiki/spaces/CAR/pages/2317942893)
- [ ] Os logs da aplicação são escritos no stdout e são enviados para New Relic ? 
- [ ] As aplicações estão com monitoramento ativo na New Relic ?

*Infraestrutura*
- [ ] A infraestrutura está automatizada com Terraform ?
- [ ] A comunicação entre as aplicações ocorre apenas entre mesmo ambiente ?


*Segurança*
- [ ] As chaves secretas que a aplicação utiliza são consumidas apenas via Secret Manager ?
- [ ] Os parâmetros de ambiente estão configurados usando o Parameter Store ?
- [ ] A exposição é feita pelo Kong ? **quando necessário expor**
    
*Testes*
- [ ] Documentação (histórias,épicos e etc) disponíveis e devidamente detalhados (utilizando o formato BDD - Behavior driven development ou User Story) ?
- [ ] Documentação (histórias, épicos e etc) apresentam impacto na jornada ?
- [ ] Todos os itens (épicos,histórias, tarefas e etc) possuem regras e critérios de aceite ?
- [ ] Todos os itens (épicos,histórias, tarefas e etc) apresentam o registro de execução de testes unitários apresentados nas tarefas ?
- [ ] Todos os itens (épicos,histórias, tarefas e etc) apresentam o registro de execução de testes de integração e/ou caminho feliz, apresentados nas tarefas ?
- [ ] As aplicações possuem teste unitário ?
- [ ] As APIs possuem teste de contrato?  
- [ ] As aplicações estão usando mocks ?
    
*Processo*
- [ ] Está usando as práticas de source control ? [Source Control](https://madeiramadeira.atlassian.net/wiki/spaces/SRE/pages/2482995227)
- [ ] O repositório de código está configurado de forma a promover práticas inner source ? 
- [ ] As histórias do time usam formato de BDD ou user story ?
- [ ] Os épicos do time apresentam o problema pro time quebrar em ações (histórias) ?
