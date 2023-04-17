package repositories

import (
	"api/internal/domain/entities"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

// NewUsersRepository cria um repositório de usuários
func NewUsersRepository(db *sql.DB) *users {
	return &users{db}
}

// Create insere um usuário no banco de dados
func (repository users) Create(user entities.User) (uint64, error) {
	statement, err := repository.db.Prepare("insert into users (name, nick, email, password) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertID), nil
}

// Get traz todos os usuários que atendem um filtro de nome ou nick
func (repository users) Get(nameOrNick string) ([]entities.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) // %nomeOuNick%

	rows, err := repository.db.Query(
		"select id, name, nick, email, createdAt from users where name LIKE ? or nick LIKE ?", nameOrNick, nameOrNick,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var user entities.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Trás um usuário do banco de dados
func (repository users) GetByID(id uint64) (entities.User, error) {
	rows, err := repository.db.Query(
		"select id, name, nick, email, createdAt from users where id = ?", id,
	)
	if err != nil {
		return entities.User{}, err
	}
	defer rows.Close()

	var user entities.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return entities.User{}, err
		}
	}
	return user, nil
}

// Update altera as informações de um usuário no banco de dados
func (repository users) Update(id uint64, user entities.User) error {
	statement, err := repository.db.Prepare("update users set name = ?, nick = ?, email = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nick, user.Email, id); err != nil {
		return err
	}
	return nil
}

// Deleta um usuario do banco de dados
func (repository users) Delete(id uint64) error {
	statement, err := repository.db.Query(
		"delete from users where id = ?", id,
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	return nil
}

// Busca um usuário por email e retorna seu id e senha com hash
func (repository users) SearchByEmail(email string) (entities.User, error) {
	row, err := repository.db.Query("select id, password from users where email = ?", email)
	if err != nil {
		return entities.User{}, err
	}
	defer row.Close()

	var user entities.User

	if row.Next() {
		if err = row.Scan(&user.ID, &user.Password); err != nil {
			return entities.User{}, err
		}
	}

	return user, nil
}

// Follow permite que um usuário siga outro
func (repository users) Follow(userID, followedID uint64) error {
	// o ignore, não permite que caso já haja a chave primária (combinação dos ids), não dê err... ele simplesmente ignora
	statement, err := repository.db.Prepare("insert ignore into followers (user_id, follower_id) values (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(userID, followedID); err != nil {
		return err
	}
	return nil
}

// Deixar de seguir um usuário
func (repository users) StopFollow(userID, followedID uint64) error {
	// o ignore, não permite que caso já haja a chave primária (combinação dos ids), não dê err... ele simplesmente ignora
	statement, err := repository.db.Prepare("delete from followers where user_id = ? and follower_id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(userID, followedID); err != nil {
		return err
	}
	return nil
}

func (repository users) GetFollowers(userID uint64) ([]entities.User, error) {
	rows, err := repository.db.Query(`
		select u.id, u.name, u.nick, u.email, u.createdAt
		from users u inner join followers s on u.id = s.follower_id where s.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repository users) GetFollowing(userID uint64) ([]entities.User, error) {
	rows, err := repository.db.Query(`
		select u.id, u.name, u.nick, u.email, u.createdAt
		from users u inner join followers s on u.id = s.user_id where s.follower_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repository users) GetPassword(userID uint64) (string, error) {
	row, err := repository.db.Query("select password from users where id = ?", userID)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var user entities.User

	if row.Next() {
		if err = row.Scan(&user.Password); err != nil {
			return "", err
		}
	}
	return user.Password, nil
}

func (repository users) UpdatePassword(userID uint64, password string) error {
	statement, err := repository.db.Prepare("update users set password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(password, userID); err != nil {
		return err
	}
	return nil
}

// o repositório simplesmente recebe um dado e altera o banco
