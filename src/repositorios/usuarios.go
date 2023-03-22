package repositorios

import (
	"api/internal/domain/entities"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// Criar insere um usuário no banco de dados
func (repositorio usuarios) Criar(user entities.User) (uint64, error) {
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

// Buscar traz todos os usuários que atendem um filtro de nome ou nick
func (repositorio usuarios) Buscar(nomeOuNick string) ([]entities.User, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?", nomeOuNick, nomeOuNick,
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
func (repositorio usuarios) BuscarPorId(id uint64) (entities.User, error) {
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

// Atualizar altera as informações de um usuário no banco de dados
func (repositorio usuarios) Atualizar(id uint64, user entities.User) error {
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
func (repositorio usuarios) Deletar(id uint64) error {
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

// Seguir permite que um usuário siga outro
func (repositorio usuarios) Seguir(usuarioID, seguidorID uint64) error {
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
func (repositorio usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	// o ignore, não permite que caso já haja a chave primária (combinação dos ids), não dê erro... ele simplesmente ignora
	statement, erro := repositorio.db.Prepare("delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio usuarios) BuscarSeguidores(usuarioID uint64) ([]entities.User, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.seguidor_id where s.usuario_id = ?
	`, usuarioID)
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

func (repositorio usuarios) BuscarSeguindo(usuarioID uint64) ([]entities.User, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.usuario_id where s.seguidor_id = ?
	`, usuarioID)
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

func (repositorio usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", usuarioID)
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

func (repositorio usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}
	return nil
}

// o repositório simplesmente recebe um dado e altera o banco
