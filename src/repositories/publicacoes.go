package repositories

import (
	"api/internal/domain/entities"
	"database/sql"
)

type Posts struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) *Posts {
	return &Posts{db}
}

func (repositorio Posts) Create(post entities.Post) (uint64, error) {
	statement, erro := repositorio.db.Prepare("insert into publicacoes (titulo, conteudo, autor_id) values (?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(post.Title, post.Content, post.AuthorID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

func (repositorio Posts) GetWithID(postID uint64) (entities.Post, error) {
	linha, erro := repositorio.db.Query(`
		select p.*, u.nick from
		publicacoes p inner join usuarios u
		on u.id = p.autor_id where p.id = ?
	`, postID)
	if erro != nil {
		return entities.Post{}, erro
	}

	defer linha.Close()

	var post entities.Post

	if linha.Next() {
		if erro = linha.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return entities.Post{}, erro
		}
	}

	return post, nil
}

func (repositorio Posts) Get(userID uint64) ([]entities.Post, error) {
	linhas, erro := repositorio.db.Query(`
		select distinct p.*, u.nick from publicacoes p
		inner join usuarios u on u.id = p.autor_id
		inner join seguidores s on p.autor_id = s.usuario_id 
		where u.id = ? or s.seguidor_id = ?
		order by 1 desc
	`, userID, userID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var posts []entities.Post

	for linhas.Next() {
		var post entities.Post
		if erro = linhas.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (repositorio Posts) Update(postID uint64, post entities.Post) error {
	statement, erro := repositorio.db.Prepare("update publicacoes set titulo = ?, conteudo = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(post.Title, post.Content, postID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Posts) Delete(postID uint64) error {
	statement, erro := repositorio.db.Query(
		"delete from publicacoes where id = ?", postID,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	return nil
}

func (repositorio Posts) GetByUser(userID uint64) ([]entities.Post, error) {
	linhas, erro := repositorio.db.Query(`
		select p.*, u.nick from publicacoes p
		join usuarios u on u.id = p.autor_id
		where p.autor_id = ?
	`, userID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var posts []entities.Post

	for linhas.Next() {
		var post entities.Post
		if erro = linhas.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (repositorio Posts) Like(postID uint64) error {
	statement, erro := repositorio.db.Prepare("update publicacoes set curtidas = curtidas + 1 where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(postID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Posts) Dislike(postID uint64) error {
	statement, erro := repositorio.db.Prepare(`
		update publicacoes set curtidas = 
		CASE 
			WHEN curtidas > 0 THEN curtidas - 1 
			ELSE 0 
		END 
		where id = ?
	 `)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(postID); erro != nil {
		return erro
	}

	return nil
}
