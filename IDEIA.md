Este arquivo contém o prompt e as diretrizes principais para a criação automática do projeto `gotodolist` usando o modo agente.

---
**Instruções iniciais para o Agente:**
Primeiro, faça um breve plano de implementação para si mesmo e, em seguida, prossiga com a implementação de todos os itens detalhados abaixo. Siga estritamente a estrutura de monorepo especificada e as diretrizes de segurança. Permaneça dentro do espaço de trabalho atual e modifique apenas o projeto atual

## gotodolist
O `gotodolist` é um aplicativo simples de gerenciamento de tarefas (to-do list). A ideia principal do projeto é servir como uma ferramenta didática para ensinar o deployment de uma aplicação web que acessa um banco de dados. O projeto é inspirado no modelo do `https://github.com/docker/getting-started-todo-app`. O aplicativo não terá autenticação, focando inteiramente no fluxo de dados de Listar, Criar, Atualizar e Deletar tarefas.

# Implantação
O projeto foi pensado para ser implantado de duas formas distintas, dependendo do perfil da turma de alunos:
1. **VMs na Oracle Cloud Infrastructure (OCI):** A aplicação (backend/monolito) rodará em uma VM em uma subrede pública, enquanto o banco de dados rodará separadamente em uma VM em uma subrede privada.
2. **Ambiente com Contêineres:** A aplicação e o banco de dados rodarão em uma única VM, mas orquestrados utilizando Docker e containers.

A aplicação estará acessível diretamente por IP (ex: `http://ip-da-vm:8080`), mas deve ser projetada de forma que possa trabalhar atrás de um proxy reverso (como Nginx) no futuro.

# Tecnologias utilizadas
- **Backend Monolito:** Go (Golang) renderizando templates HTML no servidor (server-side rendering). Utilize as versões mais novas e estáveis das bibliotecas e frameworks Go para web (ex: `net/http`, `html/template`).
- **Backend Desacoplado (API Rest):** Go (Golang). A API deve seguir boas práticas RESTful, utilizando rotas claras e retornando respostas JSON.Utilize as versões mais novas e estáveis das bibliotecas e frameworks Go para web (ex: `net/http`, `gorilla/mux`).
- **Frontend Desacoplado:** Framework JS simples (ex: Vue.js, Svelte ou HTML/Vanilla JS com Fetch API).
- **Banco de Dados:** PostgreSQL.
- **Ambiente Local e Deployment:** Docker e Docker Compose.
- **CI/CD:** GitHub Actions.

# Casos de Uso
- O usuário pode listar todas as tarefas cadastradas.
- O usuário pode adicionar uma nova tarefa informando o nome/descrição da atividade.
- O usuário pode marcar uma tarefa como concluída ou desfazer a conclusão.
- O usuário pode editar uma tarefa da lista.
- O usuário pode excluir uma tarefa da lista.

# Design
Para suportar as diferentes abordagens de aula, o repositório utilizará uma arquitetura de Monorepo com a seguinte estrutura de diretórios:
- `/monolito`: Conterá o código em Go que atua tanto como backend conectando no banco, quanto retornando o HTML renderizado.
- `/desacoplado`: Conterá duas subpastas: `/backend` (API em Go) e `/frontend` (aplicação JS simples).
- Na raiz, scripts comuns, o Dockerfile, o docker-compose para rodar localmente e testes.

**Requisitos Técnicos e de Segurança:**
- O aplicativo deve utilizar Variaveis de Ambiente para todas as configurações (Host do BD, porta, user, password).
- **Extrema importância:** Não imprima valores de segredos, não copie o conteúdo de arquivos de segredos para o repositório e não inclua segredos em arquivos de código-fonte, arquivos `.env`, Dockerfiles, arquivos compose, logs ou exemplos de `README.md`. Use variáveis genéricas como placeholders.
- A aplicação principal escutará requisições e servirá os recursos necessários, pronta para rodar solta ou via contêiner.

**Testes e Documentação:**
- Cada funcionalidade importante (criar, listar, deletar, atualizar) deve ter um teste unitário associado, com boa cobertura.
- Deve existir um arquivo `README.md` detalhado que explique o que o projeto faz, os requisitos de banco de dados, como configurar o ambiente (via vars) e como executá-lo localmente via fonte e via Docker.
- O `README.md` deve ser atualizado sempre que um aspecto importante da configuração for adicionado ou alterado.

**CI/CD:**
- Crie um workflow de integração e entrega contínua (CI/CD) para o GitHub Actions, similar à referência: `https://github.com/elton-bt/nova/blob/main/.github/workflows/publish-docker.yml`.
- Adicione também passos de lint e scanner da imagem com o Trivy

# Novas ideias
Agente, implemente os requisitos acima e leve em consideração os seguintes recursos extras e boas práticas para um serviço em produção/didático:
1. **Health Check Endpoint:** Crie um endpoint (ex: `/health`) para permitir que orquestradores e proxies verifiquem se a aplicação e o banco estão no ar.
2. **Graceful Shutdown:** Implemente o desligamento gracioso na aplicação Go, garantindo que conexões com o banco sejam fechadas corretamente ao receber sinais do SO (ex: SIGTERM).
3. **Multi-stage build no Docker:** Use compilação em múltiplos estágios nos `Dockerfile`s da aplicação Go para gerar imagens finais enxutas e seguras.
4. **Tratamento de erros genéricos:** Em caso de perda de conexão com o banco, retorne um erro tratável na interface, avisando ao instrutor e aluno sem expor as entranhas do erro (Stack Trace).