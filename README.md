# Video Sync

Este aplicativo permite que os usuários criem salas virtuais onde até duas pessoas podem compartilhar e assistir vídeos do YouTube de forma sincronizada. O app inclui recursos como controle de reprodução sincronizada (play, pause, avançar), troca de vídeo, um chat integrado e notificações em tempo real usando Server-Sent Events (SSE).

## Funcionalidades

- **Salas Virtuais**: Crie salas para até dois usuários.
- **Reprodução Sincronizada de Vídeos**: Assista vídeos do YouTube em sincronia com controles compartilhados (play, pause, avançar).
- **Troca de Vídeo**: Altere o vídeo que está sendo assistido em tempo real.
- **Chat Integrado**: Comunique-se com o outro usuário na sala por meio de um chat embutido.
- **Notificações em Tempo Real**: Receba notificações instantâneas usando Server-Sent Events (SSE).

## Como Funciona

1. **Criar uma Sala**: Um usuário pode criar uma nova sala fornecendo um nome e uma senha.
2. **Entrar em uma Sala**: Outro usuário pode entrar na sala usando o nome da sala e a senha.
3. **Reprodução Sincronizada de Vídeos**: Ambos os usuários podem assistir ao mesmo vídeo do YouTube com controles sincronizados.
4. **Chat**: Use o chat integrado para se comunicar com o outro usuário na sala.
5. **Notificações**: Notificações em tempo real mantêm os usuários informados sobre as atividades na sala.

## Como Usar

1. **Criar uma Sala**: Navegue até a página principal e crie uma nova sala inserindo um nome e uma senha.
2. **Entrar em uma Sala**: Compartilhe o link da sala com outro usuário, que pode entrar fornecendo o nome da sala e a senha.
3. **Assistir Vídeos**: Quando ambos os usuários estiverem na sala, eles podem começar a assistir vídeos juntos com controles sincronizados.
4. **Chat**: Use o recurso de chat para se comunicar com o outro usuário na sala.

## Requisitos

- Linguagem de programação Go
- Dependências listadas no arquivo `go.mod`

## Instalação

1. Clone o repositório.
2. Instale as dependências usando `go mod tidy`.
3. Execute o aplicativo com `go run main.go`.

## Contribuição

Contribuições são bem-vindas! Abra uma issue ou envie um pull request para melhorias ou correções de bugs.

## Licença

Este projeto está licenciado sob a Licença MIT. Consulte o arquivo LICENSE para mais detalhes.
