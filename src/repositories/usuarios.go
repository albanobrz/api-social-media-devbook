package repositories

import (
	"api/internal/domain/entities"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

// NewUsersRepository cria um repositório de usuários
func NewUsersRepository(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// Create insere um usuário no banco de dados
func (repositorio usuarios) Create(user entities.User) (uint64, error) {
	statement, erro := repositorio.db.Prepare("insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Get traz todos os usuários que atendem um filtro de nome ou nick
func (repositorio usuarios) Get(nameOrNick string) ([]entities.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) // %nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?", nameOrNick, nameOrNick,
	)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []entities.User

	for linhas.Next() {
		var user entities.User

		if erro = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, user)
	}

	return usuarios, nil
}

// Trás um usuário do banco de dados
func (repositorio usuarios) GetByID(id uint64) (entities.User, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where id = ?", id,
	)
	if erro != nil {
		return entities.User{}, erro
	}
	defer linhas.Close()

	var user entities.User

	if linhas.Next() {
		if erro = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return entities.User{}, erro
		}
	}
	return user, nil
}

// Update altera as informações de um usuário no banco de dados
func (repositorio usuarios) Update(id uint64, user entities.User) error {
	statement, erro := repositorio.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(user.Name, user.Nick, user.Email, id); erro != nil {
		return erro
	}
	return nil
}

// Deleta um usuario do banco de dados
func (repositorio usuarios) Delete(id uint64) error {
	statement, erro := repositorio.db.Query(
		"delete from usuarios where id = ?", id,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	return nil
}

// Busca um usuário por email e retorna seu id e senha com hash
func (repositorio usuarios) BuscarPorEmail(email string) (entities.User, error) {
	linha, erro := repositorio.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return entities.User{}, erro
	}
	defer linha.Close()

	var user entities.User

	if linha.Next() {
		if erro = linha.Scan(&user.ID, &user.Password); erro != nil {
			return entities.User{}, erro
		}
	}

	return user, nil
}

// Follow permite que um usuário siga outro
func (repositorio usuarios) Follow(usuarioID, seguidorID uint64) error {
	// o ignore, não permite que caso já haja a chave primária (combinação dos ids), não dê erro... ele simplesmente ignora
	statement, erro := repositorio.db.Prepare("insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}
	return nil
}

// Deixar de seguir um usuário
func (repositorio usuarios) StopFollow(userID, followedID uint64) error {
	// o ignore, não permite que caso já haja a chave primária (combinação dos ids), não dê erro... ele simplesmente ignora
	statement, erro := repositorio.db.Prepare("delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(userID, followedID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio usuarios) GetFollowers(userID uint64) ([]entities.User, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.seguidor_id where s.usuario_id = ?
	`, userID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var users []entities.User
	for linhas.Next() {
		var user entities.User

		if erro = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}
	return users, nil
}

func (repositorio usuarios) GetFollowing(userID uint64) ([]entities.User, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.usuario_id where s.seguidor_id = ?
	`, userID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var users []entities.User
	for linhas.Next() {
		var user entities.User

		if erro = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}
	return users, nil
}

func (repositorio usuarios) GetPassword(userID uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", userID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var user entities.User

	if linha.Next() {
		if erro = linha.Scan(&user.Password); erro != nil {
			return "", erro
		}
	}
	return user.Password, nil
}

func (repositorio usuarios) UpdatePassword(userID uint64, password string) error {
	statement, erro := repositorio.db.Prepare("update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(password, userID); erro != nil {
		return erro
	}
	return nil
}

// o repositório simplesmente recebe um dado e altera o banco
