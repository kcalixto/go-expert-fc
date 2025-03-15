# Documentação explicando como rodar o projeto em ambiente dev.

1. Clone o repositório
```bash
git clone github.com/kcalixto/go-expert-fc
```

2. Abra a pasta contendo o desafio
```bash
code go-expert-fc/challenges/auction
```

3. Execute o projeto
```bash
docker-compose up -d
```

4. Crie um usuário no mongodb através do mongosh com o ID: B9821642-EECA-4D66-9921-524352355089
```bash
docker exec -it mongodb mongosh --username admin --password admin --authenticationDatabase admin --eval 'db = db.getSiblingDB("auctions"); db.users.insertOne({ "_id": "B9821642-EECA-4D66-9921-524352355089", "name": "test-user" })'
```

5. Realize os testes
```bash
# Cria um leilão
curl --location 'http://localhost:8080/auction' \
--data '{
    "product_name": "My product",
    "category": "product",
    "description": "Some new product",
    "condition": 1
}'

# Buscar o leião criado
curl --location 'http://localhost:8080/auction?status=0'

# Cria uma oferta
curl --location 'http://localhost:8080/bid' \
--data '{
    "user_id": "B9821642-EECA-4D66-9921-524352355089",
    "auction_id": "< coloque aqui o id do seu novo leilão >",
    "amount": 1500
}'

# Verifica status do leilão
curl --location 'http://localhost:8080/auction/winner/< coloque aqui o id do seu novo leilão >'
```

6. Para finalizar o projeto
```bash
docker-compose down
```

7. O tempo de expiração do leilão é de 1 minuto, portanto, após criar o seu leilão, espere 1 minuto e verifique se o status foi alterado para 1 (fechado).

-------------------------------------------------

# descrição do desafio
Objetivo: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

- [X] Clone o seguinte repositório: clique para acessar o repositório.

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria:
- [X] adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

Você deverá desenvolver:

- [X] Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- [X] Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- [x] Um teste para validar se o fechamento está acontecendo de forma automatizada;

Dicas:

Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
[ ] Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;

Entrega:

O código-fonte completo da implementação.
Documentação explicando como rodar o projeto em ambiente dev.
Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.
